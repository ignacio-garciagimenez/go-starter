package domain

import (
	"errors"

	"github.com/google/uuid"
)

type Product struct {
	*baseEntity[uuid.UUID]
	name      string
	unitPrice float64
}

func NewProduct(name string, price float64) (*Product, error) {
	if len(name) < 10 || price <= 0.00 {
		return nil, errors.New("invalid arguments")
	}

	product := &Product{
		baseEntity: &baseEntity[uuid.UUID]{
			id: uuid.New(),
		},
		name:      name,
		unitPrice: price,
	}

	product.addDomainEvent(ProductCreated{
		ProductId:        product.id,
		ProductName:      product.name,
		ProductUnitPrice: product.unitPrice,
	})

	return product, nil
}

func (p Product) GetName() string {
	return p.name
}

func (p Product) GetPrice() float64 {
	return p.unitPrice
}
