package content

import (
	"context"

	"github.com/google/uuid"
	"github.com/namchokGithub/vocabunny-core-api/internal/core/domain"
	"github.com/namchokGithub/vocabunny-core-api/internal/core/helper"
	"github.com/namchokGithub/vocabunny-core-api/internal/core/port"
)

type LessonServiceDependencies struct {
	LessonRepository  port.LessonRepository
	SectionRepository port.SectionRepository
	TxManager         port.TransactionManager
}

type lessonService struct {
	lessonRepository  port.LessonRepository
	sectionRepository port.SectionRepository
	txManager         port.TransactionManager
}

func NewLessonService(deps LessonServiceDependencies) port.LessonService {
	return &lessonService{lessonRepository: deps.LessonRepository, sectionRepository: deps.SectionRepository, txManager: deps.TxManager}
}

func (s *lessonService) Create(ctx context.Context, input domain.LessonCreateInput) (domain.Lesson, error) {
	input.Slug = normalizeSlug(input.Slug)
	input.Title = normalizeText(input.Title)
	if input.SectionID == uuid.Nil {
		return domain.Lesson{}, helper.BadRequest("invalid_section_id", "section_id is required", nil)
	}
	if err := validateRequired(input.Slug, "invalid_slug", "slug is required"); err != nil {
		return domain.Lesson{}, err
	}
	if err := validateRequired(input.Title, "invalid_title", "title is required"); err != nil {
		return domain.Lesson{}, err
	}

	var created domain.Lesson
	err := s.txManager.RunInTx(ctx, func(txCtx context.Context) error {
		if _, err := s.sectionRepository.FindByID(txCtx, input.SectionID); err != nil {
			return helper.BadRequest("invalid_section_id", "section does not exist", err)
		}
		exists, err := s.lessonRepository.ExistsBySlug(txCtx, input.SectionID, input.Slug, nil)
		if err != nil {
			return err
		}
		if exists {
			return helper.Conflict("duplicate_lesson_slug", "lesson slug already exists in this section", nil)
		}
		created, err = s.lessonRepository.Create(txCtx, domain.Lesson{
			SectionID:   input.SectionID,
			Slug:        input.Slug,
			Title:       input.Title,
			Description: input.Description,
			OrderNo:     input.OrderNo,
			IsPublished: input.IsPublished,
			AuditFields: domain.AuditFields{CreatedBy: input.ActorID, UpdatedBy: input.ActorID},
		})
		return err
	})
	return created, err
}

func (s *lessonService) Update(ctx context.Context, input domain.LessonUpdateInput) (domain.Lesson, error) {
	var updated domain.Lesson
	err := s.txManager.RunInTx(ctx, func(txCtx context.Context) error {
		current, err := s.lessonRepository.FindByID(txCtx, input.ID)
		if err != nil {
			return err
		}
		targetSectionID := current.SectionID
		targetSlug := current.Slug
		if input.SectionID.Set {
			if input.SectionID.Value == uuid.Nil {
				return helper.BadRequest("invalid_section_id", "section_id cannot be empty", nil)
			}
			if _, err := s.sectionRepository.FindByID(txCtx, input.SectionID.Value); err != nil {
				return helper.BadRequest("invalid_section_id", "section does not exist", err)
			}
			targetSectionID = input.SectionID.Value
		}
		if input.Slug.Set {
			input.Slug.Value = normalizeSlug(input.Slug.Value)
			if err := validateRequired(input.Slug.Value, "invalid_slug", "slug cannot be empty"); err != nil {
				return err
			}
			targetSlug = input.Slug.Value
		}
		if input.Title.Set {
			input.Title.Value = normalizeText(input.Title.Value)
			if err := validateRequired(input.Title.Value, "invalid_title", "title cannot be empty"); err != nil {
				return err
			}
		}
		if targetSectionID != current.SectionID || targetSlug != current.Slug {
			exists, err := s.lessonRepository.ExistsBySlug(txCtx, targetSectionID, targetSlug, &input.ID)
			if err != nil {
				return err
			}
			if exists {
				return helper.Conflict("duplicate_lesson_slug", "lesson slug already exists in this section", nil)
			}
		}
		updated, err = s.lessonRepository.Update(txCtx, input)
		return err
	})
	return updated, err
}

func (s *lessonService) Delete(ctx context.Context, id uuid.UUID, actorID string) error {
	return s.txManager.RunInTx(ctx, func(txCtx context.Context) error {
		if _, err := s.lessonRepository.FindByID(txCtx, id); err != nil {
			return err
		}
		return s.lessonRepository.Delete(txCtx, id, actorID)
	})
}

func (s *lessonService) FindByID(ctx context.Context, id uuid.UUID) (domain.Lesson, error) {
	return s.lessonRepository.FindByID(ctx, id)
}

func (s *lessonService) FindAll(ctx context.Context, query domain.LessonQuery) (domain.PageResult[domain.Lesson], error) {
	return s.lessonRepository.FindAll(ctx, query)
}
