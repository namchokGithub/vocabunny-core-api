package content

import (
	"time"

	"github.com/namchokGithub/vocabunny-core-api/internal/core/domain"
)

func toQuestionSetResponse(item domain.QuestionSet) QuestionSetResponse {
	resp := QuestionSetResponse{
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
	if item.Unit != nil {
		resp.Unit = &UnitSummaryDTO{
			ID:       item.Unit.ID.String(),
			LessonID: item.Unit.LessonID.String(),
			Slug:     item.Unit.Slug,
			Title:    item.Unit.Title,
		}
	}
	if item.Lesson != nil {
		resp.Lesson = &LessonSummaryDTO{
			ID:        item.Lesson.ID.String(),
			SectionID: item.Lesson.SectionID.String(),
			Slug:      item.Lesson.Slug,
			Title:     item.Lesson.Title,
		}
	}
	return resp
}
