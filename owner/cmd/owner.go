package main

import (
	_ "embed"
	"log"
	"net/http"

	"github.com/dapr/go-sdk/client"
	"github.com/labstack/echo/v4"
	"github.com/spolab/petclinic/owner/pkg/api"
	"github.com/spolab/petclinic/owner/pkg/model"
)

//go:embed version.txt
var version string

func main() {
	log.Printf("Starting owner microservice %s\n", version)
	// Initialize the DAPR Client
	dapr, err := client.NewClient()
	if err != nil {
		panic(err)
	}
	// Initialize the Owner repository
	repo, err := model.NewRepository(dapr, "owner.petclinic")
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
		if err := c.Bind(cmd); err != nil {
			c.String(http.StatusBadRequest, err.Error())
			return nil
		}
		if err := cmd.Validate(); err != nil {
			c.String(http.StatusBadRequest, err.Error())
		}
		owner := model.Owner{}
		owner.Register(ctx, cmd)
		if err := repo.Save(ctx, &owner); err != nil {
			c.String(http.StatusInternalServerError, err.Error())
		}
		c.NoContent(http.StatusCreated)
		return nil
	}
}
