package storage

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"github.com/decentralized-hse/go-log-gossip/domain"
	"github.com/decentralized-hse/go-log-gossip/storage/types"
	"log"
	"sync"
)

type InMemoryStorage struct {
	mutex  *sync.Mutex
	trees  map[string]*types.MerkleTree[*LogNodeValue]
	queues map[string]map[int]*queueRecord
}

type queueRecord struct {
	Message string
	NodeId  string
	Hash    []byte
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
	return &InMemoryStorage{
		trees:  make(map[string]*types.MerkleTree[*LogNodeValue]),
		queues: make(map[string]map[int]*queueRecord),
		mutex:  new(sync.Mutex)}
}

func (m *InMemoryStorage) Append(message string, nodeId string) (*domain.Log, error) {
	nodeValue := LogNodeValue{message: message}
	_, _ = nodeValue.CalculateHash()
	m.mutex.Lock()
	tree, ok := m.trees[nodeId]

	if !ok {
		tree = types.NewMerkleTree[*LogNodeValue]()
		m.trees[nodeId] = tree
	}

	err := tree.Append(&nodeValue)
	log.Printf("Current tree root hash = %v", tree.Root.Hash)
	lastNode := tree.LastNode

	m.mutex.Unlock()

	if err != nil {
		return nil, err
	}

	return &domain.Log{
		Hash:     lastNode.Hash,
		NodeId:   nodeId,
		Message:  message,
		Position: lastNode.Position,
	}, nil
}

func (m *InMemoryStorage) InsertAt(logRecord domain.Log, nodeId string, position int) (error, int) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	tree, ok := m.trees[nodeId]
	if !ok {
		tree = types.NewMerkleTree[*LogNodeValue]()
		m.trees[nodeId] = tree
	}

	if position == 0 && len(tree.Leafs) == 0 {
		logNodeValue := LogNodeValue{message: logRecord.Message}
		hash, _ := logNodeValue.CalculateHash()
		if !hashEquals(hash, logRecord.Hash) {
			return errors.New(fmt.Sprintf("actual hash does not equals to expected")), 0
		}
		m.trees[nodeId] = tree
		log.Printf("Log %v inserted at position %d", logRecord, position)
		return tree.Append(&logNodeValue), 1
	}
	if tree.LastNode != nil && tree.LastNode.Position >= position {
		return errors.New(
			fmt.Sprintf("cannot insert message at position = %d, last inserted position = %d",
				position, tree.LastNode.Position)), -1
	}
	queue, ok := m.queues[nodeId]
	if !ok {
		queue = make(map[int]*queueRecord)
		m.queues[nodeId] = queue
	}

	lastPosition := -1
	if tree.LastNode != nil {
		lastPosition = tree.LastNode.Position
	}

	for lastPosition+1 != position {
		record, ok := queue[lastPosition+1]
		if !ok {
			queue[position] = &queueRecord{logRecord.Message, nodeId, logRecord.Hash}
			return errors.New("not all records are present in the queue"), lastPosition + 1
		}

		hash := sha256.New()

		logNodeValue := LogNodeValue{record.Message}
		logHash, _ := logNodeValue.CalculateHash()
		hash.Write(logHash)
		hash.Write(tree.LastNode.Hash)
		h := hash.Sum(nil)
		if !hashEquals(h, record.Hash) {
			return errors.New("actual hash does not equals to expected"), lastPosition
		}
		err := tree.Append(&logNodeValue)

		if err != nil {
			queue[position] = &queueRecord{logRecord.Message, nodeId, logRecord.Hash}
			return err, -1
		}
	}

	hash := sha256.New()

	logNodeValue := LogNodeValue{message: logRecord.Message}
	logHash, _ := logNodeValue.CalculateHash()
	hash.Write(logHash)
	hash.Write(tree.LastNode.Hash)
	h := hash.Sum(nil)
	if !hashEquals(h, logRecord.Hash) {
		return errors.New("actual hash does not equals to expected"), tree.LastNode.Position
	}
	err := tree.Append(&logNodeValue)

	if err != nil {
		queue[position] = &queueRecord{logRecord.Message, nodeId, logRecord.Hash}
		return err, -1
	}
	return nil, position + 1
}

func (m *InMemoryStorage) GetNodeLogs(id string) ([]*domain.Log, error) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	tree, ok := m.trees[id]
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

func (m *InMemoryStorage) GetNodeLog(nodeId string, logPosition int) (*domain.Log, error) {
	tree, ok := m.trees[nodeId]
	if !ok {
		return nil, errors.New(fmt.Sprintf("cannot find node with id = %s", nodeId))
	}

	if len(tree.Leafs) <= logPosition {
		return nil, errors.New(fmt.Sprintf("cannot find log at position = %d", logPosition))
	}

	node := tree.Leafs[logPosition]

	if node == nil {
		return nil, errors.New(fmt.Sprintf("cannot find log at position = %d", logPosition))
	}

	return &domain.Log{
		Hash:     node.Hash,
		NodeId:   nodeId,
		Message:  (*node.Value).message,
		Position: node.Position,
	}, nil
}

func hashEquals(f []byte, s []byte) bool {
	if s == nil || f == nil {
		return false
	}

	if len(f) != len(s) {
		return false
	}

	for i := 0; i < len(f); i++ {
		if f[i] != s[i] {
			return false
		}
	}

	return true
}
