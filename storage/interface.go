package storage

import "github.com/decentralized-hse/go-log-gossip/domain"

type LogStorage interface {
	Append(log string, nodeId domain.NodeId) (*domain.Log, error)
	TryInsertAt(log string, nodeId domain.NodeId, position int) bool
	GetNodeLogs(nodeId domain.NodeId) ([]*string, error)
}
