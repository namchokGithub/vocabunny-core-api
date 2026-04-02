package protocol

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/namchokGithub/vocabunny-core-api/internal/core/domain"
)

func RegisterHTTP(app *App) {
	e := app.Echo

	e.Use(middleware.CORS())

	e.GET("/health", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{
			"status":  "ok",
			"service": app.Config.App.Name,
		})
	})

	api := e.Group("/api/v1")
	registerAuthRoutes(api, app)
	registerIdentityRoutes(api, app)
}

func registerAuthRoutes(group *echo.Group, app *App) {
	auth := group.Group("/auth")
	auth.POST("/login/password", app.Handlers.AuthIdentity.LoginWithPassword)
}

func registerIdentityRoutes(group *echo.Group, app *App) {
	authenticated := group.Group("")
	authenticated.Use(app.Middleware.Authenticate())

	users := authenticated.Group("/users")
	users.GET("", app.Handlers.User.FindAll, app.Middleware.RequirePermissions(domain.PermissionUserRead))
	users.POST("", app.Handlers.User.Create, app.Middleware.RequireRoles(domain.RoleNameAdmin, domain.RoleNameContentAdmin))
	users.GET("/by-email", app.Handlers.User.FindByEmail, app.Middleware.RequirePermissions(domain.PermissionUserRead))
	users.GET("/by-username", app.Handlers.User.FindByUsername, app.Middleware.RequirePermissions(domain.PermissionUserRead))
	users.GET("/:id", app.Handlers.User.FindByID, app.Middleware.RequirePermissions(domain.PermissionUserRead))
	users.PUT("/:id", app.Handlers.User.Update, app.Middleware.RequireRoles(domain.RoleNameAdmin, domain.RoleNameContentAdmin))
	users.DELETE("/:id", app.Handlers.User.Delete, app.Middleware.RequirePermissions(domain.PermissionUserBan))

	roles := authenticated.Group("/roles", app.Middleware.RequireRoles(domain.RoleNameAdmin))
	roles.POST("", app.Handlers.Role.Create)
	roles.GET("", app.Handlers.Role.FindAll)
	roles.GET("/:id", app.Handlers.Role.FindByID)
	roles.PUT("/:id", app.Handlers.Role.Update)
	roles.DELETE("/:id", app.Handlers.Role.Delete)

	authIdentities := authenticated.Group("/auth-identities", app.Middleware.RequireRoles(domain.RoleNameAdmin))
	authIdentities.POST("", app.Handlers.AuthIdentity.Create)
	authIdentities.GET("", app.Handlers.AuthIdentity.FindAll)
	authIdentities.GET("/:id", app.Handlers.AuthIdentity.FindByID)
	authIdentities.PUT("/:id", app.Handlers.AuthIdentity.Update)
	authIdentities.DELETE("/:id", app.Handlers.AuthIdentity.Delete)
}
