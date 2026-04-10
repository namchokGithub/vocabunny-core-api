package content

import (
	"time"

	"github.com/namchokGithub/vocabunny-core-api/internal/core/domain"
)

func toSectionResponse(item domain.Section) SectionResponse {
	return SectionResponse{
		ID:          item.ID.String(),
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
