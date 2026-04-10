package content

import "github.com/namchokGithub/vocabunny-core-api/internal/core/port"

type Dependencies struct {
	SectionRepository        port.SectionRepository
	LessonRepository         port.LessonRepository
	UnitRepository           port.UnitRepository
	QuestionSetRepository    port.QuestionSetRepository
	QuestionRepository       port.QuestionRepository
	QuestionChoiceRepository port.QuestionChoiceRepository
	TagRepository            port.TagRepository
	MediaAssetRepository     port.MediaAssetRepository
	TxManager                port.TransactionManager
}

type Service struct {
	Section        port.SectionService
	Lesson         port.LessonService
	Unit           port.UnitService
	QuestionSet    port.QuestionSetService
	Question       port.QuestionService
	QuestionChoice port.QuestionChoiceService
	Tag            port.TagService
	MediaAsset     port.MediaAssetService
}

func NewService(deps Dependencies) *Service {
	return &Service{
		Section: NewSectionService(SectionServiceDependencies{
			SectionRepository: deps.SectionRepository,
			TxManager:         deps.TxManager,
		}),
		Lesson: NewLessonService(LessonServiceDependencies{
			LessonRepository:  deps.LessonRepository,
			SectionRepository: deps.SectionRepository,
			TxManager:         deps.TxManager,
		}),
		Unit: NewUnitService(UnitServiceDependencies{
			UnitRepository:   deps.UnitRepository,
			LessonRepository: deps.LessonRepository,
			TxManager:        deps.TxManager,
		}),
		QuestionSet: NewQuestionSetService(QuestionSetServiceDependencies{
			QuestionSetRepository: deps.QuestionSetRepository,
			UnitRepository:        deps.UnitRepository,
			TxManager:             deps.TxManager,
		}),
		Question: NewQuestionService(QuestionServiceDependencies{
			QuestionRepository:    deps.QuestionRepository,
			QuestionSetRepository: deps.QuestionSetRepository,
			TagRepository:         deps.TagRepository,
			TxManager:             deps.TxManager,
		}),
		QuestionChoice: NewQuestionChoiceService(QuestionChoiceServiceDependencies{
			QuestionChoiceRepository: deps.QuestionChoiceRepository,
			QuestionRepository:       deps.QuestionRepository,
			TxManager:                deps.TxManager,
		}),
		Tag: NewTagService(TagServiceDependencies{
			TagRepository: deps.TagRepository,
			TxManager:     deps.TxManager,
		}),
		MediaAsset: NewMediaAssetService(MediaAssetServiceDependencies{
			MediaAssetRepository: deps.MediaAssetRepository,
			TxManager:            deps.TxManager,
		}),
	}
}
