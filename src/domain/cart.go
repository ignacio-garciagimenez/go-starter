package domain

import (
	"errors"

	"github.com/google/uuid"
)

type Cart struct {
	*baseEntity[uuid.UUID]
	customerId uuid.UUID
	items      map[uuid.UUID]item
}

type item struct {
	ValueObject
	productId uuid.UUID
	price     float64
	quantity  int
}

func NewCart(customer *Customer) (*Cart, error) {
	if customer == nil {
		return nil, errors.New("no customer provided")
	}

	return &Cart{
		baseEntity: &baseEntity[uuid.UUID]{
			id: uuid.New(),
		},
		customerId: customer.GetID(),
		items:      map[uuid.UUID]item{},
	}, nil
}

func (c Cart) Size() int {
	var cartSize int
	for _, v := range c.items {
		cartSize += v.quantity
	}

	return cartSize
}

func (c *Cart) AddItem(product *Product, quantity int) (item, error) {
	if product == nil {
		return item{}, errors.New("invalid product")
	}

	if quantity < 1 {
		return item{}, errors.New("invalid quantity")
	}

	cartItem, err := c.findItem(product.GetID())
	if err != nil {
		return c.addNewItem(product, quantity), nil
	}

	return c.updateItemQuantity(cartItem, quantity), nil
}

func (c Cart) findItem(productId uuid.UUID) (item, error) {
	cartItem, found := c.items[productId]
	if !found {
		return item{}, errors.New("item not found")
	}

	return cartItem, nil
}

func (c *Cart) addNewItem(product *Product, quantity int) item {
	item := newItem(product, quantity)

	c.items[product.GetID()] = item
	return item
}

func (c *Cart) updateItemQuantity(cartItem item, quantityToAdd int) item {
	cartItem.quantity += quantityToAdd
	c.items[cartItem.productId] = cartItem

	return cartItem
}

func (c Cart) GetTotal() float64 {
	var total float64
	for _, item := range c.items {
		total += item.getTotal()
	}

	return total
}

func (c Cart) GetCustomerID() uuid.UUID {
	return c.customerId
}

func (c Cart) GetItems() []item {
	var output []item
	for _, cartItem := range c.items {
		output = append(output, cartItem)
	}

	return output
}

func newItem(product *Product, quantity int) item {
	return item{
		productId: product.GetID(),
		price:     product.GetPrice(),
		quantity:  quantity,
	}
}

func (i item) GetProductId() uuid.UUID {
	return i.productId
}

func (i item) GetUnitPrice() float64 {
	return i.price
}

func (i item) GetQuantity() int {
	return i.quantity
}

func (i item) getTotal() float64 {
	return i.price * float64(i.quantity)
}
