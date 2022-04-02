package controllers

import (
	"errors"

	"github.com/bitlogic/go-startup/src/application"
	"github.com/labstack/echo/v4"
)

type ProductService interface {
	CreateNewProduct(application.CreateProductCommand) (application.ProductDto, error)
}

type ProductController struct {
	service ProductService
}

func NewProductController(service ProductService) (*ProductController, error) {
	if service == nil {
		return nil, errors.New("product service was nil")
	}

	return &ProductController{
		service: service,
	}, nil
}

func (pc *ProductController) CreateNewProduct(c echo.Context) error {
	var command application.CreateProductCommand
	if err := c.Bind(&command); err != nil {
		return err
	}

	productDto, err := pc.service.CreateNewProduct(command)
	if err != nil {
		return err
	}

	return c.JSON(201, productDto)
}
