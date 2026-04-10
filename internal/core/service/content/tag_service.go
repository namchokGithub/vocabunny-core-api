package content

import (
	"context"

	"github.com/google/uuid"
	"github.com/namchokGithub/vocabunny-core-api/internal/core/domain"
	"github.com/namchokGithub/vocabunny-core-api/internal/core/helper"
	"github.com/namchokGithub/vocabunny-core-api/internal/core/port"
)

type TagServiceDependencies struct {
	TagRepository port.TagRepository
	TxManager     port.TransactionManager
}

type tagService struct {
	tagRepository port.TagRepository
	txManager     port.TransactionManager
}

func NewTagService(deps TagServiceDependencies) port.TagService {
	return &tagService{tagRepository: deps.TagRepository, txManager: deps.TxManager}
}

func (s *tagService) Create(ctx context.Context, input domain.TagCreateInput) (domain.Tag, error) {
	input.Name = normalizeText(input.Name)
	if err := validateRequired(input.Name, "invalid_name", "name is required"); err != nil {
		return domain.Tag{}, err
	}

	var created domain.Tag
	err := s.txManager.RunInTx(ctx, func(txCtx context.Context) error {
		exists, err := s.tagRepository.ExistsByName(txCtx, input.Name, nil)
		if err != nil {
			return err
		}
		if exists {
			return helper.Conflict("duplicate_tag_name", "tag name already exists", nil)
		}
		created, err = s.tagRepository.Create(txCtx, domain.Tag{
			Name:        input.Name,
			AuditFields: domain.AuditFields{CreatedBy: input.ActorID, UpdatedBy: input.ActorID},
		})
		return err
	})
	return created, err
}

func (s *tagService) Update(ctx context.Context, input domain.TagUpdateInput) (domain.Tag, error) {
	var updated domain.Tag
	err := s.txManager.RunInTx(ctx, func(txCtx context.Context) error {
		current, err := s.tagRepository.FindByID(txCtx, input.ID)
		if err != nil {
			return err
		}
		if input.Name.Set {
			input.Name.Value = normalizeText(input.Name.Value)
			if err := validateRequired(input.Name.Value, "invalid_name", "name cannot be empty"); err != nil {
				return err
			}
			if current.Name != input.Name.Value {
				exists, err := s.tagRepository.ExistsByName(txCtx, input.Name.Value, &input.ID)
				if err != nil {
					return err
				}
				if exists {
					return helper.Conflict("duplicate_tag_name", "tag name already exists", nil)
				}
			}
		}
		updated, err = s.tagRepository.Update(txCtx, input)
		return err
	})
	return updated, err
}

func (s *tagService) Delete(ctx context.Context, id uuid.UUID, actorID string) error {
	return s.txManager.RunInTx(ctx, func(txCtx context.Context) error {
		if _, err := s.tagRepository.FindByID(txCtx, id); err != nil {
			return err
		}
		return s.tagRepository.Delete(txCtx, id, actorID)
	})
}

func (s *tagService) FindByID(ctx context.Context, id uuid.UUID) (domain.Tag, error) {
	return s.tagRepository.FindByID(ctx, id)
}

func (s *tagService) FindAll(ctx context.Context, query domain.TagQuery) (domain.PageResult[domain.Tag], error) {
	return s.tagRepository.FindAll(ctx, query)
}
