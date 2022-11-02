package main

import (
	"fmt"
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
	app.GET("/:name", helloWorld(dapr))
	app.Start(":3000")
}

func helloWorld(dapr client.Client) echo.HandlerFunc {
	return func(c echo.Context) error {
		name := c.Param("name")
		message := fmt.Sprintf("Hello, %s", name)
		c.String(http.StatusOK, message)
		return nil
	}
}
