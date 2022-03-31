package cart

import (
	"github.com/google/uuid"
)

type item struct {
	productId uuid.UUID
	price     float64
	quantity  int
}
