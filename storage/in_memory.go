package storage

import "github.com/decentralized-hse/go-log-gossip/log"

type InMemoryStorage struct {
	_logs []*log.Log
}

func NewInMemoryStorage(_logs []*log.Log) *InMemoryStorage {
	return &InMemoryStorage{_logs: _logs}
}

func (storage *InMemoryStorage) Append(log *log.Log) error {
	storage._logs = append(storage._logs, log)

	return nil
}

func (storage *InMemoryStorage) Load() []*log.Log {
	return storage._logs
}
