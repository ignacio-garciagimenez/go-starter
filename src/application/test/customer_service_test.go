package test

import (
	"errors"
	"testing"

	application "github.com/bitlogic/go-startup/src/application/customer"
	"github.com/bitlogic/go-startup/src/domain/customer"
	"github.com/stretchr/testify/assert"
)

func Test_GivenANilCustomerRepository_WhenNewCustomerService_ThenReturnError(t *testing.T) {
	service, err := application.NewCustomerService(nil)

	if assert.Error(t, err) {
		assert.Equal(t, "customer repository was nil", err.Error())
	}
	assert.Nil(t, service)
}

func Test_GivenACustomerRepository_WhenNewCustomerService_ThenReturnACustomerService(t *testing.T) {
	service, err := application.NewCustomerService(&customerRepositoryMock{})

	assert.Nil(t, err)
	assert.NotEmpty(t, service)
}

func Test_GivenAValidCreateCustomerCommand_WhenCreateNewCustomer_ThenReturnACustomerDto(t *testing.T) {
	repository := &customerRepositoryMock{
		save: func(customer *customer.Customer) error {
			return nil
		},
	}
	service, _ := application.NewCustomerService(repository)
	customerToSave := application.CreateCustomerCommand{
		CustomerName: "Robert Smith Jr.",
	}

	result, err := service.CreateNewCustomer(customerToSave)

	assert.Nil(t, err)
	if assert.NotEmpty(t, result) {
		assert.NotNil(t, result.Id)
		assert.Equal(t, "Robert Smith Jr.", result.Name)
	}
	assert.Equal(t, 1, repository.callCount)
}

func Test_GivenACreateCustomerCommandWithInvalidName_WhenCreateNewCustomer_ThenReturnError(t *testing.T) {
	repository := &customerRepositoryMock{
		save: func(customer *customer.Customer) error {
			return nil
		},
	}
	service, _ := application.NewCustomerService(repository)
	customerToSave := application.CreateCustomerCommand{
		CustomerName: "bob",
	}

	result, err := service.CreateNewCustomer(customerToSave)

	assert.Empty(t, result)
	if assert.Error(t, err) {
		assert.Equal(t, "invalid name", err.Error())
	}
	assert.Equal(t, 0, repository.callCount)
}

func Test_GivenRepositoryFailsToSaveNewCustomer_WhenCreateNewCustomer_ThenReturnError(t *testing.T) {
	repository := &customerRepositoryMock{
		save: func(customer *customer.Customer) error {
			return errors.New("failed to save entity")
		},
	}
	service, _ := application.NewCustomerService(repository)
	customerToSave := application.CreateCustomerCommand{
		CustomerName: "Uncle Bob",
	}

	result, err := service.CreateNewCustomer(customerToSave)

	assert.Empty(t, result)
	if assert.Error(t, err) {
		assert.Equal(t, "failed to save entity", err.Error())
	}
	assert.Equal(t, 1, repository.callCount)
}
