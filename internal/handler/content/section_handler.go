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

type SectionHandler struct{ service port.SectionService }

func NewSectionHandler(service port.SectionService) *SectionHandler {
	return &SectionHandler{service: service}
}

func (h *SectionHandler) Create(c echo.Context) error {
	var req CreateSectionRequest
	if err := helper.BindAndValidate(c, &req); err != nil {
		return helper.RespondError(c, err)
	}
	item, err := h.service.Create(c.Request().Context(), domain.SectionCreateInput{Slug: req.Slug, Title: req.Title, Description: req.Description, OrderNo: req.OrderNo, IsPublished: req.IsPublished, ActorID: helper.ActorIDFromContext(c)})
	if err != nil {
		return helper.RespondError(c, err)
	}
	return helper.RespondSuccess(c, http.StatusCreated, toSectionResponse(item))
}

func (h *SectionHandler) Update(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return helper.RespondError(c, helper.BadRequest("invalid_section_id", "section id must be a valid uuid", err))
	}
	var req UpdateSectionRequest
	if err := helper.BindAndValidate(c, &req); err != nil {
		return helper.RespondError(c, err)
	}
	input := domain.SectionUpdateInput{ID: id, ActorID: helper.ActorIDFromContext(c)}
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
	return helper.RespondSuccess(c, http.StatusOK, toSectionResponse(item))
}

func (h *SectionHandler) Delete(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return helper.RespondError(c, helper.BadRequest("invalid_section_id", "section id must be a valid uuid", err))
	}
	if err := h.service.Delete(c.Request().Context(), id, helper.ActorIDFromContext(c)); err != nil {
		return helper.RespondError(c, err)
	}
	return helper.RespondSuccess(c, http.StatusOK, map[string]string{"id": id.String(), "status": "deleted"})
}

func (h *SectionHandler) FindByID(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return helper.RespondError(c, helper.BadRequest("invalid_section_id", "section id must be a valid uuid", err))
	}
	item, err := h.service.FindByID(c.Request().Context(), id)
	if err != nil {
		return helper.RespondError(c, err)
	}
	return helper.RespondSuccess(c, http.StatusOK, toSectionResponse(item))
}

func (h *SectionHandler) FindAll(c echo.Context) error {
	query := domain.SectionQuery{Paging: helper.BuildPaging(c), Search: strings.TrimSpace(c.QueryParam("search"))}
	query.SortBy, query.SortOrder = helper.BuildSort(c)
	if value := strings.TrimSpace(c.QueryParam("is_published")); value != "" {
		parsed := strings.EqualFold(value, "true")
		query.IsPublished = &parsed
	}
	result, err := h.service.FindAll(c.Request().Context(), query)
	if err != nil {
		return helper.RespondError(c, err)
	}
	items := make([]SectionResponse, 0, len(result.Items))
	for _, item := range result.Items {
		items = append(items, toSectionResponse(item))
	}
	return helper.RespondSuccess(c, http.StatusOK, ListResponse[SectionResponse, domain.SectionQuery]{Items: items, Paging: PagingResponse{Page: result.Paging.Page, Limit: result.Paging.Limit, Total: result.Paging.Total}, Query: query})
}
