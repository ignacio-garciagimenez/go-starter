package repositories

import (
	"errors"
	"github.com/bitlogic/go-startup/src/domain"
	"github.com/bitlogic/go-startup/src/domain/customer"
	"github.com/google/uuid"
	"golang.org/x/exp/constraints"
)

type InMemoryBaseRepository[K constraints.Ordered, E domain.Entity[K]] struct {
	entities map[K]E[K]
}

type InMemoryCustomerRepository struct {
	customers map[uuid.UUID]*customer.Customer
}

func (i *InMemoryBaseRepository[K, E]) FindByID(key K) (E, K, error) {
	entity, found := i.entities[key]
	if !found {
		return nil, nil, errors.New("cart not found")
	}

	return entity, entity.GetID(), nil
}

func (i *InMemoryBaseRepository[K, E]) Save(entity E) error {
	id := entity.GetID()
	i.entities[id] = entity
	return nil
}

func NewInMemoryCustomerRepository() customer.Repository {
	return &InMemoryCustomerRepository{
		customers: map[uuid.UUID]*customer.Customer{},
	}
}
