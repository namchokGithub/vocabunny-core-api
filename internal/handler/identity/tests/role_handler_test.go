package identity_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/namchokGithub/vocabunny-core-api/internal/core/domain"
	identity "github.com/namchokGithub/vocabunny-core-api/internal/handler/identity"
)

func TestRoleHandlerFindAllBuildsQuery(t *testing.T) {
	t.Parallel()

	var captured domain.RoleQuery
	service := &roleServiceStub{
		findAllFn: func(ctx context.Context, query domain.RoleQuery) (domain.PageResult[domain.Role], error) {
			captured = query
			return domain.PageResult[domain.Role]{
				Items:  []domain.Role{},
				Paging: domain.Paging{Page: query.Paging.Page, Limit: query.Paging.Limit, Total: 0},
			}, nil
		},
	}

	handler := identity.NewRoleHandler(service)
	rec := performJSONRequest(t, http.MethodGet, "/roles?search=%20admin%20&permission=USER_READ&sort_by=name&sort_order=desc&page=3&limit=7", "", func(c echo.Context) error {
		return handler.FindAll(c)
	})

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rec.Code)
	}
	if captured.Search != "admin" {
		t.Fatalf("expected search admin, got %q", captured.Search)
	}
	if captured.Permission == nil || *captured.Permission != domain.PermissionUserRead {
		t.Fatalf("unexpected permission: %#v", captured.Permission)
	}
	if captured.SortBy != "name" || captured.SortOrder != "desc" {
		t.Fatalf("unexpected sort: %s %s", captured.SortBy, captured.SortOrder)
	}
	if captured.Paging.Page != 3 || captured.Paging.Limit != 7 {
		t.Fatalf("unexpected paging: %#v", captured.Paging)
	}
}
