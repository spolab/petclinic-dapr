package main

import (
	"net/http"

	"github.com/dapr/go-sdk/client"
	"github.com/labstack/echo/v4"
)

func main() {
	dapr, err := client.NewClient()
	if err != nil {
		panic(err)
	}
	app := echo.New()
	app.POST("/", helloWorld(dapr))
	app.Start("127.0.0.1:3000")
}

func helloWorld(dapr client.Client) echo.HandlerFunc {
	return func(c echo.Context) error {
		c.String(http.StatusOK, "Hello world")
		return nil
	}
}
