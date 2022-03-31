package customer

import (
	"errors"
	"github.com/bitlogic/go-startup/src/domain"
	"github.com/google/uuid"
	"strings"
)

type Customer struct {
	domain.Entity
	id   uuid.UUID
	name string
}

func (c Customer) GetID() uuid.UUID {
	return c.id
}

func NewCustomer(name string) (*Customer, error) {
	trimmedName := strings.TrimSpace(name)
	if len(trimmedName) < 8 {
		return nil, errors.New("invalid name")
	}

	return &Customer{
		id:   uuid.New(),
		name: trimmedName,
	}, nil
}
