package repositories

import (
	"github.com/bitlogic/go-startup/src/domain"
	"github.com/google/uuid"
)

type InMemoryProductRepository struct {
	*inMemoryBaseRepository[uuid.UUID, *domain.Product]
}

func NewInMemoryProductRepository() domain.ProductRepository {
	return &InMemoryProductRepository{
		inMemoryBaseRepository: &inMemoryBaseRepository[uuid.UUID, *domain.Product]{
			entities: map[uuid.UUID]*domain.Product{},
		},
	}
}
