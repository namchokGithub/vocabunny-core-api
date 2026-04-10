package content

import (
	"time"

	"github.com/namchokGithub/vocabunny-core-api/internal/core/domain"
)

func toQuestionChoiceResponse(item domain.QuestionChoice) QuestionChoiceResponse {
	return QuestionChoiceResponse{
		ID:          item.ID.String(),
		QuestionID:  item.QuestionID.String(),
		ChoiceText:  item.ChoiceText,
		ChoiceOrder: item.ChoiceOrder,
		IsCorrect:   item.IsCorrect,
		CreatedAt:   item.CreatedAt.Format(time.RFC3339),
		UpdatedAt:   item.UpdatedAt.Format(time.RFC3339),
		CreatedBy:   item.CreatedBy,
		UpdatedBy:   item.UpdatedBy,
	}
}
