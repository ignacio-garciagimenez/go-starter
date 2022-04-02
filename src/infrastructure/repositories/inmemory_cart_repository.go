package repositories

import (
	"github.com/bitlogic/go-startup/src/domain"
	"github.com/google/uuid"
)

type InMemoryCartRepository struct {
	*inMemoryBaseRepository[uuid.UUID, *domain.Cart]
	customerIndex map[uuid.UUID][]*domain.Cart
}

func (i *InMemoryCartRepository) Save(entity *domain.Cart) error {
	if err := i.inMemoryBaseRepository.Save(entity); err != nil {
		return err
	}

	i.customerIndex[entity.GetCustomerID()] = append(i.customerIndex[entity.GetCustomerID()], entity)
	return nil
}

func (i *InMemoryCartRepository) GetCustomerCarts(customerId uuid.UUID) []*domain.Cart {
	return i.customerIndex[customerId]
}

func NewInMemoryCartRepository() domain.CartRepository {
	return &InMemoryCartRepository{
		inMemoryBaseRepository: &inMemoryBaseRepository[uuid.UUID, *domain.Cart]{
			entities: map[uuid.UUID]*domain.Cart{},
		},
		customerIndex: map[uuid.UUID][]*domain.Cart{},
	}
}
