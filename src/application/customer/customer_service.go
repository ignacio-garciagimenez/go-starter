package application

import (
	"errors"

	"github.com/bitlogic/go-startup/src/application"
	"github.com/bitlogic/go-startup/src/domain/customer"
)

type CustomerService struct {
	repository customer.Repository
}

func NewCustomerService(repository customer.Repository) (*CustomerService, error) {
	if repository == nil {
		return nil, errors.New("customer repository was nil")
	}

	return &CustomerService{
		repository: repository,
	}, nil
}

func (s *CustomerService) CreateNewCustomer(command CreateCustomerCommand) (application.CustomerDto, error) {
	newCustomer, err := customer.NewCustomer(command.CustomerName)
	if err != nil {
		return application.CustomerDto{}, err
	}

	if err := s.repository.Save(newCustomer); err != nil {
		return application.CustomerDto{}, err
	}

	return application.CustomerDto{
		Id:   newCustomer.GetID(),
		Name: newCustomer.GetName(),
	}, nil
}
