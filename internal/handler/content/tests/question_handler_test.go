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

func TestQuestionHandlerCreate(t *testing.T) {
	t.Parallel()

	questionSetID := uuid.New()
	tagID := uuid.New()
	questionID := uuid.New()
	service := &questionServiceStub{
		createFn: func(ctx context.Context, input domain.QuestionCreateInput) (domain.Question, error) {
			if input.QuestionSetID != questionSetID {
				t.Fatalf("expected question set id %s, got %s", questionSetID, input.QuestionSetID)
			}
			if input.ActorID != "actor-123" {
				t.Fatalf("expected actor id actor-123, got %q", input.ActorID)
			}
			if len(input.Choices) != 2 {
				t.Fatalf("expected 2 choices, got %d", len(input.Choices))
			}
			if len(input.TagIDs) != 1 || input.TagIDs[0] != tagID {
				t.Fatalf("unexpected tag ids: %#v", input.TagIDs)
			}
			return domain.Question{
				ID:            questionID,
				QuestionSetID: input.QuestionSetID,
				Type:          input.Type,
				QuestionText:  input.QuestionText,
				BlankPosition: input.BlankPosition,
				Explanation:   input.Explanation,
				Difficulty:    input.Difficulty,
				OrderNo:       input.OrderNo,
				IsActive:      input.IsActive,
				AuditFields:   testAuditFields("actor-123"),
			}, nil
		},
	}

	handler := content.NewQuestionHandler(service)
	body := `{
		"question_set_id":"` + questionSetID.String() + `",
		"type":"MULTIPLE_CHOICE",
		"question_text":"Which one is correct?",
		"blank_position":1,
		"explanation":"Because it is.",
		"difficulty":2,
		"order_no":3,
		"is_active":true,
		"choices":[
			{"choice_text":"A","choice_order":1,"is_correct":true},
			{"choice_text":"B","choice_order":2,"is_correct":false}
		],
		"tag_ids":["` + tagID.String() + `"]
	}`
	rec := performJSONRequest(t, http.MethodPost, "/questions", body, func(c echo.Context) error {
		c.Request().Header.Set("X-Actor-ID", "actor-123")
		return handler.Create(c)
	})

	if rec.Code != http.StatusCreated {
		t.Fatalf("expected status %d, got %d", http.StatusCreated, rec.Code)
	}
}

func TestQuestionHandlerUpdateInvalidID(t *testing.T) {
	t.Parallel()

	called := false
	service := &questionServiceStub{
		updateFn: func(ctx context.Context, input domain.QuestionUpdateInput) (domain.Question, error) {
			called = true
			return domain.Question{}, nil
		},
	}

	handler := content.NewQuestionHandler(service)
	rec := performJSONRequest(t, http.MethodPut, "/questions/not-a-uuid", `{"question_text":"Updated"}`, func(c echo.Context) error {
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

func TestQuestionHandlerFindAllBuildsQuery(t *testing.T) {
	t.Parallel()

	questionSetID := uuid.New()
	var captured domain.QuestionQuery
	service := &questionServiceStub{
		findAllFn: func(ctx context.Context, query domain.QuestionQuery) (domain.PageResult[domain.Question], error) {
			captured = query
			return domain.PageResult[domain.Question]{
				Items:  []domain.Question{},
				Paging: domain.Paging{Page: query.Paging.Page, Limit: query.Paging.Limit, Total: 0},
			}, nil
		},
	}

	handler := content.NewQuestionHandler(service)
	target := "/questions?search=%20correct%20&question_set_id=" + questionSetID.String() + "&type=MULTIPLE_CHOICE&is_active=true&include_choices=false&include_tags=false&sort_by=order_no&sort_order=desc&page=2&limit=8"
	rec := performJSONRequest(t, http.MethodGet, target, "", func(c echo.Context) error {
		return handler.FindAll(c)
	})

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rec.Code)
	}
	if captured.QuestionSetID == nil || *captured.QuestionSetID != questionSetID {
		t.Fatalf("unexpected question set id: %#v", captured.QuestionSetID)
	}
	if captured.Type == nil || *captured.Type != "MULTIPLE_CHOICE" {
		t.Fatalf("unexpected type: %#v", captured.Type)
	}
	if captured.IsActive == nil || !*captured.IsActive {
		t.Fatalf("unexpected active filter: %#v", captured.IsActive)
	}
	if captured.IncludeChoices || captured.IncludeTags {
		t.Fatalf("unexpected includes: choices=%v tags=%v", captured.IncludeChoices, captured.IncludeTags)
	}
}
