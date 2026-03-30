package handler

import (
	"github.com/namchokGithub/vocabunny-core-api/internal/core/helper"
	"github.com/namchokGithub/vocabunny-core-api/internal/core/service"
)

type Dependencies struct {
	Services  *service.Service
	Validator *helper.RequestValidator
}

type Handler struct {
	User *UserHandler
}

func NewHandler(deps Dependencies) *Handler {
	return &Handler{
		User: NewUserHandler(deps.Services.User),
	}
}
