package identity

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/namchokGithub/vocabunny-core-api/internal/core/domain"
)

func TestAuthIdentityHandlerLoginWithPasswordDefaultsToAppScope(t *testing.T) {
	t.Parallel()

	userID := uuid.New()
	service := &authIdentityServiceStub{
		loginWithPasswordFn: func(ctx context.Context, input domain.PasswordLoginInput) (domain.AuthToken, error) {
			if input.Scope != domain.TokenScopeApp {
				t.Fatalf("expected scope %q, got %q", domain.TokenScopeApp, input.Scope)
			}
			if input.EmailOrUsername != "bunny@example.com" {
				t.Fatalf("expected login bunny@example.com, got %q", input.EmailOrUsername)
			}
			return domain.AuthToken{
				AccessToken:      "access-token",
				RefreshToken:     "refresh-token",
				TokenType:        domain.TokenTypeBearer,
				ExpiresIn:        3600,
				RefreshExpiresIn: 7200,
				User: domain.User{
					ID:          userID,
					DisplayName: "Bunny",
					Status:      domain.UserStatusActive,
					AuditFields: domain.AuditFields{
						CreatedAt: time.Date(2026, 4, 9, 10, 0, 0, 0, time.UTC),
						UpdatedAt: time.Date(2026, 4, 9, 10, 0, 0, 0, time.UTC),
					},
				},
			}, nil
		},
	}

	handler := NewAuthIdentityHandler(service)
	rec := performJSONRequest(t, http.MethodPost, "/auth/login/password", `{"email_or_username":"bunny@example.com","password":"secret"}`, func(c echo.Context) error {
		return handler.LoginWithPassword(c)
	})

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rec.Code)
	}

	var resp struct {
		Success bool          `json:"success"`
		Data    LoginResponse `json:"data"`
	}
	decodeResponse(t, rec, &resp)

	if !resp.Success {
		t.Fatal("expected success response")
	}
	if resp.Data.AccessToken != "access-token" {
		t.Fatalf("expected access token access-token, got %q", resp.Data.AccessToken)
	}
}

func TestAuthIdentityHandlerLoginBOWithPasswordUsesBOScope(t *testing.T) {
	t.Parallel()

	service := &authIdentityServiceStub{
		loginWithPasswordFn: func(ctx context.Context, input domain.PasswordLoginInput) (domain.AuthToken, error) {
			if input.Scope != domain.TokenScopeBO {
				t.Fatalf("expected scope %q, got %q", domain.TokenScopeBO, input.Scope)
			}
			return domain.AuthToken{
				AccessToken: "bo-access",
				TokenType:   domain.TokenTypeBearer,
				User: domain.User{
					ID: uuid.New(),
					AuditFields: domain.AuditFields{
						CreatedAt: time.Date(2026, 4, 9, 10, 0, 0, 0, time.UTC),
						UpdatedAt: time.Date(2026, 4, 9, 10, 0, 0, 0, time.UTC),
					},
				},
			}, nil
		},
	}

	handler := NewAuthIdentityHandler(service)
	rec := performJSONRequest(t, http.MethodPost, "/bo/auth/login/password?scope=app", `{"email_or_username":"admin","password":"secret"}`, func(c echo.Context) error {
		return handler.LoginBOWithPassword(c)
	})

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rec.Code)
	}
}

func TestAuthIdentityHandlerFindAllRejectsInvalidUserID(t *testing.T) {
	t.Parallel()

	called := false
	service := &authIdentityServiceStub{
		findAllFn: func(ctx context.Context, query domain.AuthIdentityQuery) (domain.PageResult[domain.AuthIdentity], error) {
			called = true
			return domain.PageResult[domain.AuthIdentity]{}, nil
		},
	}

	handler := NewAuthIdentityHandler(service)
	rec := performJSONRequest(t, http.MethodGet, "/auth-identities?user_id=bad-uuid", "", func(c echo.Context) error {
		return handler.FindAll(c)
	})

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, rec.Code)
	}
	if called {
		t.Fatal("expected service find all not to be called")
	}

	var resp struct {
		Success bool `json:"success"`
		Error   struct {
			Code string `json:"code"`
		} `json:"error"`
	}
	decodeResponse(t, rec, &resp)

	if resp.Error.Code != "invalid_user_id" {
		t.Fatalf("expected error code invalid_user_id, got %q", resp.Error.Code)
	}
}
