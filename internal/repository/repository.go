package repository

import (
	"github.com/namchokGithub/vocabunny-core-api/internal/core/port"
	contentrepo "github.com/namchokGithub/vocabunny-core-api/internal/repository/content"
	identityrepo "github.com/namchokGithub/vocabunny-core-api/internal/repository/identity"
	mediarepo "github.com/namchokGithub/vocabunny-core-api/internal/repository/media"
	"gorm.io/gorm"
)

type Dependencies struct {
	DB *gorm.DB
}

type Repository struct {
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

func NewRepository(deps Dependencies) *Repository {
	identityRepositories := identityrepo.NewRepository(identityrepo.Dependencies{
		DB: deps.DB,
	})
	contentRepositories := contentrepo.NewRepository(contentrepo.Dependencies{
		DB: deps.DB,
	})
	mediaRepositories := mediarepo.NewRepository(mediarepo.Dependencies{
		DB: deps.DB,
	})

	return &Repository{
		User:           identityRepositories.User,
		Role:           identityRepositories.Role,
		AuthIdentity:   identityRepositories.AuthIdentity,
		Section:        contentRepositories.Section,
		Lesson:         contentRepositories.Lesson,
		Unit:           contentRepositories.Unit,
		QuestionSet:    contentRepositories.QuestionSet,
		Question:       contentRepositories.Question,
		QuestionChoice: contentRepositories.QuestionChoice,
		Tag:            contentRepositories.Tag,
		MediaAsset:     mediaRepositories.MediaAsset,
	}
}
