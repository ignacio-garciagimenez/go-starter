package application

import "github.com/google/uuid"

type CreateCartCommand struct {
	CustomerId uuid.UUID
}
