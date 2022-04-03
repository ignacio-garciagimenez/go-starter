package test

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/bitlogic/go-startup/src/application"
	"github.com/bitlogic/go-startup/src/infrastructure/config"
	"github.com/bitlogic/go-startup/src/infrastructure/controllers"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func Test_GivenANilCustomerService_WhenNewCustomerController_ThenReturnError(t *testing.T) {
	customerController, err := controllers.NewCustomerController(nil)

	assert.Nil(t, customerController)
	if assert.Error(t, err) {
		assert.Equal(t, "customer service was nil", err.Error())
	}
}

func Test_GivenACustomerService_WhenNewCustomerController_ThenReturnCustomerController(t *testing.T) {
	customerController, err := controllers.NewCustomerController(&customerServiceMock{})

	assert.Nil(t, err)
	assert.NotEmpty(t, customerController)
}

func Test_GivenACreateCustomerRequest_WhenCreateNewCustomer_ThenReturn200AndCustomerDto(t *testing.T) {
	newCustomerId := uuid.New()
	customerServiceMock := &customerServiceMock{
		createNewCustomer: func(_ application.CreateCustomerCommand) (application.CustomerDto, error) {
			return application.CustomerDto{
				Id:   newCustomerId,
				Name: "Martin Fowler",
			}, nil
		},
	}
	customerController, _ := controllers.NewCustomerController(customerServiceMock)

	e := echo.New()
	e.Validator = config.NewRequestValidator()
	request := httptest.NewRequest(http.MethodPost, "/products", strings.NewReader(`{"customer_name":"Martin Fowler"}`))
	request.Header.Add(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(request, rec)

	if assert.NoError(t, customerController.CreateNewCustomer(c)) {
		assert.Equal(t, http.StatusCreated, rec.Code)
		assert.Equal(t, fmt.Sprintf("{\"id\":\"%s\",\"name\":\"Martin Fowler\"}\n", newCustomerId.String()), rec.Body.String())
	}
	assert.Equal(t, 1, customerServiceMock.callCount)
}

func Test_GivenACreateCustomerRequestWithNoBody_WhenCreateNewCustomer_ThenReturn400(t *testing.T) {
	customerServiceMock := &customerServiceMock{}
	customerController, _ := controllers.NewCustomerController(customerServiceMock)

	e := echo.New()
	e.Validator = config.NewRequestValidator()
	request := httptest.NewRequest(http.MethodPost, "/products", strings.NewReader(``))
	request.Header.Add(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(request, rec)

	err := customerController.CreateNewCustomer(c)

	if assert.Error(t, err) {
		err := err.(*echo.HTTPError)
		assert.Equal(t, http.StatusBadRequest, err.Code)
		assert.Equal(t, &config.ValidationErrorsResponse{
			Message: "there were validation errors",
			Errors: []config.FieldError{
				{
					Field: "CustomerName",
					Error: "CustomerName is a required field",
				},
			},
		}, err.Message)
	}
	assert.Equal(t, 0, customerServiceMock.callCount)
}

func Test_GivenACreateCustomerRequestWithInvalidCustomerName_WhenCreateNewCustomer_ThenReturn400(t *testing.T) {
	customerServiceMock := &customerServiceMock{}
	customerController, _ := controllers.NewCustomerController(customerServiceMock)

	e := echo.New()
	e.Validator = config.NewRequestValidator()
	request := httptest.NewRequest(http.MethodPost, "/products", strings.NewReader(`{"customer_name":"Martino"}`))
	request.Header.Add(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(request, rec)

	err := customerController.CreateNewCustomer(c)

	if assert.Error(t, err) {
		err := err.(*echo.HTTPError)
		assert.Equal(t, http.StatusBadRequest, err.Code)
		assert.Equal(t, &config.ValidationErrorsResponse{
			Message: "there were validation errors",
			Errors: []config.FieldError{
				{
					Field: "CustomerName",
					Error: "CustomerName must be at least 8 characters in length",
				},
			},
		}, err.Message)
	}
	assert.Equal(t, 0, customerServiceMock.callCount)
}

func Test_GivenACreateCustomerRequestWithInvalidCustomerNameType_WhenCreateNewCustomer_ThenReturn400(t *testing.T) {
	customerServiceMock := &customerServiceMock{}
	customerController, _ := controllers.NewCustomerController(customerServiceMock)

	e := echo.New()
	e.Validator = config.NewRequestValidator()
	request := httptest.NewRequest(http.MethodPost, "/products", strings.NewReader(`{"customer_name":123}`))
	request.Header.Add(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(request, rec)

	err := customerController.CreateNewCustomer(c)

	if assert.Error(t, err) {
		err := err.(*echo.HTTPError)
		assert.Equal(t, http.StatusBadRequest, err.Code)
		assert.Equal(t, string("Unmarshal type error: expected=string, got=number, field=customer_name, offset=20"), err.Message)
	}
	assert.Equal(t, 0, customerServiceMock.callCount)
}

func Test_GivenAValidCreateCustomerRequestButCustomerServiceFailsToCreateNewCustomer_WhenCreateNewCustomer_ThenReturn500Error(t *testing.T) {
	customerServiceMock := &customerServiceMock{
		createNewCustomer: func(_ application.CreateCustomerCommand) (application.CustomerDto, error) {
			return application.CustomerDto{}, errors.New("failed to create new customer")
		},
	}
	customerController, _ := controllers.NewCustomerController(customerServiceMock)

	e := echo.New()
	e.Validator = config.NewRequestValidator()
	request := httptest.NewRequest(http.MethodPost, "/products", strings.NewReader(`{"customer_name":"Martin Fowler"}`))
	request.Header.Add(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(request, rec)

	err := customerController.CreateNewCustomer(c)

	if assert.Error(t, err) {
		err := err.(*echo.HTTPError)
		assert.Equal(t, http.StatusInternalServerError, err.Code)
		assert.Equal(t, "failed to create new customer", err.Message)
	}
	assert.Equal(t, 1, customerServiceMock.callCount)
}

type customerServiceMock struct {
	callCount         int
	createNewCustomer func(application.CreateCustomerCommand) (application.CustomerDto, error)
}

func (c *customerServiceMock) CreateNewCustomer(command application.CreateCustomerCommand) (application.CustomerDto, error) {
	c.callCount++
	return c.createNewCustomer(command)
}
