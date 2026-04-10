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

func TestSectionHandlerCreate(t *testing.T) {
	t.Parallel()

	sectionID := uuid.New()
	service := &sectionServiceStub{
		createFn: func(ctx context.Context, input domain.SectionCreateInput) (domain.Section, error) {
			if input.ActorID != "actor-123" {
				t.Fatalf("expected actor id actor-123, got %q", input.ActorID)
			}
			if input.Slug != "greetings" || input.Title != "Greetings" {
				t.Fatalf("unexpected input: %#v", input)
			}
			return domain.Section{
				ID:          sectionID,
				Slug:        input.Slug,
				Title:       input.Title,
				Description: input.Description,
				OrderNo:     input.OrderNo,
				IsPublished: input.IsPublished,
				AuditFields: testAuditFields("actor-123"),
			}, nil
		},
	}

	handler := content.NewSectionHandler(service)
	rec := performJSONRequest(t, http.MethodPost, "/sections", `{"slug":"greetings","title":"Greetings","description":"Intro","order_no":1,"is_published":true}`, func(c echo.Context) error {
		c.Request().Header.Set("X-Actor-ID", "actor-123")
		return handler.Create(c)
	})

	if rec.Code != http.StatusCreated {
		t.Fatalf("expected status %d, got %d", http.StatusCreated, rec.Code)
	}

	var resp struct {
		Success bool                    `json:"success"`
		Data    content.SectionResponse `json:"data"`
	}
	decodeResponse(t, rec, &resp)

	if !resp.Success || resp.Data.ID != sectionID.String() {
		t.Fatalf("unexpected response: %#v", resp)
	}
}

func TestSectionHandlerUpdateInvalidID(t *testing.T) {
	t.Parallel()

	called := false
	service := &sectionServiceStub{
		updateFn: func(ctx context.Context, input domain.SectionUpdateInput) (domain.Section, error) {
			called = true
			return domain.Section{}, nil
		},
	}

	handler := content.NewSectionHandler(service)
	rec := performJSONRequest(t, http.MethodPut, "/sections/not-a-uuid", `{"title":"Updated"}`, func(c echo.Context) error {
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

func TestSectionHandlerFindAllBuildsQuery(t *testing.T) {
	t.Parallel()

	var captured domain.SectionQuery
	service := &sectionServiceStub{
		findAllFn: func(ctx context.Context, query domain.SectionQuery) (domain.PageResult[domain.Section], error) {
			captured = query
			return domain.PageResult[domain.Section]{
				Items:  []domain.Section{},
				Paging: domain.Paging{Page: query.Paging.Page, Limit: query.Paging.Limit, Total: 0},
			}, nil
		},
	}

	handler := content.NewSectionHandler(service)
	rec := performJSONRequest(t, http.MethodGet, "/sections?search=%20greet%20&is_published=true&sort_by=title&sort_order=desc&page=2&limit=5", "", func(c echo.Context) error {
		return handler.FindAll(c)
	})

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rec.Code)
	}
	if captured.Search != "greet" {
		t.Fatalf("expected search greet, got %q", captured.Search)
	}
	if captured.IsPublished == nil || !*captured.IsPublished {
		t.Fatalf("unexpected published filter: %#v", captured.IsPublished)
	}
	if captured.SortBy != "title" || captured.SortOrder != "desc" {
		t.Fatalf("unexpected sort: %s %s", captured.SortBy, captured.SortOrder)
	}
	if captured.Paging.Page != 2 || captured.Paging.Limit != 5 {
		t.Fatalf("unexpected paging: %#v", captured.Paging)
	}
}
