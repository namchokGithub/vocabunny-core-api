package handler

import (
	"github.com/namchokGithub/vocabunny-core-api/internal/core/helper"
	"github.com/namchokGithub/vocabunny-core-api/internal/core/service"
	contenthandler "github.com/namchokGithub/vocabunny-core-api/internal/handler/content"
	identityhandler "github.com/namchokGithub/vocabunny-core-api/internal/handler/identity"
)

type Dependencies struct {
	Services  *service.Service
	Validator *helper.RequestValidator
}

type Handler struct {
	User           *identityhandler.UserHandler
	Role           *identityhandler.RoleHandler
	AuthIdentity   *identityhandler.AuthIdentityHandler
	Section        *contenthandler.SectionHandler
	Lesson         *contenthandler.LessonHandler
	Unit           *contenthandler.UnitHandler
	QuestionSet    *contenthandler.QuestionSetHandler
	Question       *contenthandler.QuestionHandler
	QuestionChoice *contenthandler.QuestionChoiceHandler
	Tag            *contenthandler.TagHandler
}

func NewHandler(deps Dependencies) *Handler {
	identityHandlers := identityhandler.NewHandler(identityhandler.Dependencies{
		Services:  deps.Services,
		Validator: deps.Validator,
	})
	contentHandlers := contenthandler.NewHandler(contenthandler.Dependencies{
		Services:  deps.Services,
		Validator: deps.Validator,
	})

	return &Handler{
		User:           identityHandlers.User,
		Role:           identityHandlers.Role,
		AuthIdentity:   identityHandlers.AuthIdentity,
		Section:        contentHandlers.Section,
		Lesson:         contentHandlers.Lesson,
		Unit:           contentHandlers.Unit,
		QuestionSet:    contentHandlers.QuestionSet,
		Question:       contentHandlers.Question,
		QuestionChoice: contentHandlers.QuestionChoice,
		Tag:            contentHandlers.Tag,
	}
}
