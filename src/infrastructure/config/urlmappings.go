package config

import (
	"net/http"

	"github.com/bitlogic/go-startup/src/application"
	"github.com/bitlogic/go-startup/src/infrastructure/controllers"
	"github.com/bitlogic/go-startup/src/infrastructure/repositories"
	"github.com/labstack/echo/v4"
)

var productController *controllers.ProductController

func init() {
	productRepository := repositories.NewInMemoryProductRepository()
	service, _ := application.NewProductService(productRepository)
	controller, _ := controllers.NewProductController(service)

	productController = controller
}

func MapEndpoints(e *echo.Echo) {
	e.Validator = NewRequestValidator()

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})

	e.POST("/products", productController.CreateNewProduct)
}
