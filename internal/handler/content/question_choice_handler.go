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

type QuestionChoiceHandler struct{ service port.QuestionChoiceService }

func NewQuestionChoiceHandler(service port.QuestionChoiceService) *QuestionChoiceHandler {
	return &QuestionChoiceHandler{service: service}
}
func (h *QuestionChoiceHandler) Create(c echo.Context) error {
	var req CreateQuestionChoiceRequest
	if err := helper.BindAndValidate(c, &req); err != nil {
		return helper.RespondError(c, err)
	}
	questionID, err := parseUUID(req.QuestionID, "question_id")
	if err != nil {
		return helper.RespondError(c, err)
	}
	item, err := h.service.Create(c.Request().Context(), domain.QuestionChoiceCreateInput{QuestionID: questionID, ChoiceText: req.ChoiceText, ChoiceOrder: req.ChoiceOrder, IsCorrect: req.IsCorrect, ActorID: helper.ActorIDFromContext(c)})
	if err != nil {
		return helper.RespondError(c, err)
	}
	return helper.RespondSuccess(c, http.StatusCreated, toQuestionChoiceResponse(item))
}
func (h *QuestionChoiceHandler) Update(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return helper.RespondError(c, helper.BadRequest("invalid_question_choice_id", "question choice id must be a valid uuid", err))
	}
	var req UpdateQuestionChoiceRequest
	if err := helper.BindAndValidate(c, &req); err != nil {
		return helper.RespondError(c, err)
	}
	input := domain.QuestionChoiceUpdateInput{ID: id, ActorID: helper.ActorIDFromContext(c)}
	if req.ChoiceText != nil {
		input.ChoiceText = domain.NewEntityField(strings.TrimSpace(*req.ChoiceText))
	}
	if req.ChoiceOrder != nil {
		input.ChoiceOrder = domain.NewEntityField(*req.ChoiceOrder)
	}
	if req.IsCorrect != nil {
		input.IsCorrect = domain.NewEntityField(*req.IsCorrect)
	}
	item, err := h.service.Update(c.Request().Context(), input)
	if err != nil {
		return helper.RespondError(c, err)
	}
	return helper.RespondSuccess(c, http.StatusOK, toQuestionChoiceResponse(item))
}
func (h *QuestionChoiceHandler) Delete(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return helper.RespondError(c, helper.BadRequest("invalid_question_choice_id", "question choice id must be a valid uuid", err))
	}
	if err := h.service.Delete(c.Request().Context(), id, helper.ActorIDFromContext(c)); err != nil {
		return helper.RespondError(c, err)
	}
	return helper.RespondSuccess(c, http.StatusOK, map[string]string{"id": id.String(), "status": "deleted"})
}
func (h *QuestionChoiceHandler) FindByID(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return helper.RespondError(c, helper.BadRequest("invalid_question_choice_id", "question choice id must be a valid uuid", err))
	}
	item, err := h.service.FindByID(c.Request().Context(), id)
	if err != nil {
		return helper.RespondError(c, err)
	}
	return helper.RespondSuccess(c, http.StatusOK, toQuestionChoiceResponse(item))
}
func (h *QuestionChoiceHandler) FindAll(c echo.Context) error {
	query := domain.QuestionChoiceQuery{Paging: helper.BuildPaging(c)}
	query.SortBy, query.SortOrder = helper.BuildSort(c)
	if value := strings.TrimSpace(c.QueryParam("question_id")); value != "" {
		parsed, err := parseUUID(value, "question_id")
		if err != nil {
			return helper.RespondError(c, err)
		}
		query.QuestionID = &parsed
	}
	result, err := h.service.FindAll(c.Request().Context(), query)
	if err != nil {
		return helper.RespondError(c, err)
	}
	items := make([]QuestionChoiceResponse, 0, len(result.Items))
	for _, item := range result.Items {
		items = append(items, toQuestionChoiceResponse(item))
	}
	return helper.RespondSuccess(c, http.StatusOK, ListResponse[QuestionChoiceResponse, domain.QuestionChoiceQuery]{Items: items, Paging: PagingResponse{Page: result.Paging.Page, Limit: result.Paging.Limit, Total: result.Paging.Total}, Query: query})
}
