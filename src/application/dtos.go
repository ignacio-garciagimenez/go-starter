package application

import (
	"github.com/bitlogic/go-startup/src/domain"
)

type CartDto struct {
	Id         domain.CartId
	CustomerId domain.CustomerId
	Items      []ItemDto
}

type ItemDto struct {
	ProductId domain.ProductId
	UnitPrice float64
	Quantity  int
}

type CustomerDto struct {
	Id   domain.CustomerId
	Name string
}

type ProductDto struct {
	Id        domain.ProductId
	Name      string
	UnitPrice float64
}
