package protocol

import (
	"github.com/labstack/echo/v4"
	"github.com/namchokGithub/vocabunny-core-api/internal/core/domain"
)

func registerAppRoutes(group *echo.Group, app *App) {
	appGroup := group.Group("/app")

	auth := appGroup.Group("/auth")
	auth.POST("/login/password", app.Handlers.AuthIdentity.LoginAppWithPassword)

	appProtected := appGroup.Group("")
	appProtected.Use(app.Middleware.Authenticate(), app.Middleware.RequireTokenScope(domain.TokenScopeApp))
	_ = appProtected
}
