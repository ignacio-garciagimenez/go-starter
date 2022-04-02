package repositories

import (
	"github.com/bitlogic/go-startup/src/domain"
	"github.com/google/uuid"
)

type InMemoryCustomerRepository struct {
	*inMemoryBaseRepository[uuid.UUID, *domain.Customer]
}

func NewInMemoryCustomerRepository() domain.CustomerRepository {
	return &InMemoryCustomerRepository{
		inMemoryBaseRepository: &inMemoryBaseRepository[uuid.UUID, *domain.Customer]{
			entities: map[uuid.UUID]*domain.Customer{},
		},
	}
}
