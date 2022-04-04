package application

import "fmt"

type NotFoundError struct {
	entityId   string
	entityType string
}

func (e NotFoundError) Error() string {
	return fmt.Sprintf(`%s with id %s not found`, e.entityId, e.entityType)
}

func NewNotFoundError(entityId string, entityType string) error {
	return &NotFoundError{
		entityId:   entityId,
		entityType: entityType,
	}
}
