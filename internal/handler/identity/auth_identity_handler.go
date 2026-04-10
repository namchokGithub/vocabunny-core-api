package identity

import (
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/namchokGithub/vocabunny-core-api/internal/core/domain"
	"github.com/namchokGithub/vocabunny-core-api/internal/core/helper"
	"github.com/namchokGithub/vocabunny-core-api/internal/core/port"
)

type AuthIdentityHandler struct {
	authIdentityService port.AuthIdentityService
}

func NewAuthIdentityHandler(authIdentityService port.AuthIdentityService) *AuthIdentityHandler {
	return &AuthIdentityHandler{authIdentityService: authIdentityService}
}

func (h *AuthIdentityHandler) Create(c echo.Context) error {
	var req CreateAuthIdentityRequest
	if err := helper.BindAndValidate(c, &req); err != nil {
		return helper.RespondError(c, err)
	}
	input, err := toCreateAuthIdentityDomain(req, helper.ActorIDFromContext(c))
	if err != nil {
		return helper.RespondError(c, err)
	}
	identity, err := h.authIdentityService.Create(c.Request().Context(), input)
	if err != nil {
		return helper.RespondError(c, err)
	}
	return helper.RespondSuccess(c, http.StatusCreated, toAuthIdentityResponse(identity))
}

func (h *AuthIdentityHandler) Update(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return helper.RespondError(c, helper.BadRequest("invalid_auth_identity_id", "auth identity id must be a valid uuid", err))
	}
	var req UpdateAuthIdentityRequest
	if err := helper.BindAndValidate(c, &req); err != nil {
		return helper.RespondError(c, err)
	}
	identity, err := h.authIdentityService.Update(c.Request().Context(), toUpdateAuthIdentityDomain(id, req, helper.ActorIDFromContext(c)))
	if err != nil {
		return helper.RespondError(c, err)
	}
	return helper.RespondSuccess(c, http.StatusOK, toAuthIdentityResponse(identity))
}

func (h *AuthIdentityHandler) Delete(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return helper.RespondError(c, helper.BadRequest("invalid_auth_identity_id", "auth identity id must be a valid uuid", err))
	}
	if err := h.authIdentityService.Delete(c.Request().Context(), id, helper.ActorIDFromContext(c)); err != nil {
		return helper.RespondError(c, err)
	}
	return helper.RespondSuccess(c, http.StatusOK, map[string]string{"id": id.String(), "status": "deleted"})
}

func (h *AuthIdentityHandler) FindByID(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return helper.RespondError(c, helper.BadRequest("invalid_auth_identity_id", "auth identity id must be a valid uuid", err))
	}
	identity, err := h.authIdentityService.FindByID(c.Request().Context(), id)
	if err != nil {
		return helper.RespondError(c, err)
	}
	return helper.RespondSuccess(c, http.StatusOK, toAuthIdentityResponse(identity))
}

func (h *AuthIdentityHandler) FindAll(c echo.Context) error {
	query, err := buildAuthIdentityQuery(c)
	if err != nil {
		return helper.RespondError(c, err)
	}
	result, err := h.authIdentityService.FindAll(c.Request().Context(), query)
	if err != nil {
		return helper.RespondError(c, err)
	}
	return helper.RespondSuccess(c, http.StatusOK, toAuthIdentitiesListResponse(result, query))
}

func (h *AuthIdentityHandler) LoginWithPassword(c echo.Context) error {
	return h.loginWithPassword(c, "")
}

func (h *AuthIdentityHandler) LoginAppWithPassword(c echo.Context) error {
	return h.loginWithPassword(c, domain.TokenScopeApp)
}

func (h *AuthIdentityHandler) LoginBOWithPassword(c echo.Context) error {
	return h.loginWithPassword(c, domain.TokenScopeBO)
}

func (h *AuthIdentityHandler) loginWithPassword(c echo.Context, scope string) error {
	var req PasswordLoginRequest
	if err := helper.BindAndValidate(c, &req); err != nil {
		return helper.RespondError(c, err)
	}

	if scope == "" {
		scope = strings.TrimSpace(c.QueryParam("scope"))
	}
	if scope == "" {
		scope = domain.TokenScopeApp
	}

	token, err := h.authIdentityService.LoginWithPassword(c.Request().Context(), domain.PasswordLoginInput{
		EmailOrUsername: req.EmailOrUsername,
		Password:        req.Password,
		Scope:           scope,
	})
	if err != nil {
		return helper.RespondError(c, err)
	}
	return helper.RespondSuccess(c, http.StatusOK, toLoginResponse(token))
}

func buildAuthIdentityQuery(c echo.Context) (domain.AuthIdentityQuery, error) {
	query := domain.AuthIdentityQuery{Paging: helper.BuildPaging(c)}
	query.SortBy, query.SortOrder = helper.BuildSort(c)

	if userID := strings.TrimSpace(c.QueryParam("user_id")); userID != "" {
		parsed, err := uuid.Parse(userID)
		if err != nil {
			return domain.AuthIdentityQuery{}, helper.BadRequest("invalid_user_id", "user_id must be a valid uuid", err)
		}
		query.UserID = &parsed
	}

	if provider := strings.TrimSpace(c.QueryParam("provider")); provider != "" {
		value := domain.AuthProvider(provider)
		query.Provider = &value
	}

	return query, nil
}
