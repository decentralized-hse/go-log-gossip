package commands

import (
	"fmt"
	"github.com/decentralized-hse/go-log-gossip/domain"
)

type CreateSelfLogCommand struct {
	Message string
}

func (c *CreateSelfLogCommand) String() string {
	return fmt.Sprintf("CreateSelfLogCommand{%v}", c.Message)
}

func NewCreateSelfLogCommand(message string) *CreateSelfLogCommand {
	return &CreateSelfLogCommand{Message: message}
}

type CreateSelfLogResponse struct {
	NewLog *domain.Log
}
