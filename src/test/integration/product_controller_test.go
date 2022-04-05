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

func Test_GivenAValidNewProductRequest_WhenPOSTNewProduct_ThenReturn200(t *testing.T) {
	productRepository := repositories.NewInMemoryProductRepository()
	productService, _ := application.NewProductService(productRepository)
	productController, _ := controllers.NewProductController(productService)

	request := httptest.NewRequest(http.MethodPost, "/products", strings.NewReader(`{"product_name":"Pepsi Light 2.5Lt","unit_price":1.10}`))
	request.Header.Add(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()

	e := echo.New()
	e.POST("/products", productController.CreateNewProduct)
	e.Validator = config.NewRequestValidator()
	e.ServeHTTP(rec, request)

	var productDto application.ProductDto
	json.Unmarshal(rec.Body.Bytes(), &productDto)

	assert.Equal(t, http.StatusCreated, rec.Code)
	assert.NotEmpty(t, productDto.Id)
	assert.Equal(t, "Pepsi Light 2.5Lt", productDto.Name)
	assert.Equal(t, application.PriceDto(1.10), productDto.UnitPrice)
	savedProduct, _ := productRepository.FindByID(domain.ProductId(productDto.Id))
	if assert.NotNil(t, savedProduct) {
		assert.Equal(t, "Pepsi Light 2.5Lt", savedProduct.GetName())
		assert.Equal(t, 1.10, savedProduct.GetPrice())
	}
}

func Test_GivenAnInvalidNewProductRequest_WhenPOSTNewProduct_ThenReturn400ErrorResponse(t *testing.T) {
	tests := []struct {
		testName             string
		requestBody          string
		expectedResponseBody string
	}{
		{
			testName:             "product name too short",
			requestBody:          `{"product_name":"PepsiPeps","unit_price":1.10}`,
			expectedResponseBody: `{"message":"there were validation errors","validation_errors":[{"field":"ProductName","error":"ProductName must be at least 10 characters in length"}]}`,
		},
		{
			testName:             "product name is nil",
			requestBody:          `{"unit_price":1.10}`,
			expectedResponseBody: `{"message":"there were validation errors","validation_errors":[{"field":"ProductName","error":"ProductName is a required field"}]}`,
		},
		{
			testName:             "product name is of invalid type",
			requestBody:          `{"product_name":123,"unit_price":1.10}`,
			expectedResponseBody: `{"message":"Unmarshal type error: expected=string, got=number, field=product_name, offset=19"}`,
		},
		{
			testName:             "unit price is nil",
			requestBody:          `{"product_name":"Pepsi Light 2.5Lt"}`,
			expectedResponseBody: `{"message":"there were validation errors","validation_errors":[{"field":"UnitPrice","error":"UnitPrice is a required field"}]}`,
		},
		{
			testName:             "unit price is negative",
			requestBody:          `{"product_name":"Pepsi Light 2.5Lt","unit_price":-1.10}`,
			expectedResponseBody: `{"message":"there were validation errors","validation_errors":[{"field":"UnitPrice","error":"UnitPrice must be greater than 0"}]}`,
		},
		{
			testName:             "unit price is of invalid type",
			requestBody:          `{"product_name":"Pepsi Light 2.5Lt","unit_price":"-1.10"}`,
			expectedResponseBody: `{"message":"Unmarshal type error: expected=float64, got=string, field=unit_price, offset=56"}`,
		},
		{
			testName:             "multiple validation errors",
			requestBody:          `{}`,
			expectedResponseBody: `{"message":"there were validation errors","validation_errors":[{"field":"ProductName","error":"ProductName is a required field"},{"field":"UnitPrice","error":"UnitPrice is a required field"}]}`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.testName, func(t *testing.T) {
			productRepository := repositories.NewInMemoryProductRepository()
			productService, _ := application.NewProductService(productRepository)
			productController, _ := controllers.NewProductController(productService)

			request := httptest.NewRequest(http.MethodPost, "/products", strings.NewReader(tc.requestBody))
			request.Header.Add(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()

			e := echo.New()
			e.POST("/products", productController.CreateNewProduct)
			e.Validator = config.NewRequestValidator()
			e.ServeHTTP(rec, request)

			assert.Equal(t, http.StatusBadRequest, rec.Code)
			assert.Equal(t, tc.expectedResponseBody, strings.Trim(rec.Body.String(), "\n"))
		})

	}

}
