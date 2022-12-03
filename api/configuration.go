package api

import (
	"context"
	"sync"
)

type ServerConfiguration struct {
	Addr      string
	Context   context.Context
	WaitGroup *sync.WaitGroup
}

func NewServerConfiguration(addr string, context context.Context, wg *sync.WaitGroup) *ServerConfiguration {
	return &ServerConfiguration{Addr: addr, Context: context, WaitGroup: wg}
}
