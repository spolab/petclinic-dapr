package main

import (
	"io"
	"log"
	"net/http"

	"github.com/dapr/go-sdk/client"
	"github.com/labstack/echo/v4"
)

func main() {
	dapr, err := client.NewClient()
	if err != nil {
		log.Fatal("initializing dapr client", err)
	}
	app := echo.New()
	app.PUT("/:id", registerOwner(dapr))
	app.Start("localhost:3000")
}

func registerOwner(dapr client.Client) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		id := c.Param("id")
		cmd, err := io.ReadAll(c.Request().Body)
		if err != nil {
			c.String(http.StatusBadRequest, err.Error())
			return nil
		}
		req := &client.InvokeActorRequest{
			ActorType: "owner",
			ActorID:   id,
			Method:    "Register",
			Data:      cmd,
		}
		dapr.InvokeActor(ctx, req)
		return nil
	}
}
