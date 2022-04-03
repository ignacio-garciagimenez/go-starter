package controllers

import (
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

type validationErrorsResponse struct {
	Message string       `json:"message"`
	Errors  []fieldError `json:"validation_errors"`
}

type fieldError struct {
	Field string `json:"field"`
	Error string `json:"error"`
}

type requestValidator struct {
	validator *validator.Validate
}

func NewRequestValidator() *requestValidator {
	return &requestValidator{validator.New()}
}

func (rv *requestValidator) Validate(i interface{}) error {
	if err := rv.validator.Struct(i); err != nil {
		//resp := createValidationErrorResponse(err.(validator.ValidationErrors))

		return echo.NewHTTPError(http.StatusBadRequest, err)
	}
	return nil
}

func createValidationErrorResponse(validationErrors validator.ValidationErrors) *validationErrorsResponse {
	errorResponse := &validationErrorsResponse{Message: "there were validation errors"}

	for _, fieldErr := range validationErrors {
		errorResponse.Errors = append(errorResponse.Errors, fieldError{
			Field: fieldErr.Field(),
			Error: fieldErr.Error(),
		})
	}

	return errorResponse
}

func CustomValidationErrorHanlder(errorHandler echo.HTTPErrorHandler) echo.HTTPErrorHandler {
	return func(err error, c echo.Context) {
		if httpError, ok := err.(*echo.HTTPError); ok {
			if validationError, ok := httpError.Message.(*validationErrorsResponse); ok {
				c.JSON(httpError.Code, validationError)
				return
			}
		}
		errorHandler(err, c)
	}
}
