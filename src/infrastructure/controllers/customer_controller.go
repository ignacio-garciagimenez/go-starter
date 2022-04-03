package controllers

import (
	"errors"

	"github.com/bitlogic/go-startup/src/application"
	"github.com/labstack/echo/v4"
)

type CustomerService interface {
	CreateNewCustomer(application.CreateCustomerCommand) (application.CustomerDto, error)
}

type CustomerController struct {
	customerService CustomerService
}

func NewCustomerController(customerService CustomerService) (*CustomerController, error) {
	if customerService == nil {
		return nil, errors.New("customer service was nil")
	}

	return &CustomerController{
		customerService: customerService,
	}, nil
}

func (cc *CustomerController) CreateNewCustomer(c echo.Context) error {
	var command application.CreateCustomerCommand
	if err := c.Bind(&command); err != nil {
		return err
	}

	if err := c.Validate(command); err != nil {
		return err
	}

	customerDto, err := cc.customerService.CreateNewCustomer(command)
	if err != nil {
		return echo.NewHTTPError(500, err.Error())
	}

	return c.JSON(201, customerDto)
}
