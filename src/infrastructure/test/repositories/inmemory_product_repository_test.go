package repositories

import (
	"testing"

	"github.com/bitlogic/go-startup/src/domain"
	"github.com/bitlogic/go-startup/src/infrastructure/repositories"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func Test_GivenNothing_WhenNewInMemoryProductRepository_ThenReturnAProductRepository(t *testing.T) {
	repo := repositories.NewInMemoryProductRepository()

	assert.NotNil(t, repo)
	assert.NotEmpty(t, repo)
}

func Test_GivenAProductRepository_WhenSave_ThenSaves(t *testing.T) {
	repo := repositories.NewInMemoryProductRepository()
	productToSave, _ := domain.NewProduct("Arroz con mani", 10.00)

	repo.Save(productToSave)
	productSaved, err := repo.FindByID(productToSave.GetID())

	assert.Nil(t, err)
	assert.NotEmpty(t, productSaved)

}

func Test_GivenAProductRepositoryWithItems_WhenFindByIDWithUnexistingID_ThenReturnesError(t *testing.T) {
	repo := repositories.NewInMemoryProductRepository()
	productToSave, _ := domain.NewProduct("Arroz con mani", 10.00)
	repo.Save(productToSave)

	productSaved, err := repo.FindByID(domain.ProductId(uuid.New()))

	assert.NotNil(t, err)
	assert.Equal(t, "entity not found", err.Error())
	assert.Nil(t, productSaved)
}
