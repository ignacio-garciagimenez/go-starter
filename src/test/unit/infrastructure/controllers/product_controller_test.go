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

func Test_GivenANilProductService_WhenNewProductController_ThenReturnError(t *testing.T) {
	controller, err := controllers.NewProductController(nil)

	assert.Nil(t, controller)
	if assert.Error(t, err) {
		assert.Equal(t, "product service was nil", err.Error())
	}
}

func Test_GivenAProductService_WhenNewProductController_ThenReturnController(t *testing.T) {
	controller, err := controllers.NewProductController(&productServiceMock{})

	assert.Nil(t, err)
	assert.NotEmpty(t, controller)

}

func Test_GivenANewProductRequest_WhenCreateNewProduct_ThenReturn201AndANewProductDto(t *testing.T) {
	newProductId := uuid.New()
	productServiceMock := &productServiceMock{
		createNewProduct: func(command application.CreateProductCommand) (application.ProductDto, error) {
			return application.ProductDto{
				Id:        newProductId,
				Name:      command.ProductName,
				UnitPrice: application.PriceDto(command.UnitPrice),
			}, nil
		},
	}
	controller, _ := controllers.NewProductController(productServiceMock)

	e := echo.New()
	e.Validator = config.NewRequestValidator()
	request := httptest.NewRequest(http.MethodPost, "/products", strings.NewReader(`{"product_name":"Pepsi Light 2.5Lt","unit_price":0.01}`))
	request.Header.Add(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(request, rec)

	if assert.NoError(t, controller.CreateNewProduct(c)) {
		assert.Equal(t, http.StatusCreated, rec.Code)
		assert.Equal(t, fmt.Sprintf("{\"id\":\"%s\",\"name\":\"Pepsi Light 2.5Lt\",\"unit_price\":0.01}\n", newProductId.String()), rec.Body.String())
	}
	assert.Equal(t, 1, productServiceMock.callCount)
}

func Test_GivenANewProductRequestWithNoPrice_WhenCreateNewProduct_ThenReturn400Error(t *testing.T) {
	productServiceMock := &productServiceMock{}
	controller, _ := controllers.NewProductController(productServiceMock)

	e := echo.New()
	e.Validator = config.NewRequestValidator()
	request := httptest.NewRequest(http.MethodPost, "/products", strings.NewReader(`{"product_name":"Pepsi Light 2.5Lt"}`))
	request.Header.Add(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(request, rec)

	err := controller.CreateNewProduct(c)
	if assert.Error(t, err) {
		err := err.(*echo.HTTPError)
		assert.Equal(t, http.StatusBadRequest, err.Code)
		assert.Equal(t, &config.ValidationErrorsResponse{
			Message: "there were validation errors",
			Errors: []config.FieldError{
				{
					Field: "UnitPrice",
					Error: "UnitPrice is a required field",
				},
			},
		}, err.Message)
	}
	assert.Equal(t, 0, productServiceMock.callCount)
}

func Test_GivenANewProductRequestWithInvalidPriceAndInvalidName_WhenCreateNewProduct_ThenReturn400Error(t *testing.T) {
	productServiceMock := &productServiceMock{}
	controller, _ := controllers.NewProductController(productServiceMock)

	e := echo.New()
	e.Validator = config.NewRequestValidator()
	request := httptest.NewRequest(http.MethodPost, "/products", strings.NewReader(`{"product_name":"Pepsi","unit_price":0}`))
	request.Header.Add(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(request, rec)

	err := controller.CreateNewProduct(c)
	if assert.Error(t, err) {
		err := err.(*echo.HTTPError)
		assert.Equal(t, http.StatusBadRequest, err.Code)
		assert.Equal(t, &config.ValidationErrorsResponse{
			Message: "there were validation errors",
			Errors: []config.FieldError{
				{
					Field: "ProductName",
					Error: "ProductName must be at least 10 characters in length",
				},
				{
					Field: "UnitPrice",
					Error: "UnitPrice is a required field",
				},
			},
		}, err.Message)
	}
	assert.Equal(t, 0, productServiceMock.callCount)
}

func Test_GivenANewProductRequestWithInvalidPriceAndNoName_WhenCreateNewProduct_ThenReturn400Error(t *testing.T) {
	productServiceMock := &productServiceMock{}
	controller, _ := controllers.NewProductController(productServiceMock)

	e := echo.New()
	e.Validator = config.NewRequestValidator()
	request := httptest.NewRequest(http.MethodPost, "/products", strings.NewReader(`{"unit_price":0}`))
	request.Header.Add(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(request, rec)

	err := controller.CreateNewProduct(c)
	if assert.Error(t, err) {
		err := err.(*echo.HTTPError)
		assert.Equal(t, http.StatusBadRequest, err.Code)
		assert.Equal(t, &config.ValidationErrorsResponse{
			Message: "there were validation errors",
			Errors: []config.FieldError{
				{
					Field: "ProductName",
					Error: "ProductName is a required field",
				},
				{
					Field: "UnitPrice",
					Error: "UnitPrice is a required field",
				},
			},
		}, err.Message)
	}
	assert.Equal(t, 0, productServiceMock.callCount)
}

func Test_GivenANewProductRequestWithNegativePriceAndNoName_WhenCreateNewProduct_ThenReturn400Error(t *testing.T) {
	productServiceMock := &productServiceMock{}
	controller, _ := controllers.NewProductController(productServiceMock)

	e := echo.New()
	e.Validator = config.NewRequestValidator()
	request := httptest.NewRequest(http.MethodPost, "/products", strings.NewReader(`{"unit_price":-1}`))
	request.Header.Add(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(request, rec)

	err := controller.CreateNewProduct(c)
	if assert.Error(t, err) {
		err := err.(*echo.HTTPError)
		assert.Equal(t, http.StatusBadRequest, err.Code)
		assert.Equal(t, &config.ValidationErrorsResponse{
			Message: "there were validation errors",
			Errors: []config.FieldError{
				{
					Field: "ProductName",
					Error: "ProductName is a required field",
				},
				{
					Field: "UnitPrice",
					Error: "UnitPrice must be greater than 0",
				},
			},
		}, err.Message)
	}
	assert.Equal(t, 0, productServiceMock.callCount)

}

func Test_GivenANewProductRequestWithInvalidUnitPriceType_WhenCreateNewProduct_ThenReturn400Error(t *testing.T) {
	productServiceMock := &productServiceMock{}
	controller, _ := controllers.NewProductController(productServiceMock)

	e := echo.New()
	e.Validator = config.NewRequestValidator()
	request := httptest.NewRequest(http.MethodPost, "/products", strings.NewReader(`{"unit_price":"123"}`))
	request.Header.Add(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(request, rec)

	err := controller.CreateNewProduct(c)
	if assert.Error(t, err) {
		err := err.(*echo.HTTPError)
		assert.Equal(t, http.StatusBadRequest, err.Code)
		assert.Equal(t, string("Unmarshal type error: expected=float64, got=string, field=unit_price, offset=19"), err.Message)
	}
	assert.Equal(t, 0, productServiceMock.callCount)

}

func Test_GivenANewProductRequestWithInvalidNameType_WhenCreateNewProduct_ThenReturn400Error(t *testing.T) {
	productServiceMock := &productServiceMock{}
	controller, _ := controllers.NewProductController(productServiceMock)

	e := echo.New()
	e.Validator = config.NewRequestValidator()
	request := httptest.NewRequest(http.MethodPost, "/products", strings.NewReader(`{"product_name":123}`))
	request.Header.Add(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(request, rec)

	err := controller.CreateNewProduct(c)
	if assert.Error(t, err) {
		err := err.(*echo.HTTPError)
		assert.Equal(t, http.StatusBadRequest, err.Code)
		assert.Equal(t, string("Unmarshal type error: expected=string, got=number, field=product_name, offset=19"), err.Message)
	}
	assert.Equal(t, 0, productServiceMock.callCount)

}

func Test_GivenProductServiceFailsToCreateNewProduct_WhenCreateNewProduct_ThenReturnError(t *testing.T) {
	productServiceMock := &productServiceMock{
		createNewProduct: func(application.CreateProductCommand) (application.ProductDto, error) {
			return application.ProductDto{}, errors.New("failed to create new product")
		},
	}
	controller, _ := controllers.NewProductController(productServiceMock)

	e := echo.New()
	e.Validator = config.NewRequestValidator()
	request := httptest.NewRequest(http.MethodPost, "/products", strings.NewReader(`{"product_name":"Pepsi Light 2.5Lt","unit_price":0.01}`))
	request.Header.Add(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(request, rec)

	err := controller.CreateNewProduct(c)
	if assert.Error(t, err) {
		err := err.(*echo.HTTPError)
		assert.Equal(t, http.StatusInternalServerError, err.Code)
		assert.Equal(t, "failed to create new product", err.Message)
	}
	assert.Equal(t, 1, productServiceMock.callCount)

}

type productServiceMock struct {
	callCount        int
	createNewProduct func(application.CreateProductCommand) (application.ProductDto, error)
}

func (s *productServiceMock) CreateNewProduct(command application.CreateProductCommand) (application.ProductDto, error) {
	s.callCount++
	return s.createNewProduct(command)
}
