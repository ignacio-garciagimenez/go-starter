package repositories

import (
	"errors"
	"github.com/bitlogic/go-startup/src/domain/cart"
	"github.com/google/uuid"
)

type InMemoryCartRepository struct {
	carts         map[uuid.UUID]*cart.Cart
	customerIndex map[uuid.UUID][]*cart.Cart
}

func (i *InMemoryCartRepository) FindByID(key uuid.UUID) (*cart.Cart, error) {
	entity, found := i.carts[key]
	if !found {
		return nil, errors.New("cart not found")
	}

	return entity, nil
}

func (i *InMemoryCartRepository) Save(entity *cart.Cart) error {
	i.carts[entity.GetID()] = entity
	i.customerIndex[entity.GetCustomerID()] = append(i.customerIndex[entity.GetCustomerID()], entity)
	return nil
}

func (i *InMemoryCartRepository) GetCustomerCarts(customerId uuid.UUID) []*cart.Cart {
	return i.customerIndex[customerId]
}

func NewInMemoryCartRepository() cart.Repository {
	return &InMemoryCartRepository{
		carts:         map[uuid.UUID]*cart.Cart{},
		customerIndex: map[uuid.UUID][]*cart.Cart{},
	}
}
