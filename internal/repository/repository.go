package repository

import (
	"gorm.io/gorm"

	"github.com/namchokGithub/vocabunny-core-api/internal/core/port"
)

type Dependencies struct {
	DB *gorm.DB
}

type Repository struct {
	User port.UserRepository
}

func NewRepository(deps Dependencies) *Repository {
	return &Repository{
		User: NewUserRepository(deps.DB),
	}
}
