package repositories

import (
	"github.com/bitlogic/go-startup/src/domain/cart"
	"github.com/google/uuid"
)

type InMemoryCartRepository struct {
	*inMemoryBaseRepository[uuid.UUID, *cart.Cart]
	customerIndex map[uuid.UUID][]*cart.Cart
}

func (i *InMemoryCartRepository) Save(entity *cart.Cart) error {
	if err := i.inMemoryBaseRepository.Save(entity); err != nil {
		return err
	}

	i.customerIndex[entity.GetCustomerID()] = append(i.customerIndex[entity.GetCustomerID()], entity)
	return nil
}

func (i *InMemoryCartRepository) GetCustomerCarts(customerId uuid.UUID) []*cart.Cart {
	return i.customerIndex[customerId]
}

func NewInMemoryCartRepository() cart.Repository {
	return &InMemoryCartRepository{
		inMemoryBaseRepository: &inMemoryBaseRepository[uuid.UUID, *cart.Cart]{
			entities: map[uuid.UUID]*cart.Cart{},
		},
		customerIndex: map[uuid.UUID][]*cart.Cart{},
	}
}
