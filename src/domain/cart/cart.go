package cart

import (
	"errors"
	"github.com/bitlogic/go-startup/src/domain/customer"
	"github.com/bitlogic/go-startup/src/domain/product"
	"github.com/google/uuid"
)

type Cart struct {
	id         uuid.UUID
	customerId uuid.UUID
	items      []Item
}

func NewCart(customer *customer.Customer) (*Cart, error) {
	if customer == nil {
		return nil, errors.New("no customer provided")
	}

	return &Cart{
		id:         uuid.New(),
		customerId: customer.GetId(),
		items:      []Item{},
	}, nil
}

func newItem(product *product.Product) Item {
	return Item{
		productId: product.GetId(),
		price:     product.GetPrice(),
	}
}

func (c Cart) Size() int {
	return len(c.items)
}

func (c *Cart) AddItem(product *product.Product) (Item, error) {
	if product == nil {
		return Item{}, errors.New("invalid product")
	}

	item := newItem(product)
	c.items = append(c.items, item)

	return item, nil
}
