package identity

import (
	"github.com/namchokGithub/vocabunny-core-api/internal/core/helper"
	"github.com/namchokGithub/vocabunny-core-api/internal/core/service"
)

type Dependencies struct {
	Services  *service.Service
	Validator *helper.RequestValidator
}

type Handler struct {
	User         *UserHandler
	Role         *RoleHandler
	AuthIdentity *AuthIdentityHandler
}

func NewHandler(deps Dependencies) *Handler {
	return &Handler{
		User:         NewUserHandler(deps.Services.User),
		Role:         NewRoleHandler(deps.Services.Role),
		AuthIdentity: NewAuthIdentityHandler(deps.Services.AuthIdentity),
	}
}
