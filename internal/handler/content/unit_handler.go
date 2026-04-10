package content

import (
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/namchokGithub/vocabunny-core-api/internal/core/domain"
	"github.com/namchokGithub/vocabunny-core-api/internal/core/helper"
	"github.com/namchokGithub/vocabunny-core-api/internal/core/port"
)

type UnitHandler struct{ service port.UnitService }

func NewUnitHandler(service port.UnitService) *UnitHandler { return &UnitHandler{service: service} }
func (h *UnitHandler) Create(c echo.Context) error {
	var req CreateUnitRequest
	if err := helper.BindAndValidate(c, &req); err != nil {
		return helper.RespondError(c, err)
	}
	lessonID, err := parseUUID(req.LessonID, "lesson_id")
	if err != nil {
		return helper.RespondError(c, err)
	}
	item, err := h.service.Create(c.Request().Context(), domain.UnitCreateInput{LessonID: lessonID, Slug: req.Slug, Title: req.Title, Description: req.Description, OrderNo: req.OrderNo, IsPublished: req.IsPublished, ActorID: helper.ActorIDFromContext(c)})
	if err != nil {
		return helper.RespondError(c, err)
	}
	return helper.RespondSuccess(c, http.StatusCreated, toUnitResponse(item))
}
func (h *UnitHandler) Update(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return helper.RespondError(c, helper.BadRequest("invalid_unit_id", "unit id must be a valid uuid", err))
	}
	var req UpdateUnitRequest
	if err := helper.BindAndValidate(c, &req); err != nil {
		return helper.RespondError(c, err)
	}
	input := domain.UnitUpdateInput{ID: id, ActorID: helper.ActorIDFromContext(c)}
	if req.LessonID != nil {
		lessonID, err := parseUUID(*req.LessonID, "lesson_id")
		if err != nil {
			return helper.RespondError(c, err)
		}
		input.LessonID = domain.NewEntityField(lessonID)
	}
	if req.Slug != nil {
		input.Slug = domain.NewEntityField(strings.TrimSpace(*req.Slug))
	}
	if req.Title != nil {
		input.Title = domain.NewEntityField(strings.TrimSpace(*req.Title))
	}
	if req.Description != nil {
		input.Description = domain.NewEntityField(req.Description)
	}
	if req.OrderNo != nil {
		input.OrderNo = domain.NewEntityField(*req.OrderNo)
	}
	if req.IsPublished != nil {
		input.IsPublished = domain.NewEntityField(*req.IsPublished)
	}
	item, err := h.service.Update(c.Request().Context(), input)
	if err != nil {
		return helper.RespondError(c, err)
	}
	return helper.RespondSuccess(c, http.StatusOK, toUnitResponse(item))
}
func (h *UnitHandler) Delete(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return helper.RespondError(c, helper.BadRequest("invalid_unit_id", "unit id must be a valid uuid", err))
	}
	if err := h.service.Delete(c.Request().Context(), id, helper.ActorIDFromContext(c)); err != nil {
		return helper.RespondError(c, err)
	}
	return helper.RespondSuccess(c, http.StatusOK, map[string]string{"id": id.String(), "status": "deleted"})
}
func (h *UnitHandler) FindByID(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return helper.RespondError(c, helper.BadRequest("invalid_unit_id", "unit id must be a valid uuid", err))
	}
	item, err := h.service.FindByID(c.Request().Context(), id)
	if err != nil {
		return helper.RespondError(c, err)
	}
	return helper.RespondSuccess(c, http.StatusOK, toUnitResponse(item))
}
func (h *UnitHandler) FindAll(c echo.Context) error {
	query := domain.UnitQuery{Paging: helper.BuildPaging(c), Search: strings.TrimSpace(c.QueryParam("search"))}
	query.SortBy, query.SortOrder = helper.BuildSort(c)
	if value := strings.TrimSpace(c.QueryParam("lesson_id")); value != "" {
		parsed, err := parseUUID(value, "lesson_id")
		if err != nil {
			return helper.RespondError(c, err)
		}
		query.LessonID = &parsed
	}
	if value := strings.TrimSpace(c.QueryParam("is_published")); value != "" {
		parsed := strings.EqualFold(value, "true")
		query.IsPublished = &parsed
	}
	result, err := h.service.FindAll(c.Request().Context(), query)
	if err != nil {
		return helper.RespondError(c, err)
	}
	items := make([]UnitResponse, 0, len(result.Items))
	for _, item := range result.Items {
		items = append(items, toUnitResponse(item))
	}
	return helper.RespondSuccess(c, http.StatusOK, ListResponse[UnitResponse, domain.UnitQuery]{Items: items, Paging: PagingResponse{Page: result.Paging.Page, Limit: result.Paging.Limit, Total: result.Paging.Total}, Query: query})
}
