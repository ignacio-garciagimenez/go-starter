package repositories

import (
	"github.com/bitlogic/go-startup/src/domain/product"
	"github.com/bitlogic/go-startup/src/infrastructure/repositories"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_GivenNothing_WhenNewInMemoryProductRepository_ThenReturnAProductRepository(t *testing.T) {
	repo := repositories.NewInMemoryProductRepository()

	assert.NotNil(t, repo)
	assert.NotEmpty(t, repo)
}

func Test_GivenAProductRepository_WhenSave_ThenSaves(t *testing.T) {
	repo := repositories.NewInMemoryProductRepository()
	productToSave, _ := product.NewProduct("Arroz con mani", 10.00)

	repo.Save(productToSave)
	customerSaved, err := repo.FindByID(productToSave.GetID())

	assert.Nil(t, err)
	assert.NotEmpty(t, customerSaved)

}