package domain

import (
	"errors"
	"strings"

	"github.com/google/uuid"
)

type Customer struct {
	*baseEntity[uuid.UUID]
	name string
}

func (c Customer) GetName() string {
	return c.name
}

func NewCustomer(name string) (*Customer, error) {
	trimmedName := strings.TrimSpace(name)
	if len(trimmedName) < 8 {
		return nil, errors.New("invalid name")
	}

	return &Customer{
		baseEntity: &baseEntity[uuid.UUID]{
			id: uuid.New(),
		},
		name: trimmedName,
	}, nil
}
