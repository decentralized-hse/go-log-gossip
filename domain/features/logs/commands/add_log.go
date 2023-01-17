package commands

import (
	"fmt"
	"github.com/decentralized-hse/go-log-gossip/domain"
)

type AddLogCommand struct {
	Log *domain.Log
}

func (c *AddLogCommand) String() string {
	return fmt.Sprintf("AddLogCommand{%v}", c.Log.String())
}

func NewAddLogCommand(log *domain.Log) *AddLogCommand {
	return &AddLogCommand{Log: log}
}

type AddLogResponse struct {
}
