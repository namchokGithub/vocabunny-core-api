package content_test

import (
	"context"
	"encoding/json"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/namchokGithub/vocabunny-core-api/internal/core/domain"
	"github.com/namchokGithub/vocabunny-core-api/internal/core/helper"
)

type sectionServiceStub struct {
	createFn   func(ctx context.Context, input domain.SectionCreateInput) (domain.Section, error)
	updateFn   func(ctx context.Context, input domain.SectionUpdateInput) (domain.Section, error)
	deleteFn   func(ctx context.Context, id uuid.UUID, actorID string) error
	findByIDFn func(ctx context.Context, id uuid.UUID) (domain.Section, error)
	findAllFn  func(ctx context.Context, query domain.SectionQuery) (domain.PageResult[domain.Section], error)
}

func (s *sectionServiceStub) Create(ctx context.Context, input domain.SectionCreateInput) (domain.Section, error) {
	return s.createFn(ctx, input)
}

func (s *sectionServiceStub) Update(ctx context.Context, input domain.SectionUpdateInput) (domain.Section, error) {
	return s.updateFn(ctx, input)
}

func (s *sectionServiceStub) Delete(ctx context.Context, id uuid.UUID, actorID string) error {
	return s.deleteFn(ctx, id, actorID)
}

func (s *sectionServiceStub) FindByID(ctx context.Context, id uuid.UUID) (domain.Section, error) {
	return s.findByIDFn(ctx, id)
}

func (s *sectionServiceStub) FindAll(ctx context.Context, query domain.SectionQuery) (domain.PageResult[domain.Section], error) {
	return s.findAllFn(ctx, query)
}

type lessonServiceStub struct {
	createFn   func(ctx context.Context, input domain.LessonCreateInput) (domain.Lesson, error)
	updateFn   func(ctx context.Context, input domain.LessonUpdateInput) (domain.Lesson, error)
	deleteFn   func(ctx context.Context, id uuid.UUID, actorID string) error
	findByIDFn func(ctx context.Context, id uuid.UUID) (domain.Lesson, error)
	findAllFn  func(ctx context.Context, query domain.LessonQuery) (domain.PageResult[domain.Lesson], error)
}

func (s *lessonServiceStub) Create(ctx context.Context, input domain.LessonCreateInput) (domain.Lesson, error) {
	return s.createFn(ctx, input)
}

func (s *lessonServiceStub) Update(ctx context.Context, input domain.LessonUpdateInput) (domain.Lesson, error) {
	return s.updateFn(ctx, input)
}

func (s *lessonServiceStub) Delete(ctx context.Context, id uuid.UUID, actorID string) error {
	return s.deleteFn(ctx, id, actorID)
}

func (s *lessonServiceStub) FindByID(ctx context.Context, id uuid.UUID) (domain.Lesson, error) {
	return s.findByIDFn(ctx, id)
}

func (s *lessonServiceStub) FindAll(ctx context.Context, query domain.LessonQuery) (domain.PageResult[domain.Lesson], error) {
	return s.findAllFn(ctx, query)
}

type unitServiceStub struct {
	createFn   func(ctx context.Context, input domain.UnitCreateInput) (domain.Unit, error)
	updateFn   func(ctx context.Context, input domain.UnitUpdateInput) (domain.Unit, error)
	deleteFn   func(ctx context.Context, id uuid.UUID, actorID string) error
	findByIDFn func(ctx context.Context, id uuid.UUID) (domain.Unit, error)
	findAllFn  func(ctx context.Context, query domain.UnitQuery) (domain.PageResult[domain.Unit], error)
}

func (s *unitServiceStub) Create(ctx context.Context, input domain.UnitCreateInput) (domain.Unit, error) {
	return s.createFn(ctx, input)
}

func (s *unitServiceStub) Update(ctx context.Context, input domain.UnitUpdateInput) (domain.Unit, error) {
	return s.updateFn(ctx, input)
}

func (s *unitServiceStub) Delete(ctx context.Context, id uuid.UUID, actorID string) error {
	return s.deleteFn(ctx, id, actorID)
}

func (s *unitServiceStub) FindByID(ctx context.Context, id uuid.UUID) (domain.Unit, error) {
	return s.findByIDFn(ctx, id)
}

func (s *unitServiceStub) FindAll(ctx context.Context, query domain.UnitQuery) (domain.PageResult[domain.Unit], error) {
	return s.findAllFn(ctx, query)
}

type questionSetServiceStub struct {
	createFn   func(ctx context.Context, input domain.QuestionSetCreateInput) (domain.QuestionSet, error)
	updateFn   func(ctx context.Context, input domain.QuestionSetUpdateInput) (domain.QuestionSet, error)
	deleteFn   func(ctx context.Context, id uuid.UUID, actorID string) error
	findByIDFn func(ctx context.Context, id uuid.UUID) (domain.QuestionSet, error)
	findAllFn  func(ctx context.Context, query domain.QuestionSetQuery) (domain.PageResult[domain.QuestionSet], error)
}

func (s *questionSetServiceStub) Create(ctx context.Context, input domain.QuestionSetCreateInput) (domain.QuestionSet, error) {
	return s.createFn(ctx, input)
}

func (s *questionSetServiceStub) Update(ctx context.Context, input domain.QuestionSetUpdateInput) (domain.QuestionSet, error) {
	return s.updateFn(ctx, input)
}

func (s *questionSetServiceStub) Delete(ctx context.Context, id uuid.UUID, actorID string) error {
	return s.deleteFn(ctx, id, actorID)
}

