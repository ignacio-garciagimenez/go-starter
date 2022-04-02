package domain

import "github.com/google/uuid"

type ProductRepository interface {
	Repository[uuid.UUID, *Product]
}

type CustomerRepository interface {
	Repository[uuid.UUID, *Customer]
}

type CartRepository interface {
	Repository[uuid.UUID, *Cart]
	GetCustomerCarts(customerId uuid.UUID) []*Cart
}
