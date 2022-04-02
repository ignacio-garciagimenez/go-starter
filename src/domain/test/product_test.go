package test

import (
	"testing"

	"github.com/bitlogic/go-startup/src/domain"
	"github.com/stretchr/testify/assert"
)

func Test_GivenAnEmptyName_WhenNewProduct_ThenReturnError(t *testing.T) {
	product, err := domain.NewProduct("", 10.00)

	assert.Nil(t, product)
	assert.Error(t, err)
	assert.Equal(t, "invalid arguments", err.Error())
}

func Test_GivenANameWith9Characters_WhenNewProduct_ThenReturnError(t *testing.T) {
	product, err := domain.NewProduct("123456789", 10.00)

	assert.Nil(t, product)
	assert.Error(t, err)
	assert.Equal(t, "invalid arguments", err.Error())
}

func Test_GivenAnInvalidPrice_WhenNewProduct_ThenReturnError(t *testing.T) {
	product, err := domain.NewProduct("Arroz yamani", 0.00)

	assert.Nil(t, product)
	assert.Error(t, err)
	assert.Equal(t, "invalid arguments", err.Error())
}

func Test_GivenValidParameters_WhenNewProduct_ThenReturnAProduct(t *testing.T) {
	product, err := domain.NewProduct("Arroz yamani", 0.01)

	assert.Nil(t, err)
	assert.NotNil(t, product)
}
