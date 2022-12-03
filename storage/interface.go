package storage

import (
	"github.com/decentralized-hse/go-log-gossip/domain"
	"time"
)

type LogStorage interface {
	Append(log *domain.Log) error
	GetLast(time time.Time, nodeId domain.NodeId) (*domain.Log, error)
	Get(time time.Time, nodeId domain.NodeId) ([]*domain.Log, error)
}
