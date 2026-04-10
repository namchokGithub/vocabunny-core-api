package middleware

import (
	"slices"

	"github.com/labstack/echo/v4"
	"github.com/namchokGithub/vocabunny-core-api/infrastructure"
	"github.com/namchokGithub/vocabunny-core-api/internal/core/domain"
	"github.com/namchokGithub/vocabunny-core-api/internal/core/helper"
	"github.com/namchokGithub/vocabunny-core-api/internal/core/port"
)

type Dependencies struct {
	JWTManager  *infrastructure.JWTManager
	UserService port.UserService
}

type Middleware struct {
	jwtManager  *infrastructure.JWTManager
	userService port.UserService
}

func New(deps Dependencies) *Middleware {
	return &Middleware{
		jwtManager:  deps.JWTManager,
		userService: deps.UserService,
	}
}

func (m *Middleware) Authenticate() echo.MiddlewareFunc {
	jwtMiddleware := m.jwtManager.Middleware()
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return jwtMiddleware(func(c echo.Context) error {
			subject, _ := c.Get("jwt_subject").(string)
			user, err := m.userService.FindBySubject(c.Request().Context(), subject)
			if err != nil {
				return helper.RespondError(c, helper.Unauthorized("invalid_subject", "user for token subject was not found", err))
			}
			helper.SetCurrentUser(c, user)
			return next(c)
		})
	}
}

func (m *Middleware) RequireTokenScope(scopes ...string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			scope, _ := c.Get("jwt_scope").(string)
			if scope == "" {
				return helper.RespondError(c, helper.Unauthorized("invalid_scope", "token scope is not available", nil))
			}
			if !slices.Contains(scopes, scope) {
				return helper.RespondError(c, helper.Forbidden("insufficient_scope", "token scope is not allowed for this route", nil))
			}
			return next(c)
		}
	}
}

func (m *Middleware) RequireRoles(roleNames ...domain.RoleName) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			user, ok := helper.CurrentUser(c)
			if !ok {
				return helper.RespondError(c, helper.Unauthorized("unauthorized", "current user is not available", nil))
			}

			for _, role := range user.Roles {
				if slices.Contains(roleNames, role.Name) {
					return next(c)
				}
			}

			return helper.RespondError(c, helper.Forbidden("insufficient_role", "current user does not have required role", nil))
		}
	}
}

func (m *Middleware) RequirePermissions(codes ...domain.PermissionCode) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			user, ok := helper.CurrentUser(c)
			if !ok {
				return helper.RespondError(c, helper.Unauthorized("unauthorized", "current user is not available", nil))
			}

			for _, role := range user.Roles {
				for _, permission := range role.Permissions {
					if slices.Contains(codes, permission.PermissionCode) {
						return next(c)
					}
				}
			}

			return helper.RespondError(c, helper.Forbidden("insufficient_permission", "current user does not have required permission", nil))
		}
	}
}
