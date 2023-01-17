package commands

import (
	"context"
	"github.com/decentralized-hse/go-log-gossip/domain/features/dtos"
	"github.com/decentralized-hse/go-log-gossip/infra/gossip"
	"github.com/decentralized-hse/go-log-gossip/storage"
)

type CreateSelfLogHandler struct {
	storage    storage.LogStorage
	gossiper   *gossip.Gossiper
	selfNodeId string
}

func NewCreateSelfLogHandler(storage storage.LogStorage, gossiper *gossip.Gossiper, selfNodeId string) *CreateSelfLogHandler {
	return &CreateSelfLogHandler{storage: storage, gossiper: gossiper, selfNodeId: selfNodeId}
}

func (c *CreateSelfLogHandler) Handle(_ context.Context, command *CreateSelfLogCommand) (*CreateSelfLogResponse, error) {
	log, err := c.storage.Append(command.Message, c.selfNodeId)
	if err != nil {
		return nil, err
	}

	dto := dtos.NewLogDTO(log)
	err = c.gossiper.BroadcastMessage(gossip.Push, dto)

	return &CreateSelfLogResponse{NewLog: log}, err
}
