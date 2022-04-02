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
	if assert.NotEmpty(t, cart.GetDomainEvents()) {
		assert.Equal(t, 1, len(cart.GetDomainEvents()))
	}
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
	cart.ClearDomainEvents()

	cartItem, err := cart.AddItem(nil, 1)

	assert.Error(t, err)
	assert.Equal(t, "invalid product", err.Error())
	assert.Empty(t, cartItem)
	assert.Empty(t, cart.GetDomainEvents())
}

func Test_GivenAValidProduct_WhenAddProductToCart_ThenReturnAnItem(t *testing.T) {
	cartCustomer, _ := domain.NewCustomer("John Mayer")
	cart, _ := domain.NewCart(cartCustomer)
	cart.ClearDomainEvents()

	productToAdd, _ := domain.NewProduct("Arroz Blanco Gallo", 8.00)

	cartItem, err := cart.AddItem(productToAdd, 1)

	assert.Nil(t, err)
	assert.NotEmpty(t, cartItem)
	assert.Equal(t, 1, cart.Size())
	if assert.NotEmpty(t, cart.GetDomainEvents()) {
		assert.Equal(t, 1, len(cart.GetDomainEvents()))
	}
}

func Test_GivenAnInvalidQuantity_WhenAddProductToCart_ThenReturnError(t *testing.T) {
	cartCustomer, _ := domain.NewCustomer("John Mayer")
	cart, _ := domain.NewCart(cartCustomer)
	cart.ClearDomainEvents()

	productToAdd, _ := domain.NewProduct("Arroz Blanco Gallo", 8.00)

	cartItem, err := cart.AddItem(productToAdd, 0)

	assert.Error(t, err)
	assert.Equal(t, "invalid quantity", err.Error())
	assert.Empty(t, cartItem)
	assert.Equal(t, 0, cart.Size())
	assert.Empty(t, cart.GetDomainEvents())
}

func Test_GivenACartWithAnItem_WhenAddTheSameItemToCart_ThenUpdateQuantity(t *testing.T) {
	cartCustomer, _ := domain.NewCustomer("John Mayer")
	cart, _ := domain.NewCart(cartCustomer)
	productToAdd, _ := domain.NewProduct("Arroz Blanco Gallo", 8.00)
	cart.AddItem(productToAdd, 1)
	cart.ClearDomainEvents()

	cartItem, err := cart.AddItem(productToAdd, 2)

	assert.Nil(t, err)
	assert.NotEmpty(t, cartItem)
	assert.Equal(t, 3, cart.Size())
	if assert.NotEmpty(t, cart.GetDomainEvents()) {
		assert.Equal(t, 1, len(cart.GetDomainEvents()))
	}
}

func Test_GivenAnEmptyCart_WhenGetTotal_ThenReturnZero(t *testing.T) {
	cartCustomer, _ := domain.NewCustomer("John Mayer")
	cart, _ := domain.NewCart(cartCustomer)
	cart.ClearDomainEvents()

	price := cart.GetTotal()

	assert.Equal(t, float64(0), price)
	assert.Empty(t, cart.GetDomainEvents())
}

func Test_GivenANonEmptyCart_WhenGetTotal_ThenReturnCorrectTotal(t *testing.T) {
	cartCustomer, _ := domain.NewCustomer("John Mayer")
	productToAdd, _ := domain.NewProduct("Arroz Blanco Gallo", 8.10)
	anotherProductToAdd, _ := domain.NewProduct("Pepsi 2.5Lt", 12.00)

	cart, _ := domain.NewCart(cartCustomer)
	cart.AddItem(productToAdd, 1)
	cart.AddItem(productToAdd, 2)
	cart.AddItem(anotherProductToAdd, 2)
	cart.ClearDomainEvents()

	total := cart.GetTotal()

	assert.Equal(t, 48.30, total)
	assert.Empty(t, cart.GetDomainEvents())
}

func Test_GivenACart_WhenEqualsToItself_ThenReturnsTrue(t *testing.T) {
	cartCustomer, _ := domain.NewCustomer("John Mayer")
	cart, _ := domain.NewCart(cartCustomer)

	assert.True(t, cart.EqualsTo(cart))
}

func Test_GivenACart_WhenEqualsToAnotherCart_ThenReturnsFalse(t *testing.T) {
	cartCustomer, _ := domain.NewCustomer("John Mayer")
	cart, _ := domain.NewCart(cartCustomer)
	cart2, _ := domain.NewCart(cartCustomer)

	assert.False(t, cart.EqualsTo(cart2))
}

func Test_GivenACart_WhenEqualsToAnotherTypeOfEntity_ThenReturnsFalse(t *testing.T) {
	cartCustomer, _ := domain.NewCustomer("John Mayer")
	product, _ := domain.NewProduct("Pepsi Ligh", 10.00)
	cart, _ := domain.NewCart(cartCustomer)

	assert.False(t, cart.EqualsTo(cartCustomer))
	assert.False(t, cart.EqualsTo(product))
}
