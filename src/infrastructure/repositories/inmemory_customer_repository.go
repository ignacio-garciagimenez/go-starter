package repositories

import (
	"github.com/bitlogic/go-startup/src/domain/customer"
	"github.com/google/uuid"
)

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
