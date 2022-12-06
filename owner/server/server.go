package server

import (
	"github.com/dapr/go-sdk/actor"
	"github.com/dapr/go-sdk/client"
	"go.uber.org/zap"
)

const (
	ErrCodeOK = iota
	ErrCodeInvalidCommand
)

type OwnerActor struct {
	actor.ServerImplBase
	logger *zap.Logger
	dapr   client.Client
}

func OwnerActorFactory(logger *zap.Logger, dapr client.Client) actor.Factory {
	return func() actor.Server {
		return &OwnerActor{logger: logger, dapr: dapr}
	}
}

// Type returns the name that will be used to bind the actor to the URL, i.e. `/actor/<type>/method/<methodName>`
func (a *OwnerActor) Type() string {
	return "owner"
}
