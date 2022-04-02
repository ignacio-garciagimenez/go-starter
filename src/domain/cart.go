package domain

import (
	"errors"
	"reflect"

	"github.com/google/uuid"
)

type CartId uuid.UUID

type Cart struct {
	*baseEntity[CartId]
	customerId CustomerId
	items      map[ProductId]item
}

type item struct {
	productId ProductId
	price     float64
	quantity  int
}

func NewCart(customer *Customer) (*Cart, error) {
	if customer == nil {
		return nil, errors.New("no customer provided")
	}

	cart := &Cart{
		baseEntity: &baseEntity[CartId]{
			id: CartId(uuid.New()),
		},
		customerId: customer.GetID(),
		items:      map[ProductId]item{},
	}

	cart.addDomainEvent(CartCreated{
		CartId:     cart.id,
		CustomerId: cart.customerId,
	})

	return cart, nil
}

func (c *Cart) EqualsTo(entity Entity[CartId]) bool {
	return reflect.TypeOf(c) == reflect.TypeOf(entity) && c.GetID() == entity.GetID()
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

	productId := product.GetID()
	if cartItem, found := c.items[productId]; found {
		c.items[productId] = cartItem.addQuantity(quantity)
	} else {
		c.items[productId] = item{
			productId: product.GetID(),
			price:     product.GetPrice(),
			quantity:  quantity,
		}
	}

	c.addDomainEvent(ItemAddedToCart{
		CartId:    c.id,
		ProductId: productId,
		Quantity:  quantity,
	})

	return c.items[productId], nil
}

func (c Cart) GetTotal() float64 {
	var total float64
	for _, item := range c.items {
		total += item.getTotal()
	}

	return total
}

func (c Cart) GetCustomerID() CustomerId {
	return c.customerId
}

func (c Cart) GetItems() []item {
	var output []item
	for _, cartItem := range c.items {
		output = append(output, cartItem)
	}

	return output
}

func (i item) GetProductId() ProductId {
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

func (i item) addQuantity(quantityToAdd int) item {
	return item{
		productId: i.productId,
		price:     i.price,
		quantity:  i.quantity + quantityToAdd,
	}
}

func (v item) EqualsTo(other ValueObject) bool {
	return reflect.DeepEqual(v, other)
}
