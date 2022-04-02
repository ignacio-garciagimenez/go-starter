package application

import (
	"errors"

	"github.com/bitlogic/go-startup/src/domain"
	"github.com/google/uuid"
)

type CustomerService struct {
	repository domain.CustomerRepository
}

func NewCustomerService(repository domain.CustomerRepository) (*CustomerService, error) {
	if repository == nil {
		return nil, errors.New("customer repository was nil")
	}

	return &CustomerService{
		repository: repository,
	}, nil
}

func (s *CustomerService) CreateNewCustomer(command CreateCustomerCommand) (CustomerDto, error) {
	newCustomer, err := domain.NewCustomer(command.CustomerName)
	if err != nil {
		return CustomerDto{}, err
	}

	if err := s.repository.Save(newCustomer); err != nil {
		return CustomerDto{}, err
	}

	return CustomerDto{
		Id:   uuid.UUID(newCustomer.GetID()),
		Name: newCustomer.GetName(),
	}, nil
}
