package tests

import (
	"testing"

	"github.com/bitlogic/go-startup/src/domain"
	"github.com/bitlogic/go-startup/src/infrastructure/repositories"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func Test_GivenNothing_WhenNewInMemoryCustomerRepository_ThenReturnACustomerRepository(t *testing.T) {
	repo := repositories.NewInMemoryCustomerRepository()

	assert.NotNil(t, repo)
	assert.NotEmpty(t, repo)
}

func Test_GivenACustomerRepository_WhenSave_ThenSaves(t *testing.T) {
	repo := repositories.NewInMemoryCustomerRepository()
	customerToSave, _ := domain.NewCustomer("John Mayer")

	repo.Save(customerToSave)
	customerSaved, err := repo.FindByID(customerToSave.GetID())

	assert.Nil(t, err)
	assert.NotEmpty(t, customerSaved)

}

func Test_GivenACustomerRepositoryWithOneItem_WhenFinByIDWithUnexistingID_ThenReturnsError(t *testing.T) {
	repo := repositories.NewInMemoryCustomerRepository()
	customerToSave, _ := domain.NewCustomer("John Mayer")
	repo.Save(customerToSave)

	customerSaved, err := repo.FindByID(domain.CustomerId(uuid.New()))

	if assert.Error(t, err) {
		assert.Equal(t, "entity not found", err.Error())
	}
	assert.Nil(t, customerSaved)

}
