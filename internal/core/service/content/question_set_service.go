package content

import (
	"context"

	"github.com/google/uuid"
	"github.com/namchokGithub/vocabunny-core-api/internal/core/domain"
	"github.com/namchokGithub/vocabunny-core-api/internal/core/helper"
	"github.com/namchokGithub/vocabunny-core-api/internal/core/port"
)

type QuestionSetServiceDependencies struct {
	QuestionSetRepository port.QuestionSetRepository
	UnitRepository        port.UnitRepository
	TxManager             port.TransactionManager
}

type questionSetService struct {
	questionSetRepository port.QuestionSetRepository
	unitRepository        port.UnitRepository
	txManager             port.TransactionManager
}

func NewQuestionSetService(deps QuestionSetServiceDependencies) port.QuestionSetService {
	return &questionSetService{questionSetRepository: deps.QuestionSetRepository, unitRepository: deps.UnitRepository, txManager: deps.TxManager}
}

func (s *questionSetService) Create(ctx context.Context, input domain.QuestionSetCreateInput) (domain.QuestionSet, error) {
	input.Slug = normalizeSlug(input.Slug)
	input.Title = normalizeText(input.Title)
	if input.UnitID == uuid.Nil {
		return domain.QuestionSet{}, helper.BadRequest("invalid_unit_id", "unit_id is required", nil)
	}
	if err := validateRequired(input.Slug, "invalid_slug", "slug is required"); err != nil {
		return domain.QuestionSet{}, err
	}
	if err := validateRequired(input.Title, "invalid_title", "title is required"); err != nil {
		return domain.QuestionSet{}, err
	}
	if input.Version <= 0 {
		input.Version = 1
	}

	var created domain.QuestionSet
	err := s.txManager.RunInTx(ctx, func(txCtx context.Context) error {
		if _, err := s.unitRepository.FindByID(txCtx, input.UnitID); err != nil {
			return helper.BadRequest("invalid_unit_id", "unit does not exist", err)
		}
		exists, err := s.questionSetRepository.ExistsBySlugVersion(txCtx, input.UnitID, input.Slug, input.Version, nil)
		if err != nil {
			return err
		}
		if exists {
			return helper.Conflict("duplicate_question_set_slug", "question set slug and version already exist in this unit", nil)
		}
		created, err = s.questionSetRepository.Create(txCtx, domain.QuestionSet{
			UnitID:           input.UnitID,
			Slug:             input.Slug,
			Title:            input.Title,
			Description:      input.Description,
			OrderNo:          input.OrderNo,
			EstimatedSeconds: input.EstimatedSeconds,
			IsPublished:      input.IsPublished,
			Version:          input.Version,
			AuditFields:      domain.AuditFields{CreatedBy: input.ActorID, UpdatedBy: input.ActorID},
		})
		return err
	})
	return created, err
}

func (s *questionSetService) Update(ctx context.Context, input domain.QuestionSetUpdateInput) (domain.QuestionSet, error) {
	var updated domain.QuestionSet
	err := s.txManager.RunInTx(ctx, func(txCtx context.Context) error {
		current, err := s.questionSetRepository.FindByID(txCtx, input.ID)
		if err != nil {
			return err
		}
		targetUnitID := current.UnitID
		targetSlug := current.Slug
		targetVersion := current.Version
		if input.UnitID.Set {
			if input.UnitID.Value == uuid.Nil {
				return helper.BadRequest("invalid_unit_id", "unit_id cannot be empty", nil)
			}
			if _, err := s.unitRepository.FindByID(txCtx, input.UnitID.Value); err != nil {
				return helper.BadRequest("invalid_unit_id", "unit does not exist", err)
			}
			targetUnitID = input.UnitID.Value
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
		if input.Version.Set {
			if input.Version.Value <= 0 {
				return helper.BadRequest("invalid_version", "version must be greater than zero", nil)
			}
			targetVersion = input.Version.Value
		}
		if targetUnitID != current.UnitID || targetSlug != current.Slug || targetVersion != current.Version {
			exists, err := s.questionSetRepository.ExistsBySlugVersion(txCtx, targetUnitID, targetSlug, targetVersion, &input.ID)
			if err != nil {
				return err
			}
			if exists {
				return helper.Conflict("duplicate_question_set_slug", "question set slug and version already exist in this unit", nil)
			}
		}
		updated, err = s.questionSetRepository.Update(txCtx, input)
		return err
	})
	return updated, err
}

func (s *questionSetService) Delete(ctx context.Context, id uuid.UUID, actorID string) error {
	return s.txManager.RunInTx(ctx, func(txCtx context.Context) error {
		if _, err := s.questionSetRepository.FindByID(txCtx, id); err != nil {
			return err
		}
		return s.questionSetRepository.Delete(txCtx, id, actorID)
	})
}

func (s *questionSetService) FindByID(ctx context.Context, id uuid.UUID) (domain.QuestionSet, error) {
	return s.questionSetRepository.FindByID(ctx, id)
}

func (s *questionSetService) FindAll(ctx context.Context, query domain.QuestionSetQuery) (domain.PageResult[domain.QuestionSet], error) {
	return s.questionSetRepository.FindAll(ctx, query)
}
