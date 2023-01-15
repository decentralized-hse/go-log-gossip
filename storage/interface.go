package storage

import "github.com/decentralized-hse/go-log-gossip/domain"

type LogStorage interface {
	Append(log string, nodeId string) (*domain.Log, error)
	InsertAt(log domain.Log, nodeId string, position int) (error, int)
	GetNodeLogs(nodeId string) ([]*domain.Log, error)
	GetNodeLog(nodeId string, logPosition int) (*domain.Log, error)
}
