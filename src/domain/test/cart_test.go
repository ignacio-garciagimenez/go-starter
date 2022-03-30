package test

import (
	"github.com/bitlogic/go-startup/src/domain/cart"
	"github.com/bitlogic/go-startup/src/domain/product"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_GivenNoItems_WhenNewCart_ThenReturnError(t *testing.T) {
	cart, err := cart.NewCart(nil)

	assert.Error(t, err, "err should not be nil")
	assert.Equal(t, "no items provided", err.Error())
	assert.Nil(t, cart, "cart should be nil")
}

func Test_GivenAnEmptyItemsSlice_WhenNewCart_ThenReturnError(t *testing.T) {
	cart, err := cart.NewCart([]*product.Product{})

	assert.Error(t, err)
	assert.Equal(t, "no items provided", err.Error())
	assert.Nil(t, cart, "Cart should be nil")
}

func Test_GivenAnItem_WhenNewCart_ThenReturnNewCart(t *testing.T) {
	startingProduct, _ := product.NewProduct("Arroz Yamaní", 10.00)
	cart, err := cart.NewCart([]*product.Product{startingProduct})

	assert.Nil(t, err)
	assert.NotNil(t, cart)
}

func Test_GivenACartWithOneItem_WhenSize_ThenReturnOne(t *testing.T) {
	startingProduct, _ := product.NewProduct("Arroz Yamaní", 10.00)
	cart, _ := cart.NewCart([]*product.Product{startingProduct})

	cartSize := cart.Size()

	assert.Equal(t, 1, cartSize)
}

func Test_GivenANilProduct_WhenAddProductToCart_ThenReturnError(t *testing.T) {
	startingProduct, _ := product.NewProduct("Arroz Yamaní", 10.00)
	cart, _ := cart.NewCart([]*product.Product{startingProduct})

	cartItem, err := cart.AddItem(nil)

	assert.Error(t, err)
	assert.Equal(t, "invalid product", err.Error())
	assert.Empty(t, cartItem)
}

func Test_GivenAValidProduct_WhenAddProductToCart_ThenReturnAnItem(t *testing.T) {
	startingProduct, _ := product.NewProduct("Arroz Yamaní", 10.00)
	cart, _ := cart.NewCart([]*product.Product{startingProduct})

	productToAdd, _ := product.NewProduct("Arroz Blanco Gallo", 8.00)

	cartItem, err := cart.AddItem(productToAdd)

	assert.Nil(t, err)
	assert.NotEmpty(t, cartItem)
	assert.Equal(t, 2, cart.Size())
}
