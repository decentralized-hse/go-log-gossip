package storage

import (
	"crypto/sha256"
	"github.com/decentralized-hse/go-log-gossip/domain"
	"github.com/decentralized-hse/go-log-gossip/storage/types"
	"sync"
)

type InMemoryStorage struct {
	trees map[domain.NodeId]*types.MerkleTree
	mutex *sync.Mutex
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

func NewInMemoryStorage(config FileStorageConfiguration) *InMemoryStorage {
	return &InMemoryStorage{trees: make(map[domain.NodeId]*types.MerkleTree), mutex: new(sync.Mutex)}
}

func (storage *InMemoryStorage) Append(message string, nodeId domain.NodeId) (*domain.Log, error) {
	nodeValue := LogNodeValue{message: message}

	tree, ok := storage.trees[nodeId]

	if !ok {
		tree = types.NewMerkleTree()
		storage.trees[nodeId] = tree
	}

	err := tree.Append(&nodeValue)

	if err != nil {
		return nil, err
	}

	nodeHash, err := nodeValue.CalculateHash()

	if err != nil {
		return nil, err
	}

	return &domain.Log{
		Hash:    domain.Hash(nodeHash[:]),
		NodeId:  nodeId,
		Message: message,
	}, nil
}
