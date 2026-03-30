package helper

import (
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

type RequestValidator struct {
	validate *validator.Validate
}

func NewRequestValidator() *RequestValidator {
	return &RequestValidator{
		validate: validator.New(),
	}
}

func (v *RequestValidator) Validate(i interface{}) error {
	if err := v.validate.Struct(i); err != nil {
		return BadRequest("validation_error", err.Error(), err)
	}

	return nil
}

func BindAndValidate(c echo.Context, payload interface{}) error {
	if err := c.Bind(payload); err != nil {
		return BadRequest("invalid_request", "invalid request payload", err)
	}

	validatorInstance, ok := c.Echo().Validator.(*RequestValidator)
	if !ok {
		return Internal("validator_unavailable", "validator is not configured", nil)
	}

	if err := validatorInstance.Validate(payload); err != nil {
		return err
	}

	return nil
}

func ParseUUIDParam(c echo.Context, key string) (string, error) {
	value := c.Param(key)
	if value == "" {
		return "", echo.NewHTTPError(http.StatusBadRequest, key+" is required")
	}

	return value, nil
}
