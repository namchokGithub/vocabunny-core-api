package content

import (
	"time"

	"github.com/namchokGithub/vocabunny-core-api/internal/core/domain"
)

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
