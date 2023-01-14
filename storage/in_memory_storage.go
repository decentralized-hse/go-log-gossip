package storage

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"github.com/decentralized-hse/go-log-gossip/domain"
	"github.com/decentralized-hse/go-log-gossip/storage/types"
	"sync"
)

type InMemoryStorage struct {
	mutex  *sync.Mutex
	trees  map[domain.NodeId]*types.MerkleTree[*LogNodeValue]
	queues map[domain.NodeId]map[int]*queueRecord
}

type queueRecord struct {
	Message string
	NodeId  domain.NodeId
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

func (m *InMemoryStorage) Append(message string, nodeId domain.NodeId) (*domain.Log, error) {
	nodeValue := LogNodeValue{message: message}

	m.mutex.Lock()
	tree, ok := m.trees[nodeId]

	if !ok {
		tree = types.NewMerkleTree[*LogNodeValue]()
		m.trees[nodeId] = tree
	}

	err := tree.Append(&nodeValue)
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

func (m *InMemoryStorage) InsertAt(message string, nodeId domain.NodeId, position int) (error, int) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	tree, ok := m.trees[nodeId]
	if position == 0 {
		tree = types.NewMerkleTree[*LogNodeValue]()
		m.trees[nodeId] = tree
		return tree.Append(&LogNodeValue{message: message}), 0
	}
	if tree.LastNode.Position >= position {
		return errors.New(
			fmt.Sprintf("cannot insert message at position = %d, last inserted position = %d",
				position, tree.LastNode.Position)), -1
	}
	queue, ok := m.queues[nodeId]
	if !ok {
		queue = make(map[int]*queueRecord)
		m.queues[nodeId] = queue
	}

	for lastPosition := tree.LastNode.Position; lastPosition+1 != position; {
		record, ok := queue[lastPosition+1]
		if !ok {
			queue[position] = &queueRecord{message, nodeId}
			return errors.New("not all records are present in the queue"), lastPosition + 1
		}
		err := tree.Append(&LogNodeValue{record.Message})

		if err != nil {
			queue[position] = &queueRecord{message, nodeId}
			return err, -1
		}
	}

	err := tree.Append(&LogNodeValue{message})

	if err != nil {
		queue[position] = &queueRecord{message, nodeId}
		return err, -1
	}
	return nil, -1
}

func (m *InMemoryStorage) GetNodeLogs(id domain.NodeId) ([]*domain.Log, error) {
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

func (m *InMemoryStorage) GetNodeLog(nodeId domain.NodeId, logPosition int) (*domain.Log, error) {
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
