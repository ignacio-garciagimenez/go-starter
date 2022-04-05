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

func Test_GivenANilCartService_WhenNewCartController_ThenReturnError(t *testing.T) {
	controller, err := controllers.NewCartController(nil)

	assert.Nil(t, controller)
	if assert.Error(t, err) {
		assert.Equal(t, "cart service was nil", err.Error())
	}
}

func Test_GivenACartService_WhenNewCartController_ThenReturnACartController(t *testing.T) {
	controller, err := controllers.NewCartController(&cartServiceMock{})

	assert.NoError(t, err)
	assert.NotEmpty(t, controller)
	assert.IsType(t, &controllers.CartController{}, controller)
}

func Test_GivenAValidCreateNewCartRequest_WhenCreateNewCart_ThenReturn201AndACartDto(t *testing.T) {
	newCartId := uuid.New()
	customerId := uuid.New()
	cartServiceMock := &cartServiceMock{
		createNewCart: func(_ application.CreateCartCommand) (application.CartDto, error) {
			return application.CartDto{
				Id:         newCartId,
				CustomerId: customerId,
				Items:      []application.ItemDto{},
			}, nil
		},
	}
	controller, _ := controllers.NewCartController(cartServiceMock)

	e := echo.New()
	e.Validator = config.NewRequestValidator()
	request := httptest.NewRequest(http.MethodPost, "/carts", strings.NewReader(fmt.Sprintf(`{"customer_id":"%s"}`, customerId.String())))
	request.Header.Add(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(request, rec)

	if assert.NoError(t, controller.CreateNewCart(c)) {
		assert.Equal(t, http.StatusCreated, rec.Code)
		assert.Equal(t, fmt.Sprintf("{\"id\":\"%s\",\"customer_id\":\"%s\",\"items\":[]}\n", newCartId.String(), customerId.String()), rec.Body.String())
	}
	assert.Equal(t, 1, cartServiceMock.callCount)

}

func Test_GivenAValidCreateNewCartRequestButCartServiceFailsToCreateCart_WhenCreateNewCart_ThenReturn500(t *testing.T) {
	customerId := uuid.New()
	cartServiceMock := &cartServiceMock{
		createNewCart: func(_ application.CreateCartCommand) (application.CartDto, error) {
			return application.CartDto{}, errors.New("failed to create cart")
		},
	}
	controller, _ := controllers.NewCartController(cartServiceMock)

	e := echo.New()
	e.Validator = config.NewRequestValidator()
	request := httptest.NewRequest(http.MethodPost, "/carts", strings.NewReader(fmt.Sprintf(`{"customer_id":"%s"}`, customerId.String())))
	request.Header.Add(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(request, rec)

	err := controller.CreateNewCart(c)
	if assert.Error(t, err) {
		err := err.(*echo.HTTPError)
		assert.Equal(t, http.StatusInternalServerError, err.Code)
		assert.Equal(t, "failed to create cart", err.Message)
	}
	assert.Equal(t, 1, cartServiceMock.callCount)

}

func Test_GivenAValidCreateNewCartRequestButCustomerDoesExist_WhenCreateNewCart_ThenReturn500(t *testing.T) {
	customerId := uuid.New()
	cartServiceMock := &cartServiceMock{
		createNewCart: func(_ application.CreateCartCommand) (application.CartDto, error) {
			return application.CartDto{}, application.NewNotFoundError(customerId.String(), "customer")
		},
	}
	controller, _ := controllers.NewCartController(cartServiceMock)

	e := echo.New()
	e.Validator = config.NewRequestValidator()
	request := httptest.NewRequest(http.MethodPost, "/carts", strings.NewReader(fmt.Sprintf(`{"customer_id":"%s"}`, customerId.String())))
	request.Header.Add(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(request, rec)

	err := controller.CreateNewCart(c)
	if assert.Error(t, err) {
		err := err.(*echo.HTTPError)
		assert.Equal(t, http.StatusNotFound, err.Code)
		assert.Equal(t, fmt.Sprintf("customer with id %s not found", customerId.String()), err.Message)
	}
	assert.Equal(t, 1, cartServiceMock.callCount)

}

func Test_GivenACreateCartRequestWithWrongCustomerIdType_WhenCreateNewCart_ThenReturn400(t *testing.T) {
	cartServiceMock := &cartServiceMock{}
	controller, _ := controllers.NewCartController(cartServiceMock)

	e := echo.New()
	e.Validator = config.NewRequestValidator()
	request := httptest.NewRequest(http.MethodPost, "/carts", strings.NewReader(`{"customer_id":123}`))
	request.Header.Add(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(request, rec)

	err := controller.CreateNewCart(c)
	if assert.Error(t, err) {
		err := err.(*echo.HTTPError)
		assert.Equal(t, http.StatusBadRequest, err.Code)
		assert.Equal(t, "Unmarshal type error: expected=uuid.UUID, got=number, field=customer_id, offset=18", err.Message)
	}
	assert.Equal(t, 0, cartServiceMock.callCount)
}

func Test_GivenACreateCartRequestWithNoCustomerId_WhenCreateNewCart_ThenReturn400(t *testing.T) {
	cartServiceMock := &cartServiceMock{}
	controller, _ := controllers.NewCartController(cartServiceMock)

	e := echo.New()
	e.Validator = config.NewRequestValidator()
	request := httptest.NewRequest(http.MethodPost, "/carts", strings.NewReader(``))
	request.Header.Add(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(request, rec)

	err := controller.CreateNewCart(c)
	if assert.Error(t, err) {
		err := err.(*echo.HTTPError)
		assert.Equal(t, http.StatusBadRequest, err.Code)
		assert.Equal(t, &config.ValidationErrorsResponse{
			Message: "there were validation errors",
			Errors: []config.FieldError{
				{
					Field: "CustomerId",
					Error: "CustomerId is a required field",
				},
			},
		}, err.Message)
	}
	assert.Equal(t, 0, cartServiceMock.callCount)
}

func Test_GivenACreateCartRequestWithCustomerIdShorterThanUUIDLength_WhenCreateNewCart_ThenReturn400(t *testing.T) {
	cartServiceMock := &cartServiceMock{}
	controller, _ := controllers.NewCartController(cartServiceMock)

	e := echo.New()
	e.Validator = config.NewRequestValidator()
	request := httptest.NewRequest(http.MethodPost, "/carts", strings.NewReader(`{"customer_id":"asdasd-asd"}`))
	request.Header.Add(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(request, rec)

	err := controller.CreateNewCart(c)
	if assert.Error(t, err) {
		err := err.(*echo.HTTPError)
		assert.Equal(t, http.StatusBadRequest, err.Code)
		assert.Equal(t, "invalid UUID length: 10", err.Message)
	}
	assert.Equal(t, 0, cartServiceMock.callCount)
}

func Test_GivenACreateCartRequestWithCustomerIdLongerThanUUIDLength_WhenCreateNewCart_ThenReturn400(t *testing.T) {
	cartServiceMock := &cartServiceMock{}
	controller, _ := controllers.NewCartController(cartServiceMock)

	e := echo.New()
	e.Validator = config.NewRequestValidator()
	request := httptest.NewRequest(http.MethodPost, "/carts", strings.NewReader(`{"customer_id":"asdasd-asdasdasd-asdnasdkasd-asda"}`))
	request.Header.Add(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(request, rec)

	err := controller.CreateNewCart(c)
	if assert.Error(t, err) {
		err := err.(*echo.HTTPError)
		assert.Equal(t, http.StatusBadRequest, err.Code)
		assert.Equal(t, "invalid UUID length: 33", err.Message)
	}
	assert.Equal(t, 0, cartServiceMock.callCount)
}

func Test_GivenACreateCartRequestWithMalformedUUID_WhenCreateNewCart_ThenReturn400(t *testing.T) {
	cartServiceMock := &cartServiceMock{}
	controller, _ := controllers.NewCartController(cartServiceMock)

	e := echo.New()
	e.Validator = config.NewRequestValidator()
	request := httptest.NewRequest(http.MethodPost, "/carts", strings.NewReader(`{"customer_id":"12313W12331231231233123123123322"}`))
	request.Header.Add(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(request, rec)

	err := controller.CreateNewCart(c)
	if assert.Error(t, err) {
		err := err.(*echo.HTTPError)
		assert.Equal(t, http.StatusBadRequest, err.Code)
		assert.Equal(t, "invalid UUID format", err.Message)
	}
	assert.Equal(t, 0, cartServiceMock.callCount)
}

func Test_GivenAValidAddItemToCartRequestAndAnExistingEmptyCart_WhenAddItemToCart_ThenReturn200AndCartDtoWithNewItem(t *testing.T) {
	cartId := uuid.New()
	customerId := uuid.New()
	productId := uuid.New()
	cartServiceMock := &cartServiceMock{
		addItemToCart: func(command application.AddItemToCartCommand) (application.CartDto, error) {
			if command.CartId == cartId {
				return application.CartDto{
					Id:         cartId,
					CustomerId: customerId,
					Items: []application.ItemDto{
						{
							ProductId: command.ProductId,
							UnitPrice: 10.10,
							Quantity:  command.Quantity,
						},
					},
				}, nil
			}

			return application.CartDto{}, errors.New("cart doesnt exist")
		},
	}
	controller, _ := controllers.NewCartController(cartServiceMock)

	e := echo.New()
	e.Validator = config.NewRequestValidator()
	request := httptest.NewRequest(http.MethodPost, "/carts", strings.NewReader(
		fmt.Sprintf(`{"product_id":"%s","quantity":2}`, productId.String())))
	request.Header.Add(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(request, rec)
	c.SetPath("/carts/:cartId")
	c.SetParamNames("cartId")
	c.SetParamValues(cartId.String())

	if assert.NoError(t, controller.AddItemToCart(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t,
			fmt.Sprintf("{\"id\":\"%s\",\"customer_id\":\"%s\",\"items\":[{\"product_id\":\"%s\",\"unit_price\":10.10,\"quantity\":2}]}\n",
				cartId.String(), customerId.String(), productId.String()),
			rec.Body.String())
	}
	assert.Equal(t, 1, cartServiceMock.callCount)
}

func Test_GivenANonExistantCart_WhenAddItemToCart_ThenReturn404(t *testing.T) {
	cartId := uuid.New()
	customerId := uuid.New()
	productId := uuid.New()
	cartServiceMock := &cartServiceMock{
		addItemToCart: func(command application.AddItemToCartCommand) (application.CartDto, error) {
			if command.CartId == cartId {
				return application.CartDto{
					Id:         cartId,
					CustomerId: customerId,
					Items: []application.ItemDto{
						{
							ProductId: command.ProductId,
							UnitPrice: 10.10,
							Quantity:  command.Quantity,
						},
					},
				}, nil
			}

			return application.CartDto{}, application.NewNotFoundError(cartId.String(), "cart")
		},
	}
	controller, _ := controllers.NewCartController(cartServiceMock)

	e := echo.New()
	e.Validator = config.NewRequestValidator()
	request := httptest.NewRequest(http.MethodPost, "/carts", strings.NewReader(
		fmt.Sprintf(`{"product_id":"%s","quantity":2}`, productId.String())))
	request.Header.Add(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(request, rec)
	c.SetPath("/carts/:cartId")
	c.SetParamNames("cartId")
	c.SetParamValues(uuid.New().String())

	err := controller.AddItemToCart(c)
	if assert.Error(t, err) {
		err := err.(*echo.HTTPError)
		assert.Equal(t, http.StatusNotFound, err.Code)
		assert.Equal(t, fmt.Sprintf("cart with id %s not found", cartId.String()), err.Message)
	}
	assert.Equal(t, 1, cartServiceMock.callCount)
}

func Test_GivenACartServiceFailsToAddItemToCart_WhenAddItemToCart_ThenReturn500(t *testing.T) {
	cartId := uuid.New()
	productId := uuid.New()
	cartServiceMock := &cartServiceMock{
		addItemToCart: func(command application.AddItemToCartCommand) (application.CartDto, error) {
			return application.CartDto{}, errors.New("failed to add item to cart")
		},
	}
	controller, _ := controllers.NewCartController(cartServiceMock)

	e := echo.New()
	e.Validator = config.NewRequestValidator()
	request := httptest.NewRequest(http.MethodPost, "/carts", strings.NewReader(
		fmt.Sprintf(`{"product_id":"%s","quantity":2}`, productId.String())))
	request.Header.Add(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(request, rec)
	c.SetPath("/carts/:cartId")
	c.SetParamNames("cartId")
	c.SetParamValues(cartId.String())

	err := controller.AddItemToCart(c)
	if assert.Error(t, err) {
		err := err.(*echo.HTTPError)
		assert.Equal(t, http.StatusInternalServerError, err.Code)
		assert.Equal(t, "failed to add item to cart", err.Message)
	}
	assert.Equal(t, 1, cartServiceMock.callCount)
}

func Test_GivenANilCartId_WhenAddItemToCart_ThenReturn400(t *testing.T) {
	productId := uuid.New()
	cartServiceMock := &cartServiceMock{}
	controller, _ := controllers.NewCartController(cartServiceMock)

	e := echo.New()
	e.Validator = config.NewRequestValidator()
	request := httptest.NewRequest(http.MethodPost, "/carts", strings.NewReader(
		fmt.Sprintf(`{"product_id":"%s","quantity":2}`, productId.String())))
	request.Header.Add(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(request, rec)
	c.SetPath("/carts/:cartId")
	c.SetParamNames("cartId")

	err := controller.AddItemToCart(c)
	if assert.Error(t, err) {
		err := err.(*echo.HTTPError)
		assert.Equal(t, http.StatusBadRequest, err.Code)
		assert.Equal(t, &config.ValidationErrorsResponse{
			Message: "there were validation errors",
			Errors: []config.FieldError{
				{
					Field: "CartId",
					Error: "CartId is a required field",
				},
			},
		}, err.Message)
	}
	assert.Equal(t, 0, cartServiceMock.callCount)
}

func Test_GivenAnInvalidGuid_WhenAddItemToCart_ThenReturn400(t *testing.T) {
	productId := uuid.New()
	cartServiceMock := &cartServiceMock{}
	controller, _ := controllers.NewCartController(cartServiceMock)

	e := echo.New()
	e.Validator = config.NewRequestValidator()
	request := httptest.NewRequest(http.MethodPost, "/carts", strings.NewReader(
		fmt.Sprintf(`{"product_id":"%s","quantity":2}`, productId.String())))
	request.Header.Add(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(request, rec)
	c.SetPath("/carts/:cartId")
	c.SetParamNames("cartId")
	c.SetParamValues("1234567890")

	err := controller.AddItemToCart(c)
	if assert.Error(t, err) {
		err := err.(*echo.HTTPError)
		assert.Equal(t, http.StatusBadRequest, err.Code)
		assert.Equal(t, &config.ValidationErrorsResponse{
			Message: "there were validation errors",
			Errors: []config.FieldError{
				{
					Field: "CartId",
					Error: "CartId is a required field",
				},
			},
		}, err.Message)
	}
	assert.Equal(t, 0, cartServiceMock.callCount)
}

func Test_GivenAnInvalidGuidFormat_WhenAddItemToCart_ThenReturn400(t *testing.T) {
	productId := uuid.New()
	cartServiceMock := &cartServiceMock{}
	controller, _ := controllers.NewCartController(cartServiceMock)

	e := echo.New()
	e.Validator = config.NewRequestValidator()
	request := httptest.NewRequest(http.MethodPost, "/carts", strings.NewReader(
		fmt.Sprintf(`{"product_id":"%s","quantity":2}`, productId.String())))
	request.Header.Add(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(request, rec)
	c.SetPath("/carts/:cartId")
	c.SetParamNames("cartId")
	c.SetParamValues("12W456789012W456789012W456789022")

	err := controller.AddItemToCart(c)
	if assert.Error(t, err) {
		err := err.(*echo.HTTPError)
		assert.Equal(t, http.StatusBadRequest, err.Code)
		assert.Equal(t, &config.ValidationErrorsResponse{
			Message: "there were validation errors",
			Errors: []config.FieldError{
				{
					Field: "CartId",
					Error: "CartId is a required field",
				},
			},
		}, err.Message)
	}
	assert.Equal(t, 0, cartServiceMock.callCount)
}

func Test_GivenAnNilProductId_WhenAddItemToCart_ThenReturn400(t *testing.T) {
	cartServiceMock := &cartServiceMock{}
	controller, _ := controllers.NewCartController(cartServiceMock)

	e := echo.New()
	e.Validator = config.NewRequestValidator()
	request := httptest.NewRequest(http.MethodPost, "/carts", strings.NewReader(`{"quantity":2}`))
	request.Header.Add(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(request, rec)
	c.SetPath("/carts/:cartId")
	c.SetParamNames("cartId")
	c.SetParamValues(uuid.New().String())

	err := controller.AddItemToCart(c)
	if assert.Error(t, err) {
		err := err.(*echo.HTTPError)
		assert.Equal(t, http.StatusBadRequest, err.Code)
		assert.Equal(t, &config.ValidationErrorsResponse{
			Message: "there were validation errors",
			Errors: []config.FieldError{
				{
					Field: "ProductId",
					Error: "ProductId is a required field",
				},
			},
		}, err.Message)
	}
	assert.Equal(t, 0, cartServiceMock.callCount)
}

func Test_GivenAnInvalidFormatUUIDForProductId_WhenAddItemToCart_ThenReturn400(t *testing.T) {
	cartServiceMock := &cartServiceMock{}
	controller, _ := controllers.NewCartController(cartServiceMock)

	e := echo.New()
	e.Validator = config.NewRequestValidator()
	request := httptest.NewRequest(http.MethodPost, "/carts", strings.NewReader(`{"product_id":123,"quantity":2}`))
	request.Header.Add(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(request, rec)
	c.SetPath("/carts/:cartId")
	c.SetParamNames("cartId")
	c.SetParamValues(uuid.New().String())

	err := controller.AddItemToCart(c)
	if assert.Error(t, err) {
		err := err.(*echo.HTTPError)
		assert.Equal(t, http.StatusBadRequest, err.Code)
		assert.Equal(t, "Unmarshal type error: expected=uuid.UUID, got=number, field=product_id, offset=17", err.Message)
	}
	assert.Equal(t, 0, cartServiceMock.callCount)
}

func Test_GivenAnShortUUIDForProductId_WhenAddItemToCart_ThenReturn400(t *testing.T) {
	cartServiceMock := &cartServiceMock{}
	controller, _ := controllers.NewCartController(cartServiceMock)

	e := echo.New()
	e.Validator = config.NewRequestValidator()
	request := httptest.NewRequest(http.MethodPost, "/carts", strings.NewReader(`{"product_id":"123123","quantity":2}`))
	request.Header.Add(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(request, rec)
	c.SetPath("/carts/:cartId")
	c.SetParamNames("cartId")
	c.SetParamValues(uuid.New().String())

	err := controller.AddItemToCart(c)
	if assert.Error(t, err) {
		err := err.(*echo.HTTPError)
		assert.Equal(t, http.StatusBadRequest, err.Code)
		assert.Equal(t, "invalid UUID length: 6", err.Message)
	}
	assert.Equal(t, 0, cartServiceMock.callCount)
}

func Test_GivenAnLongUUIDForProductId_WhenAddItemToCart_ThenReturn400(t *testing.T) {
	cartServiceMock := &cartServiceMock{}
	controller, _ := controllers.NewCartController(cartServiceMock)

	e := echo.New()
	e.Validator = config.NewRequestValidator()
	request := httptest.NewRequest(http.MethodPost, "/carts", strings.NewReader(`{"product_id":"123123123123123123123123123123323","quantity":2}`))
	request.Header.Add(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(request, rec)
	c.SetPath("/carts/:cartId")
	c.SetParamNames("cartId")
	c.SetParamValues(uuid.New().String())

	err := controller.AddItemToCart(c)
	if assert.Error(t, err) {
		err := err.(*echo.HTTPError)
		assert.Equal(t, http.StatusBadRequest, err.Code)
		assert.Equal(t, "invalid UUID length: 33", err.Message)
	}
	assert.Equal(t, 0, cartServiceMock.callCount)
}

func Test_GivenANilQuantity_WhenAddItemToCart_ThenReturn400(t *testing.T) {
	cartServiceMock := &cartServiceMock{}
	controller, _ := controllers.NewCartController(cartServiceMock)

	e := echo.New()
	e.Validator = config.NewRequestValidator()
	request := httptest.NewRequest(http.MethodPost, "/carts", strings.NewReader(fmt.Sprintf(`{"product_id":"%s"}`, uuid.New().String())))
	request.Header.Add(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(request, rec)
	c.SetPath("/carts/:cartId")
	c.SetParamNames("cartId")
	c.SetParamValues(uuid.New().String())

	err := controller.AddItemToCart(c)
	if assert.Error(t, err) {
		err := err.(*echo.HTTPError)
		assert.Equal(t, http.StatusBadRequest, err.Code)
		assert.Equal(t, &config.ValidationErrorsResponse{
			Message: "there were validation errors",
			Errors: []config.FieldError{
				{
					Field: "Quantity",
					Error: "Quantity is a required field",
				},
			},
		}, err.Message)
	}
	assert.Equal(t, 0, cartServiceMock.callCount)
}

func Test_GivenAZeroQuantity_WhenAddItemToCart_ThenReturn400(t *testing.T) {
	cartServiceMock := &cartServiceMock{}
	controller, _ := controllers.NewCartController(cartServiceMock)

	e := echo.New()
	e.Validator = config.NewRequestValidator()
	request := httptest.NewRequest(http.MethodPost, "/carts", strings.NewReader(fmt.Sprintf(`{"product_id":"%s","quantity":0}`, uuid.New().String())))
	request.Header.Add(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(request, rec)
	c.SetPath("/carts/:cartId")
	c.SetParamNames("cartId")
	c.SetParamValues(uuid.New().String())

	err := controller.AddItemToCart(c)
	if assert.Error(t, err) {
		err := err.(*echo.HTTPError)
		assert.Equal(t, http.StatusBadRequest, err.Code)
		assert.Equal(t, &config.ValidationErrorsResponse{
			Message: "there were validation errors",
			Errors: []config.FieldError{
				{
					Field: "Quantity",
					Error: "Quantity is a required field",
				},
			},
		}, err.Message)
	}
	assert.Equal(t, 0, cartServiceMock.callCount)
}

func Test_GivenANegativeQuantity_WhenAddItemToCart_ThenReturn400(t *testing.T) {
	cartServiceMock := &cartServiceMock{}
	controller, _ := controllers.NewCartController(cartServiceMock)

	e := echo.New()
	e.Validator = config.NewRequestValidator()
	request := httptest.NewRequest(http.MethodPost, "/carts", strings.NewReader(fmt.Sprintf(`{"product_id":"%s","quantity":-1}`, uuid.New().String())))
	request.Header.Add(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(request, rec)
	c.SetPath("/carts/:cartId")
	c.SetParamNames("cartId")
	c.SetParamValues(uuid.New().String())

	err := controller.AddItemToCart(c)
	if assert.Error(t, err) {
		err := err.(*echo.HTTPError)
		assert.Equal(t, http.StatusBadRequest, err.Code)
		assert.Equal(t, &config.ValidationErrorsResponse{
			Message: "there were validation errors",
			Errors: []config.FieldError{
				{
					Field: "Quantity",
					Error: "Quantity must be greater than 0",
				},
			},
		}, err.Message)
	}
	assert.Equal(t, 0, cartServiceMock.callCount)
}

type cartServiceMock struct {
	callCount     int
	createNewCart func(application.CreateCartCommand) (application.CartDto, error)
	addItemToCart func(application.AddItemToCartCommand) (application.CartDto, error)
}

func (c *cartServiceMock) CreateNewCart(command application.CreateCartCommand) (application.CartDto, error) {
	c.callCount++
	return c.createNewCart(command)
}

func (c *cartServiceMock) AddItemToCart(command application.AddItemToCartCommand) (application.CartDto, error) {
	c.callCount++
	return c.addItemToCart(command)
}
