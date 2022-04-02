package application

import (
	"errors"

	"github.com/bitlogic/go-startup/src/application"
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

func (s *ProductService) CreateNewProduct(command CreateProductCommand) (application.ProductDto, error) {
	newProduct, err := product.NewProduct(command.ProductName, command.UnitPrice)
	if err != nil {
		return application.ProductDto{}, err
	}

	if err := s.repository.Save(newProduct); err != nil {
		return application.ProductDto{}, err
	}

	return application.ProductDto{
		Id:        newProduct.GetID(),
		Name:      newProduct.GetName(),
		UnitPrice: newProduct.GetPrice(),
	}, nil
}
