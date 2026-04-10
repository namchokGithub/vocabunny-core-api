package content

import (
	"context"

	"github.com/google/uuid"
	"github.com/namchokGithub/vocabunny-core-api/internal/core/domain"
	"github.com/namchokGithub/vocabunny-core-api/internal/core/helper"
	"github.com/namchokGithub/vocabunny-core-api/internal/core/port"
)

type QuestionServiceDependencies struct {
	QuestionRepository    port.QuestionRepository
	QuestionSetRepository port.QuestionSetRepository
	TagRepository         port.TagRepository
	TxManager             port.TransactionManager
}

type questionService struct {
	questionRepository    port.QuestionRepository
	questionSetRepository port.QuestionSetRepository
	tagRepository         port.TagRepository
	txManager             port.TransactionManager
}

func NewQuestionService(deps QuestionServiceDependencies) port.QuestionService {
	return &questionService{
		questionRepository:    deps.QuestionRepository,
		questionSetRepository: deps.QuestionSetRepository,
		tagRepository:         deps.TagRepository,
		txManager:             deps.TxManager,
	}
}

func (s *questionService) Create(ctx context.Context, input domain.QuestionCreateInput) (domain.Question, error) {
	input.Type = normalizeText(input.Type)
	input.QuestionText = normalizeText(input.QuestionText)
	if input.QuestionSetID == uuid.Nil {
		return domain.Question{}, helper.BadRequest("invalid_question_set_id", "question_set_id is required", nil)
	}
	if err := validateRequired(input.Type, "invalid_type", "type is required"); err != nil {
		return domain.Question{}, err
	}
	if err := validateRequired(input.QuestionText, "invalid_question_text", "question_text is required"); err != nil {
		return domain.Question{}, err
	}
	if input.Difficulty <= 0 {
		input.Difficulty = 1
	}

	var created domain.Question
	err := s.txManager.RunInTx(ctx, func(txCtx context.Context) error {
		if _, err := s.questionSetRepository.FindByID(txCtx, input.QuestionSetID); err != nil {
			return helper.BadRequest("invalid_question_set_id", "question set does not exist", err)
		}
		for _, choice := range input.Choices {
			if normalizeText(choice.ChoiceText) == "" {
				return helper.BadRequest("invalid_choice_text", "choice_text is required", nil)
			}
		}
		if err := ensureTagIDsExist(s.tagRepository, txCtx, input.TagIDs); err != nil {
			return err
		}
		question, err := s.questionRepository.Create(txCtx, domain.Question{
			QuestionSetID: input.QuestionSetID,
			Type:          input.Type,
			QuestionText:  input.QuestionText,
			BlankPosition: input.BlankPosition,
			Explanation:   input.Explanation,
			ImageURL:      input.ImageURL,
			Difficulty:    input.Difficulty,
			OrderNo:       input.OrderNo,
			IsActive:      input.IsActive,
			AuditFields:   domain.AuditFields{CreatedBy: input.ActorID, UpdatedBy: input.ActorID},
		})
		if err != nil {
			return err
		}
		if err := s.questionRepository.ReplaceChoices(txCtx, question.ID, input.Choices, input.ActorID); err != nil {
			return err
		}
		if err := s.questionRepository.ReplaceTags(txCtx, question.ID, input.TagIDs, input.ActorID); err != nil {
			return err
		}
		created, err = s.questionRepository.FindByID(txCtx, question.ID)
		return err
	})
	return created, err
}

func (s *questionService) Update(ctx context.Context, input domain.QuestionUpdateInput) (domain.Question, error) {
	var updated domain.Question
	err := s.txManager.RunInTx(ctx, func(txCtx context.Context) error {
		current, err := s.questionRepository.FindByID(txCtx, input.ID)
		if err != nil {
			return err
		}
		if input.QuestionSetID.Set {
			if input.QuestionSetID.Value == uuid.Nil {
				return helper.BadRequest("invalid_question_set_id", "question_set_id cannot be empty", nil)
			}
			if _, err := s.questionSetRepository.FindByID(txCtx, input.QuestionSetID.Value); err != nil {
				return helper.BadRequest("invalid_question_set_id", "question set does not exist", err)
			}
		}
		if input.Type.Set {
			input.Type.Value = normalizeText(input.Type.Value)
			if err := validateRequired(input.Type.Value, "invalid_type", "type cannot be empty"); err != nil {
				return err
			}
		}
		if input.QuestionText.Set {
			input.QuestionText.Value = normalizeText(input.QuestionText.Value)
			if err := validateRequired(input.QuestionText.Value, "invalid_question_text", "question_text cannot be empty"); err != nil {
				return err
			}
		}
		if input.Difficulty.Set && input.Difficulty.Value <= 0 {
			return helper.BadRequest("invalid_difficulty", "difficulty must be greater than zero", nil)
		}
		if input.Choices.Set {
			for _, choice := range input.Choices.Value {
				if normalizeText(choice.ChoiceText) == "" {
					return helper.BadRequest("invalid_choice_text", "choice_text is required", nil)
				}
			}
		}
		if input.TagIDs.Set {
			if err := ensureTagIDsExist(s.tagRepository, txCtx, input.TagIDs.Value); err != nil {
				return err
			}
		}

		if _, err := s.questionRepository.Update(txCtx, input); err != nil {
			return err
		}
		if input.Choices.Set {
			if err := s.questionRepository.ReplaceChoices(txCtx, input.ID, input.Choices.Value, input.ActorID); err != nil {
				return err
			}
		}
		if input.TagIDs.Set {
			if err := s.questionRepository.ReplaceTags(txCtx, input.ID, input.TagIDs.Value, input.ActorID); err != nil {
				return err
			}
		}

		if !input.Choices.Set && !input.TagIDs.Set {
			updated = current
		}
		updated, err = s.questionRepository.FindByID(txCtx, input.ID)
		return err
	})
	return updated, err
}

func (s *questionService) Delete(ctx context.Context, id uuid.UUID, actorID string) error {
	return s.txManager.RunInTx(ctx, func(txCtx context.Context) error {
		if _, err := s.questionRepository.FindByID(txCtx, id); err != nil {
			return err
		}
		return s.questionRepository.Delete(txCtx, id, actorID)
	})
}

func (s *questionService) FindByID(ctx context.Context, id uuid.UUID) (domain.Question, error) {
	return s.questionRepository.FindByID(ctx, id)
}

func (s *questionService) FindAll(ctx context.Context, query domain.QuestionQuery) (domain.PageResult[domain.Question], error) {
	return s.questionRepository.FindAll(ctx, query)
}
