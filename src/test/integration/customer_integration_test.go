package test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/bitlogic/go-startup/src/application"
	"github.com/bitlogic/go-startup/src/domain"
	"github.com/bitlogic/go-startup/src/infrastructure/config"
	"github.com/bitlogic/go-startup/src/infrastructure/controllers"
	"github.com/bitlogic/go-startup/src/infrastructure/repositories"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func Test_GivenAValidNewCustomerRequest_WhenPOSTNewCustomer_ThenReturn200(t *testing.T) {
	customerRepository := repositories.NewInMemoryCustomerRepository()
	customerService, _ := application.NewCustomerService(customerRepository)
	customerController, _ := controllers.NewCustomerController(customerService)

	request := httptest.NewRequest(http.MethodPost, "/customers", strings.NewReader(`{"customer_name":"Linus Torvalds"}`))
	request.Header.Add(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()

	e := echo.New()
	e.POST("/customers", customerController.CreateNewCustomer)
	e.Validator = config.NewRequestValidator()
	e.ServeHTTP(rec, request)

	var customerDto application.CustomerDto
	json.Unmarshal(rec.Body.Bytes(), &customerDto)

	assert.Equal(t, http.StatusCreated, rec.Code)
	assert.NotEmpty(t, customerDto.Id)
	assert.Equal(t, "Linus Torvalds", customerDto.Name)
	savedCustomer, _ := customerRepository.FindByID(domain.CustomerId(customerDto.Id))
	if assert.NotNil(t, savedCustomer) {
		assert.Equal(t, "Linus Torvalds", savedCustomer.GetName())
	}
}

func Test_GivenAnInvalidNewCustomerRequest_WhenPOSTNewCustomer_ThenReturn400ErrorResponse(t *testing.T) {
	tests := []struct {
		testName             string
		requestBody          string
		expectedResponseBody string
		expectedResponseCode int
	}{
		{
			testName:             "customer name too short",
			requestBody:          `{"customer_name":"Linus"}`,
			expectedResponseBody: `{"message":"there were validation errors","validation_errors":[{"field":"CustomerName","error":"CustomerName must be at least 8 characters in length"}]}`,
			expectedResponseCode: http.StatusBadRequest,
		},
		{
			testName:             "customer name filled with whitespaces",
			requestBody:          `{"customer_name":"Linus       "}`,
			expectedResponseBody: `{"message":"invalid name"}`,
			expectedResponseCode: http.StatusInternalServerError,
		},
		{
			testName:             "customer name is nil",
			requestBody:          `{}`,
			expectedResponseBody: `{"message":"there were validation errors","validation_errors":[{"field":"CustomerName","error":"CustomerName is a required field"}]}`,
			expectedResponseCode: http.StatusBadRequest,
		},
		{
			testName:             "customer name is of invalid type",
			requestBody:          `{"customer_name":123}`,
			expectedResponseBody: `{"message":"Unmarshal type error: expected=string, got=number, field=customer_name, offset=20"}`,
			expectedResponseCode: http.StatusBadRequest,
		},
	}

	for _, tc := range tests {
		t.Run(tc.testName, func(t *testing.T) {
			customerRepository := repositories.NewInMemoryCustomerRepository()
			customerService, _ := application.NewCustomerService(customerRepository)
			customerController, _ := controllers.NewCustomerController(customerService)

			request := httptest.NewRequest(http.MethodPost, "/customers", strings.NewReader(tc.requestBody))
			request.Header.Add(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()

			e := echo.New()
			e.POST("/customers", customerController.CreateNewCustomer)
			e.Validator = config.NewRequestValidator()
			e.ServeHTTP(rec, request)

			assert.Equal(t, tc.expectedResponseCode, rec.Code)
			assert.Equal(t, tc.expectedResponseBody, strings.Trim(rec.Body.String(), "\n"))
		})

	}

}
