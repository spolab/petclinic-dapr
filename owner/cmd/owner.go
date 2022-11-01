package main

import "github.com/labstack/echo/v4"

func main() {
	app := echo.New()
	app.Start("localhost:3000")
}
