package domain

type CartCreated struct {
	DomainEvent
	CartId     CartId
	CustomerId CustomerId
}

type ItemAddedToCart struct {
	DomainEvent
	CartId    CartId
	ProductId ProductId
	Quantity  int
}

type CustomerCreated struct {
	CustomerId   CustomerId
	CustomerName string
}

type ProductCreated struct {
	ProductId        ProductId
	ProductName      string
	ProductUnitPrice float64
}
