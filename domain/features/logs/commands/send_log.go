package commands

import (
	"fmt"
	"github.com/decentralized-hse/go-log-gossip/domain"
)

type SendLogCommand struct {
	NodeId      domain.NodeId
	LogPosition int
	SenderId    domain.NodeId
}

func (c *SendLogCommand) String() string {
	return fmt.Sprintf("SendLogCommand{%s, %d}", c.NodeId, c.LogPosition)
}

func NewSendLogCommand(senderId domain.NodeId, nodeId domain.NodeId, logPosition int) *SendLogCommand {
	return &SendLogCommand{SenderId: senderId, NodeId: nodeId, LogPosition: logPosition}
}

type SendLogResponse struct {
	Log *domain.Log
}
