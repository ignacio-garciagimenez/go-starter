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
