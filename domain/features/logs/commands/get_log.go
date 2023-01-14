package commands

import (
	"fmt"
	"github.com/decentralized-hse/go-log-gossip/domain"
)

type GetLogQuery struct {
	NodeId      domain.NodeId
	LogPosition int
}

func (c *GetLogQuery) String() string {
	return fmt.Sprintf("GetLogQuery{%s, %d}", c.NodeId, c.LogPosition)
}

func NewGetLogQuery(nodeId domain.NodeId, logPosition int) *GetLogQuery {
	return &GetLogQuery{NodeId: nodeId, LogPosition: logPosition}
}

type GetLogResponse struct {
	Log *domain.Log
}
