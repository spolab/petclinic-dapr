package main

import (
	_ "embed"
	"log"
	"net/http"
	"os"

	"github.com/dapr/go-sdk/client"
	"github.com/labstack/echo/v4"
	"github.com/spolab/petclinic/owner/pkg/api"
	"github.com/spolab/petclinic/owner/pkg/model"
)

//go:embed version.txt
var version string

func main() {
	stateStoreName := os.Getenv("STATESTORE_NAME")
	pubsubName := os.Getenv("PUBSUB_NAME")
	pubsubTopic := os.Getenv("PUBSUB_TOPIC")
	log.Printf("Launch parameters: STATESTORE_NAME='%s', PUBSUB_NAME='%s', PUBSUBTOPIC='%s'", stateStoreName, pubsubName, pubsubTopic)
	log.Printf("Starting owner microservice %s\n", version)
	// Initialize the DAPR Client
	dapr, err := client.NewClient()
	if err != nil {
		panic(err)
	}
	// Initialize the Owner repository
	repo, err := model.NewRepository(dapr, stateStoreName, pubsubName, pubsubTopic)
	if err != nil {
		panic(err)
	}
	// Setup the API Endpoint
	app := echo.New()
	app.HideBanner = true
	app.POST("/register", register(repo))
	if err := app.Start(":3000"); err != nil {
		panic(err)
	}
}

// register registers a new owner
func register(repo *model.Repository) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		var cmd *api.RegisterOwner
		if err := c.Bind(&cmd); err != nil {
			c.String(http.StatusBadRequest, err.Error())
			return nil
		}
		if err := cmd.Validate(); err != nil {
			c.String(http.StatusBadRequest, err.Error())
		}
		owner := model.Owner{}
		owner.Register(ctx, cmd)
		if err := repo.Save(ctx, &owner, client.WithConcurrency(client.StateConcurrencyFirstWrite)); err != nil {
			c.String(http.StatusInternalServerError, err.Error())
		}
		c.NoContent(http.StatusCreated)
		return nil
	}
}
