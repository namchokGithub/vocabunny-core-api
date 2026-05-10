package content

import (
	"context"

	"github.com/namchokGithub/vocabunny-core-api/internal/core/domain"
	"github.com/namchokGithub/vocabunny-core-api/internal/core/port"
)

type ContentOrderServiceDependencies struct {
	SectionRepository     port.SectionRepository
	LessonRepository      port.LessonRepository
	UnitRepository        port.UnitRepository
	QuestionSetRepository port.QuestionSetRepository
	QuestionRepository    port.QuestionRepository
}

type contentOrderService struct {
	sectionRepository     port.SectionRepository
	lessonRepository      port.LessonRepository
	unitRepository        port.UnitRepository
	questionSetRepository port.QuestionSetRepository
	questionRepository    port.QuestionRepository
}

func NewContentOrderService(deps ContentOrderServiceDependencies) port.ContentOrderService {
	return &contentOrderService{
		sectionRepository:     deps.SectionRepository,
		lessonRepository:      deps.LessonRepository,
		unitRepository:        deps.UnitRepository,
		questionSetRepository: deps.QuestionSetRepository,
		questionRepository:    deps.QuestionRepository,
	}
}

func (s *contentOrderService) GetLastOrderNos(ctx context.Context) (domain.ContentOrderNoSummary, error) {
	sections, err := s.sectionRepository.FindLastOrderNo(ctx)
	if err != nil {
		return domain.ContentOrderNoSummary{}, err
	}

	lessons, err := s.lessonRepository.FindLastOrderNo(ctx)
	if err != nil {
		return domain.ContentOrderNoSummary{}, err
	}

	units, err := s.unitRepository.FindLastOrderNo(ctx)
	if err != nil {
		return domain.ContentOrderNoSummary{}, err
	}

	questionSets, err := s.questionSetRepository.FindLastOrderNo(ctx)
	if err != nil {
		return domain.ContentOrderNoSummary{}, err
	}

	questions, err := s.questionRepository.FindLastOrderNo(ctx)
	if err != nil {
		return domain.ContentOrderNoSummary{}, err
	}

	return domain.ContentOrderNoSummary{
		Sections:     sections,
		Lessons:      lessons,
		Units:        units,
		QuestionSets: questionSets,
		Questions:    questions,
	}, nil
}
