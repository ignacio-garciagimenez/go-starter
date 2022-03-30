package cart

import (
	"github.com/google/uuid"
)

type Item struct {
	productId uuid.UUID
	price     float64
}
