package content

import (
	"github.com/namchokGithub/vocabunny-core-api/internal/core/helper"
	"github.com/namchokGithub/vocabunny-core-api/internal/core/service"
)

type Dependencies struct {
	Services  *service.Service
	Validator *helper.RequestValidator
}

type Handler struct {
	Section        *SectionHandler
	Lesson         *LessonHandler
	Unit           *UnitHandler
	QuestionSet    *QuestionSetHandler
	Question       *QuestionHandler
	QuestionChoice *QuestionChoiceHandler
	Tag            *TagHandler
	MediaAsset     *MediaAssetHandler
}

func NewHandler(deps Dependencies) *Handler {
	return &Handler{
		Section:        NewSectionHandler(deps.Services.Section),
		Lesson:         NewLessonHandler(deps.Services.Lesson),
		Unit:           NewUnitHandler(deps.Services.Unit),
		QuestionSet:    NewQuestionSetHandler(deps.Services.QuestionSet),
		Question:       NewQuestionHandler(deps.Services.Question),
		QuestionChoice: NewQuestionChoiceHandler(deps.Services.QuestionChoice),
		Tag:            NewTagHandler(deps.Services.Tag),
		MediaAsset:     NewMediaAssetHandler(deps.Services.MediaAsset),
	}
}
