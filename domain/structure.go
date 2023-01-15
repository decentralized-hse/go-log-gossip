package domain

import (
	"fmt"
)

type Log struct {
	Hash     []byte
	Position int
	NodeId   string
	Message  string
}

func (log *Log) String() string {
	return fmt.Sprintf("[%s] [%d] %s", log.Hash, log.Position, log.Message)
}
