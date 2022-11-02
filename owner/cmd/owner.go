package main

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func main() {
	app := echo.New()
	app.GET("/", helloWorld)
	app.Start(":3000")
}

func helloWorld(c echo.Context) error {
	c.String(http.StatusOK, "Hello Peter!")
	return nil
}
