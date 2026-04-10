package service

import (
	"github.com/namchokGithub/vocabunny-core-api/internal/core/port"
	contentservice "github.com/namchokGithub/vocabunny-core-api/internal/core/service/content"
	identityservice "github.com/namchokGithub/vocabunny-core-api/internal/core/service/identity"
)

type Dependencies struct {
	Repositories *RepositoryPorts
	TxManager    port.TransactionManager
	Storage      port.FileStorage
	TokenManager port.TokenManager
}

type RepositoryPorts struct {
	User           port.UserRepository
	Role           port.RoleRepository
	AuthIdentity   port.AuthIdentityRepository
	Section        port.SectionRepository
	Lesson         port.LessonRepository
	Unit           port.UnitRepository
	QuestionSet    port.QuestionSetRepository
	Question       port.QuestionRepository
	QuestionChoice port.QuestionChoiceRepository
	Tag            port.TagRepository
	MediaAsset     port.MediaAssetRepository
}

type Service struct {
	User           port.UserService
	Role           port.RoleService
	AuthIdentity   port.AuthIdentityService
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
	identityServices := identityservice.NewService(identityservice.Dependencies{
		UserRepository:         deps.Repositories.User,
		RoleRepository:         deps.Repositories.Role,
		AuthIdentityRepository: deps.Repositories.AuthIdentity,
		TxManager:              deps.TxManager,
		TokenManager:           deps.TokenManager,
	})
	contentServices := contentservice.NewService(contentservice.Dependencies{
		SectionRepository:        deps.Repositories.Section,
		LessonRepository:         deps.Repositories.Lesson,
		UnitRepository:           deps.Repositories.Unit,
		QuestionSetRepository:    deps.Repositories.QuestionSet,
		QuestionRepository:       deps.Repositories.Question,
		QuestionChoiceRepository: deps.Repositories.QuestionChoice,
		TagRepository:            deps.Repositories.Tag,
		MediaAssetRepository:     deps.Repositories.MediaAsset,
		TxManager:                deps.TxManager,
	})

	return &Service{
		User:           identityServices.User,
		Role:           identityServices.Role,
		AuthIdentity:   identityServices.AuthIdentity,
		Section:        contentServices.Section,
		Lesson:         contentServices.Lesson,
		Unit:           contentServices.Unit,
		QuestionSet:    contentServices.QuestionSet,
		Question:       contentServices.Question,
		QuestionChoice: contentServices.QuestionChoice,
		Tag:            contentServices.Tag,
		MediaAsset:     contentServices.MediaAsset,
	}
}
