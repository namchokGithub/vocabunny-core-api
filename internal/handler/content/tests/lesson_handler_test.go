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

func TestLessonHandlerCreate(t *testing.T) {
	t.Parallel()

	sectionID := uuid.New()
	lessonID := uuid.New()
	service := &lessonServiceStub{
		createFn: func(ctx context.Context, input domain.LessonCreateInput) (domain.Lesson, error) {
			if input.ActorID != "actor-123" {
				t.Fatalf("expected actor id actor-123, got %q", input.ActorID)
			}
			if input.SectionID != sectionID {
				t.Fatalf("expected section id %s, got %s", sectionID, input.SectionID)
			}
			return domain.Lesson{
				ID:          lessonID,
				SectionID:   input.SectionID,
				Slug:        input.Slug,
				Title:       input.Title,
				Description: input.Description,
				OrderNo:     input.OrderNo,
				IsPublished: input.IsPublished,
				AuditFields: testAuditFields("actor-123"),
			}, nil
		},
	}

	handler := content.NewLessonHandler(service)
	body := `{"section_id":"` + sectionID.String() + `","slug":"lesson-1","title":"Lesson 1","description":"Desc","order_no":1,"is_published":true}`
	rec := performJSONRequest(t, http.MethodPost, "/lessons", body, func(c echo.Context) error {
		c.Request().Header.Set("X-Actor-ID", "actor-123")
		return handler.Create(c)
	})

	if rec.Code != http.StatusCreated {
		t.Fatalf("expected status %d, got %d", http.StatusCreated, rec.Code)
	}
}

func TestLessonHandlerFindAllBuildsQuery(t *testing.T) {
	t.Parallel()

	sectionID := uuid.New()
	var captured domain.LessonQuery
	service := &lessonServiceStub{
		findAllFn: func(ctx context.Context, query domain.LessonQuery) (domain.PageResult[domain.Lesson], error) {
			captured = query
			return domain.PageResult[domain.Lesson]{
				Items:  []domain.Lesson{},
				Paging: domain.Paging{Page: query.Paging.Page, Limit: query.Paging.Limit, Total: 0},
			}, nil
		},
	}

	handler := content.NewLessonHandler(service)
	target := "/lessons?search=%20beginner%20&section_id=" + sectionID.String() + "&is_published=false&sort_by=order_no&sort_order=asc&page=3&limit=10"
	rec := performJSONRequest(t, http.MethodGet, target, "", func(c echo.Context) error {
		return handler.FindAll(c)
	})

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rec.Code)
	}
	if captured.Search != "beginner" {
		t.Fatalf("expected search beginner, got %q", captured.Search)
	}
	if captured.SectionID == nil || *captured.SectionID != sectionID {
		t.Fatalf("unexpected section id: %#v", captured.SectionID)
	}
	if captured.IsPublished == nil || *captured.IsPublished {
		t.Fatalf("unexpected published filter: %#v", captured.IsPublished)
	}
}

func TestLessonHandlerFindAllRejectsInvalidSectionID(t *testing.T) {
	t.Parallel()

	called := false
	service := &lessonServiceStub{
		findAllFn: func(ctx context.Context, query domain.LessonQuery) (domain.PageResult[domain.Lesson], error) {
			called = true
			return domain.PageResult[domain.Lesson]{}, nil
		},
	}

	handler := content.NewLessonHandler(service)
	rec := performJSONRequest(t, http.MethodGet, "/lessons?section_id=bad-uuid", "", func(c echo.Context) error {
		return handler.FindAll(c)
	})

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, rec.Code)
	}
	if called {
		t.Fatal("expected service find all not to be called")
	}
}
