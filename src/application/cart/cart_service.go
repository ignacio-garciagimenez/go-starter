package application

import (
	"errors"

	"github.com/bitlogic/go-startup/src/application"
	"github.com/bitlogic/go-startup/src/domain/cart"
	"github.com/bitlogic/go-startup/src/domain/customer"
)

type CartService struct {
	cartRepository     cart.Repository
	customerRepository customer.Repository
}

func NewCartService(cartRepository cart.Repository, customerRepository customer.Repository) (*CartService, error) {
	if cartRepository == nil {
		return nil, errors.New("cart repository was nil")
	}

	if customerRepository == nil {
		return nil, errors.New("customer repository was nil")
	}

	return &CartService{
		cartRepository:     cartRepository,
		customerRepository: customerRepository,
	}, nil

}

func (s *CartService) CreateNewCart(command CreateCartCommand) (application.CartDto, error) {
	customer, _ := s.customerRepository.FindByID(command.CustomerId)
	cart, _ := cart.NewCart(customer)

	s.cartRepository.Save(cart)

	cartDto := mapCartToDto(cart)

	return cartDto, nil
}

func mapCartToDto(cart *cart.Cart) application.CartDto {
	var itemDtos []application.ItemDto

	for _, item := range cart.GetItems() {
		itemDtos = append(itemDtos, application.ItemDto{
			ProductId: item.GetProductId(),
			UnitPrice: item.GetUnitPrice(),
			Quantity:  item.GetQuantity(),
		})
	}

	return application.CartDto{
		Id:         cart.GetID(),
		CustomerId: cart.GetCustomerID(),
		Items:      itemDtos,
	}
}
