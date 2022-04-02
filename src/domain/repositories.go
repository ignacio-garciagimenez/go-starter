package domain

type ProductRepository interface {
	Repository[ProductId, *Product]
}

type CustomerRepository interface {
	Repository[CustomerId, *Customer]
}

type CartRepository interface {
	Repository[CartId, *Cart]
	GetCustomerCarts(customerId CustomerId) []*Cart
}
