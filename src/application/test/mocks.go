package test

import (
	"github.com/bitlogic/go-startup/src/domain/cart"
	"github.com/bitlogic/go-startup/src/domain/customer"
	"github.com/bitlogic/go-startup/src/domain/product"
	"github.com/google/uuid"
)

type cartRepositoryMock struct {
	callCount int
	findById  func(uuid.UUID) (*cart.Cart, error)
	save      func(*cart.Cart) error
}

func (r *cartRepositoryMock) FindByID(cartId uuid.UUID) (*cart.Cart, error) {
	r.callCount++
	return r.findById(cartId)
}

func (r *cartRepositoryMock) Save(cart *cart.Cart) error {
	r.callCount++
	return r.save(cart)
}

func (r *cartRepositoryMock) GetCustomerCarts(customerId uuid.UUID) []*cart.Cart {
	return nil
}

type productRepositoryMock struct {
	callCount int
	findByID  func(uuid.UUID) (*product.Product, error)
	save      func(*product.Product) error
}

func (m *productRepositoryMock) FindByID(productId uuid.UUID) (*product.Product, error) {
	m.callCount++
	return m.findByID(productId)
}

func (m *productRepositoryMock) Save(newProduct *product.Product) error {
	m.callCount++
	return m.save(newProduct)
}

type customerRepositoryMock struct {
	callCount int
	findById  func(uuid.UUID) (*customer.Customer, error)
	save      func(*customer.Customer) error
}

func (r *customerRepositoryMock) FindByID(customerId uuid.UUID) (*customer.Customer, error) {
	r.callCount++
	return r.findById(customerId)
}

func (r *customerRepositoryMock) Save(customer *customer.Customer) error {
	r.callCount++
	return r.save(customer)
}
