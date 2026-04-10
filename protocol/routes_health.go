package protocol

import (
	"net/http"
	"os"
	"time"

	"github.com/labstack/echo/v4"
)

func registerHealthRoutes(e *echo.Echo, app *App) {
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
}
