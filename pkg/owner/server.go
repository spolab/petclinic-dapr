package owner

import (
	"github.com/dapr/go-sdk/actor"
	"github.com/dapr/go-sdk/client"
	"github.com/go-playground/validator/v10"
	"github.com/rs/zerolog/log"
)

const (
	ErrCodeOK = iota
	ErrCodeInvalidCommand
)

type OwnerActor struct {
	actor.ServerImplBase
	dapr     client.Client
	validate *validator.Validate
	broker   string
	topic    string
}

// OwnerActorFactory is a
func OwnerActorFactory(dapr client.Client, broker string, topic string) actor.Factory {
	return func() actor.Server {
		log.Info().Interface("client", dapr).Str("broker", broker).Str("topic", topic).Msg("initializing actor server")
		return &OwnerActor{dapr: dapr, validate: validator.New(), broker: broker, topic: topic}
	}
}

// Type returns the name that will be used to bind the actor to the URL, i.e. `/actor/<type>/method/<methodName>`
func (a *OwnerActor) Type() string {
	return "owner"
}
