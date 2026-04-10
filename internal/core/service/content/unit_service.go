package content

import (
	"context"

	"github.com/google/uuid"
	"github.com/namchokGithub/vocabunny-core-api/internal/core/domain"
	"github.com/namchokGithub/vocabunny-core-api/internal/core/helper"
	"github.com/namchokGithub/vocabunny-core-api/internal/core/port"
)

type UnitServiceDependencies struct {
	UnitRepository   port.UnitRepository
	LessonRepository port.LessonRepository
	TxManager        port.TransactionManager
}

type unitService struct {
	unitRepository   port.UnitRepository
	lessonRepository port.LessonRepository
	txManager        port.TransactionManager
}

func NewUnitService(deps UnitServiceDependencies) port.UnitService {
	return &unitService{unitRepository: deps.UnitRepository, lessonRepository: deps.LessonRepository, txManager: deps.TxManager}
}

func (s *unitService) Create(ctx context.Context, input domain.UnitCreateInput) (domain.Unit, error) {
	input.Slug = normalizeSlug(input.Slug)
	input.Title = normalizeText(input.Title)
	if input.LessonID == uuid.Nil {
		return domain.Unit{}, helper.BadRequest("invalid_lesson_id", "lesson_id is required", nil)
	}
	if err := validateRequired(input.Slug, "invalid_slug", "slug is required"); err != nil {
		return domain.Unit{}, err
	}
	if err := validateRequired(input.Title, "invalid_title", "title is required"); err != nil {
		return domain.Unit{}, err
	}

	var created domain.Unit
	err := s.txManager.RunInTx(ctx, func(txCtx context.Context) error {
		if _, err := s.lessonRepository.FindByID(txCtx, input.LessonID); err != nil {
			return helper.BadRequest("invalid_lesson_id", "lesson does not exist", err)
		}
		exists, err := s.unitRepository.ExistsBySlug(txCtx, input.LessonID, input.Slug, nil)
		if err != nil {
			return err
		}
		if exists {
			return helper.Conflict("duplicate_unit_slug", "unit slug already exists in this lesson", nil)
		}
		created, err = s.unitRepository.Create(txCtx, domain.Unit{
			LessonID:    input.LessonID,
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

func (s *unitService) Update(ctx context.Context, input domain.UnitUpdateInput) (domain.Unit, error) {
	var updated domain.Unit
	err := s.txManager.RunInTx(ctx, func(txCtx context.Context) error {
		current, err := s.unitRepository.FindByID(txCtx, input.ID)
		if err != nil {
			return err
		}
		targetLessonID := current.LessonID
		targetSlug := current.Slug
		if input.LessonID.Set {
			if input.LessonID.Value == uuid.Nil {
				return helper.BadRequest("invalid_lesson_id", "lesson_id cannot be empty", nil)
			}
			if _, err := s.lessonRepository.FindByID(txCtx, input.LessonID.Value); err != nil {
				return helper.BadRequest("invalid_lesson_id", "lesson does not exist", err)
			}
			targetLessonID = input.LessonID.Value
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
		if targetLessonID != current.LessonID || targetSlug != current.Slug {
			exists, err := s.unitRepository.ExistsBySlug(txCtx, targetLessonID, targetSlug, &input.ID)
			if err != nil {
				return err
			}
			if exists {
				return helper.Conflict("duplicate_unit_slug", "unit slug already exists in this lesson", nil)
			}
		}
		updated, err = s.unitRepository.Update(txCtx, input)
		return err
	})
	return updated, err
}

func (s *unitService) Delete(ctx context.Context, id uuid.UUID, actorID string) error {
	return s.txManager.RunInTx(ctx, func(txCtx context.Context) error {
		if _, err := s.unitRepository.FindByID(txCtx, id); err != nil {
			return err
		}
		return s.unitRepository.Delete(txCtx, id, actorID)
	})
}

func (s *unitService) FindByID(ctx context.Context, id uuid.UUID) (domain.Unit, error) {
	return s.unitRepository.FindByID(ctx, id)
}

func (s *unitService) FindAll(ctx context.Context, query domain.UnitQuery) (domain.PageResult[domain.Unit], error) {
	return s.unitRepository.FindAll(ctx, query)
}
