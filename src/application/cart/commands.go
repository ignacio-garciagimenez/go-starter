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
