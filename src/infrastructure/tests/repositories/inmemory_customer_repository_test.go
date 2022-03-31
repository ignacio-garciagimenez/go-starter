package repositories

import (
	"github.com/bitlogic/go-startup/src/domain/customer"
	"github.com/bitlogic/go-startup/src/infrastructure/repositories"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_GivenNothing_WhenNewInMemoryCustomerRepository_ThenReturnACustomerRepository(t *testing.T) {
	repo := repositories.NewInMemoryCustomerRepository()

	assert.NotNil(t, repo)
	assert.NotEmpty(t, repo)
}

func Test_GivenACustomerRepository_WhenSave_ThenSaves(t *testing.T) {
	repo := repositories.NewInMemoryCustomerRepository()
	customerToSave, _ := customer.NewCustomer("John Mayer")

	repo.Save(customerToSave)
	customerSaved, err := repo.FindByID(customerToSave.GetID())

	assert.Nil(t, err)
	assert.NotEmpty(t, customerSaved)

}
