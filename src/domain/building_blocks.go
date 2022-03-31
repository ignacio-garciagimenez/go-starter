package domain

import (
	"github.com/google/uuid"
	"golang.org/x/exp/constraints"
)

type Entity[K EntityKey] interface {
	GetID() K
}

type ValueObject interface{}

type EntityKey interface {
	constraints.Ordered | uuid.UUID
}

type Repository[K EntityKey, E Entity[K]] interface {
	FindByID(key K) (E, error)
	Save(entity E) error
}
