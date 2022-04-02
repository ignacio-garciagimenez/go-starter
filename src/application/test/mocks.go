package test

import (
	"github.com/bitlogic/go-startup/src/domain"
	"github.com/google/uuid"
)

type cartRepositoryMock struct {
	callCount int
	findById  func(uuid.UUID) (*domain.Cart, error)
	save      func(*domain.Cart) error
}

func (r *cartRepositoryMock) FindByID(cartId uuid.UUID) (*domain.Cart, error) {
	r.callCount++
	return r.findById(cartId)
}

func (r *cartRepositoryMock) Save(cart *domain.Cart) error {
	r.callCount++
	return r.save(cart)
}

func (r *cartRepositoryMock) GetCustomerCarts(customerId uuid.UUID) []*domain.Cart {
	return nil
}

type productRepositoryMock struct {
	callCount int
	findByID  func(uuid.UUID) (*domain.Product, error)
	save      func(*domain.Product) error
}

func (m *productRepositoryMock) FindByID(productId uuid.UUID) (*domain.Product, error) {
	m.callCount++
	return m.findByID(productId)
}

func (m *productRepositoryMock) Save(newProduct *domain.Product) error {
	m.callCount++
	return m.save(newProduct)
}

type customerRepositoryMock struct {
	callCount int
	findById  func(uuid.UUID) (*domain.Customer, error)
	save      func(*domain.Customer) error
}

func (r *customerRepositoryMock) FindByID(customerId uuid.UUID) (*domain.Customer, error) {
	r.callCount++
	return r.findById(customerId)
}

func (r *customerRepositoryMock) Save(customer *domain.Customer) error {
	r.callCount++
	return r.save(customer)
}
