package server

import (
	"github.com/dapr/go-sdk/actor"
	"github.com/dapr/go-sdk/client"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
)

const (
	ErrCodeOK = iota
	ErrCodeInvalidCommand
)

type OwnerActor struct {
	actor.ServerImplBase
	logger   *zap.Logger
	dapr     client.Client
	validate *validator.Validate
	broker   string
	topic    string
}

// OwnerActorFactory is a
func OwnerActorFactory(logger *zap.Logger, dapr client.Client, broker string, topic string) actor.Factory {
	return func() actor.Server {
		logger.Info("initializing actor server", zap.Any("client", dapr), zap.String("broker", broker), zap.String("topic", topic))
		return &OwnerActor{logger: logger, dapr: dapr, validate: validator.New(), broker: broker, topic: topic}
	}
}

// Type returns the name that will be used to bind the actor to the URL, i.e. `/actor/<type>/method/<methodName>`
func (a *OwnerActor) Type() string {
	return "owner"
}
