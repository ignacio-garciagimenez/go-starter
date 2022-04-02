package test

import (
	"testing"

	"github.com/bitlogic/go-startup/src/domain"
	"github.com/stretchr/testify/assert"
)

func Test_GivenNoCustomer_WhenNewCart_ThenReturnError(t *testing.T) {
	cart, err := domain.NewCart(nil)

	assert.Error(t, err, "err should not be nil")
	assert.Equal(t, "no customer provided", err.Error())
	assert.Nil(t, cart, "cart should be nil")
}

func Test_GivenACustomer_WhenNewCart_ThenReturnNewCart(t *testing.T) {
	cartCustomer, _ := domain.NewCustomer("John Mayer")
	cart, err := domain.NewCart(cartCustomer)

	assert.Nil(t, err)
	assert.NotNil(t, cart)
	assert.NotEmpty(t, cart)
}

func Test_GivenAnEmptyCart_WhenSize_ThenReturnZero(t *testing.T) {
	cartCustomer, _ := domain.NewCustomer("John Mayer")
	cart, _ := domain.NewCart(cartCustomer)

	cartSize := cart.Size()

	assert.Equal(t, 0, cartSize)
}

func Test_GivenANilProduct_WhenAddProductToCart_ThenReturnError(t *testing.T) {
	cartCustomer, _ := domain.NewCustomer("John Mayer")
	cart, _ := domain.NewCart(cartCustomer)

	cartItem, err := cart.AddItem(nil, 1)

	assert.Error(t, err)
	assert.Equal(t, "invalid product", err.Error())
	assert.Empty(t, cartItem)
}

func Test_GivenAValidProduct_WhenAddProductToCart_ThenReturnAnItem(t *testing.T) {
	cartCustomer, _ := domain.NewCustomer("John Mayer")
	cart, _ := domain.NewCart(cartCustomer)

	productToAdd, _ := domain.NewProduct("Arroz Blanco Gallo", 8.00)

	cartItem, err := cart.AddItem(productToAdd, 1)

	assert.Nil(t, err)
	assert.NotEmpty(t, cartItem)
	assert.Equal(t, 1, cart.Size())
}

func Test_GivenAnInvalidQuantity_WhenAddProductToCart_ThenReturnError(t *testing.T) {
	cartCustomer, _ := domain.NewCustomer("John Mayer")
	cart, _ := domain.NewCart(cartCustomer)

	productToAdd, _ := domain.NewProduct("Arroz Blanco Gallo", 8.00)

	cartItem, err := cart.AddItem(productToAdd, 0)

	assert.Error(t, err)
	assert.Equal(t, "invalid quantity", err.Error())
	assert.Empty(t, cartItem)
	assert.Equal(t, 0, cart.Size())
}

func Test_GivenACartWithAnItem_WhenAddTheSameItemToCart_ThenUpdateQuantity(t *testing.T) {
	cartCustomer, _ := domain.NewCustomer("John Mayer")
	cart, _ := domain.NewCart(cartCustomer)
	productToAdd, _ := domain.NewProduct("Arroz Blanco Gallo", 8.00)
	cart.AddItem(productToAdd, 1)

	cartItem, err := cart.AddItem(productToAdd, 2)

	assert.Nil(t, err)
	assert.NotEmpty(t, cartItem)
	assert.Equal(t, 3, cart.Size())
}

func Test_GivenAnEmptyCart_WhenGetTotal_ThenReturnZero(t *testing.T) {
	cartCustomer, _ := domain.NewCustomer("John Mayer")
	cart, _ := domain.NewCart(cartCustomer)

	price := cart.GetTotal()

	assert.Equal(t, float64(0), price)
}

func Test_GivenANonEmptyCart_WhenGetTotal_ThenReturnCorrectTotal(t *testing.T) {
	cartCustomer, _ := domain.NewCustomer("John Mayer")
	productToAdd, _ := domain.NewProduct("Arroz Blanco Gallo", 8.10)
	anotherProductToAdd, _ := domain.NewProduct("Pepsi 2.5Lt", 12.00)

	cart, _ := domain.NewCart(cartCustomer)
	cart.AddItem(productToAdd, 1)
	cart.AddItem(productToAdd, 2)
	cart.AddItem(anotherProductToAdd, 2)

	total := cart.GetTotal()

	assert.Equal(t, 48.30, total)
}
