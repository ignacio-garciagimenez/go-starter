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

type cartServiceMock struct {
	callCount     int
	createNewCart func(application.CreateCartCommand) (application.CartDto, error)
}

func (c *cartServiceMock) CreateNewCart(command application.CreateCartCommand) (application.CartDto, error) {
	c.callCount++
	return c.createNewCart(command)
}
