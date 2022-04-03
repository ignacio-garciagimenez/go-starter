package test

import (
	"github.com/bitlogic/go-startup/src/domain"
)

type cartRepositoryMock struct {
	callCount int
	findById  func(domain.CartId) (*domain.Cart, error)
	save      func(*domain.Cart) error
}

func (r *cartRepositoryMock) FindByID(cartId domain.CartId) (*domain.Cart, error) {
	r.callCount++
	return r.findById(cartId)
}

func (r *cartRepositoryMock) Save(cart *domain.Cart) error {
	r.callCount++
	return r.save(cart)
}

func (r *cartRepositoryMock) GetCustomerCarts(customerId domain.CustomerId) []*domain.Cart {
	return nil
}

type productRepositoryMock struct {
	callCount int
	findByID  func(domain.ProductId) (*domain.Product, error)
	save      func(*domain.Product) error
}

func (m *productRepositoryMock) FindByID(productId domain.ProductId) (*domain.Product, error) {
	m.callCount++
	return m.findByID(productId)
}

func (m *productRepositoryMock) Save(newProduct *domain.Product) error {
	m.callCount++
	return m.save(newProduct)
}

type customerRepositoryMock struct {
	callCount int
	findById  func(domain.CustomerId) (*domain.Customer, error)
	save      func(*domain.Customer) error
}

func (r *customerRepositoryMock) FindByID(customerId domain.CustomerId) (*domain.Customer, error) {
	r.callCount++
	return r.findById(customerId)
}

func (r *customerRepositoryMock) Save(customer *domain.Customer) error {
	r.callCount++
	return r.save(customer)
}
