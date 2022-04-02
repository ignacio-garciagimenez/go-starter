package application

import "github.com/google/uuid"

type CreateCustomerCommand struct {
	CustomerName string
}

type CreateCustomerResult struct {
	Id           uuid.UUID
	CustomerName string
}
