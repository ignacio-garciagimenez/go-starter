package application

import (
	"github.com/bitlogic/go-startup/src/domain"
)

type CreateCartCommand struct {
	CustomerId domain.CustomerId
}

type AddItemToCartCommand struct {
	CartId    domain.CartId
	ProductId domain.ProductId
	Quantity  int
}

type CreateCustomerCommand struct {
	CustomerName string `json:"customer_name" validate:"required,gte=8"`
}

type CreateProductCommand struct {
	ProductName string  `json:"product_name" validate:"required,gte=10"`
	UnitPrice   float64 `json:"unit_price" validate:"required,gt=0"`
}
