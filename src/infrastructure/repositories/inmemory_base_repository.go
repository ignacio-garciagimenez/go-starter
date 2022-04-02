package repositories

import (
	"errors"

	"github.com/bitlogic/go-startup/src/domain"
)

type inMemoryBaseRepository[K comparable, E domain.Entity[K]] struct {
	entities map[K]E
}

func (i *inMemoryBaseRepository[K, E]) FindByID(key K) (E, error) {
	var entity E
	if entity, found := i.entities[key]; found {
		return entity, nil
	}

	return entity, errors.New("entity not found")
}

func (i *inMemoryBaseRepository[K, E]) Save(entity E) error {
	i.entities[entity.GetID()] = entity
	return nil
}
