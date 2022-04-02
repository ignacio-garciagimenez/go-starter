package application

import (
	"errors"

	"github.com/bitlogic/go-startup/src/domain"
	"github.com/google/uuid"
)

type CartService struct {
	cartRepository     domain.CartRepository
	customerRepository domain.CustomerRepository
	productRepository  domain.ProductRepository
}

func NewCartService(cartRepository domain.CartRepository, customerRepository domain.CustomerRepository, productRepository domain.ProductRepository) (*CartService, error) {
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

func (s *CartService) CreateNewCart(command CreateCartCommand) (CartDto, error) {
	customer, err := s.customerRepository.FindByID(command.CustomerId)
	if err != nil {
		return CartDto{}, err
	}

	cart, err := domain.NewCart(customer)
	if err != nil {
		return CartDto{}, err
	}

	if err = s.cartRepository.Save(cart); err != nil {
		return CartDto{}, err
	}

	return mapCartToDto(cart), nil
}

func (s *CartService) AddItemToCart(command AddItemToCartCommand) (CartDto, error) {
	product, err := s.productRepository.FindByID(command.ProductId)
	if err != nil {
		return CartDto{}, err
	}

	cart, err := s.cartRepository.FindByID(command.CartId)
	if err != nil {
		return CartDto{}, err
	}

	if _, err = cart.AddItem(product, command.Quantity); err != nil {
		return CartDto{}, err
	}

	if err = s.cartRepository.Save(cart); err != nil {
		return CartDto{}, err
	}

	return mapCartToDto(cart), nil
}

func mapCartToDto(cart *domain.Cart) CartDto {
	var itemDtos []ItemDto

	for _, item := range cart.GetItems() {
		itemDtos = append(itemDtos, ItemDto{
			ProductId: uuid.UUID(item.GetProductId()),
			UnitPrice: item.GetUnitPrice(),
			Quantity:  item.GetQuantity(),
		})
	}

	return CartDto{
		Id:         uuid.UUID(cart.GetID()),
		CustomerId: uuid.UUID(cart.GetCustomerID()),
		Items:      itemDtos,
	}
}
