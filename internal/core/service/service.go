package service

import (
	"github.com/namchokGithub/vocabunny-core-api/internal/core/port"
	identityservice "github.com/namchokGithub/vocabunny-core-api/internal/core/service/identity"
)

type Dependencies struct {
	Repositories *RepositoryPorts
	TxManager    port.TransactionManager
	Storage      port.FileStorage
	TokenManager port.TokenManager
}

type RepositoryPorts struct {
	User         port.UserRepository
	Role         port.RoleRepository
	AuthIdentity port.AuthIdentityRepository
}

type Service struct {
	User         port.UserService
	Role         port.RoleService
	AuthIdentity port.AuthIdentityService
}

func NewService(deps Dependencies) *Service {
	identityServices := identityservice.NewService(identityservice.Dependencies{
		UserRepository:         deps.Repositories.User,
		RoleRepository:         deps.Repositories.Role,
		AuthIdentityRepository: deps.Repositories.AuthIdentity,
		TxManager:              deps.TxManager,
		TokenManager:           deps.TokenManager,
	})

	return &Service{
		User:         identityServices.User,
		Role:         identityServices.Role,
		AuthIdentity: identityServices.AuthIdentity,
	}
}
