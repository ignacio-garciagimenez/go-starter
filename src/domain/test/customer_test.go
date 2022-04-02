package test

import (
	"testing"

	"github.com/bitlogic/go-startup/src/domain"
	"github.com/stretchr/testify/assert"
)

func Test_GivenAnEmptyName_WhenNewCustomer_ThenReturnError(t *testing.T) {
	customer, err := domain.NewCustomer("")

	assert.Error(t, err)
	assert.Equal(t, "invalid name", err.Error())
	assert.Nil(t, customer)
}

func Test_GivenABlankName_WhenNewCustomer_ThenReturnError(t *testing.T) {
	customer, err := domain.NewCustomer("                         ")

	assert.Error(t, err)
	assert.Equal(t, "invalid name", err.Error())
	assert.Nil(t, customer)
}

func Test_GivenAShortName_WhenNewCustomer_ThenReturnError(t *testing.T) {
	customer, err := domain.NewCustomer("John                 ")

	assert.Error(t, err)
	assert.Equal(t, "invalid name", err.Error())
	assert.Nil(t, customer)
}

func Test_GivenAValidName_WhenNewCustomer_ThenReturnANonEmptyCustomer(t *testing.T) {
	customer, err := domain.NewCustomer("     John Mayer       ")

	assert.Nil(t, err)
	assert.NotNil(t, customer)
	assert.NotEmpty(t, customer)
	if assert.NotEmpty(t, customer.GetDomainEvents()) {
		assert.Equal(t, 1, len(customer.GetDomainEvents()))
	}
}

func Test_GivenACustomer_WhenEqualsToItself_ThenReturnsTrue(t *testing.T) {
	customer, _ := domain.NewCustomer("John Mayer")

	assert.True(t, customer.EqualsTo(customer))
}

func Test_GivenACustomer_WhenEqualsToAnotherCustomer_ThenReturnsFalse(t *testing.T) {
	customer, _ := domain.NewCustomer("John Mayer")
	customer2, _ := domain.NewCustomer("John Mayer")

	assert.False(t, customer.EqualsTo(customer2))
}

func Test_GivenACustomer_WhenEqualsToAnotherTypeOfEntity_ThenReturnsFalse(t *testing.T) {
	customer, _ := domain.NewCustomer("John Mayer")
	product, _ := domain.NewProduct("Pepsi Ligh", 10.00)

	assert.False(t, customer.EqualsTo(product))
}
