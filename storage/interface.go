package storage

import "github.com/decentralized-hse/go-log-gossip/log"

type LogStorage interface {
	Append(log *log.Log) error
	Load() []*log.Log
}
