package repository

import (
	"github.com/namchokGithub/vocabunny-core-api/internal/core/port"
	identityrepo "github.com/namchokGithub/vocabunny-core-api/internal/repository/identity"
	"gorm.io/gorm"
)

type Dependencies struct {
	DB *gorm.DB
}

type Repository struct {
	User         port.UserRepository
	Role         port.RoleRepository
	AuthIdentity port.AuthIdentityRepository
}

func NewRepository(deps Dependencies) *Repository {
	identityRepositories := identityrepo.NewRepository(identityrepo.Dependencies{
		DB: deps.DB,
	})

	return &Repository{
		User:         identityRepositories.User,
		Role:         identityRepositories.Role,
		AuthIdentity: identityRepositories.AuthIdentity,
	}
}
