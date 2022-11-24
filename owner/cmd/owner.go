package main

import (
	_ "embed"
	"net/http"
	"os"

	ce "github.com/cloudevents/sdk-go/v2"
	"github.com/dapr/go-sdk/client"
	"github.com/labstack/echo/v4"
	"github.com/spolab/petclinic/owner/pkg/api"
	"github.com/spolab/petclinic/owner/pkg/model"
	"go.uber.org/zap"
)

//go:embed version.txt
var version string

func main() {
	stateStoreName := os.Getenv("STATESTORE_NAME")
	pubsubName := os.Getenv("PUBSUB_NAME")
	pubsubTopic := os.Getenv("PUBSUB_TOPIC")
	//
	// Initialize logging
	logger, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}
	defer logger.Sync()
	//
	// Announces the bootstrap of the microservice
	logger.Info("Starting Owner Microservice", zap.String("version", version), zap.String("STATESTORE_NAM£", stateStoreName), zap.String("PUBSUB_NAME", pubsubName), zap.String("PUBSUBTOPIC", pubsubTopic))
	//
	// Initialize the DAPR Client
	dapr, err := client.NewClient()
	if err != nil {
		logger.Panic("Error initializing the Dapr client", zap.Error(err))
	}
	//
	// Initialize the Owner repository
	repo, err := model.NewRepository(dapr, stateStoreName, pubsubName, pubsubTopic)
	if err != nil {
		logger.Panic("Error initializing the Owner repository", zap.Error(err))
	}
	//
	// Setup the API Endpoint
	app := echo.New()
	app.HideBanner = true
	app.POST("/register", register(logger, repo))
	app.POST("/registered", registered(logger, dapr))
	if err := app.Start("localhost:3000"); err != nil {
		panic(err)
	}
}

// register registers a new owner
func register(logger *zap.Logger, repo *model.Repository) echo.HandlerFunc {
	return func(c echo.Context) error {
		logger.Debug("Received owner creation request")
		ctx := c.Request().Context()
		// Try to map the request against the expected command
		var cmd *api.RegisterOwner
		if err := c.Bind(&cmd); err != nil {
			c.String(http.StatusBadRequest, err.Error())
			return nil
		}
		// Validate the command
		if err := cmd.Validate(); err != nil {
			c.String(http.StatusBadRequest, err.Error())
		}
		// Execute the command
		// TODO - this should not happen. It should be a clerk that registers - we will change the code when we will have authentication in place
		owner := model.Owner{}
		owner.Register(ctx, cmd)
		if err := repo.Create(ctx, &owner); err != nil {
			logger.Error("Saving the owner state", zap.Error(err))
			c.String(http.StatusInternalServerError, err.Error())
		}
		logger.Debug("Owner registration complete", zap.String("id", owner.Id))
		c.NoContent(http.StatusCreated)
		return nil
	}
}

// this method listens for the ownercreated event and updates the read-only cache
func registered(logger *zap.Logger, dapr client.Client) echo.HandlerFunc {
	return func(c echo.Context) error {
		logger.Info("I am getting here")
		if event, err := ce.NewEventFromHTTPRequest(c.Request()); err == nil {
			logger.Info("Response from callback", zap.Any("event", event))
		} else {
			logger.Warn("Receiving incorrect data from the subscription", zap.Error(err))
		}
		c.NoContent(http.StatusOK)
		return nil
	}
}
