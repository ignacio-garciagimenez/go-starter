package domain

import "golang.org/x/exp/constraints"

type Entity[V constraints.Ordered] interface {
	GetID() V
}
type ValueObject interface{}

//
//type Repository[K any, E *Entity] interface {
//	FindByID(key K) (E, error)
//	Save(entity E) error
//}
