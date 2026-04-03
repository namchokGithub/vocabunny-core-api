package protocol

import (
	"github.com/labstack/echo/v4/middleware"
)

func RegisterHTTP(app *App) {
	e := app.Echo

	e.Use(middleware.CORS())
	registerHealthRoutes(e, app)

	api := e.Group("/api/v1")
	registerAppRoutes(api, app)
	registerBORoutes(api, app)
}
