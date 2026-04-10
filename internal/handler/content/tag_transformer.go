package content

import (
	"time"

	"github.com/namchokGithub/vocabunny-core-api/internal/core/domain"
)

func toTagResponse(item domain.Tag) TagResponse {
	return TagResponse{
		ID:        item.ID.String(),
		Name:      item.Name,
		CreatedAt: item.CreatedAt.Format(time.RFC3339),
		UpdatedAt: item.UpdatedAt.Format(time.RFC3339),
		CreatedBy: item.CreatedBy,
		UpdatedBy: item.UpdatedBy,
	}
}
