package commands

import (
	"fmt"
	"github.com/decentralized-hse/go-log-gossip/domain"
)

type SendLogCommand struct {
	NodeId      string
	LogPosition int
	SenderId    string
}

func (c *SendLogCommand) String() string {
	return fmt.Sprintf("SendLogCommand{%s, %d}", c.NodeId, c.LogPosition)
}

func NewSendLogCommand(senderId string, nodeId string, logPosition int) *SendLogCommand {
	return &SendLogCommand{SenderId: senderId, NodeId: nodeId, LogPosition: logPosition}
}

type SendLogResponse struct {
	Log *domain.Log
}
