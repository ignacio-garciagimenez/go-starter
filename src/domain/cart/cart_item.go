package cart

import (
	"github.com/bitlogic/go-startup/src/domain"
	"github.com/google/uuid"
)

type item struct {
	domain.ValueObject
	productId uuid.UUID
	price     float64
	quantity  int
}

func (i item) getTotal() float64 {
	return i.price * float64(i.quantity)
}
