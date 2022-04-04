package application

import (
	"github.com/google/uuid"
)

type CreateCartCommand struct {
	CustomerId uuid.UUID `json:"customer_id" validate:"required"`
}

type AddItemToCartCommand struct {
	CartId    uuid.UUID
	ProductId uuid.UUID
	Quantity  int
}

type CreateCustomerCommand struct {
	CustomerName string `json:"customer_name" validate:"required,gte=8"`
}

type CreateProductCommand struct {
	ProductName string  `json:"product_name" validate:"required,gte=10"`
	UnitPrice   float64 `json:"unit_price" validate:"required,gt=0"`
}
