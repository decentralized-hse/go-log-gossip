package storage

import (
	"crypto/sha256"
	"fmt"
	"github.com/decentralized-hse/go-log-gossip/domain"
	"github.com/decentralized-hse/go-log-gossip/storage/types"
	"path"
	"sync"
)

type FileStorageConfiguration struct {
	rootPath string
}

type FileStorage struct {
	config FileStorageConfiguration
	mutex  *sync.Mutex
}

func NewFileStorage(config FileStorageConfiguration) *FileStorage {
	return &FileStorage{config: config, mutex: new(sync.Mutex)}
}

func (storage *FileStorage) Append(log string, nodeId domain.NodeId) (*domain.Log, error) {
	logsFilePath := storage.getFilePath(nodeId, "log")
	logsFile, err := types.OpenAppendOnlyFile(logsFilePath)
	defer logsFile.Close()

	if err != nil {
		return nil, err
	}

	indexFilePath := storage.getFilePath(nodeId, "index")
	indexFile, err := types.OpenAppendOnlyFile(indexFilePath)
	defer indexFile.Close()

	if err != nil {
		return nil, err
	}

	storage.mutex.TryLock()
	defer storage.mutex.Unlock()

	logFileSize, err := logsFile.Size()

	if err != nil {
		return nil, err
	}

	_, err = indexFile.Append(fmt.Sprintf("%016x", logFileSize))
	if err != nil {
		return nil, err
	}

	_, err = logsFile.AppendLine(log)
	if err != nil {
		return nil, err
	}

	err = storage.updateHashTree(log)

	if err != nil {
		return nil, err
	}

	return &domain.Log{
		Hash:    "",
		Id:      0,
		NodeId:  nodeId,
		Message: log,
	}, nil
}

func (storage *FileStorage) getFilePath(nodeId domain.NodeId, fileName string) string {
	return path.Join(storage.config.rootPath, fmt.Sprintf("%s", nodeId), fileName)
}

func (storage *FileStorage) updateHashTree(data string) error {
	hashTreeRoot := path.Join(storage.config.rootPath, "hash")

	bytes := []byte(data)
	layerNumber := 0

	for {
		hash := sha256.Sum256(bytes)
		hashFile, err := types.OpenAppendOnlyFile(path.Join(hashTreeRoot, fmt.Sprintf("%d.hash", layerNumber)))
		if err != nil {
			panic(err)
		}
		_, err = hashFile.AppendBytes(hash[:])

		if err != nil {
			panic(err)
		}

		hashFileLength, err := hashFile.Size()

		if err != nil {
			panic(err)
		}

		layerHashesCount := hashFileLength / sha256.Size

		if layerHashesCount%2 == 1 {
			break
		}

		_, err = hashFile.ReadAt(bytes, hashFileLength-sha256.Size*2)

		hashFile.Close()
	}

	return nil
}
