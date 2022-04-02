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
	CustomerName string
}

type CreateProductCommand struct {
	ProductName string  `json:"product_name"`
	UnitPrice   float64 `json:"unit_price"`
}
