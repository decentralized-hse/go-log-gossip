package commands

import (
	"context"
	"github.com/decentralized-hse/go-log-gossip/infra/gossip"
	"github.com/decentralized-hse/go-log-gossip/storage"
)

const SelfNodeId = "self"

type CreateSelfLogHandler struct {
	storage  storage.LogStorage
	gossiper *gossip.Gossiper
}

func NewCreateSelfLogHandler(storage storage.LogStorage, gossiper *gossip.Gossiper) *CreateSelfLogHandler {
	return &CreateSelfLogHandler{storage: storage, gossiper: gossiper}
}

func (c *CreateSelfLogHandler) Handle(_ context.Context, command *CreateSelfLogCommand) (response *CreateSelfLogResponse, err error) {
	//log, err := c.storage.Append(command.Message, SelfNodeId)
	if err != nil {
		return nil, err
	}
	_ = c.gossiper.BroadcastMessage("new", command.Message)
	response = &CreateSelfLogResponse{NewLog: nil}
	return
}
