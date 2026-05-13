package content

import (
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/namchokGithub/vocabunny-core-api/internal/constants"
	"github.com/namchokGithub/vocabunny-core-api/internal/core/domain"
	"github.com/namchokGithub/vocabunny-core-api/internal/core/helper"
	"github.com/namchokGithub/vocabunny-core-api/internal/core/port"
)

type TagHandler struct{ service port.TagService }

func NewTagHandler(service port.TagService) *TagHandler { return &TagHandler{service: service} }
func (h *TagHandler) Create(c echo.Context) error {
	var req CreateTagRequest
	if err := helper.BindAndValidate(c, &req); err != nil {
		return helper.RespondError(c, err)
	}
	item, err := h.service.Create(c.Request().Context(), domain.TagCreateInput{
		Name:    req.Name,
		Color:   strings.TrimSpace(req.Color),
		ActorID: helper.ActorIDFromContext(c),
	})
	if err != nil {
		return helper.RespondError(c, err)
	}
	return helper.RespondSuccess(c, http.StatusCreated, toTagResponse(item), constants.CodeCreated)
}
func (h *TagHandler) Update(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return helper.RespondError(c, helper.BadRequest(constants.CodeInvalidQueryParam, "tag id must be a valid uuid", err))
	}
	var req UpdateTagRequest
	if err := helper.BindAndValidate(c, &req); err != nil {
		return helper.RespondError(c, err)
	}
	input := domain.TagUpdateInput{ID: id, ActorID: helper.ActorIDFromContext(c)}
	if req.Name != nil {
		input.Name = domain.NewEntityField(strings.TrimSpace(*req.Name))
	}
	if req.Color != nil {
		input.Color = domain.NewEntityField(strings.TrimSpace(*req.Color))
	}
	item, err := h.service.Update(c.Request().Context(), input)
	if err != nil {
		return helper.RespondError(c, err)
	}
	return helper.RespondSuccess(c, http.StatusOK, toTagResponse(item), constants.CodeUpdated)
}
func (h *TagHandler) Delete(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return helper.RespondError(c, helper.BadRequest(constants.CodeInvalidQueryParam, "tag id must be a valid uuid", err))
	}
	if err := h.service.Delete(c.Request().Context(), id, helper.ActorIDFromContext(c)); err != nil {
		return helper.RespondError(c, err)
	}
	return helper.RespondSuccess(c, http.StatusOK, map[string]string{"id": id.String(), "status": "deleted"}, constants.CodeDeleted)
}
func (h *TagHandler) FindByID(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return helper.RespondError(c, helper.BadRequest(constants.CodeInvalidQueryParam, "tag id must be a valid uuid", err))
	}
	item, err := h.service.FindByID(c.Request().Context(), id)
	if err != nil {
		return helper.RespondError(c, err)
	}
	return helper.RespondSuccess(c, http.StatusOK, toTagResponse(item), constants.CodeSuccess)
}
func (h *TagHandler) FindAll(c echo.Context) error {
	query := domain.TagQuery{Paging: helper.BuildPaging(c), Search: strings.TrimSpace(c.QueryParam("search"))}
	query.SortBy, query.SortOrder = helper.BuildSort(c)
	result, err := h.service.FindAll(c.Request().Context(), query)
	if err != nil {
		return helper.RespondError(c, err)
	}
	items := make([]TagResponse, 0, len(result.Items))
	for _, item := range result.Items {
		items = append(items, toTagResponse(item))
	}
	return helper.RespondSuccess(c, http.StatusOK, ListResponse[TagResponse, domain.TagQuery]{Items: items, Paging: helper.NewPagingResponse(result.Paging), Query: query}, constants.CodeSuccess)
}
