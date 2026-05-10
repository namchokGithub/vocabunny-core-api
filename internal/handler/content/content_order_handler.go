package content

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/namchokGithub/vocabunny-core-api/internal/core/helper"
	"github.com/namchokGithub/vocabunny-core-api/internal/core/port"
)

type ContentOrderHandler struct{ service port.ContentOrderService }

func NewContentOrderHandler(service port.ContentOrderService) *ContentOrderHandler {
	return &ContentOrderHandler{service: service}
}

func (h *ContentOrderHandler) GetLastOrderNos(c echo.Context) error {
	item, err := h.service.GetLastOrderNos(c.Request().Context())
	if err != nil {
		return helper.RespondError(c, err)
	}
	return helper.RespondSuccess(c, http.StatusOK, toContentOrderNoResponse(item))
}
