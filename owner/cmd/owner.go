package main

import (
	_ "embed"
	"log"
	"net/http"

	"github.com/dapr/go-sdk/client"
	"github.com/labstack/echo/v4"
	"github.com/spolab/petclinic/owner/pkg/api/command"
)

//go:embed version.txt
var version string

func main() {
	log.Printf("Starting owner microservice %s\n", version)
	dapr, err := client.NewClient()
	if err != nil {
		panic(err)
	}
	app := echo.New()
	app.HideBanner = true
	app.POST("/register/:id", register(dapr))
	if err := app.Start(":3000"); err != nil {
		panic(err)
	}
}

// register registers a new owner
func register(dapr client.Client) echo.HandlerFunc {
	return func(c echo.Context) error {
		log.Printf("Processing registration for ID %s", c.Param("id"))
		var cmd command.RegisterOwner
		if err := c.Bind(&cmd); err != nil {
			c.String(http.StatusBadRequest, err.Error())
			return nil
		}
		c.JSON(http.StatusOK, &cmd)
		return nil
	}
}
