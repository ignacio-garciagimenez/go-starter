package application

import "github.com/google/uuid"

type CreateCartCommand struct {
	CustomerId uuid.UUID
}

type AddItemToCartCommand struct {
	CartId    uuid.UUID
	ProductId uuid.UUID
	Quantity  int
}

type CreateCustomerCommand struct {
	CustomerName string
}

type CreateCustomerResult struct {
	Id           uuid.UUID
	CustomerName string
}

type CreateProductCommand struct {
	ProductName string
	UnitPrice   float64
}

type CreateProductResult struct {
	Id          uuid.UUID
	ProductName string
	UnitPrice   float64
}
