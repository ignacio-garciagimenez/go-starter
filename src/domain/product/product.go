package product

import (
	"errors"
	"github.com/google/uuid"
)

type Product struct {
	id    uuid.UUID
	name  string
	price float64
}

func NewProduct(name string, price float64) (*Product, error) {
	if len(name) < 10 || price <= 0.00 {
		return nil, errors.New("invalid arguments")
	}

	return &Product{
		id:    uuid.New(),
		name:  name,
		price: price,
	}, nil
}

func (p Product) GetId() uuid.UUID {
	return p.id
}

func (p Product) GetName() string {
	return p.name
}

func (p Product) GetPrice() float64 {
	return p.price
}
