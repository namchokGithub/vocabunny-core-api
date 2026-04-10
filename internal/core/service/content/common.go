package content

import (
	"context"
	"strings"

	"github.com/google/uuid"
	"github.com/namchokGithub/vocabunny-core-api/internal/core/domain"
	"github.com/namchokGithub/vocabunny-core-api/internal/core/helper"
)

func normalizeSlug(value string) string {
	return strings.TrimSpace(strings.ToLower(value))
}

func normalizeText(value string) string {
	return strings.TrimSpace(value)
}

func validateRequired(value string, code string, message string) error {
	if strings.TrimSpace(value) == "" {
		return helper.BadRequest(code, message, nil)
	}
	return nil
}

func ensureTagIDsExist(repo interface {
	FindByID(ctx context.Context, id uuid.UUID) (domain.Tag, error)
}, ctx context.Context, tagIDs []uuid.UUID) error {
	for _, tagID := range tagIDs {
		if tagID == uuid.Nil {
			return helper.BadRequest("invalid_tag_id", "tag_id must be a valid uuid", nil)
		}
		if _, err := repo.FindByID(ctx, tagID); err != nil {
			return helper.BadRequest("invalid_tag_id", "one or more tags do not exist", err)
		}
	}
	return nil
}
