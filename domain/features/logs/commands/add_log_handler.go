package commands

import (
	"context"
	"github.com/decentralized-hse/go-log-gossip/infra/gossip"
	"github.com/decentralized-hse/go-log-gossip/storage"
)

type AddLogHandler struct {
	storage  storage.LogStorage
	gossiper *gossip.Gossiper
}

func NewAddLogHandler(storage storage.LogStorage, gossiper *gossip.Gossiper) *AddLogHandler {
	return &AddLogHandler{storage: storage, gossiper: gossiper}
}

func (c *AddLogHandler) Handle(_ context.Context, command *AddLogCommand) (response *AddLogResponse, err error) {
	err, _ = c.storage.InsertAt(command.Log.Message, command.Log.NodeId, command.Log.Position)
	if err != nil {
		// TODO: gossiper pull record with id = lastInserted
		return nil, err
	}
	_ = c.gossiper.BroadcastMessage("sync", command.Log)
	response = &AddLogResponse{}
	return
}
