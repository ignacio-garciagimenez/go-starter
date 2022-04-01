package application

import "github.com/google/uuid"

type CreateProductCommand struct {
	ProductName string
	UnitPrice   float64
}

type CreateProductResult struct {
	Id          uuid.UUID
	ProductName string
	UnitPrice   float64
}
