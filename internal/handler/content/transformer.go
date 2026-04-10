package content

import (
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/namchokGithub/vocabunny-core-api/internal/core/domain"
	"github.com/namchokGithub/vocabunny-core-api/internal/core/helper"
)

func toSectionResponse(item domain.Section) SectionResponse {
	return SectionResponse{ID: item.ID.String(), Slug: item.Slug, Title: item.Title, Description: item.Description, OrderNo: item.OrderNo, IsPublished: item.IsPublished, CreatedAt: item.CreatedAt.Format(time.RFC3339), UpdatedAt: item.UpdatedAt.Format(time.RFC3339), CreatedBy: item.CreatedBy, UpdatedBy: item.UpdatedBy}
}

func toLessonResponse(item domain.Lesson) LessonResponse {
	return LessonResponse{ID: item.ID.String(), SectionID: item.SectionID.String(), Slug: item.Slug, Title: item.Title, Description: item.Description, OrderNo: item.OrderNo, IsPublished: item.IsPublished, CreatedAt: item.CreatedAt.Format(time.RFC3339), UpdatedAt: item.UpdatedAt.Format(time.RFC3339), CreatedBy: item.CreatedBy, UpdatedBy: item.UpdatedBy}
}

func toUnitResponse(item domain.Unit) UnitResponse {
	return UnitResponse{ID: item.ID.String(), LessonID: item.LessonID.String(), Slug: item.Slug, Title: item.Title, Description: item.Description, OrderNo: item.OrderNo, IsPublished: item.IsPublished, CreatedAt: item.CreatedAt.Format(time.RFC3339), UpdatedAt: item.UpdatedAt.Format(time.RFC3339), CreatedBy: item.CreatedBy, UpdatedBy: item.UpdatedBy}
}

func toQuestionSetResponse(item domain.QuestionSet) QuestionSetResponse {
	return QuestionSetResponse{ID: item.ID.String(), UnitID: item.UnitID.String(), Slug: item.Slug, Title: item.Title, Description: item.Description, OrderNo: item.OrderNo, EstimatedSeconds: item.EstimatedSeconds, IsPublished: item.IsPublished, Version: item.Version, CreatedAt: item.CreatedAt.Format(time.RFC3339), UpdatedAt: item.UpdatedAt.Format(time.RFC3339), CreatedBy: item.CreatedBy, UpdatedBy: item.UpdatedBy}
}

func toQuestionChoiceResponse(item domain.QuestionChoice) QuestionChoiceResponse {
	return QuestionChoiceResponse{ID: item.ID.String(), QuestionID: item.QuestionID.String(), ChoiceText: item.ChoiceText, ChoiceOrder: item.ChoiceOrder, IsCorrect: item.IsCorrect, CreatedAt: item.CreatedAt.Format(time.RFC3339), UpdatedAt: item.UpdatedAt.Format(time.RFC3339), CreatedBy: item.CreatedBy, UpdatedBy: item.UpdatedBy}
}

func toTagResponse(item domain.Tag) TagResponse {
	return TagResponse{ID: item.ID.String(), Name: item.Name, CreatedAt: item.CreatedAt.Format(time.RFC3339), UpdatedAt: item.UpdatedAt.Format(time.RFC3339), CreatedBy: item.CreatedBy, UpdatedBy: item.UpdatedBy}
}

func toQuestionResponse(item domain.Question) QuestionResponse {
	response := QuestionResponse{
		ID:            item.ID.String(),
		QuestionSetID: item.QuestionSetID.String(),
		Type:          item.Type,
		QuestionText:  item.QuestionText,
		BlankPosition: item.BlankPosition,
		Explanation:   item.Explanation,
		ImageURL:      item.ImageURL,
		Difficulty:    item.Difficulty,
		OrderNo:       item.OrderNo,
		IsActive:      item.IsActive,
		CreatedAt:     item.CreatedAt.Format(time.RFC3339),
		UpdatedAt:     item.UpdatedAt.Format(time.RFC3339),
		CreatedBy:     item.CreatedBy,
		UpdatedBy:     item.UpdatedBy,
	}
	for _, choice := range item.Choices {
		response.Choices = append(response.Choices, toQuestionChoiceResponse(choice))
	}
	for _, tag := range item.Tags {
		response.Tags = append(response.Tags, toTagResponse(tag))
	}
	return response
}

func parseUUID(value string, field string) (uuid.UUID, error) {
	parsed, err := uuid.Parse(strings.TrimSpace(value))
	if err != nil {
		return uuid.Nil, helper.BadRequest("invalid_"+field, field+" must be a valid uuid", err)
	}
	return parsed, nil
}

func parseOptionalUUIDField(value *string, field string) (*uuid.UUID, error) {
	if value == nil {
		return nil, nil
	}
	trimmed := strings.TrimSpace(*value)
	if trimmed == "" {
		return nil, nil
	}
	parsed, err := parseUUID(trimmed, field)
	if err != nil {
		return nil, err
	}
	return &parsed, nil
}

func parseUUIDList(values []string, field string) ([]uuid.UUID, error) {
	items := make([]uuid.UUID, 0, len(values))
	for _, value := range values {
		if strings.TrimSpace(value) == "" {
			continue
		}
		parsed, err := parseUUID(value, field)
		if err != nil {
			return nil, err
		}
		items = append(items, parsed)
	}
	return items, nil
}

func toChoiceInputs(values []QuestionChoicePayload) ([]domain.QuestionChoiceInput, error) {
	items := make([]domain.QuestionChoiceInput, 0, len(values))
	for _, value := range values {
		item := domain.QuestionChoiceInput{
			ChoiceText:  strings.TrimSpace(value.ChoiceText),
			ChoiceOrder: value.ChoiceOrder,
			IsCorrect:   value.IsCorrect,
		}
		if value.ID != nil && strings.TrimSpace(*value.ID) != "" {
			parsed, err := parseUUID(*value.ID, "choice_id")
			if err != nil {
				return nil, err
			}
			item.ID = parsed
		}
		items = append(items, item)
	}
	return items, nil
}
