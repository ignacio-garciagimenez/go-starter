package repositories

import (
	"testing"

	"github.com/bitlogic/go-startup/src/domain/cart"
	"github.com/bitlogic/go-startup/src/domain/customer"
	"github.com/bitlogic/go-startup/src/domain/product"
	"github.com/bitlogic/go-startup/src/infrastructure/repositories"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func Test_GivenNothing_WhenNewInMemoryCartRepository_ThenReturnACartRepository(t *testing.T) {
	repo := repositories.NewInMemoryCartRepository()

	assert.NotNil(t, repo)
	assert.NotEmpty(t, repo)
}

func Test_GivenACartRepository_WhenSave_ThenSaves(t *testing.T) {
	repo := repositories.NewInMemoryCartRepository()
	aCustomer, _ := customer.NewCustomer("John Mayer")
	aProduct, _ := product.NewProduct("Arroz con leche", 10.00)
	cartToSave, _ := cart.NewCart(aCustomer)
	cartToSave.AddItem(aProduct, 1)

	repo.Save(cartToSave)
	cartSaved, err := repo.FindByID(cartToSave.GetID())

	assert.Nil(t, err)
	assert.NotEmpty(t, cartSaved)
	assert.Equal(t, 1, cartSaved.Size())
	assert.Equal(t, 10.00, cartSaved.GetTotal())
}

func Test_GivenACartRepository_WhenGetByCustomer_ThenReturnsTheCustomersCart(t *testing.T) {
	repo := repositories.NewInMemoryCartRepository()
	aCustomer, _ := customer.NewCustomer("John Mayer")
	aProduct, _ := product.NewProduct("Arroz con leche", 10.00)
	cartToSave, _ := cart.NewCart(aCustomer)
	cartToSave.AddItem(aProduct, 1)

	repo.Save(cartToSave)
	cartsSaved := repo.GetCustomerCarts(aCustomer.GetID())

	assert.NotEmpty(t, cartsSaved)
	assert.Equal(t, 1, len(cartsSaved))
	assert.Equal(t, 1, cartsSaved[0].Size())
	assert.Equal(t, 10.00, cartsSaved[0].GetTotal())
}

func Test_GivenACartRepositoryWithOneCart_WhenFindByIDWithUnexistingID_ThenReturnsError(t *testing.T) {
	repo := repositories.NewInMemoryCartRepository()
	aCustomer, _ := customer.NewCustomer("John Mayer")
	aProduct, _ := product.NewProduct("Arroz con leche", 10.00)
	cartToSave, _ := cart.NewCart(aCustomer)
	cartToSave.AddItem(aProduct, 1)

	repo.Save(cartToSave)
	cartSaved, err := repo.FindByID(uuid.New())

	if assert.Error(t, err) {
		assert.Equal(t, "entity not found", err.Error())
	}
	assert.Nil(t, cartSaved)

}
