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

type QuestionHandler struct{ service port.QuestionService }

func NewQuestionHandler(service port.QuestionService) *QuestionHandler {
	return &QuestionHandler{service: service}
}

func (h *QuestionHandler) Create(c echo.Context) error {
	var req CreateQuestionRequest
	if err := helper.BindAndValidate(c, &req); err != nil {
		return helper.RespondError(c, err)
	}
	questionSetID, err := parseUUID(req.QuestionSetID, "question_set_id")
	if err != nil {
		return helper.RespondError(c, err)
	}
	choices, err := toChoiceInputs(req.Choices)
	if err != nil {
		return helper.RespondError(c, err)
	}
	tagIDs, err := parseUUIDList(req.TagIDs, "tag_id")
	if err != nil {
		return helper.RespondError(c, err)
	}
	item, err := h.service.Create(c.Request().Context(), domain.QuestionCreateInput{QuestionSetID: questionSetID, Type: req.Type, QuestionText: req.QuestionText, BlankPosition: req.BlankPosition, Explanation: req.Explanation, ImageURL: req.ImageURL, Difficulty: req.Difficulty, OrderNo: req.OrderNo, IsActive: req.IsActive, Choices: choices, TagIDs: tagIDs, ActorID: helper.ActorIDFromContext(c)})
	if err != nil {
		return helper.RespondError(c, err)
	}
	return helper.RespondSuccess(c, http.StatusCreated, toQuestionResponse(item))
}

func (h *QuestionHandler) Update(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return helper.RespondError(c, helper.BadRequest("invalid_question_id", "question id must be a valid uuid", err))
	}
	var req UpdateQuestionRequest
	if err := helper.BindAndValidate(c, &req); err != nil {
		return helper.RespondError(c, err)
	}
	input := domain.QuestionUpdateInput{ID: id, ActorID: helper.ActorIDFromContext(c)}
	if req.QuestionSetID != nil {
		questionSetID, err := parseUUID(*req.QuestionSetID, "question_set_id")
		if err != nil {
			return helper.RespondError(c, err)
		}
		input.QuestionSetID = domain.NewEntityField(questionSetID)
	}
	if req.Type != nil {
		input.Type = domain.NewEntityField(strings.TrimSpace(*req.Type))
	}
	if req.QuestionText != nil {
		input.QuestionText = domain.NewEntityField(strings.TrimSpace(*req.QuestionText))
	}
	if req.BlankPosition != nil {
		input.BlankPosition = domain.NewEntityField(req.BlankPosition)
	}
	if req.Explanation != nil {
		input.Explanation = domain.NewEntityField(req.Explanation)
	}
	if req.ImageURL != nil {
		input.ImageURL = domain.NewEntityField(req.ImageURL)
	}
	if req.Difficulty != nil {
		input.Difficulty = domain.NewEntityField(*req.Difficulty)
	}
	if req.OrderNo != nil {
		input.OrderNo = domain.NewEntityField(*req.OrderNo)
	}
	if req.IsActive != nil {
		input.IsActive = domain.NewEntityField(*req.IsActive)
	}
	if req.Choices != nil {
		choices, err := toChoiceInputs(req.Choices)
		if err != nil {
			return helper.RespondError(c, err)
		}
		input.Choices = domain.NewEntityField(choices)
	}
	if req.TagIDs != nil {
		tagIDs, err := parseUUIDList(req.TagIDs, "tag_id")
		if err != nil {
			return helper.RespondError(c, err)
		}
		input.TagIDs = domain.NewEntityField(tagIDs)
	}
	item, err := h.service.Update(c.Request().Context(), input)
	if err != nil {
		return helper.RespondError(c, err)
	}
	return helper.RespondSuccess(c, http.StatusOK, toQuestionResponse(item))
}

func (h *QuestionHandler) Delete(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return helper.RespondError(c, helper.BadRequest("invalid_question_id", "question id must be a valid uuid", err))
	}
	if err := h.service.Delete(c.Request().Context(), id, helper.ActorIDFromContext(c)); err != nil {
		return helper.RespondError(c, err)
	}
	return helper.RespondSuccess(c, http.StatusOK, map[string]string{"id": id.String(), "status": "deleted"})
}
func (h *QuestionHandler) FindByID(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return helper.RespondError(c, helper.BadRequest("invalid_question_id", "question id must be a valid uuid", err))
	}
	item, err := h.service.FindByID(c.Request().Context(), id)
	if err != nil {
		return helper.RespondError(c, err)
	}
	return helper.RespondSuccess(c, http.StatusOK, toQuestionResponse(item))
}
func (h *QuestionHandler) FindAll(c echo.Context) error {
	query := domain.QuestionQuery{Paging: helper.BuildPaging(c), Search: strings.TrimSpace(c.QueryParam("search")), IncludeChoices: !strings.EqualFold(strings.TrimSpace(c.QueryParam("include_choices")), "false"), IncludeTags: !strings.EqualFold(strings.TrimSpace(c.QueryParam("include_tags")), "false")}
	query.SortBy, query.SortOrder = helper.BuildSort(c)
	if value := strings.TrimSpace(c.QueryParam("question_set_id")); value != "" {
		parsed, err := parseUUID(value, "question_set_id")
		if err != nil {
			return helper.RespondError(c, err)
		}
		query.QuestionSetID = &parsed
	}
	if value := strings.TrimSpace(c.QueryParam("type")); value != "" {
		query.Type = &value
	}
	if value := strings.TrimSpace(c.QueryParam("is_active")); value != "" {
		parsed := strings.EqualFold(value, "true")
		query.IsActive = &parsed
	}
	result, err := h.service.FindAll(c.Request().Context(), query)
	if err != nil {
		return helper.RespondError(c, err)
	}
	items := make([]QuestionResponse, 0, len(result.Items))
	for _, item := range result.Items {
		items = append(items, toQuestionResponse(item))
	}
	return helper.RespondSuccess(c, http.StatusOK, ListResponse[QuestionResponse, domain.QuestionQuery]{Items: items, Paging: PagingResponse{Page: result.Paging.Page, Limit: result.Paging.Limit, Total: result.Paging.Total}, Query: query})
}
