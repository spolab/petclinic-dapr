package user

import (
	"github.com/dapr/go-sdk/actor"
	"github.com/dapr/go-sdk/client"
	"github.com/go-playground/validator/v10"
	"github.com/rs/zerolog/log"
)

const actorName = "user"

type UserActor struct {
	actor.ServerImplBase
	dapr     client.Client
	validate *validator.Validate
	broker   string
	topic    string
}

// UserActorFactory is a
func UserActorFactory(dapr client.Client, broker string, topic string) actor.Factory {
	return func() actor.Server {
		log.Debug().Str("broker", broker).Str("topic", topic).Msg("initializing actor factory")
		return &UserActor{dapr: dapr, validate: validator.New(), broker: broker, topic: topic}
	}
}

// Type returns the name that will be used to bind the actor to the URL, i.e. `/actor/<type>/method/<methodName>`
func (a *UserActor) Type() string {
	return actorName
}
