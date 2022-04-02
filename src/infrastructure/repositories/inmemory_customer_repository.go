package repositories

import (
	"github.com/bitlogic/go-startup/src/domain"
)

type InMemoryCustomerRepository struct {
	*inMemoryBaseRepository[domain.CustomerId, *domain.Customer]
}

func NewInMemoryCustomerRepository() domain.CustomerRepository {
	return &InMemoryCustomerRepository{
		inMemoryBaseRepository: &inMemoryBaseRepository[domain.CustomerId, *domain.Customer]{
			entities: map[domain.CustomerId]*domain.Customer{},
		},
	}
}
