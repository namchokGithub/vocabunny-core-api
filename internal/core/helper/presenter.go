package helper

import (
	"fmt"
	"os"

	"github.com/labstack/echo/v4"
)

type SuccessResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Meta    *Meta       `json:"meta,omitempty"`
}

type ErrorResponse struct {
	Success bool     `json:"success"`
	Error   APIError `json:"error"`
}

type Meta struct {
	Code string `json:"code"`
}

type APIError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// BuildCode formats a business/system code as vocab-{env}-{code}.
// Reads APP_ENV from the environment; falls back to "local".
func BuildCode(code string) string {
	env := os.Getenv("APP_ENV")
	if env == "" {
		env = "local"
	}
	return fmt.Sprintf("vocab-%s-%s", env, code)
}

// RespondSuccess writes a JSON success response.
// Pass an optional business code (e.g. constants.CodeCreated) to include meta.code.
func RespondSuccess(c echo.Context, statusCode int, data interface{}, code ...string) error {
	resp := SuccessResponse{
		Success: true,
		Data:    data,
	}
	if len(code) > 0 && code[0] != "" {
		resp.Meta = &Meta{Code: BuildCode(code[0])}
	}
	return c.JSON(statusCode, resp)
}

func RespondError(c echo.Context, err error) error {
	appErr := AsAppError(err)
	return c.JSON(appErr.StatusCode, ErrorResponse{
		Success: false,
		Error: APIError{
			Code:    BuildCode(appErr.Code),
			Message: appErr.Message,
		},
	})
}
