package test

import (
	"errors"
	"testing"

	"github.com/bitlogic/go-startup/src/application"
	"github.com/bitlogic/go-startup/src/domain"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func Test_GivenANilCartRepository_WhenNewCartService_ThenReturnError(t *testing.T) {
	service, err := application.NewCartService(nil, &customerRepositoryMock{}, &productRepositoryMock{})

	if assert.Error(t, err) {
		assert.Equal(t, "cart repository was nil", err.Error())
	}
	assert.Nil(t, service)
}

func Test_GivenANilCustomerRepository_WhenNewCartService_ThenReturnError(t *testing.T) {
	service, err := application.NewCartService(&cartRepositoryMock{}, nil, &productRepositoryMock{})

	if assert.Error(t, err) {
		assert.Equal(t, "customer repository was nil", err.Error())
	}
	assert.Nil(t, service)
}

func Test_GivenANilProductRepository_WhenNewCartService_ThenReturnError(t *testing.T) {
	service, err := application.NewCartService(&cartRepositoryMock{}, &customerRepositoryMock{}, nil)

	if assert.Error(t, err) {
		assert.Equal(t, "product repository was nil", err.Error())
	}
	assert.Nil(t, service)
}

func Test_GivenAllRepositories_WhenNewCartService_ThenReturnACartService(t *testing.T) {
	service, err := application.NewCartService(&cartRepositoryMock{}, &customerRepositoryMock{}, &productRepositoryMock{})

	assert.Nil(t, err)
	assert.NotEmpty(t, service)
}

func Test_GivenAValidCreateCartCommand_WhenCreateNewCart_ThenReturnACreateCartResponse(t *testing.T) {
	cartRepository := &cartRepositoryMock{
		save: func(cart *domain.Cart) error {
			return nil
		},
	}
	savedCustomer, _ := domain.NewCustomer("Grady Booch")
	customerRepository := &customerRepositoryMock{
		findById: func(customerId uuid.UUID) (*domain.Customer, error) {
			return savedCustomer, nil
		},
	}
	service, _ := application.NewCartService(cartRepository, customerRepository, &productRepositoryMock{})
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

func Test_GivenACreateCartCommandWithNonExistinantCustomerId_WhenCreateNewCart_ThenReturnError(t *testing.T) {
	cartRepository := &cartRepositoryMock{
		save: func(cart *domain.Cart) error {
			return nil
		},
	}
	customerRepository := &customerRepositoryMock{
		findById: func(customerId uuid.UUID) (*domain.Customer, error) {
			return nil, errors.New("entity not found")
		},
	}
	service, _ := application.NewCartService(cartRepository, customerRepository, &productRepositoryMock{})
	command := application.CreateCartCommand{
		CustomerId: uuid.New(),
	}

	result, err := service.CreateNewCart(command)

	assert.Empty(t, result)
	if assert.Error(t, err) {
		assert.Equal(t, "entity not found", err.Error())
	}
	assert.Equal(t, 1, customerRepository.callCount)
	assert.Equal(t, 0, cartRepository.callCount)
}

func Test_GivenCartRepositoryFailsToSaveNewlyCreatedCart_WhenCreateNewCart_ThenReturnError(t *testing.T) {
	cartRepository := &cartRepositoryMock{
		save: func(cart *domain.Cart) error {
			return errors.New("failed to save entity")
		},
	}
	savedCustomer, _ := domain.NewCustomer("Grady Booch")
	customerRepository := &customerRepositoryMock{
		findById: func(customerId uuid.UUID) (*domain.Customer, error) {
			return savedCustomer, nil
		},
	}
	service, _ := application.NewCartService(cartRepository, customerRepository, &productRepositoryMock{})
	command := application.CreateCartCommand{
		CustomerId: savedCustomer.GetID(),
	}

	result, err := service.CreateNewCart(command)

	assert.Empty(t, result)
	if assert.Error(t, err) {
		assert.Equal(t, "failed to save entity", err.Error())
	}
	assert.Equal(t, 1, customerRepository.callCount)
	assert.Equal(t, 1, cartRepository.callCount)
}

func Test_GivenCustomerRepositoryReturnsNilCustomerAndNilError_WhenCreateNewCart_ThenReturnError(t *testing.T) {
	cartRepository := &cartRepositoryMock{
		save: func(cart *domain.Cart) error {
			return nil
		},
	}
	customerRepository := &customerRepositoryMock{
		findById: func(customerId uuid.UUID) (*domain.Customer, error) {
			return nil, nil
		},
	}
	service, _ := application.NewCartService(cartRepository, customerRepository, &productRepositoryMock{})
	command := application.CreateCartCommand{
		CustomerId: uuid.New(),
	}

	result, err := service.CreateNewCart(command)

	assert.Empty(t, result)
	if assert.Error(t, err) {
		assert.Equal(t, "no customer provided", err.Error())
	}
	assert.Equal(t, 1, customerRepository.callCount)
	assert.Equal(t, 0, cartRepository.callCount)
}

func Test_GivenACart_WhenAddItemToCart_ThenTheItemIsAddedToTheCart(t *testing.T) {
	vaughnVernon, _ := domain.NewCustomer("Vaughn Vernon")
	vaughnVernonsCart, _ := domain.NewCart(vaughnVernon)
	productVaughnVernonWantsToAdd, _ := domain.NewProduct("Implementing Domain Driven Design Book", 50.00)

	productRepository := &productRepositoryMock{
		findByID: func(productId uuid.UUID) (*domain.Product, error) {
			return productVaughnVernonWantsToAdd, nil
		},
	}

	cartRepository := &cartRepositoryMock{
		findById: func(cartId uuid.UUID) (*domain.Cart, error) {
			return vaughnVernonsCart, nil
		},
		save: func(cart *domain.Cart) error {
			return nil
		},
	}
	service, _ := application.NewCartService(cartRepository, &customerRepositoryMock{}, productRepository)
	command := application.AddItemToCartCommand{
		CartId:    vaughnVernonsCart.GetID(),
		ProductId: productVaughnVernonWantsToAdd.GetID(),
		Quantity:  1,
	}

	result, err := service.AddItemToCart(command)

	assert.Nil(t, err)
	if assert.NotEmpty(t, result) {
		assert.Equal(t, vaughnVernonsCart.GetID(), result.Id)
		assert.Equal(t, vaughnVernonsCart.GetCustomerID(), result.CustomerId)
		if assert.NotEmpty(t, result.Items) {
			assert.Equal(t, 1, result.Items[0].Quantity)
			assert.Equal(t, 50.00, result.Items[0].UnitPrice)
			assert.Equal(t, productVaughnVernonWantsToAdd.GetID(), result.Items[0].ProductId)
		}
	}
	assert.Equal(t, 2, cartRepository.callCount)
	assert.Equal(t, 1, productRepository.callCount)
}

func Test_GivenAnInvalidQuantity_WhenAddItemToCart_ThenReturnError(t *testing.T) {
	vaughnVernon, _ := domain.NewCustomer("Vaughn Vernon")
	vaughnVernonsCart, _ := domain.NewCart(vaughnVernon)
	productVaughnVernonWantsToAdd, _ := domain.NewProduct("Implementing Domain Driven Design Book", 50.00)

	productRepository := &productRepositoryMock{
		findByID: func(productId uuid.UUID) (*domain.Product, error) {
			return productVaughnVernonWantsToAdd, nil
		},
	}

	cartRepository := &cartRepositoryMock{
		findById: func(cartId uuid.UUID) (*domain.Cart, error) {
			return vaughnVernonsCart, nil
		},
		save: func(cart *domain.Cart) error {
			return nil
		},
	}
	service, _ := application.NewCartService(cartRepository, &customerRepositoryMock{}, productRepository)
	command := application.AddItemToCartCommand{
		CartId:    vaughnVernonsCart.GetID(),
		ProductId: productVaughnVernonWantsToAdd.GetID(),
		Quantity:  0,
	}

	result, err := service.AddItemToCart(command)

	assert.Empty(t, result)
	if assert.Error(t, err) {
		assert.Equal(t, "invalid quantity", err.Error())
	}
	assert.Equal(t, 1, productRepository.callCount)
	assert.Equal(t, 1, cartRepository.callCount)
}

func Test_GivenANonExistantProduct_WhenAddItemToCart_ThenReturnError(t *testing.T) {
	vaughnVernon, _ := domain.NewCustomer("Vaughn Vernon")
	vaughnVernonsCart, _ := domain.NewCart(vaughnVernon)

	productRepository := &productRepositoryMock{
		findByID: func(productId uuid.UUID) (*domain.Product, error) {
			return nil, errors.New("product not found")
		},
	}

	cartRepository := &cartRepositoryMock{
		findById: func(cartId uuid.UUID) (*domain.Cart, error) {
			return vaughnVernonsCart, nil
		},
		save: func(cart *domain.Cart) error {
			return nil
		},
	}
	service, _ := application.NewCartService(cartRepository, &customerRepositoryMock{}, productRepository)
	command := application.AddItemToCartCommand{
		CartId:    vaughnVernonsCart.GetID(),
		ProductId: uuid.New(),
		Quantity:  0,
	}

	result, err := service.AddItemToCart(command)

	assert.Empty(t, result)
	if assert.Error(t, err) {
		assert.Equal(t, "product not found", err.Error())
	}
	assert.Equal(t, 1, productRepository.callCount)
	assert.Equal(t, 0, cartRepository.callCount)
}

func Test_GivenANonExistantCart_WhenAddItemToCart_ThenReturnError(t *testing.T) {
	productVaughnVernonWantsToAdd, _ := domain.NewProduct("Implementing Domain Driven Design Book", 50.00)

	productRepository := &productRepositoryMock{
		findByID: func(productId uuid.UUID) (*domain.Product, error) {
			return productVaughnVernonWantsToAdd, nil
		},
	}

	cartRepository := &cartRepositoryMock{
		findById: func(cartId uuid.UUID) (*domain.Cart, error) {
			return nil, errors.New("cart not found")
		},
		save: func(cart *domain.Cart) error {
			return nil
		},
	}
	service, _ := application.NewCartService(cartRepository, &customerRepositoryMock{}, productRepository)
	command := application.AddItemToCartCommand{
		CartId:    uuid.New(),
		ProductId: productVaughnVernonWantsToAdd.GetID(),
		Quantity:  1,
	}

	result, err := service.AddItemToCart(command)

	assert.Empty(t, result)
	if assert.Error(t, err) {
		assert.Equal(t, "cart not found", err.Error())
	}
	assert.Equal(t, 1, productRepository.callCount)
	assert.Equal(t, 1, cartRepository.callCount)
}

func Test_GivenCartRepositoryFailsToSave_WhenAddItemToCart_ThenReturnError(t *testing.T) {
	vaughnVernon, _ := domain.NewCustomer("Vaughn Vernon")
	vaughnVernonsCart, _ := domain.NewCart(vaughnVernon)
	productVaughnVernonWantsToAdd, _ := domain.NewProduct("Implementing Domain Driven Design Book", 50.00)

	productRepository := &productRepositoryMock{
		findByID: func(productId uuid.UUID) (*domain.Product, error) {
			return productVaughnVernonWantsToAdd, nil
		},
	}

	cartRepository := &cartRepositoryMock{
		findById: func(cartId uuid.UUID) (*domain.Cart, error) {
			return vaughnVernonsCart, nil
		},
		save: func(cart *domain.Cart) error {
			return errors.New("failed to save cart")
		},
	}
	service, _ := application.NewCartService(cartRepository, &customerRepositoryMock{}, productRepository)
	command := application.AddItemToCartCommand{
		CartId:    vaughnVernonsCart.GetID(),
		ProductId: productVaughnVernonWantsToAdd.GetID(),
		Quantity:  1,
	}

	result, err := service.AddItemToCart(command)

	assert.Empty(t, result)
	if assert.Error(t, err) {
		assert.Equal(t, "failed to save cart", err.Error())
	}
	assert.Equal(t, 1, productRepository.callCount)
	assert.Equal(t, 2, cartRepository.callCount)
}
