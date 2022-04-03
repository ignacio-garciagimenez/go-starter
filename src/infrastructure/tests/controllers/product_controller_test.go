package tests

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/bitlogic/go-startup/src/application"
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
	controller, _ := controllers.NewProductController(&productServiceMock{
		createNewProduct: func(command application.CreateProductCommand) (application.ProductDto, error) {
			return application.ProductDto{
				Id:        newProductId,
				Name:      command.ProductName,
				UnitPrice: application.PriceDto(command.UnitPrice),
			}, nil
		},
	})

	e := echo.New()
	e.Validator = controllers.NewRequestValidator()
	request := httptest.NewRequest(http.MethodPost, "/products", strings.NewReader(`{"product_name":"Pepsi Light 2.5Lt","unit_price":10.00}`))
	request.Header.Add(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(request, rec)

	if assert.NoError(t, controller.CreateNewProduct(c)) {
		assert.Equal(t, http.StatusCreated, rec.Code)
		assert.Equal(t, fmt.Sprintf("{\"id\":\"%s\",\"name\":\"Pepsi Light 2.5Lt\",\"unit_price\":10.00}\n", newProductId.String()), rec.Body.String())
	}
}

func Test_GivenANewProductRequesTWithInvalidPrice_WhenCreateNewProduct_ThenReturn400Error(t *testing.T) {
	controller, _ := controllers.NewProductController(&productServiceMock{})

	e := echo.New()
	e.Validator = controllers.NewRequestValidator()
	//e.HTTPErrorHandler = controllers.CustomValidationErrorHanlder(e.DefaultHTTPErrorHandler)
	request := httptest.NewRequest(http.MethodPost, "/products", strings.NewReader(`{"product_name":"Pepsi Light 2.5Lt","unit_price":0}`))
	request.Header.Add(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(request, rec)

	if assert.Error(t, controller.CreateNewProduct(c)) {
		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.Equal(t, `{"validation_errors":[{"field":"unit_price","errors":["price should be greater than 0.00"]}]}`, rec.Body.String())
	}
}

type productServiceMock struct {
	callCount        int
	createNewProduct func(application.CreateProductCommand) (application.ProductDto, error)
}

func (s *productServiceMock) CreateNewProduct(command application.CreateProductCommand) (application.ProductDto, error) {
	s.callCount++
	return s.createNewProduct(command)
}
