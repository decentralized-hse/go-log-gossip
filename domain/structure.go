package domain

import (
	"fmt"
)

type NodeId string

type Log struct {
	Hash     []byte
	Position int
	NodeId   NodeId
	Message  string
}

func (log *Log) String() string {
	return fmt.Sprintf("[%s] [%d] %s", log.Hash, log.Position, log.Message)
}
