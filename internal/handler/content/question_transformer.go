package content

import (
	"time"

	"github.com/namchokGithub/vocabunny-core-api/internal/core/domain"
)

func toQuestionResponse(item domain.Question) QuestionResponse {
	resp := QuestionResponse{
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
	if item.QuestionSet != nil {
		resp.QuestionSet = &QuestionSetSummaryDTO{
			ID:      item.QuestionSet.ID.String(),
			UnitID:  item.QuestionSet.UnitID.String(),
			Slug:    item.QuestionSet.Slug,
			Title:   item.QuestionSet.Title,
			Version: item.QuestionSet.Version,
		}
	}
	for _, choice := range item.Choices {
		resp.Choices = append(resp.Choices, toQuestionChoiceResponse(choice))
	}
	for _, tag := range item.Tags {
		resp.Tags = append(resp.Tags, toTagResponse(tag))
	}
	return resp
}
