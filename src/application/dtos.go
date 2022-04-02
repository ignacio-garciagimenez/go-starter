package application

import "github.com/google/uuid"

type CartDto struct {
	Id         uuid.UUID
	CustomerId uuid.UUID
	Items      []ItemDto
}

type ItemDto struct {
	ProductId uuid.UUID
	UnitPrice float64
	Quantity  int
}

type CustomerDto struct {
	Id   uuid.UUID
	Name string
}

type ProductDto struct {
	Id        uuid.UUID
	Name      string
	UnitPrice float64
}
