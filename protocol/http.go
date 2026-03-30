package protocol

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
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
	registerUserRoutes(api, app)
}

func registerUserRoutes(group *echo.Group, app *App) {
	users := group.Group("/users")

	users.POST("", app.Handlers.User.Create)
	users.GET("", app.Handlers.User.FindAll)
	users.GET("/by-email", app.Handlers.User.FindByEmail)
	users.GET("/:id", app.Handlers.User.FindByID)
	users.PUT("/:id", app.Handlers.User.Update)
	users.DELETE("/:id", app.Handlers.User.Delete)
}
