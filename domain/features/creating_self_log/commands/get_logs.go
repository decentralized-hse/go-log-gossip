package commands

import (
	"fmt"
	"github.com/decentralized-hse/go-log-gossip/domain"
)

type GetLogsCommand struct {
	NodeId domain.NodeId
}

func (c *GetLogsCommand) String() string {
	return fmt.Sprintf("GetLogsCommand{%v}", c.NodeId)
}

func NewGetLogsCommand(nodeId domain.NodeId) *GetLogsCommand {
	return &GetLogsCommand{NodeId: nodeId}
}

type GetLogsResponse struct {
	Logs []*domain.Log
}
