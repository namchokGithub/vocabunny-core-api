package content

import (
	"time"

	"github.com/namchokGithub/vocabunny-core-api/internal/core/domain"
)

func toUnitResponse(item domain.Unit) UnitResponse {
	return UnitResponse{
		ID:          item.ID.String(),
		LessonID:    item.LessonID.String(),
		Slug:        item.Slug,
		Title:       item.Title,
		Description: item.Description,
		OrderNo:     item.OrderNo,
		IsPublished: item.IsPublished,
		CreatedAt:   item.CreatedAt.Format(time.RFC3339),
		UpdatedAt:   item.UpdatedAt.Format(time.RFC3339),
		CreatedBy:   item.CreatedBy,
		UpdatedBy:   item.UpdatedBy,
	}
}
