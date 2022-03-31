package domain

type Entity interface{}
type ValueObject interface{}

type Repository[K any, E Entity] interface {
	FindByID(key K) (E, error)
	Save(entity E) error
}
