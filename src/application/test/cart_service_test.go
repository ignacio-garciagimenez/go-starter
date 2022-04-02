package test

import (
	"testing"

	application "github.com/bitlogic/go-startup/src/application/cart"
	"github.com/bitlogic/go-startup/src/domain/cart"
	"github.com/bitlogic/go-startup/src/domain/customer"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func Test_GivenANilCartRepository_WhenNewCartService_ThenReturnError(t *testing.T) {
	service, err := application.NewCartService(nil, &customerRepositoryMock{})

	if assert.Error(t, err) {
		assert.Equal(t, "cart repository was nil", err.Error())
	}
	assert.Nil(t, service)
}

func Test_GivenANilCustomerRepository_WhenNewCartService_ThenReturnError(t *testing.T) {
	service, err := application.NewCartService(&cartRepositoryMock{}, nil)

	if assert.Error(t, err) {
		assert.Equal(t, "customer repository was nil", err.Error())
	}
	assert.Nil(t, service)
}

func Test_GivenACartRepository_WhenNewCartService_ThenReturnACartService(t *testing.T) {
	service, err := application.NewCartService(&cartRepositoryMock{}, &customerRepositoryMock{})

	assert.Nil(t, err)
	assert.NotEmpty(t, service)
}

func Test_GivenAValidCreateCartCommand_WhenCreateNewCart_ThenReturnACreateCartResponse(t *testing.T) {
	cartRepository := &cartRepositoryMock{
		save: func(cart *cart.Cart) error {
			return nil
		},
	}
	savedCustomer, _ := customer.NewCustomer("Grady Booch")
	customerRepository := &customerRepositoryMock{
		findById: func(customerId uuid.UUID) (*customer.Customer, error) {
			return savedCustomer, nil
		},
	}
	service, _ := application.NewCartService(cartRepository, customerRepository)
	command := application.CreateCartCommand{
		CustomerId: savedCustomer.GetID(),
	}

	result, err := service.CreateNewCart(command)

	assert.Nil(t, err)
	if assert.NotEmpty(t, result) {
		assert.NotNil(t, result.Id)
		assert.Equal(t, savedCustomer.GetID(), result.CustomerId)
		assert.Empty(t, result.Items)
	}
	assert.Equal(t, 1, customerRepository.callCount)
	assert.Equal(t, 1, cartRepository.callCount)

}

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
