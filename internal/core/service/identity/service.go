package identity

import "github.com/namchokGithub/vocabunny-core-api/internal/core/port"

type Dependencies struct {
	UserRepository         port.UserRepository
	RoleRepository         port.RoleRepository
	AuthIdentityRepository port.AuthIdentityRepository
	TxManager              port.TransactionManager
	TokenManager           port.TokenManager
}

type Service struct {
	User         port.UserService
	Role         port.RoleService
	AuthIdentity port.AuthIdentityService
}

func NewService(deps Dependencies) *Service {
	return &Service{
		User: NewUserService(UserServiceDependencies{
			UserRepository: deps.UserRepository,
			RoleRepository: deps.RoleRepository,
			TxManager:      deps.TxManager,
		}),
		Role: NewRoleService(RoleServiceDependencies{
			RoleRepository: deps.RoleRepository,
			TxManager:      deps.TxManager,
		}),
		AuthIdentity: NewAuthIdentityService(AuthIdentityServiceDependencies{
			AuthIdentityRepository: deps.AuthIdentityRepository,
			UserRepository:         deps.UserRepository,
			TxManager:              deps.TxManager,
			TokenManager:           deps.TokenManager,
		}),
	}
}
