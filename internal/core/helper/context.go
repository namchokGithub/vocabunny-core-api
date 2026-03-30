package helper

import (
	"context"

	"github.com/labstack/echo/v4"
)

type txContextKey struct{}

func ContextWithTx(ctx context.Context, tx interface{}) context.Context {
	return context.WithValue(ctx, txContextKey{}, tx)
}

func TxFromContext(ctx context.Context) interface{} {
	return ctx.Value(txContextKey{})
}

func ActorIDFromContext(c echo.Context) string {
	if actorID := c.Request().Header.Get("X-Actor-ID"); actorID != "" {
		return actorID
	}

	if subject, ok := c.Get("jwt_subject").(string); ok && subject != "" {
		return subject
	}

	return "system"
}
