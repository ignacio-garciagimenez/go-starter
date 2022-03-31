package repositories

import (
	"errors"
	"github.com/bitlogic/go-startup/src/domain/customer"
	"github.com/google/uuid"
)

type InMemoryCustomerRepository struct {
	customers map[uuid.UUID]*customer.Customer
}

func (i *InMemoryCustomerRepository) FindByID(key uuid.UUID) (*customer.Customer, error) {
	entity, found := i.customers[key]
	if !found {
		return nil, errors.New("cart not found")
	}

	return entity, nil
}

func (i *InMemoryCustomerRepository) Save(entity *customer.Customer) error {
	i.customers[entity.GetID()] = entity
	return nil
}

func NewInMemoryCustomerRepository() customer.Repository {
	return &InMemoryCustomerRepository{
		customers: map[uuid.UUID]*customer.Customer{},
	}
}
