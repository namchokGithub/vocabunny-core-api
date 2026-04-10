package content

import (
	"context"

	"github.com/google/uuid"
	"github.com/namchokGithub/vocabunny-core-api/internal/core/domain"
	"github.com/namchokGithub/vocabunny-core-api/internal/core/helper"
	"github.com/namchokGithub/vocabunny-core-api/internal/core/port"
)

type SectionServiceDependencies struct {
	SectionRepository port.SectionRepository
	TxManager         port.TransactionManager
}

type sectionService struct {
	sectionRepository port.SectionRepository
	txManager         port.TransactionManager
}

func NewSectionService(deps SectionServiceDependencies) port.SectionService {
	return &sectionService{sectionRepository: deps.SectionRepository, txManager: deps.TxManager}
}

func (s *sectionService) Create(ctx context.Context, input domain.SectionCreateInput) (domain.Section, error) {
	input.Slug = normalizeSlug(input.Slug)
	input.Title = normalizeText(input.Title)
	if err := validateRequired(input.Slug, "invalid_slug", "slug is required"); err != nil {
		return domain.Section{}, err
	}
	if err := validateRequired(input.Title, "invalid_title", "title is required"); err != nil {
		return domain.Section{}, err
	}

	var created domain.Section
	err := s.txManager.RunInTx(ctx, func(txCtx context.Context) error {
		exists, err := s.sectionRepository.ExistsBySlug(txCtx, input.Slug, nil)
		if err != nil {
			return err
		}
		if exists {
			return helper.Conflict("duplicate_section_slug", "section slug already exists", nil)
		}
		created, err = s.sectionRepository.Create(txCtx, domain.Section{
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

func (s *sectionService) Update(ctx context.Context, input domain.SectionUpdateInput) (domain.Section, error) {
	var updated domain.Section
	err := s.txManager.RunInTx(ctx, func(txCtx context.Context) error {
		current, err := s.sectionRepository.FindByID(txCtx, input.ID)
		if err != nil {
			return err
		}
		if input.Slug.Set {
			input.Slug.Value = normalizeSlug(input.Slug.Value)
			if err := validateRequired(input.Slug.Value, "invalid_slug", "slug cannot be empty"); err != nil {
				return err
			}
			if current.Slug != input.Slug.Value {
				exists, err := s.sectionRepository.ExistsBySlug(txCtx, input.Slug.Value, &input.ID)
				if err != nil {
					return err
				}
				if exists {
					return helper.Conflict("duplicate_section_slug", "section slug already exists", nil)
				}
			}
		}
		if input.Title.Set {
			input.Title.Value = normalizeText(input.Title.Value)
			if err := validateRequired(input.Title.Value, "invalid_title", "title cannot be empty"); err != nil {
				return err
			}
		}
		updated, err = s.sectionRepository.Update(txCtx, input)
		return err
	})
	return updated, err
}

func (s *sectionService) Delete(ctx context.Context, id uuid.UUID, actorID string) error {
	return s.txManager.RunInTx(ctx, func(txCtx context.Context) error {
		if _, err := s.sectionRepository.FindByID(txCtx, id); err != nil {
			return err
		}
		return s.sectionRepository.Delete(txCtx, id, actorID)
	})
}

func (s *sectionService) FindByID(ctx context.Context, id uuid.UUID) (domain.Section, error) {
	return s.sectionRepository.FindByID(ctx, id)
}

func (s *sectionService) FindAll(ctx context.Context, query domain.SectionQuery) (domain.PageResult[domain.Section], error) {
	return s.sectionRepository.FindAll(ctx, query)
}
