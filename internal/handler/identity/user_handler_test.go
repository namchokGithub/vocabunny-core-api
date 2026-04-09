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

func TestUserHandlerCreate(t *testing.T) {
	t.Parallel()

	roleID := uuid.New()
	userID := uuid.New()
	now := time.Date(2026, 4, 9, 10, 0, 0, 0, time.UTC)

	service := &userServiceStub{
		createFn: func(ctx context.Context, input domain.UserCreateInput) (domain.User, error) {
			if input.ActorID != "actor-123" {
				t.Fatalf("expected actor id actor-123, got %q", input.ActorID)
			}
			if input.Email != "bunny@example.com" {
				t.Fatalf("expected email bunny@example.com, got %q", input.Email)
			}
			if len(input.RoleIDs) != 1 || input.RoleIDs[0] != roleID {
				t.Fatalf("unexpected role ids: %#v", input.RoleIDs)
			}

			return domain.User{
				ID:          userID,
				Email:       input.Email,
				Username:    input.Username,
				DisplayName: input.DisplayName,
				Status:      input.Status,
				Roles: []domain.Role{{
					ID:   roleID,
					Name: domain.RoleNameUser,
					AuditFields: domain.AuditFields{
						CreatedAt: now,
						UpdatedAt: now,
						CreatedBy: "seed",
						UpdatedBy: "seed",
					},
				}},
				AuditFields: domain.AuditFields{
					CreatedAt: now,
					UpdatedAt: now,
					CreatedBy: "actor-123",
					UpdatedBy: "actor-123",
				},
			}, nil
		},
	}

	handler := NewUserHandler(service)
	body := `{"email":"bunny@example.com","username":"bun","display_name":"Bun","status":"ACTIVE","role_ids":["` + roleID.String() + `"]}`
	rec := performJSONRequest(t, http.MethodPost, "/users", body, func(c echo.Context) error {
		c.Request().Header.Set("X-Actor-ID", "actor-123")
		return handler.Create(c)
	})

	if rec.Code != http.StatusCreated {
		t.Fatalf("expected status %d, got %d", http.StatusCreated, rec.Code)
	}

	var resp struct {
		Success bool         `json:"success"`
		Data    UserResponse `json:"data"`
	}
	decodeResponse(t, rec, &resp)

	if !resp.Success {
		t.Fatal("expected success response")
	}
	if resp.Data.ID != userID.String() {
		t.Fatalf("expected user id %s, got %s", userID, resp.Data.ID)
	}
	if resp.Data.Email != "bunny@example.com" {
		t.Fatalf("expected email bunny@example.com, got %q", resp.Data.Email)
	}
}

func TestUserHandlerUpdateInvalidID(t *testing.T) {
	t.Parallel()

	called := false
	service := &userServiceStub{
		updateFn: func(ctx context.Context, input domain.UserUpdateInput) (domain.User, error) {
			called = true
			return domain.User{}, nil
		},
	}

	handler := NewUserHandler(service)
	rec := performJSONRequest(t, http.MethodPut, "/users/not-a-uuid", `{"display_name":"New Name"}`, func(c echo.Context) error {
		c.SetParamNames("id")
		c.SetParamValues("not-a-uuid")
		return handler.Update(c)
	})

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, rec.Code)
	}
	if called {
		t.Fatal("expected service update not to be called")
	}

	var resp struct {
		Success bool `json:"success"`
		Error   struct {
			Code string `json:"code"`
		} `json:"error"`
	}
	decodeResponse(t, rec, &resp)

	if resp.Success {
		t.Fatal("expected error response")
	}
	if resp.Error.Code != "invalid_user_id" {
		t.Fatalf("expected error code invalid_user_id, got %q", resp.Error.Code)
	}
}

func TestUserHandlerFindAllBuildsQuery(t *testing.T) {
	t.Parallel()

	roleID := uuid.New()
	userID := uuid.New()
	var captured domain.UserQuery

	service := &userServiceStub{
		findAllFn: func(ctx context.Context, query domain.UserQuery) (domain.PageResult[domain.User], error) {
			captured = query
			return domain.PageResult[domain.User]{
				Items: []domain.User{{
					ID:          userID,
					Email:       "find@example.com",
					Username:    "finder",
					DisplayName: "Finder",
					Status:      domain.UserStatusActive,
					AuditFields: domain.AuditFields{
						CreatedAt: time.Date(2026, 4, 9, 10, 0, 0, 0, time.UTC),
						UpdatedAt: time.Date(2026, 4, 9, 10, 0, 0, 0, time.UTC),
						CreatedBy: "seed",
						UpdatedBy: "seed",
					},
				}},
				Paging: domain.Paging{Page: query.Paging.Page, Limit: query.Paging.Limit, Total: 1},
			}, nil
		},
	}

	handler := NewUserHandler(service)
	rec := performJSONRequest(t, http.MethodGet, "/users?search=%20rabbit%20&include_auth=true&page=2&limit=5&sort_by=email&sort_order=desc&status=ACTIVE&role_id="+roleID.String(), "", func(c echo.Context) error {
		return handler.FindAll(c)
	})

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rec.Code)
	}
	if captured.Search != "rabbit" {
		t.Fatalf("expected search rabbit, got %q", captured.Search)
	}
	if !captured.IncludeAuth {
		t.Fatal("expected include auth to be true")
	}
	if captured.Paging.Page != 2 || captured.Paging.Limit != 5 {
		t.Fatalf("unexpected paging: %#v", captured.Paging)
	}
	if captured.SortBy != "email" || captured.SortOrder != "desc" {
		t.Fatalf("unexpected sort: %s %s", captured.SortBy, captured.SortOrder)
	}
	if captured.RoleID == nil || *captured.RoleID != roleID {
		t.Fatalf("unexpected role id: %#v", captured.RoleID)
	}
	if captured.Status == nil || *captured.Status != domain.UserStatusActive {
		t.Fatalf("unexpected status: %#v", captured.Status)
	}
}
