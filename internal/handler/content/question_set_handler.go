package content

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/namchokGithub/vocabunny-core-api/internal/core/domain"
	"github.com/namchokGithub/vocabunny-core-api/internal/core/helper"
	"github.com/namchokGithub/vocabunny-core-api/internal/core/port"
)

type QuestionSetHandler struct{ service port.QuestionSetService }

func NewQuestionSetHandler(service port.QuestionSetService) *QuestionSetHandler {
	return &QuestionSetHandler{service: service}
}
func (h *QuestionSetHandler) Create(c echo.Context) error {
	var req CreateQuestionSetRequest
	if err := helper.BindAndValidate(c, &req); err != nil {
		return helper.RespondError(c, err)
	}
	unitID, err := parseUUID(req.UnitID, "unit_id")
	if err != nil {
		return helper.RespondError(c, err)
	}
	item, err := h.service.Create(c.Request().Context(), domain.QuestionSetCreateInput{UnitID: unitID, Slug: req.Slug, Title: req.Title, Description: req.Description, OrderNo: req.OrderNo, EstimatedSeconds: req.EstimatedSeconds, IsPublished: req.IsPublished, Version: req.Version, ActorID: helper.ActorIDFromContext(c)})
	if err != nil {
		return helper.RespondError(c, err)
	}
	return helper.RespondSuccess(c, http.StatusCreated, toQuestionSetResponse(item))
}
func (h *QuestionSetHandler) Update(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return helper.RespondError(c, helper.BadRequest("invalid_question_set_id", "question set id must be a valid uuid", err))
	}
	var req UpdateQuestionSetRequest
	if err := helper.BindAndValidate(c, &req); err != nil {
		return helper.RespondError(c, err)
	}
	input := domain.QuestionSetUpdateInput{ID: id, ActorID: helper.ActorIDFromContext(c)}
	if req.UnitID != nil {
		unitID, err := parseUUID(*req.UnitID, "unit_id")
		if err != nil {
			return helper.RespondError(c, err)
		}
		input.UnitID = domain.NewEntityField(unitID)
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
	if req.EstimatedSeconds != nil {
		input.EstimatedSeconds = domain.NewEntityField(req.EstimatedSeconds)
	}
	if req.IsPublished != nil {
		input.IsPublished = domain.NewEntityField(*req.IsPublished)
	}
	if req.Version != nil {
		input.Version = domain.NewEntityField(*req.Version)
	}
	item, err := h.service.Update(c.Request().Context(), input)
	if err != nil {
		return helper.RespondError(c, err)
	}
	return helper.RespondSuccess(c, http.StatusOK, toQuestionSetResponse(item))
}
func (h *QuestionSetHandler) Delete(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return helper.RespondError(c, helper.BadRequest("invalid_question_set_id", "question set id must be a valid uuid", err))
	}
	if err := h.service.Delete(c.Request().Context(), id, helper.ActorIDFromContext(c)); err != nil {
		return helper.RespondError(c, err)
	}
	return helper.RespondSuccess(c, http.StatusOK, map[string]string{"id": id.String(), "status": "deleted"})
}
func (h *QuestionSetHandler) FindByID(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return helper.RespondError(c, helper.BadRequest("invalid_question_set_id", "question set id must be a valid uuid", err))
	}
	item, err := h.service.FindByID(c.Request().Context(), id)
	if err != nil {
		return helper.RespondError(c, err)
	}
	return helper.RespondSuccess(c, http.StatusOK, toQuestionSetResponse(item))
}
func (h *QuestionSetHandler) FindAll(c echo.Context) error {
	query := domain.QuestionSetQuery{Paging: helper.BuildPaging(c), Search: strings.TrimSpace(c.QueryParam("search"))}
	query.SortBy, query.SortOrder = helper.BuildSort(c)
	if value := strings.TrimSpace(c.QueryParam("unit_id")); value != "" {
		parsed, err := parseUUID(value, "unit_id")
		if err != nil {
			return helper.RespondError(c, err)
		}
		query.UnitID = &parsed
	}
	if value := strings.TrimSpace(c.QueryParam("version")); value != "" {
		parsed, err := strconv.Atoi(value)
		if err != nil {
			return helper.RespondError(c, helper.BadRequest("invalid_version", "version must be an integer", err))
		}
		query.Version = &parsed
	}
	if value := strings.TrimSpace(c.QueryParam("is_published")); value != "" {
		parsed := strings.EqualFold(value, "true")
		query.IsPublished = &parsed
	}
	result, err := h.service.FindAll(c.Request().Context(), query)
	if err != nil {
		return helper.RespondError(c, err)
	}
	items := make([]QuestionSetResponse, 0, len(result.Items))
	for _, item := range result.Items {
		items = append(items, toQuestionSetResponse(item))
	}
	return helper.RespondSuccess(c, http.StatusOK, ListResponse[QuestionSetResponse, domain.QuestionSetQuery]{Items: items, Paging: PagingResponse{Page: result.Paging.Page, Limit: result.Paging.Limit, Total: result.Paging.Total}, Query: query})
}
