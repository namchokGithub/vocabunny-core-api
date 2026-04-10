package protocol

import (
	"github.com/labstack/echo/v4"
	"github.com/namchokGithub/vocabunny-core-api/internal/core/domain"
)

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

	content := authenticated.Group("/content")

	sections := content.Group("/sections")
	sections.POST("", app.Handlers.Section.Create, app.Middleware.RequirePermissions(domain.PermissionContentWrite))
	sections.GET("", app.Handlers.Section.FindAll, app.Middleware.RequirePermissions(domain.PermissionContentRead))
	sections.GET("/:id", app.Handlers.Section.FindByID, app.Middleware.RequirePermissions(domain.PermissionContentRead))
	sections.PUT("/:id", app.Handlers.Section.Update, app.Middleware.RequirePermissions(domain.PermissionContentWrite))
	sections.DELETE("/:id", app.Handlers.Section.Delete, app.Middleware.RequirePermissions(domain.PermissionContentWrite))

	lessons := content.Group("/lessons")
	lessons.POST("", app.Handlers.Lesson.Create, app.Middleware.RequirePermissions(domain.PermissionContentWrite))
	lessons.GET("", app.Handlers.Lesson.FindAll, app.Middleware.RequirePermissions(domain.PermissionContentRead))
	lessons.GET("/:id", app.Handlers.Lesson.FindByID, app.Middleware.RequirePermissions(domain.PermissionContentRead))
	lessons.PUT("/:id", app.Handlers.Lesson.Update, app.Middleware.RequirePermissions(domain.PermissionContentWrite))
	lessons.DELETE("/:id", app.Handlers.Lesson.Delete, app.Middleware.RequirePermissions(domain.PermissionContentWrite))

	units := content.Group("/units")
	units.POST("", app.Handlers.Unit.Create, app.Middleware.RequirePermissions(domain.PermissionContentWrite))
	units.GET("", app.Handlers.Unit.FindAll, app.Middleware.RequirePermissions(domain.PermissionContentRead))
	units.GET("/:id", app.Handlers.Unit.FindByID, app.Middleware.RequirePermissions(domain.PermissionContentRead))
	units.PUT("/:id", app.Handlers.Unit.Update, app.Middleware.RequirePermissions(domain.PermissionContentWrite))
	units.DELETE("/:id", app.Handlers.Unit.Delete, app.Middleware.RequirePermissions(domain.PermissionContentWrite))

	questionSets := content.Group("/question-sets")
	questionSets.POST("", app.Handlers.QuestionSet.Create, app.Middleware.RequirePermissions(domain.PermissionContentWrite))
	questionSets.GET("", app.Handlers.QuestionSet.FindAll, app.Middleware.RequirePermissions(domain.PermissionContentRead))
	questionSets.GET("/:id", app.Handlers.QuestionSet.FindByID, app.Middleware.RequirePermissions(domain.PermissionContentRead))
	questionSets.PUT("/:id", app.Handlers.QuestionSet.Update, app.Middleware.RequirePermissions(domain.PermissionContentWrite))
	questionSets.DELETE("/:id", app.Handlers.QuestionSet.Delete, app.Middleware.RequirePermissions(domain.PermissionContentWrite))

	questions := content.Group("/questions")
	questions.POST("", app.Handlers.Question.Create, app.Middleware.RequirePermissions(domain.PermissionContentWrite))
	questions.GET("", app.Handlers.Question.FindAll, app.Middleware.RequirePermissions(domain.PermissionContentRead))
	questions.GET("/:id", app.Handlers.Question.FindByID, app.Middleware.RequirePermissions(domain.PermissionContentRead))
	questions.PUT("/:id", app.Handlers.Question.Update, app.Middleware.RequirePermissions(domain.PermissionContentWrite))
	questions.DELETE("/:id", app.Handlers.Question.Delete, app.Middleware.RequirePermissions(domain.PermissionContentWrite))

	questionChoices := content.Group("/question-choices")
	questionChoices.POST("", app.Handlers.QuestionChoice.Create, app.Middleware.RequirePermissions(domain.PermissionContentWrite))
	questionChoices.GET("", app.Handlers.QuestionChoice.FindAll, app.Middleware.RequirePermissions(domain.PermissionContentRead))
	questionChoices.GET("/:id", app.Handlers.QuestionChoice.FindByID, app.Middleware.RequirePermissions(domain.PermissionContentRead))
	questionChoices.PUT("/:id", app.Handlers.QuestionChoice.Update, app.Middleware.RequirePermissions(domain.PermissionContentWrite))
	questionChoices.DELETE("/:id", app.Handlers.QuestionChoice.Delete, app.Middleware.RequirePermissions(domain.PermissionContentWrite))

	tags := content.Group("/tags")
	tags.POST("", app.Handlers.Tag.Create, app.Middleware.RequirePermissions(domain.PermissionContentWrite))
	tags.GET("", app.Handlers.Tag.FindAll, app.Middleware.RequirePermissions(domain.PermissionContentRead))
	tags.GET("/:id", app.Handlers.Tag.FindByID, app.Middleware.RequirePermissions(domain.PermissionContentRead))
	tags.PUT("/:id", app.Handlers.Tag.Update, app.Middleware.RequirePermissions(domain.PermissionContentWrite))
	tags.DELETE("/:id", app.Handlers.Tag.Delete, app.Middleware.RequirePermissions(domain.PermissionContentWrite))

	mediaAssets := content.Group("/media-assets")
	mediaAssets.POST("", app.Handlers.MediaAsset.Create, app.Middleware.RequirePermissions(domain.PermissionContentWrite))
	mediaAssets.GET("", app.Handlers.MediaAsset.FindAll, app.Middleware.RequirePermissions(domain.PermissionContentRead))
	mediaAssets.GET("/:id", app.Handlers.MediaAsset.FindByID, app.Middleware.RequirePermissions(domain.PermissionContentRead))
	mediaAssets.PUT("/:id", app.Handlers.MediaAsset.Update, app.Middleware.RequirePermissions(domain.PermissionContentWrite))
	mediaAssets.DELETE("/:id", app.Handlers.MediaAsset.Delete, app.Middleware.RequirePermissions(domain.PermissionContentWrite))
}
