package config

import (
	"net/http"

	"github.com/bitlogic/go-startup/src/application"
	"github.com/bitlogic/go-startup/src/infrastructure/controllers"
	"github.com/bitlogic/go-startup/src/infrastructure/repositories"
	"github.com/labstack/echo/v4"
)

var productController *controllers.ProductController
var customerController *controllers.CustomerController
var cartController *controllers.CartController

func init() {
	productRepository := repositories.NewInMemoryProductRepository()
	productService, _ := application.NewProductService(productRepository)
	productController, _ = controllers.NewProductController(productService)

	customerRepository := repositories.NewInMemoryCustomerRepository()
	customerService, _ := application.NewCustomerService(customerRepository)
	customerController, _ = controllers.NewCustomerController(customerService)

	cartRepository := repositories.NewInMemoryCartRepository()
	cartService, _ := application.NewCartService(cartRepository, customerRepository, productRepository)
	cartController, _ = controllers.NewCartController(cartService)
}

func MapEndpoints(e *echo.Echo) {
	e.Validator = NewRequestValidator()

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})

	e.POST("/products", productController.CreateNewProduct)
	e.POST("/customers", customerController.CreateNewCustomer)
	e.POST("/carts", cartController.CreateNewCart)
	e.POST("/carts/:cartId", cartController.AddItemToCart)
}
