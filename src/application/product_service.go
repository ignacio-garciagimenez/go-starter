package application

import (
	"errors"

	"github.com/bitlogic/go-startup/src/domain"
)

type ProductService struct {
	repository domain.ProductRepository
}

func NewProductService(repository domain.ProductRepository) (*ProductService, error) {
	if repository == nil {
		return nil, errors.New("repository was nil")
	}

	return &ProductService{
		repository: repository,
	}, nil
}

func (s *ProductService) CreateNewProduct(command CreateProductCommand) (ProductDto, error) {
	newProduct, err := domain.NewProduct(command.ProductName, command.UnitPrice)
	if err != nil {
		return ProductDto{}, err
	}

	if err := s.repository.Save(newProduct); err != nil {
		return ProductDto{}, err
	}

	return ProductDto{
		Id:        newProduct.GetID(),
		Name:      newProduct.GetName(),
		UnitPrice: newProduct.GetPrice(),
	}, nil
}
