package main

import (
	"github.com/bitlogic/go-startup/src/infrastructure/controllers"
	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()

	controllers.MapEndpoints(e)

	e.Logger.Fatal(e.Start(":8080"))
}
