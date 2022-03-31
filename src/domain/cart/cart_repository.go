package cart

import (
	"github.com/bitlogic/go-startup/src/domain"
	"github.com/google/uuid"
)

type Repository interface {
	domain.Repository[uuid.UUID, *Cart]
	GetCustomerCarts(customerId uuid.UUID) []*Cart
}
