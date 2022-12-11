package main

import (
	"io"
	"net/http"

	"github.com/dapr/go-sdk/client"
	"github.com/labstack/echo/v4"
)

const actorType = "owner"

func main() {
	// create the microservice
	app := echo.New()
	// create a dapr client
	dapr, err := client.NewClient()
	if err != nil {
		app.Logger.Fatalf("creating the dapr client", err)
	}
	// register the endpoints
	app.PUT("/:id", registerOwner(dapr))
	// start the microservice
	app.Start("127.0.0.1:3000")
}

func registerOwner(dapr client.Client) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		data, err := io.ReadAll(c.Request().Body)
		if err != nil {
			return err
		}
		req := &client.InvokeActorRequest{
			ActorType: actorType,
			Method:    "Register",
			ActorID:   c.Param("id"),
			Data:      data,
		}
		_, err = dapr.InvokeActor(ctx, req)
		if err != nil {
			return err
		}
		return c.NoContent(http.StatusAccepted)
	}
}
