package repositories

import (
	"errors"
	"github.com/bitlogic/go-startup/src/domain/product"
	"github.com/google/uuid"
)

type InMemoryProductRepository struct {
	products map[uuid.UUID]*product.Product
}

func (r *InMemoryProductRepository) FindByID(key uuid.UUID) (*product.Product, error) {
	entity, found := r.products[key]
	if !found {
		return nil, errors.New("product not found")
	}

	return entity, nil
}

func (r *InMemoryProductRepository) Save(entity *product.Product) error {
	r.products[entity.GetID()] = entity

	return nil
}

func NewInMemoryProductRepository() product.Repository {
	return &InMemoryProductRepository{
		products: map[uuid.UUID]*product.Product{},
	}
}
