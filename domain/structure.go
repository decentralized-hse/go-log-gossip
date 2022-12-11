package domain

import (
	"fmt"
)

type Hash string
type NodeId string
type Id uint64

type Log struct {
	Hash    Hash
	Id      Id
	NodeId  NodeId
	Message string
}

func (log *Log) String() string {
	return fmt.Sprintf("[%s] [%d] %s", log.Hash, log.Id, log.Message)
}
