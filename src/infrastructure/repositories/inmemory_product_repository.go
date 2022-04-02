package repositories

import (
	"github.com/bitlogic/go-startup/src/domain"
)

type InMemoryProductRepository struct {
	*inMemoryBaseRepository[domain.ProductId, *domain.Product]
}

func NewInMemoryProductRepository() domain.ProductRepository {
	return &InMemoryProductRepository{
		inMemoryBaseRepository: &inMemoryBaseRepository[domain.ProductId, *domain.Product]{
			entities: map[domain.ProductId]*domain.Product{},
		},
	}
}
