package test

import (
	"github.com/bitlogic/go-startup/src/domain/customer"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_GivenAnEmptyName_WhenNewCustomer_ThenReturnError(t *testing.T) {
	customer, err := customer.NewCustomer("")

	assert.Error(t, err)
	assert.Equal(t, "invalid name", err.Error())
	assert.Nil(t, customer)
}

func Test_GivenABlankName_WhenNewCustomer_ThenReturnError(t *testing.T) {
	customer, err := customer.NewCustomer("                         ")

	assert.Error(t, err)
	assert.Equal(t, "invalid name", err.Error())
	assert.Nil(t, customer)
}

func Test_GivenAShortName_WhenNewCustomer_ThenReturnError(t *testing.T) {
	customer, err := customer.NewCustomer("John                 ")

	assert.Error(t, err)
	assert.Equal(t, "invalid name", err.Error())
	assert.Nil(t, customer)
}

func Test_GivenAValidName_WhenNewCustomer_ThenReturnANonEmptyCustomer(t *testing.T) {
	customer, err := customer.NewCustomer("     John Mayer       ")

	assert.Nil(t, err)
	assert.NotNil(t, customer)
	assert.NotEmpty(t, customer)
}
