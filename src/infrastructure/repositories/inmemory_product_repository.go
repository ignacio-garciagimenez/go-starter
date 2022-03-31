package repositories

import (
	"github.com/bitlogic/go-startup/src/domain/product"
	"github.com/google/uuid"
)

type InMemoryProductRepository struct {
	*inMemoryBaseRepository[uuid.UUID, *product.Product]
}

func NewInMemoryProductRepository() product.Repository {
	return &InMemoryProductRepository{
		inMemoryBaseRepository: &inMemoryBaseRepository[uuid.UUID, *product.Product]{
			entities: map[uuid.UUID]*product.Product{},
		},
	}
}
