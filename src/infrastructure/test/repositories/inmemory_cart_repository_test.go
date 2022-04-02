package repositories

import (
	"testing"

	"github.com/bitlogic/go-startup/src/domain"
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
	aCustomer, _ := domain.NewCustomer("John Mayer")
	aProduct, _ := domain.NewProduct("Arroz con leche", 10.00)
	cartToSave, _ := domain.NewCart(aCustomer)
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
	aCustomer, _ := domain.NewCustomer("John Mayer")
	aProduct, _ := domain.NewProduct("Arroz con leche", 10.00)
	cartToSave, _ := domain.NewCart(aCustomer)
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
	aCustomer, _ := domain.NewCustomer("John Mayer")
	aProduct, _ := domain.NewProduct("Arroz con leche", 10.00)
	cartToSave, _ := domain.NewCart(aCustomer)
	cartToSave.AddItem(aProduct, 1)

	repo.Save(cartToSave)
	cartSaved, err := repo.FindByID(domain.CartId(uuid.New()))

	if assert.Error(t, err) {
		assert.Equal(t, "entity not found", err.Error())
	}
	assert.Nil(t, cartSaved)

}
