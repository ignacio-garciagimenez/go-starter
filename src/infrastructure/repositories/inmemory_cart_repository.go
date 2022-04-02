package repositories

import (
	"github.com/bitlogic/go-startup/src/domain"
)

type InMemoryCartRepository struct {
	*inMemoryBaseRepository[domain.CartId, *domain.Cart]
	customerIndex map[domain.CustomerId][]*domain.Cart
}

func (i *InMemoryCartRepository) Save(entity *domain.Cart) error {
	if err := i.inMemoryBaseRepository.Save(entity); err != nil {
		return err
	}

	i.customerIndex[entity.GetCustomerID()] = append(i.customerIndex[entity.GetCustomerID()], entity)
	return nil
}

func (i *InMemoryCartRepository) GetCustomerCarts(customerId domain.CustomerId) []*domain.Cart {
	return i.customerIndex[customerId]
}

func NewInMemoryCartRepository() domain.CartRepository {
	return &InMemoryCartRepository{
		inMemoryBaseRepository: &inMemoryBaseRepository[domain.CartId, *domain.Cart]{
			entities: map[domain.CartId]*domain.Cart{},
		},
		customerIndex: map[domain.CustomerId][]*domain.Cart{},
	}
}
