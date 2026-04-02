package helper

import (
	"context"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/namchokGithub/vocabunny-core-api/internal/core/domain"
)

type txContextKey struct{}
type currentUserContextKey struct{}

func ContextWithTx(ctx context.Context, tx interface{}) context.Context {
	return context.WithValue(ctx, txContextKey{}, tx)
}

func TxFromContext(ctx context.Context) interface{} {
	return ctx.Value(txContextKey{})
}

func ContextWithCurrentUser(ctx context.Context, user domain.User) context.Context {
	return context.WithValue(ctx, currentUserContextKey{}, user)
}

func CurrentUserFromContext(ctx context.Context) (domain.User, bool) {
	user, ok := ctx.Value(currentUserContextKey{}).(domain.User)
	return user, ok
}

func SetCurrentUser(c echo.Context, user domain.User) {
	c.Set("current_user", user)
	c.SetRequest(c.Request().WithContext(ContextWithCurrentUser(c.Request().Context(), user)))
}

func CurrentUser(c echo.Context) (domain.User, bool) {
	if user, ok := c.Get("current_user").(domain.User); ok {
		return user, true
	}

	return CurrentUserFromContext(c.Request().Context())
}

func ActorIDFromContext(c echo.Context) string {
	if actorID := c.Request().Header.Get("X-Actor-ID"); actorID != "" {
		return actorID
	}

	if user, ok := CurrentUser(c); ok && user.ID != uuid.Nil {
		return user.ID.String()
	}

	if subject, ok := c.Get("jwt_subject").(string); ok && subject != "" {
		return subject
	}

	return "system"
}
