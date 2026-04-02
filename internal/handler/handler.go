package handler

import (
	"github.com/namchokGithub/vocabunny-core-api/internal/core/helper"
	"github.com/namchokGithub/vocabunny-core-api/internal/core/service"
	identityhandler "github.com/namchokGithub/vocabunny-core-api/internal/handler/identity"
)

type Dependencies struct {
	Services  *service.Service
	Validator *helper.RequestValidator
}

type Handler struct {
	User         *identityhandler.UserHandler
	Role         *identityhandler.RoleHandler
	AuthIdentity *identityhandler.AuthIdentityHandler
}

func NewHandler(deps Dependencies) *Handler {
	identityHandlers := identityhandler.NewHandler(identityhandler.Dependencies{
		Services:  deps.Services,
		Validator: deps.Validator,
	})

	return &Handler{
		User:         identityHandlers.User,
		Role:         identityHandlers.Role,
		AuthIdentity: identityHandlers.AuthIdentity,
	}
}
