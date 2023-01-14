package storage

import "github.com/decentralized-hse/go-log-gossip/domain"

type LogStorage interface {
	Append(log string, nodeId domain.NodeId) (*domain.Log, error)
	InsertAt(log string, nodeId domain.NodeId, position int) (error, int)
	GetNodeLogs(nodeId domain.NodeId) ([]*domain.Log, error)
	GetNodeLog(nodeId domain.NodeId, logPosition int) (*domain.Log, error)
}
