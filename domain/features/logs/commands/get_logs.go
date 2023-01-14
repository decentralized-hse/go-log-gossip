package commands

import (
	"fmt"
	"github.com/decentralized-hse/go-log-gossip/domain"
)

type GetLogsQuery struct {
	NodeId domain.NodeId
}

func (c *GetLogsQuery) String() string {
	return fmt.Sprintf("GetLogsQuery{%v}", c.NodeId)
}

func NewGetLogsCommand(nodeId domain.NodeId) *GetLogsQuery {
	return &GetLogsQuery{NodeId: nodeId}
}

type GetLogsResponse struct {
	Logs []*domain.Log
}
