package content

import (
	"strings"

	"github.com/google/uuid"
	"github.com/namchokGithub/vocabunny-core-api/internal/core/domain"
	"github.com/namchokGithub/vocabunny-core-api/internal/core/helper"
)

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
