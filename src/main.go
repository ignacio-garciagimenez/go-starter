package main

import (
	"github.com/bitlogic/go-startup/src/infrastructure/config"
	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()

	config.MapEndpoints(e)

	e.Logger.Fatal(e.Start(":8080"))
}
