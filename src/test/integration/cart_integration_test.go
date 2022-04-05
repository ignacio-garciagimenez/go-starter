package test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/bitlogic/go-startup/src/application"
	"github.com/bitlogic/go-startup/src/domain"
	"github.com/bitlogic/go-startup/src/infrastructure/config"
	"github.com/bitlogic/go-startup/src/infrastructure/controllers"
	"github.com/bitlogic/go-startup/src/infrastructure/repositories"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func Test_GivenAValidNewCartRequest_WhenPOSTNewCart_ThenReturn201(t *testing.T) {
	existingCustomer, _ := domain.NewCustomer("Bjarne Stroustrup")
	cartRepository := repositories.NewInMemoryCartRepository()
	customerRepository := repositories.NewInMemoryCustomerRepository()
	productRepository := repositories.NewInMemoryProductRepository()
	cartService, _ := application.NewCartService(cartRepository, customerRepository, productRepository)
	cartController, _ := controllers.NewCartController(cartService)

	customerRepository.Save(existingCustomer)

	request := httptest.NewRequest(http.MethodPost, "/carts", strings.NewReader(
		fmt.Sprintf(`{"customer_id":"%s"}`, uuid.UUID(existingCustomer.GetID()).String())))
	request.Header.Add(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()

	e := echo.New()
	e.POST("/carts", cartController.CreateNewCart)
	e.Validator = config.NewRequestValidator()
	e.ServeHTTP(rec, request)

	var cartDto application.CartDto
	json.Unmarshal(rec.Body.Bytes(), &cartDto)

	assert.Equal(t, http.StatusCreated, rec.Code)
	assert.NotEmpty(t, cartDto.Id)
	assert.Equal(t, uuid.UUID(existingCustomer.GetID()), cartDto.CustomerId)
	assert.Empty(t, cartDto.Items)
	savedCart, _ := cartRepository.FindByID(domain.CartId(cartDto.Id))
	if assert.NotNil(t, savedCart) {
		assert.Equal(t, existingCustomer.GetID(), savedCart.GetCustomerID())
		assert.Equal(t, 0, len(savedCart.GetItems()))
	}
}

func Test_GivenAnInvalidNewCartRequest_WhenPOSTNewCart_ThenReturn400ErrorResponse(t *testing.T) {
	nonExistantCustomerId := uuid.New()

	tests := []struct {
		testName             string
		requestBody          string
		expectedResponseBody string
		expectedResponseCode int
	}{
		{
			testName:             "customer id too short",
			requestBody:          `{"customer_id":"1231231231231231231231231231231"}`,
			expectedResponseBody: `{"message":"invalid UUID length: 31"}`,
			expectedResponseCode: http.StatusBadRequest,
		},
		{
			testName:             "customer id too long",
			requestBody:          `{"customer_id":"123123123123123123123123123123133"}`,
			expectedResponseBody: `{"message":"invalid UUID length: 33"}`,
			expectedResponseCode: http.StatusBadRequest,
		},
		{
			testName:             "customer id bad uuid format",
			requestBody:          `{"customer_id":"123123123W2312312312312312312313"}`,
			expectedResponseBody: `{"message":"invalid UUID format"}`,
			expectedResponseCode: http.StatusBadRequest,
		},
		{
			testName:             "customer id is nil",
			requestBody:          `{}`,
			expectedResponseBody: `{"message":"there were validation errors","validation_errors":[{"field":"CustomerId","error":"CustomerId is a required field"}]}`,
			expectedResponseCode: http.StatusBadRequest,
		},
		{
			testName:             "customer id is of invalid type",
			requestBody:          `{"customer_id":123}`,
			expectedResponseBody: `{"message":"Unmarshal type error: expected=uuid.UUID, got=number, field=customer_id, offset=18"}`,
			expectedResponseCode: http.StatusBadRequest,
		},
		{
			testName:             "customer with id doesnt exist",
			requestBody:          fmt.Sprintf(`{"customer_id":"%s"}`, nonExistantCustomerId.String()),
			expectedResponseBody: fmt.Sprintf(`{"message":"customer with id %s not found"}`, nonExistantCustomerId.String()),
			expectedResponseCode: http.StatusNotFound,
		},
	}

	for _, tc := range tests {
		t.Run(tc.testName, func(t *testing.T) {
			cartRepository := repositories.NewInMemoryCartRepository()
			customerRepository := repositories.NewInMemoryCustomerRepository()
			productRepository := repositories.NewInMemoryProductRepository()
			cartService, _ := application.NewCartService(cartRepository, customerRepository, productRepository)
			cartController, _ := controllers.NewCartController(cartService)

			request := httptest.NewRequest(http.MethodPost, "/carts", strings.NewReader(tc.requestBody))
			request.Header.Add(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()

			e := echo.New()
			e.POST("/carts", cartController.CreateNewCart)
			e.Validator = config.NewRequestValidator()
			e.ServeHTTP(rec, request)

			assert.Equal(t, tc.expectedResponseCode, rec.Code)
			assert.Equal(t, tc.expectedResponseBody, strings.Trim(rec.Body.String(), "\n"))
		})

	}

}

func Test_GivenAValidAddItemToCartRequest_WhenPOSTAddItemToCart_ThenReturn200(t *testing.T) {
	existingCustomer, _ := domain.NewCustomer("Bjarne Stroustrup")
	existingProduct, _ := domain.NewProduct("Mortadela 1 Kg", 10.00)
	existingCart, _ := domain.NewCart(existingCustomer)

	cartRepository := repositories.NewInMemoryCartRepository()
	customerRepository := repositories.NewInMemoryCustomerRepository()
	productRepository := repositories.NewInMemoryProductRepository()
	cartService, _ := application.NewCartService(cartRepository, customerRepository, productRepository)
	cartController, _ := controllers.NewCartController(cartService)

	customerRepository.Save(existingCustomer)
	productRepository.Save(existingProduct)
	cartRepository.Save(existingCart)

	request := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/carts/%s", uuid.UUID(existingCart.GetID()).String()), strings.NewReader(
		fmt.Sprintf(`{"product_id":"%s","quantity":2}`, uuid.UUID(existingProduct.GetID()).String())))
	request.Header.Add(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()

	e := echo.New()
	e.POST("/carts/:cartId", cartController.AddItemToCart)
	e.Validator = config.NewRequestValidator()
	e.ServeHTTP(rec, request)

	var cartDto application.CartDto
	json.Unmarshal(rec.Body.Bytes(), &cartDto)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, uuid.UUID(existingCart.GetID()), cartDto.Id)
	assert.Equal(t, uuid.UUID(existingCustomer.GetID()), cartDto.CustomerId)
	if assert.NotEmpty(t, cartDto.Items) {
		assert.Equal(t, 1, len(cartDto.Items))
		assert.Equal(t, uuid.UUID(existingProduct.GetID()), cartDto.Items[0].ProductId)
		assert.Equal(t, 2, cartDto.Items[0].Quantity)
		assert.Equal(t, existingProduct.GetPrice(), float64(cartDto.Items[0].UnitPrice))
	}

	savedCart, _ := cartRepository.FindByID(domain.CartId(cartDto.Id))
	if assert.NotNil(t, savedCart) {
		assert.Equal(t, existingCustomer.GetID(), savedCart.GetCustomerID())
		if assert.Equal(t, 1, len(savedCart.GetItems())) {
			assert.Equal(t, uuid.UUID(existingProduct.GetID()), uuid.UUID(savedCart.GetItems()[0].GetProductId()))
			assert.Equal(t, 2, savedCart.GetItems()[0].GetQuantity())
			assert.Equal(t, 10.00, savedCart.GetItems()[0].GetUnitPrice())
		}
		assert.Equal(t, 20.00, savedCart.GetTotal())
	}
}

