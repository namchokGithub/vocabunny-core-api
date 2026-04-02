package protocol

import (
	"net/http"
	"os"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/namchokGithub/vocabunny-core-api/internal/core/domain"
)

func RegisterHTTP(app *App) {
	e := app.Echo

	e.Use(middleware.CORS())

	e.GET("/health", func(c echo.Context) error {
		now := time.Now()
		hostname, err := os.Hostname()
		if err != nil {
			hostname = "unknown"
		}

		return c.JSON(http.StatusOK, map[string]any{
			"status":      "ok",
			"service":     app.Config.App.Name,
			"environment": app.Config.App.Env,
			"server_time": now.Format(time.RFC3339),
			"timezone":    now.Location().String(),
			"utc_offset":  now.Format("-07:00"),
			"server_host": hostname,
		})
	})

	api := e.Group("/api/v1")
	registerAppRoutes(api, app)
	registerBORoutes(api, app)
}

func registerAppRoutes(group *echo.Group, app *App) {
	appGroup := group.Group("/app")

	auth := appGroup.Group("/auth")
	auth.POST("/login/password", app.Handlers.AuthIdentity.LoginAppWithPassword)

	appProtected := appGroup.Group("")
	appProtected.Use(app.Middleware.Authenticate(), app.Middleware.RequireTokenScope(domain.TokenScopeApp))
	_ = appProtected
}

func registerBORoutes(group *echo.Group, app *App) {
	boGroup := group.Group("/bo")

	auth := boGroup.Group("/auth")
	auth.POST("/login/password", app.Handlers.AuthIdentity.LoginBOWithPassword)

	authenticated := boGroup.Group("")
	authenticated.Use(app.Middleware.Authenticate(), app.Middleware.RequireTokenScope(domain.TokenScopeBO))

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
