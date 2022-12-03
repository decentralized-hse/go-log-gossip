package storage

import "github.com/decentralized-hse/go-log-gossip/domain"

type InMemoryStorage struct {
	_logs []*domain.Log
}

func NewInMemoryStorage(_logs []*domain.Log) *InMemoryStorage {
	return &InMemoryStorage{_logs: _logs}
}

func (storage *InMemoryStorage) Append(log *domain.Log) error {
	storage._logs = append(storage._logs, log)

	return nil
}

func (storage *InMemoryStorage) Load() []*domain.Log {
	return storage._logs
}