func Test_GivenAnInvalidAddItemToCartRequest_WhenPOSTAddItemToCart_ThenReturn400ErrorResponse(t *testing.T) {
	nonExistantProductId := uuid.New()
	existantProduct, _ := domain.NewProduct("Mortadela 1Kg", 10)
	cartId := uuid.New()

	tests := []struct {
		cartId               string
		testName             string
		requestBody          string
		expectedResponseBody string
		expectedResponseCode int
	}{
		{
			testName:             "cart doesnt exist",
			requestBody:          fmt.Sprintf(`{"product_id":"%s","quantity":1}`, uuid.UUID(existantProduct.GetID()).String()),
			expectedResponseBody: fmt.Sprintf(`{"message":"cart with id %s not found"}`, cartId.String()),
			expectedResponseCode: http.StatusNotFound,
		},
		{
			testName:             "product id too short",
			requestBody:          `{"product_id":"1231231231231231231231231231231","quantity":1}`,
			expectedResponseBody: `{"message":"invalid UUID length: 31"}`,
			expectedResponseCode: http.StatusBadRequest,
		},
		{
			testName:             "product id too long",
			requestBody:          `{"product_id":"123123123123123123123123123123133","quantity":1}`,
			expectedResponseBody: `{"message":"invalid UUID length: 33"}`,
			expectedResponseCode: http.StatusBadRequest,
		},
		{
			testName:             "product id bad uuid format",
			requestBody:          `{"product_id":"123123123W2312312312312312312313","quantity":1}`,
			expectedResponseBody: `{"message":"invalid UUID format"}`,
			expectedResponseCode: http.StatusBadRequest,
		},
		{
			testName:             "product id is nil",
			requestBody:          `{"quantity":1}`,
			expectedResponseBody: `{"message":"there were validation errors","validation_errors":[{"field":"ProductId","error":"ProductId is a required field"}]}`,
			expectedResponseCode: http.StatusBadRequest,
		},
		{
			testName:             "product id is of invalid type",
			requestBody:          `{"product_id":123,"quantity":1}`,
			expectedResponseBody: `{"message":"Unmarshal type error: expected=uuid.UUID, got=number, field=product_id, offset=17"}`,
			expectedResponseCode: http.StatusBadRequest,
		},
		{
			testName:             "product with id doesnt exist",
			requestBody:          fmt.Sprintf(`{"product_id":"%s","quantity":1}`, nonExistantProductId.String()),
			expectedResponseBody: fmt.Sprintf(`{"message":"product with id %s not found"}`, nonExistantProductId.String()),
			expectedResponseCode: http.StatusNotFound,
		},
		{
			testName:             "quantity nil",
			requestBody:          `{"product_id":"12312312312312312312312312312311"}`,
			expectedResponseBody: `{"message":"there were validation errors","validation_errors":[{"field":"Quantity","error":"Quantity is a required field"}]}`,
			expectedResponseCode: http.StatusBadRequest,
		},
		{
			testName:             "quantity negative",
			requestBody:          `{"product_id":"12312312312312312312312312312313","quantity":-1}`,
			expectedResponseBody: `{"message":"there were validation errors","validation_errors":[{"field":"Quantity","error":"Quantity must be greater than 0"}]}`,
			expectedResponseCode: http.StatusBadRequest,
		},
		{
			testName:             "quantity wrong type",
			requestBody:          `{"product_id":"12312312322312312312312312312313","quantity":"1"}`,
			expectedResponseBody: `{"message":"Unmarshal type error: expected=int, got=string, field=quantity, offset=63"}`,
			expectedResponseCode: http.StatusBadRequest,
		},
	}

	for _, tc := range tests {
		t.Run(tc.testName, func(t *testing.T) {
			cartRepository := repositories.NewInMemoryCartRepository()
			customerRepository := repositories.NewInMemoryCustomerRepository()
			productRepository := repositories.NewInMemoryProductRepository()
			cartService, _ := application.NewCartService(cartRepository, customerRepository, productRepository)
			cartController, _ := controllers.NewCartController(cartService)

			productRepository.Save(existantProduct)

			request := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/carts/%s", cartId.String()), strings.NewReader(tc.requestBody))
			request.Header.Add(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()

			e := echo.New()
			e.POST("/carts/:cartId", cartController.AddItemToCart)
			e.Validator = config.NewRequestValidator()
			e.ServeHTTP(rec, request)

			assert.Equal(t, tc.expectedResponseCode, rec.Code)
			assert.Equal(t, tc.expectedResponseBody, strings.Trim(rec.Body.String(), "\n"))
		})

	}

}
