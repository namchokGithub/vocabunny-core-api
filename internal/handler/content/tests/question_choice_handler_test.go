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

func TestQuestionChoiceHandlerCreate(t *testing.T) {
	t.Parallel()

	questionID := uuid.New()
	choiceID := uuid.New()
	service := &questionChoiceServiceStub{
		createFn: func(ctx context.Context, input domain.QuestionChoiceCreateInput) (domain.QuestionChoice, error) {
			if input.QuestionID != questionID {
				t.Fatalf("expected question id %s, got %s", questionID, input.QuestionID)
			}
			if input.ActorID != "actor-123" {
				t.Fatalf("expected actor id actor-123, got %q", input.ActorID)
			}
			return domain.QuestionChoice{
				ID:          choiceID,
				QuestionID:  input.QuestionID,
				ChoiceText:  input.ChoiceText,
				ChoiceOrder: input.ChoiceOrder,
				IsCorrect:   input.IsCorrect,
				AuditFields: testAuditFields("actor-123"),
			}, nil
		},
	}

	handler := content.NewQuestionChoiceHandler(service)
	body := `{"question_id":"` + questionID.String() + `","choice_text":"Choice A","choice_order":1,"is_correct":true}`
	rec := performJSONRequest(t, http.MethodPost, "/question-choices", body, func(c echo.Context) error {
		c.Request().Header.Set("X-Actor-ID", "actor-123")
		return handler.Create(c)
	})

	if rec.Code != http.StatusCreated {
		t.Fatalf("expected status %d, got %d", http.StatusCreated, rec.Code)
	}
}

func TestQuestionChoiceHandlerFindAllRejectsInvalidQuestionID(t *testing.T) {
	t.Parallel()

	called := false
	service := &questionChoiceServiceStub{
		findAllFn: func(ctx context.Context, query domain.QuestionChoiceQuery) (domain.PageResult[domain.QuestionChoice], error) {
			called = true
			return domain.PageResult[domain.QuestionChoice]{}, nil
		},
	}

	handler := content.NewQuestionChoiceHandler(service)
	rec := performJSONRequest(t, http.MethodGet, "/question-choices?question_id=bad-uuid", "", func(c echo.Context) error {
		return handler.FindAll(c)
	})

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, rec.Code)
	}
	if called {
		t.Fatal("expected service find all not to be called")
	}
}
