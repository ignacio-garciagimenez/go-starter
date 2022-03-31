package repositories

import (
	"errors"

	"github.com/bitlogic/go-startup/src/domain"
	"github.com/bitlogic/go-startup/src/domain/customer"
	"github.com/google/uuid"
)

type inMemoryBaseRepository[K domain.EntityKey, E domain.Entity[K]] struct {
	entities map[K]E
}

func (i *inMemoryBaseRepository[K, E]) FindByID(key K) (E, error) {
	entity, found := i.entities[key]
	if !found {
		return *new(E), errors.New("cart not found")
	}

	return entity, nil
}

func (i *inMemoryBaseRepository[K, E]) Save(entity E) error {
	i.entities[entity.GetID()] = entity
	return nil
}

type InMemoryCustomerRepository struct {
	*inMemoryBaseRepository[uuid.UUID, *customer.Customer]
}

func NewInMemoryCustomerRepository() customer.Repository {
	return &InMemoryCustomerRepository{
		inMemoryBaseRepository: &inMemoryBaseRepository[uuid.UUID, *customer.Customer]{
			entities: map[uuid.UUID]*customer.Customer{},
		},
	}
}
