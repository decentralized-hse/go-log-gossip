package commands

import (
	"fmt"
	"github.com/decentralized-hse/go-log-gossip/domain/features/dtos"
)

type GetLogsQuery struct {
	NodeId string
}

func (c *GetLogsQuery) String() string {
	return fmt.Sprintf("GetLogsQuery{%v}", c.NodeId)
}

func NewGetLogsCommand(nodeId string) *GetLogsQuery {
	return &GetLogsQuery{NodeId: nodeId}
}

type GetLogsResponse struct {
	Logs []*dtos.LogDTO
}
