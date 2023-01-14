package storage

import (
	"crypto/sha256"
	"errors"
	"github.com/decentralized-hse/go-log-gossip/domain"
	"github.com/decentralized-hse/go-log-gossip/storage/types"
	"sync"
)

type InMemoryStorage struct {
	mutex *sync.Mutex
	trees map[domain.NodeId]*types.MerkleTree[*LogNodeValue]
}

type LogNodeValue struct {
	message string
}

func (l *LogNodeValue) CalculateHash() ([]byte, error) {
	h := sha256.New()
	_, err := h.Write([]byte(l.message))

	if err != nil {
		return nil, err
	}

	return h.Sum(nil), nil
}

func (l *LogNodeValue) Equals(other types.NodeValue) (bool, error) {
	return l == other, nil
}

func NewInMemoryStorage() *InMemoryStorage {
	return &InMemoryStorage{trees: make(map[domain.NodeId]*types.MerkleTree[*LogNodeValue]), mutex: new(sync.Mutex)}
}

func (storage *InMemoryStorage) Append(message string, nodeId domain.NodeId) (*domain.Log, error) {
	nodeValue := LogNodeValue{message: message}

	storage.mutex.Lock()
	tree, ok := storage.trees[nodeId]

	if !ok {
		tree = types.NewMerkleTree[*LogNodeValue]()
		storage.trees[nodeId] = tree
	}

	previousNode := tree.Leafs[len(tree.Leafs)-1]
	length, err := tree.Append(&nodeValue)

	storage.mutex.Unlock()

	if err != nil {
		return nil, err
	}

	logHash, err := nodeValue.CalculateHash()

	h := sha256.New()
	_, err = h.Write(logHash)

	if err != nil {
		return nil, err
	}

	nodeHash := h.Sum(previousNode.Hash)

	if err != nil {
		return nil, err
	}

	return &domain.Log{
		Hash:     nodeHash[:],
		NodeId:   nodeId,
		Message:  message,
		Position: length,
	}, nil
}

func (storage *InMemoryStorage) TryInsertAt(message string, nodeId domain.NodeId, position int) bool {
	storage.mutex.Lock()

	storage.mutex.Unlock()
}

func (storage *InMemoryStorage) GetNodeMessages(id domain.NodeId) ([]*domain.Log, error) {
	storage.mutex.Lock()
	defer storage.mutex.Unlock()
	tree, ok := storage.trees[id]
	if !ok {
		return nil, errors.New("node not found")
	}
	logs := make([]*domain.Log, len(tree.Leafs))
	for i := 0; i < len(tree.Leafs); i++ {
		logs[i] = &domain.Log{
			Hash:     tree.Leafs[i].Hash,
			Position: i,
			NodeId:   id,
			Message:  (*tree.Leafs[i].Value).message,
		}
	}

	return logs, nil
}
