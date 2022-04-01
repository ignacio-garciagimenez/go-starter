package product

import (
	"errors"

	"github.com/bitlogic/go-startup/src/domain"
	"github.com/google/uuid"
)

type Product struct {
	domain.Entity[uuid.UUID]
	id        uuid.UUID
	name      string
	unitPrice float64
}

func NewProduct(name string, price float64) (*Product, error) {
	if len(name) < 10 || price <= 0.00 {
		return nil, errors.New("invalid arguments")
	}

	return &Product{
		id:        uuid.New(),
		name:      name,
		unitPrice: price,
	}, nil
}

func (p Product) GetID() uuid.UUID {
	return p.id
}

func (p Product) GetName() string {
	return p.name
}

func (p Product) GetPrice() float64 {
	return p.unitPrice
}