func (s *questionSetServiceStub) FindByID(ctx context.Context, id uuid.UUID) (domain.QuestionSet, error) {
	return s.findByIDFn(ctx, id)
}

func (s *questionSetServiceStub) FindAll(ctx context.Context, query domain.QuestionSetQuery) (domain.PageResult[domain.QuestionSet], error) {
	return s.findAllFn(ctx, query)
}

type questionServiceStub struct {
	createFn   func(ctx context.Context, input domain.QuestionCreateInput) (domain.Question, error)
	updateFn   func(ctx context.Context, input domain.QuestionUpdateInput) (domain.Question, error)
	deleteFn   func(ctx context.Context, id uuid.UUID, actorID string) error
	findByIDFn func(ctx context.Context, id uuid.UUID) (domain.Question, error)
	findAllFn  func(ctx context.Context, query domain.QuestionQuery) (domain.PageResult[domain.Question], error)
}

func (s *questionServiceStub) Create(ctx context.Context, input domain.QuestionCreateInput) (domain.Question, error) {
	return s.createFn(ctx, input)
}

func (s *questionServiceStub) Update(ctx context.Context, input domain.QuestionUpdateInput) (domain.Question, error) {
	return s.updateFn(ctx, input)
}

func (s *questionServiceStub) Delete(ctx context.Context, id uuid.UUID, actorID string) error {
	return s.deleteFn(ctx, id, actorID)
}

func (s *questionServiceStub) FindByID(ctx context.Context, id uuid.UUID) (domain.Question, error) {
	return s.findByIDFn(ctx, id)
}

func (s *questionServiceStub) FindAll(ctx context.Context, query domain.QuestionQuery) (domain.PageResult[domain.Question], error) {
	return s.findAllFn(ctx, query)
}

type questionChoiceServiceStub struct {
	createFn   func(ctx context.Context, input domain.QuestionChoiceCreateInput) (domain.QuestionChoice, error)
	updateFn   func(ctx context.Context, input domain.QuestionChoiceUpdateInput) (domain.QuestionChoice, error)
	deleteFn   func(ctx context.Context, id uuid.UUID, actorID string) error
	findByIDFn func(ctx context.Context, id uuid.UUID) (domain.QuestionChoice, error)
	findAllFn  func(ctx context.Context, query domain.QuestionChoiceQuery) (domain.PageResult[domain.QuestionChoice], error)
}

func (s *questionChoiceServiceStub) Create(ctx context.Context, input domain.QuestionChoiceCreateInput) (domain.QuestionChoice, error) {
	return s.createFn(ctx, input)
}

func (s *questionChoiceServiceStub) Update(ctx context.Context, input domain.QuestionChoiceUpdateInput) (domain.QuestionChoice, error) {
	return s.updateFn(ctx, input)
}

func (s *questionChoiceServiceStub) Delete(ctx context.Context, id uuid.UUID, actorID string) error {
	return s.deleteFn(ctx, id, actorID)
}

func (s *questionChoiceServiceStub) FindByID(ctx context.Context, id uuid.UUID) (domain.QuestionChoice, error) {
	return s.findByIDFn(ctx, id)
}

func (s *questionChoiceServiceStub) FindAll(ctx context.Context, query domain.QuestionChoiceQuery) (domain.PageResult[domain.QuestionChoice], error) {
	return s.findAllFn(ctx, query)
}

type tagServiceStub struct {
	createFn   func(ctx context.Context, input domain.TagCreateInput) (domain.Tag, error)
	updateFn   func(ctx context.Context, input domain.TagUpdateInput) (domain.Tag, error)
	deleteFn   func(ctx context.Context, id uuid.UUID, actorID string) error
	findByIDFn func(ctx context.Context, id uuid.UUID) (domain.Tag, error)
	findAllFn  func(ctx context.Context, query domain.TagQuery) (domain.PageResult[domain.Tag], error)
}

func (s *tagServiceStub) Create(ctx context.Context, input domain.TagCreateInput) (domain.Tag, error) {
	return s.createFn(ctx, input)
}

func (s *tagServiceStub) Update(ctx context.Context, input domain.TagUpdateInput) (domain.Tag, error) {
	return s.updateFn(ctx, input)
}

func (s *tagServiceStub) Delete(ctx context.Context, id uuid.UUID, actorID string) error {
	return s.deleteFn(ctx, id, actorID)
}

func (s *tagServiceStub) FindByID(ctx context.Context, id uuid.UUID) (domain.Tag, error) {
	return s.findByIDFn(ctx, id)
}

func (s *tagServiceStub) FindAll(ctx context.Context, query domain.TagQuery) (domain.PageResult[domain.Tag], error) {
	return s.findAllFn(ctx, query)
}

func performJSONRequest(t *testing.T, method, target, body string, run func(c echo.Context) error) *httptest.ResponseRecorder {
	t.Helper()

	e := echo.New()
	e.Validator = helper.NewRequestValidator()

	req := httptest.NewRequest(method, target, strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	if err := run(c); err != nil {
		t.Fatalf("handler returned error: %v", err)
	}

	return rec
}

func decodeResponse(t *testing.T, rec *httptest.ResponseRecorder, dest interface{}) {
	t.Helper()

	if err := json.Unmarshal(rec.Body.Bytes(), dest); err != nil {
		t.Fatalf("failed to decode response: %v; body=%s", err, rec.Body.String())
	}
}

func testAuditFields(actorID string) domain.AuditFields {
	now := time.Date(2026, 4, 9, 10, 0, 0, 0, time.UTC)
	return domain.AuditFields{
		CreatedAt: now,
		UpdatedAt: now,
		CreatedBy: actorID,
		UpdatedBy: actorID,
	}
}
