package test

import (
	"github.com/bitlogic/go-startup/src/domain/cart"
	"github.com/bitlogic/go-startup/src/domain/customer"
	"github.com/bitlogic/go-startup/src/domain/product"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_GivenNoCustomer_WhenNewCart_ThenReturnError(t *testing.T) {
	cart, err := cart.NewCart(nil)

	assert.Error(t, err, "err should not be nil")
	assert.Equal(t, "no customer provided", err.Error())
	assert.Nil(t, cart, "cart should be nil")
}

func Test_GivenACustomer_WhenNewCart_ThenReturnNewCart(t *testing.T) {
	cartCustomer, _ := customer.NewCustomer("John Mayer")
	cart, err := cart.NewCart(cartCustomer)

	assert.Nil(t, err)
	assert.NotNil(t, cart)
	assert.NotEmpty(t, cart)
}

func Test_GivenAnEmptyCart_WhenSize_ThenReturnZero(t *testing.T) {
	cartCustomer, _ := customer.NewCustomer("John Mayer")
	cart, _ := cart.NewCart(cartCustomer)

	cartSize := cart.Size()

	assert.Equal(t, 0, cartSize)
}

func Test_GivenANilProduct_WhenAddProductToCart_ThenReturnError(t *testing.T) {
	cartCustomer, _ := customer.NewCustomer("John Mayer")
	cart, _ := cart.NewCart(cartCustomer)

	cartItem, err := cart.AddItem(nil)

	assert.Error(t, err)
	assert.Equal(t, "invalid product", err.Error())
	assert.Empty(t, cartItem)
}

func Test_GivenAValidProduct_WhenAddProductToCart_ThenReturnAnItem(t *testing.T) {
	cartCustomer, _ := customer.NewCustomer("John Mayer")
	cart, _ := cart.NewCart(cartCustomer)

	productToAdd, _ := product.NewProduct("Arroz Blanco Gallo", 8.00)

	cartItem, err := cart.AddItem(productToAdd)

	assert.Nil(t, err)
	assert.NotEmpty(t, cartItem)
	assert.Equal(t, 1, cart.Size())
}
