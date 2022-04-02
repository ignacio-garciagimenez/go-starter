package domain

import (
	"errors"
	"reflect"

	"github.com/google/uuid"
)

type ProductId uuid.UUID

type Product struct {
	*baseEntity[ProductId]
	name      string
	unitPrice float64
}

func NewProduct(name string, price float64) (*Product, error) {
	if len(name) < 10 || price <= 0.00 {
		return nil, errors.New("invalid arguments")
	}

	product := &Product{
		baseEntity: &baseEntity[ProductId]{
			id: ProductId(uuid.New()),
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

func (p *Product) EqualsTo(entity Entity[ProductId]) bool {
	return reflect.TypeOf(p) == reflect.TypeOf(entity) &&
		p.GetID() == entity.GetID()
}
