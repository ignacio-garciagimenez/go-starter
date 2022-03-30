package cart

import (
	"errors"
	"github.com/bitlogic/go-startup/src/domain/product"
)

type Cart struct {
	items []Item
}

func NewCart(itemsToAdd []*product.Product) (*Cart, error) {
	if itemsToAdd == nil || len(itemsToAdd) < 1 {
		return nil, errors.New("no items provided")
	}

	var items []Item

	for _, product := range itemsToAdd {
		items = append(items, newItem(product))
	}

	return &Cart{
		items: items,
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
