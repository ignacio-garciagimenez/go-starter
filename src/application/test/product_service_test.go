package test

import (
	"errors"
	"testing"

	"github.com/bitlogic/go-startup/src/application"
	"github.com/bitlogic/go-startup/src/domain"
	"github.com/bitlogic/go-startup/src/infrastructure/repositories"
	"github.com/stretchr/testify/assert"
)

func Test_GivenANilProductRepository_WhenNewProductService_ThenReturnError(t *testing.T) {
	productService, err := application.NewProductService(nil)

	if assert.Error(t, err) {
		assert.Equal(t, "repository was nil", err.Error())
	}
	assert.Nil(t, productService)
}

func Test_GivenAProductRepository_WhenNewProductService_ThenReturnAProductService(t *testing.T) {
	repo := repositories.NewInMemoryProductRepository()

	productService, err := application.NewProductService(repo)

	assert.Nil(t, err)
	assert.NotNil(t, productService)
}

func Test_GivenAWellFormedCreateProductCommand_WhenCreateNewProduct_ThenReturnANewProductDto(t *testing.T) {
	repositoryMock := &productRepositoryMock{
		save: func(product *domain.Product) error {
			return nil
		},
	}
	productService, _ := application.NewProductService(repositoryMock)
	createProductCommand := application.CreateProductCommand{
		ProductName: "Pepsi 2.25Lt",
		UnitPrice:   10.00,
	}

	output, err := productService.CreateNewProduct(createProductCommand)

	assert.Nil(t, err)
	if assert.NotEmpty(t, output) {
		assert.NotNil(t, output.Id)
		assert.Equal(t, "Pepsi 2.25Lt", output.Name)
		assert.Equal(t, application.PriceDto(10.00), output.UnitPrice)
	}

	assert.Equal(t, 1, repositoryMock.callCount)
}

func Test_GivenACreateProductCommandWithInvalidName_WhenCreateNewProduct_ThenReturnError(t *testing.T) {
	repositoryMock := &productRepositoryMock{
		save: func(product *domain.Product) error {
			return nil
		},
	}
	productService, _ := application.NewProductService(repositoryMock)
	createProductCommand := application.CreateProductCommand{
		ProductName: "Pepsi",
		UnitPrice:   10.00,
	}

	output, err := productService.CreateNewProduct(createProductCommand)

	if assert.Error(t, err) {
		assert.Equal(t, "invalid arguments", err.Error())
	}
	assert.Empty(t, output)
	assert.Equal(t, 0, repositoryMock.callCount)
}

func Test_GivenACreateProductCommandWithInvalidPrice_WhenCreateNewProduct_ThenReturnError(t *testing.T) {
	repositoryMock := &productRepositoryMock{
		save: func(product *domain.Product) error {
			return nil
		},
	}
	productService, _ := application.NewProductService(repositoryMock)
	createProductCommand := application.CreateProductCommand{
		ProductName: "Pepsi 2.25Lts",
		UnitPrice:   0.00,
	}

	output, err := productService.CreateNewProduct(createProductCommand)

	if assert.Error(t, err) {
		assert.Equal(t, "invalid arguments", err.Error())
	}
	assert.Empty(t, output)
	assert.Equal(t, 0, repositoryMock.callCount)
}

func Test_GivenSavingProductFails_WhenCreateNewProduct_ThenReturnError(t *testing.T) {
	repositoryMock := &productRepositoryMock{
		save: func(product *domain.Product) error {
			return errors.New("failed to save entity")
		},
	}
	productService, _ := application.NewProductService(repositoryMock)
	createProductCommand := application.CreateProductCommand{
		ProductName: "Pepsi 2.25Lts",
		UnitPrice:   10.00,
	}

	output, err := productService.CreateNewProduct(createProductCommand)

	if assert.Error(t, err) {
		assert.Equal(t, "failed to save entity", err.Error())
	}
	assert.Empty(t, output)
	assert.Equal(t, 1, repositoryMock.callCount)
}
