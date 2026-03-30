package service

import (
	"github.com/namchokGithub/vocabunny-core-api/internal/core/port"
)

type Dependencies struct {
	Repositories *RepositoryPorts
	TxManager    port.TransactionManager
	Storage      port.FileStorage
}

type RepositoryPorts struct {
	User port.UserRepository
}

type Service struct {
	User port.UserService
}

func NewService(deps Dependencies) *Service {
	return &Service{
		User: NewUserService(UserServiceDependencies{
			UserRepository: deps.Repositories.User,
			TxManager:      deps.TxManager,
			Storage:        deps.Storage,
		}),
	}
}
