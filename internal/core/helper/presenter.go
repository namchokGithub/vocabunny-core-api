package helper

import "github.com/labstack/echo/v4"

type successResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data"`
}

type errorResponse struct {
	Success bool   `json:"success"`
	Error   apiErr `json:"error"`
}

type apiErr struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func RespondSuccess(c echo.Context, statusCode int, data interface{}) error {
	return c.JSON(statusCode, successResponse{
		Success: true,
		Data:    data,
	})
}

func RespondError(c echo.Context, err error) error {
	appErr := AsAppError(err)

	return c.JSON(appErr.StatusCode, errorResponse{
		Success: false,
		Error: apiErr{
			Code:    appErr.Code,
			Message: appErr.Message,
		},
	})
}
