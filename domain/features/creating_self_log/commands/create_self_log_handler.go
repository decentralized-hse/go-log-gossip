package commands

import (
	"context"
	"github.com/decentralized-hse/go-log-gossip/storage"
	"log"
)

type CreateSelfLogHandler struct {
	storage storage.LogStorage
}

func NewCreateSelfLogHandler(storage storage.LogStorage) *CreateSelfLogHandler {
	return &CreateSelfLogHandler{storage: storage}
}

func (c *CreateSelfLogHandler) Handle(_ context.Context, command *CreateSelfLogCommand) (response *CreateSelfLogResponse, err error) {
	log.Printf("Handling create self log handler, %v", command)

	response = &CreateSelfLogResponse{NewLog: nil}
	return
}
