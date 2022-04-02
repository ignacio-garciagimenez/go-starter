package domain

import (
	"errors"
	"reflect"
	"strings"

	"github.com/google/uuid"
)

type CustomerId uuid.UUID

type Customer struct {
	*baseEntity[CustomerId]
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

	customer := &Customer{
		baseEntity: &baseEntity[CustomerId]{
			id: CustomerId(uuid.New()),
		},
		name: trimmedName,
	}

	customer.addDomainEvent(CustomerCreated{
		CustomerId:   customer.id,
		CustomerName: customer.name,
	})

	return customer, nil
}

func (c *Customer) EqualsTo(entity Entity[CustomerId]) bool {
	return reflect.TypeOf(c) == reflect.TypeOf(entity) &&
		c.GetID() == entity.GetID()
}
