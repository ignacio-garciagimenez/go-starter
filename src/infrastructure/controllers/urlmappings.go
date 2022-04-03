package controllers

import (
	"net/http"

	"github.com/bitlogic/go-startup/src/application"
	"github.com/bitlogic/go-startup/src/infrastructure/repositories"
	"github.com/labstack/echo/v4"
)

var productController *ProductController

func init() {
	productRepository := repositories.NewInMemoryProductRepository()
	service, _ := application.NewProductService(productRepository)
	controller, _ := NewProductController(service)

	productController = controller
}

func MapEndpoints(e *echo.Echo) {
	e.Validator = NewRequestValidator()
	//e.HTTPErrorHandler = CustomValidationErrorHanlder(e.DefaultHTTPErrorHandler)

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})

	e.POST("/products", productController.CreateNewProduct)
}
