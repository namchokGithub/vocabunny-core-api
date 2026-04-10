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

func TestUnitHandlerCreate(t *testing.T) {
	t.Parallel()

	lessonID := uuid.New()
	unitID := uuid.New()
	service := &unitServiceStub{
		createFn: func(ctx context.Context, input domain.UnitCreateInput) (domain.Unit, error) {
			if input.LessonID != lessonID {
				t.Fatalf("expected lesson id %s, got %s", lessonID, input.LessonID)
			}
			if input.ActorID != "actor-123" {
				t.Fatalf("expected actor id actor-123, got %q", input.ActorID)
			}
			return domain.Unit{
				ID:          unitID,
				LessonID:    input.LessonID,
				Slug:        input.Slug,
				Title:       input.Title,
				Description: input.Description,
				OrderNo:     input.OrderNo,
				IsPublished: input.IsPublished,
				AuditFields: testAuditFields("actor-123"),
			}, nil
		},
	}

	handler := content.NewUnitHandler(service)
	body := `{"lesson_id":"` + lessonID.String() + `","slug":"unit-1","title":"Unit 1","description":"Desc","order_no":2,"is_published":true}`
	rec := performJSONRequest(t, http.MethodPost, "/units", body, func(c echo.Context) error {
		c.Request().Header.Set("X-Actor-ID", "actor-123")
		return handler.Create(c)
	})

	if rec.Code != http.StatusCreated {
		t.Fatalf("expected status %d, got %d", http.StatusCreated, rec.Code)
	}
}

func TestUnitHandlerFindAllBuildsQuery(t *testing.T) {
	t.Parallel()

	lessonID := uuid.New()
	var captured domain.UnitQuery
	service := &unitServiceStub{
		findAllFn: func(ctx context.Context, query domain.UnitQuery) (domain.PageResult[domain.Unit], error) {
			captured = query
			return domain.PageResult[domain.Unit]{
				Items:  []domain.Unit{},
				Paging: domain.Paging{Page: query.Paging.Page, Limit: query.Paging.Limit, Total: 0},
			}, nil
		},
	}

	handler := content.NewUnitHandler(service)
	target := "/units?search=%20core%20&lesson_id=" + lessonID.String() + "&is_published=true&sort_by=title&sort_order=desc&page=4&limit=6"
	rec := performJSONRequest(t, http.MethodGet, target, "", func(c echo.Context) error {
		return handler.FindAll(c)
	})

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rec.Code)
	}
	if captured.Search != "core" {
		t.Fatalf("expected search core, got %q", captured.Search)
	}
	if captured.LessonID == nil || *captured.LessonID != lessonID {
		t.Fatalf("unexpected lesson id: %#v", captured.LessonID)
	}
	if captured.IsPublished == nil || !*captured.IsPublished {
		t.Fatalf("unexpected published filter: %#v", captured.IsPublished)
	}
}

func TestUnitHandlerFindAllRejectsInvalidLessonID(t *testing.T) {
	t.Parallel()

	called := false
	service := &unitServiceStub{
		findAllFn: func(ctx context.Context, query domain.UnitQuery) (domain.PageResult[domain.Unit], error) {
			called = true
			return domain.PageResult[domain.Unit]{}, nil
		},
	}

	handler := content.NewUnitHandler(service)
	rec := performJSONRequest(t, http.MethodGet, "/units?lesson_id=bad-uuid", "", func(c echo.Context) error {
		return handler.FindAll(c)
	})

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, rec.Code)
	}
	if called {
		t.Fatal("expected service find all not to be called")
	}
}
