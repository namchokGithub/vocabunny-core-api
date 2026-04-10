package content

import (
	"context"

	"github.com/google/uuid"
	"github.com/namchokGithub/vocabunny-core-api/internal/core/domain"
	"github.com/namchokGithub/vocabunny-core-api/internal/core/helper"
	"github.com/namchokGithub/vocabunny-core-api/internal/core/port"
)

type QuestionChoiceServiceDependencies struct {
	QuestionChoiceRepository port.QuestionChoiceRepository
	QuestionRepository       port.QuestionRepository
	TxManager                port.TransactionManager
}

type questionChoiceService struct {
	questionChoiceRepository port.QuestionChoiceRepository
	questionRepository       port.QuestionRepository
	txManager                port.TransactionManager
}

func NewQuestionChoiceService(deps QuestionChoiceServiceDependencies) port.QuestionChoiceService {
	return &questionChoiceService{
		questionChoiceRepository: deps.QuestionChoiceRepository,
		questionRepository:       deps.QuestionRepository,
		txManager:                deps.TxManager,
	}
}

func (s *questionChoiceService) Create(ctx context.Context, input domain.QuestionChoiceCreateInput) (domain.QuestionChoice, error) {
	input.ChoiceText = normalizeText(input.ChoiceText)
	if input.QuestionID == uuid.Nil {
		return domain.QuestionChoice{}, helper.BadRequest("invalid_question_id", "question_id is required", nil)
	}
	if err := validateRequired(input.ChoiceText, "invalid_choice_text", "choice_text is required"); err != nil {
		return domain.QuestionChoice{}, err
	}

	var created domain.QuestionChoice
	var createErr error
	err := s.txManager.RunInTx(ctx, func(txCtx context.Context) error {
		if _, err := s.questionRepository.FindByID(txCtx, input.QuestionID); err != nil {
			return helper.BadRequest("invalid_question_id", "question does not exist", err)
		}
		created, createErr = s.questionChoiceRepository.Create(txCtx, domain.QuestionChoice{
			QuestionID:  input.QuestionID,
			ChoiceText:  input.ChoiceText,
			ChoiceOrder: input.ChoiceOrder,
			IsCorrect:   input.IsCorrect,
			AuditFields: domain.AuditFields{CreatedBy: input.ActorID, UpdatedBy: input.ActorID},
		})
		return createErr
	})
	return created, err
}

func (s *questionChoiceService) Update(ctx context.Context, input domain.QuestionChoiceUpdateInput) (domain.QuestionChoice, error) {
	if input.ChoiceText.Set {
		input.ChoiceText.Value = normalizeText(input.ChoiceText.Value)
		if err := validateRequired(input.ChoiceText.Value, "invalid_choice_text", "choice_text cannot be empty"); err != nil {
			return domain.QuestionChoice{}, err
		}
	}
	return s.questionChoiceRepository.Update(ctx, input)
}

func (s *questionChoiceService) Delete(ctx context.Context, id uuid.UUID, actorID string) error {
	return s.txManager.RunInTx(ctx, func(txCtx context.Context) error {
		if _, err := s.questionChoiceRepository.FindByID(txCtx, id); err != nil {
			return err
		}
		return s.questionChoiceRepository.Delete(txCtx, id, actorID)
	})
}

func (s *questionChoiceService) FindByID(ctx context.Context, id uuid.UUID) (domain.QuestionChoice, error) {
	return s.questionChoiceRepository.FindByID(ctx, id)
}

func (s *questionChoiceService) FindAll(ctx context.Context, query domain.QuestionChoiceQuery) (domain.PageResult[domain.QuestionChoice], error) {
	return s.questionChoiceRepository.FindAll(ctx, query)
}
