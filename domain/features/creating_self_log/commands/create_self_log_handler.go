package commands

import (
	"context"
	"github.com/decentralized-hse/go-log-gossip/storage"
)

const SelfNodeId = "self"

type CreateSelfLogHandler struct {
	storage storage.LogStorage
}

func NewCreateSelfLogHandler(storage storage.LogStorage) *CreateSelfLogHandler {
	return &CreateSelfLogHandler{storage: storage}
}

func (c *CreateSelfLogHandler) Handle(_ context.Context, command *CreateSelfLogCommand) (response *CreateSelfLogResponse, err error) {
	log, err := c.storage.Append(command.Message, SelfNodeId)
	if err != nil {
		return nil, err
	}
	response = &CreateSelfLogResponse{NewLog: log}
	return
}
