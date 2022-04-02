package domain

import "github.com/google/uuid"

type CartCreated struct {
	DomainEvent
	CartId     uuid.UUID
	CustomerId uuid.UUID
}

type ItemAddedToCart struct {
	DomainEvent
	CartId    uuid.UUID
	ProductId uuid.UUID
	Quantity  int
}

type CustomerCreated struct {
	CustomerId   uuid.UUID
	CustomerName string
}

type ProductCreated struct {
	ProductId        uuid.UUID
	ProductName      string
	ProductUnitPrice float64
}
