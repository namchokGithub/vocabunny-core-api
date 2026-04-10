package content_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/namchokGithub/vocabunny-core-api/internal/core/domain"
	content "github.com/namchokGithub/vocabunny-core-api/internal/handler/content"
)

func TestTagHandlerCreate(t *testing.T) {
	t.Parallel()

	tagID := uuid.New()
	service := &tagServiceStub{
		createFn: func(ctx context.Context, input domain.TagCreateInput) (domain.Tag, error) {
			if input.ActorID != "actor-123" {
				t.Fatalf("expected actor id actor-123, got %q", input.ActorID)
			}
			if input.Name != "grammar" {
				t.Fatalf("expected tag name grammar, got %q", input.Name)
			}
			return domain.Tag{
				ID:          tagID,
				Name:        input.Name,
				AuditFields: testAuditFields("actor-123"),
			}, nil
		},
	}

	handler := content.NewTagHandler(service)
	rec := performJSONRequest(t, http.MethodPost, "/tags", `{"name":"grammar"}`, func(c echo.Context) error {
		c.Request().Header.Set("X-Actor-ID", "actor-123")
		return handler.Create(c)
	})

	if rec.Code != http.StatusCreated {
		t.Fatalf("expected status %d, got %d", http.StatusCreated, rec.Code)
	}
}

func TestTagHandlerUpdateInvalidID(t *testing.T) {
	t.Parallel()

	called := false
	service := &tagServiceStub{
		updateFn: func(ctx context.Context, input domain.TagUpdateInput) (domain.Tag, error) {
			called = true
			return domain.Tag{}, nil
		},
	}

	handler := content.NewTagHandler(service)
	rec := performJSONRequest(t, http.MethodPut, "/tags/not-a-uuid", `{"name":"updated"}`, func(c echo.Context) error {
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
}

func TestTagHandlerFindAllBuildsQuery(t *testing.T) {
	t.Parallel()

	var captured domain.TagQuery
	service := &tagServiceStub{
		findAllFn: func(ctx context.Context, query domain.TagQuery) (domain.PageResult[domain.Tag], error) {
			captured = query
			return domain.PageResult[domain.Tag]{
				Items:  []domain.Tag{},
				Paging: domain.Paging{Page: query.Paging.Page, Limit: query.Paging.Limit, Total: 0},
			}, nil
		},
	}

	handler := content.NewTagHandler(service)
	rec := performJSONRequest(t, http.MethodGet, "/tags?search=%20vocab%20&sort_by=name&sort_order=asc&page=2&limit=9", "", func(c echo.Context) error {
		return handler.FindAll(c)
	})

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rec.Code)
	}
	if captured.Search != "vocab" {
		t.Fatalf("expected search vocab, got %q", captured.Search)
	}
	if captured.SortBy != "name" || captured.SortOrder != "asc" {
		t.Fatalf("unexpected sort: %s %s", captured.SortBy, captured.SortOrder)
	}
}
