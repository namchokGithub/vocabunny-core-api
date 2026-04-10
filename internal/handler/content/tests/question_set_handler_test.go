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

func TestQuestionSetHandlerCreate(t *testing.T) {
	t.Parallel()

	unitID := uuid.New()
	questionSetID := uuid.New()
	service := &questionSetServiceStub{
		createFn: func(ctx context.Context, input domain.QuestionSetCreateInput) (domain.QuestionSet, error) {
			if input.UnitID != unitID {
				t.Fatalf("expected unit id %s, got %s", unitID, input.UnitID)
			}
			if input.ActorID != "actor-123" {
				t.Fatalf("expected actor id actor-123, got %q", input.ActorID)
			}
			return domain.QuestionSet{
				ID:               questionSetID,
				UnitID:           input.UnitID,
				Slug:             input.Slug,
				Title:            input.Title,
				Description:      input.Description,
				OrderNo:          input.OrderNo,
				EstimatedSeconds: input.EstimatedSeconds,
				IsPublished:      input.IsPublished,
				Version:          input.Version,
				AuditFields:      testAuditFields("actor-123"),
			}, nil
		},
	}

	handler := content.NewQuestionSetHandler(service)
	body := `{"unit_id":"` + unitID.String() + `","slug":"set-1","title":"Set 1","description":"Desc","order_no":1,"estimated_seconds":90,"is_published":true,"version":2}`
	rec := performJSONRequest(t, http.MethodPost, "/question-sets", body, func(c echo.Context) error {
		c.Request().Header.Set("X-Actor-ID", "actor-123")
		return handler.Create(c)
	})

	if rec.Code != http.StatusCreated {
		t.Fatalf("expected status %d, got %d", http.StatusCreated, rec.Code)
	}
}

func TestQuestionSetHandlerFindAllBuildsQuery(t *testing.T) {
	t.Parallel()

	unitID := uuid.New()
	var captured domain.QuestionSetQuery
	service := &questionSetServiceStub{
		findAllFn: func(ctx context.Context, query domain.QuestionSetQuery) (domain.PageResult[domain.QuestionSet], error) {
			captured = query
			return domain.PageResult[domain.QuestionSet]{
				Items:  []domain.QuestionSet{},
				Paging: domain.Paging{Page: query.Paging.Page, Limit: query.Paging.Limit, Total: 0},
			}, nil
		},
	}

	handler := content.NewQuestionSetHandler(service)
	target := "/question-sets?search=%20set%20&unit_id=" + unitID.String() + "&version=3&is_published=false&sort_by=version&sort_order=desc&page=2&limit=4"
	rec := performJSONRequest(t, http.MethodGet, target, "", func(c echo.Context) error {
		return handler.FindAll(c)
	})

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rec.Code)
	}
	if captured.UnitID == nil || *captured.UnitID != unitID {
		t.Fatalf("unexpected unit id: %#v", captured.UnitID)
	}
	if captured.Version == nil || *captured.Version != 3 {
		t.Fatalf("unexpected version: %#v", captured.Version)
	}
	if captured.IsPublished == nil || *captured.IsPublished {
		t.Fatalf("unexpected published filter: %#v", captured.IsPublished)
	}
}

func TestQuestionSetHandlerFindAllRejectsInvalidVersion(t *testing.T) {
	t.Parallel()

	called := false
	service := &questionSetServiceStub{
		findAllFn: func(ctx context.Context, query domain.QuestionSetQuery) (domain.PageResult[domain.QuestionSet], error) {
			called = true
			return domain.PageResult[domain.QuestionSet]{}, nil
		},
	}

	handler := content.NewQuestionSetHandler(service)
	rec := performJSONRequest(t, http.MethodGet, "/question-sets?version=bad", "", func(c echo.Context) error {
		return handler.FindAll(c)
	})

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, rec.Code)
	}
	if called {
		t.Fatal("expected service find all not to be called")
	}
}
