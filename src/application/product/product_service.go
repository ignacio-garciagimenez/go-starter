package application

import (
	"errors"

	"github.com/bitlogic/go-startup/src/domain/product"
)

type ProductService struct {
	repository product.Repository
}

func NewProductService(repository product.Repository) (*ProductService, error) {
	if repository == nil {
		return nil, errors.New("repository was nil")
	}

	return &ProductService{
		repository: repository,
	}, nil
}

func (s *ProductService) CreateNewProduct(command CreateProductCommand) (CreateProductResult, error) {
	newProduct, err := product.NewProduct(command.ProductName, command.UnitPrice)
	if err != nil {
		return CreateProductResult{}, err
	}

	s.repository.Save(newProduct)

	return CreateProductResult{
		Id:          newProduct.GetID(),
		ProductName: newProduct.GetName(),
		UnitPrice:   newProduct.GetPrice(),
	}, nil
}
