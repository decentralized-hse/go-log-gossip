package domain

import (
	"fmt"
	"time"
)

type Hash string
type NodeId string
type Id uint64

type Log struct {
	Hash    Hash
	Id      Id
	NodeId  NodeId
	Message string
	Time    time.Time
}

func (log *Log) String() string {
	return fmt.Sprintf("[%s] [%d] [%v] %s", log.Hash, log.Id, log.Time, log.Message)
}
