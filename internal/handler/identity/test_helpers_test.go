package identity

import (
	"context"
	"encoding/json"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/namchokGithub/vocabunny-core-api/internal/core/domain"
	"github.com/namchokGithub/vocabunny-core-api/internal/core/helper"
)

type userServiceStub struct {
	createFn         func(ctx context.Context, input domain.UserCreateInput) (domain.User, error)
	updateFn         func(ctx context.Context, input domain.UserUpdateInput) (domain.User, error)
	deleteFn         func(ctx context.Context, id uuid.UUID, actorID string) error
	findByIDFn       func(ctx context.Context, id uuid.UUID) (domain.User, error)
	findByEmailFn    func(ctx context.Context, email string) (domain.User, error)
	findByUsernameFn func(ctx context.Context, username string) (domain.User, error)
	findBySubjectFn  func(ctx context.Context, subject string) (domain.User, error)
	findAllFn        func(ctx context.Context, query domain.UserQuery) (domain.PageResult[domain.User], error)
}

func (s *userServiceStub) Create(ctx context.Context, input domain.UserCreateInput) (domain.User, error) {
	return s.createFn(ctx, input)
}

func (s *userServiceStub) Update(ctx context.Context, input domain.UserUpdateInput) (domain.User, error) {
	return s.updateFn(ctx, input)
}

func (s *userServiceStub) Delete(ctx context.Context, id uuid.UUID, actorID string) error {
	return s.deleteFn(ctx, id, actorID)
}

func (s *userServiceStub) FindByID(ctx context.Context, id uuid.UUID) (domain.User, error) {
	return s.findByIDFn(ctx, id)
}

func (s *userServiceStub) FindByEmail(ctx context.Context, email string) (domain.User, error) {
	return s.findByEmailFn(ctx, email)
}

func (s *userServiceStub) FindByUsername(ctx context.Context, username string) (domain.User, error) {
	return s.findByUsernameFn(ctx, username)
}

func (s *userServiceStub) FindBySubject(ctx context.Context, subject string) (domain.User, error) {
	return s.findBySubjectFn(ctx, subject)
}

func (s *userServiceStub) FindAll(ctx context.Context, query domain.UserQuery) (domain.PageResult[domain.User], error) {
	return s.findAllFn(ctx, query)
}

type roleServiceStub struct {
	createFn   func(ctx context.Context, input domain.RoleCreateInput) (domain.Role, error)
	updateFn   func(ctx context.Context, input domain.RoleUpdateInput) (domain.Role, error)
	deleteFn   func(ctx context.Context, id uuid.UUID, actorID string) error
	findByIDFn func(ctx context.Context, id uuid.UUID) (domain.Role, error)
	findAllFn  func(ctx context.Context, query domain.RoleQuery) (domain.PageResult[domain.Role], error)
}

func (s *roleServiceStub) Create(ctx context.Context, input domain.RoleCreateInput) (domain.Role, error) {
	return s.createFn(ctx, input)
}

func (s *roleServiceStub) Update(ctx context.Context, input domain.RoleUpdateInput) (domain.Role, error) {
	return s.updateFn(ctx, input)
}

func (s *roleServiceStub) Delete(ctx context.Context, id uuid.UUID, actorID string) error {
	return s.deleteFn(ctx, id, actorID)
}

func (s *roleServiceStub) FindByID(ctx context.Context, id uuid.UUID) (domain.Role, error) {
	return s.findByIDFn(ctx, id)
}

func (s *roleServiceStub) FindAll(ctx context.Context, query domain.RoleQuery) (domain.PageResult[domain.Role], error) {
	return s.findAllFn(ctx, query)
}

type authIdentityServiceStub struct {
	createFn            func(ctx context.Context, input domain.AuthIdentityCreateInput) (domain.AuthIdentity, error)
	updateFn            func(ctx context.Context, input domain.AuthIdentityUpdateInput) (domain.AuthIdentity, error)
	deleteFn            func(ctx context.Context, id uuid.UUID, actorID string) error
	findByIDFn          func(ctx context.Context, id uuid.UUID) (domain.AuthIdentity, error)
	findAllFn           func(ctx context.Context, query domain.AuthIdentityQuery) (domain.PageResult[domain.AuthIdentity], error)
	loginWithPasswordFn func(ctx context.Context, input domain.PasswordLoginInput) (domain.AuthToken, error)
}

func (s *authIdentityServiceStub) Create(ctx context.Context, input domain.AuthIdentityCreateInput) (domain.AuthIdentity, error) {
	return s.createFn(ctx, input)
}

func (s *authIdentityServiceStub) Update(ctx context.Context, input domain.AuthIdentityUpdateInput) (domain.AuthIdentity, error) {
	return s.updateFn(ctx, input)
}

func (s *authIdentityServiceStub) Delete(ctx context.Context, id uuid.UUID, actorID string) error {
	return s.deleteFn(ctx, id, actorID)
}

func (s *authIdentityServiceStub) FindByID(ctx context.Context, id uuid.UUID) (domain.AuthIdentity, error) {
	return s.findByIDFn(ctx, id)
}

func (s *authIdentityServiceStub) FindAll(ctx context.Context, query domain.AuthIdentityQuery) (domain.PageResult[domain.AuthIdentity], error) {
	return s.findAllFn(ctx, query)
}

func (s *authIdentityServiceStub) LoginWithPassword(ctx context.Context, input domain.PasswordLoginInput) (domain.AuthToken, error) {
	return s.loginWithPasswordFn(ctx, input)
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
