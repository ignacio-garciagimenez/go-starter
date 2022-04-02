package application

import (
	"errors"

	"github.com/bitlogic/go-startup/src/application"
	"github.com/bitlogic/go-startup/src/domain/cart"
	"github.com/bitlogic/go-startup/src/domain/customer"
	"github.com/bitlogic/go-startup/src/domain/product"
)

type CartService struct {
	cartRepository     cart.Repository
	customerRepository customer.Repository
	productRepository  product.Repository
}

func NewCartService(cartRepository cart.Repository, customerRepository customer.Repository, productRepository product.Repository) (*CartService, error) {
	if cartRepository == nil {
		return nil, errors.New("cart repository was nil")
	}

	if customerRepository == nil {
		return nil, errors.New("customer repository was nil")
	}

	if productRepository == nil {
		return nil, errors.New("product repository was nil")
	}

	return &CartService{
		cartRepository:     cartRepository,
		customerRepository: customerRepository,
		productRepository:  productRepository,
	}, nil
}

func (s *CartService) CreateNewCart(command CreateCartCommand) (application.CartDto, error) {
	customer, err := s.customerRepository.FindByID(command.CustomerId)
	if err != nil {
		return application.CartDto{}, err
	}

	cart, err := cart.NewCart(customer)
	if err != nil {
		return application.CartDto{}, err
	}

	if err = s.cartRepository.Save(cart); err != nil {
		return application.CartDto{}, err
	}

	return mapCartToDto(cart), nil
}

func (s *CartService) AddItemToCart(command AddItemToCartCommand) (application.CartDto, error) {
	product, err := s.productRepository.FindByID(command.ProductId)
	if err != nil {
		return application.CartDto{}, err
	}

	cart, err := s.cartRepository.FindByID(command.CartId)
	if err != nil {
		return application.CartDto{}, err
	}

	if _, err = cart.AddItem(product, command.Quantity); err != nil {
		return application.CartDto{}, err
	}

	if err = s.cartRepository.Save(cart); err != nil {
		return application.CartDto{}, err
	}

	return mapCartToDto(cart), nil
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
