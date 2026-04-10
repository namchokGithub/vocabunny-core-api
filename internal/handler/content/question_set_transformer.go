package content

import (
	"time"

	"github.com/namchokGithub/vocabunny-core-api/internal/core/domain"
)

func toQuestionSetResponse(item domain.QuestionSet) QuestionSetResponse {
	return QuestionSetResponse{
		ID:               item.ID.String(),
		UnitID:           item.UnitID.String(),
		Slug:             item.Slug,
		Title:            item.Title,
		Description:      item.Description,
		OrderNo:          item.OrderNo,
		EstimatedSeconds: item.EstimatedSeconds,
		IsPublished:      item.IsPublished,
		Version:          item.Version,
		CreatedAt:        item.CreatedAt.Format(time.RFC3339),
		UpdatedAt:        item.UpdatedAt.Format(time.RFC3339),
		CreatedBy:        item.CreatedBy,
		UpdatedBy:        item.UpdatedBy,
	}
}
