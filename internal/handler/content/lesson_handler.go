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

type LessonHandler struct{ service port.LessonService }

func NewLessonHandler(service port.LessonService) *LessonHandler {
	return &LessonHandler{service: service}
}

func (h *LessonHandler) Create(c echo.Context) error {
	var req CreateLessonRequest
	if err := helper.BindAndValidate(c, &req); err != nil {
		return helper.RespondError(c, err)
	}
	sectionID, err := parseUUID(req.SectionID, "section_id")
	if err != nil {
		return helper.RespondError(c, err)
	}
	item, err := h.service.Create(c.Request().Context(), domain.LessonCreateInput{SectionID: sectionID, Slug: req.Slug, Title: req.Title, Description: req.Description, OrderNo: req.OrderNo, IsPublished: req.IsPublished, ActorID: helper.ActorIDFromContext(c)})
	if err != nil {
		return helper.RespondError(c, err)
	}
	return helper.RespondSuccess(c, http.StatusCreated, toLessonResponse(item))
}

func (h *LessonHandler) Update(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return helper.RespondError(c, helper.BadRequest("invalid_lesson_id", "lesson id must be a valid uuid", err))
	}
	var req UpdateLessonRequest
	if err := helper.BindAndValidate(c, &req); err != nil {
		return helper.RespondError(c, err)
	}
	input := domain.LessonUpdateInput{ID: id, ActorID: helper.ActorIDFromContext(c)}
	if req.SectionID != nil {
		sectionID, err := parseUUID(*req.SectionID, "section_id")
		if err != nil {
			return helper.RespondError(c, err)
		}
		input.SectionID = domain.NewEntityField(sectionID)
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
	return helper.RespondSuccess(c, http.StatusOK, toLessonResponse(item))
}

func (h *LessonHandler) Delete(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return helper.RespondError(c, helper.BadRequest("invalid_lesson_id", "lesson id must be a valid uuid", err))
	}
	if err := h.service.Delete(c.Request().Context(), id, helper.ActorIDFromContext(c)); err != nil {
		return helper.RespondError(c, err)
	}
	return helper.RespondSuccess(c, http.StatusOK, map[string]string{"id": id.String(), "status": "deleted"})
}

func (h *LessonHandler) FindByID(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return helper.RespondError(c, helper.BadRequest("invalid_lesson_id", "lesson id must be a valid uuid", err))
	}
	item, err := h.service.FindByID(c.Request().Context(), id)
	if err != nil {
		return helper.RespondError(c, err)
	}
	return helper.RespondSuccess(c, http.StatusOK, toLessonResponse(item))
}

func (h *LessonHandler) FindAll(c echo.Context) error {
	query := domain.LessonQuery{Paging: helper.BuildPaging(c), Search: strings.TrimSpace(c.QueryParam("search"))}
	query.SortBy, query.SortOrder = helper.BuildSort(c)
	if value := strings.TrimSpace(c.QueryParam("section_id")); value != "" {
		parsed, err := parseUUID(value, "section_id")
		if err != nil {
			return helper.RespondError(c, err)
		}
		query.SectionID = &parsed
	}
	if value := strings.TrimSpace(c.QueryParam("is_published")); value != "" {
		parsed := strings.EqualFold(value, "true")
		query.IsPublished = &parsed
	}
	result, err := h.service.FindAll(c.Request().Context(), query)
	if err != nil {
		return helper.RespondError(c, err)
	}
	items := make([]LessonResponse, 0, len(result.Items))
	for _, item := range result.Items {
		items = append(items, toLessonResponse(item))
	}
	return helper.RespondSuccess(c, http.StatusOK, ListResponse[LessonResponse, domain.LessonQuery]{Items: items, Paging: PagingResponse{Page: result.Paging.Page, Limit: result.Paging.Limit, Total: result.Paging.Total}, Query: query})
}
