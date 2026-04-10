package port

import (
	"context"

	"github.com/google/uuid"
	"github.com/namchokGithub/vocabunny-core-api/internal/core/domain"
)

type UserService interface {
	Create(ctx context.Context, input domain.UserCreateInput) (domain.User, error)
	Update(ctx context.Context, input domain.UserUpdateInput) (domain.User, error)
	Delete(ctx context.Context, id uuid.UUID, actorID string) error
	FindByID(ctx context.Context, id uuid.UUID) (domain.User, error)
	FindByEmail(ctx context.Context, email string) (domain.User, error)
	FindByUsername(ctx context.Context, username string) (domain.User, error)
	FindBySubject(ctx context.Context, subject string) (domain.User, error)
	FindAll(ctx context.Context, query domain.UserQuery) (domain.PageResult[domain.User], error)
}

type RoleService interface {
	Create(ctx context.Context, input domain.RoleCreateInput) (domain.Role, error)
	Update(ctx context.Context, input domain.RoleUpdateInput) (domain.Role, error)
	Delete(ctx context.Context, id uuid.UUID, actorID string) error
	FindByID(ctx context.Context, id uuid.UUID) (domain.Role, error)
	FindAll(ctx context.Context, query domain.RoleQuery) (domain.PageResult[domain.Role], error)
}

type AuthIdentityService interface {
	Create(ctx context.Context, input domain.AuthIdentityCreateInput) (domain.AuthIdentity, error)
	Update(ctx context.Context, input domain.AuthIdentityUpdateInput) (domain.AuthIdentity, error)
	Delete(ctx context.Context, id uuid.UUID, actorID string) error
	FindByID(ctx context.Context, id uuid.UUID) (domain.AuthIdentity, error)
	FindAll(ctx context.Context, query domain.AuthIdentityQuery) (domain.PageResult[domain.AuthIdentity], error)
	LoginWithPassword(ctx context.Context, input domain.PasswordLoginInput) (domain.AuthToken, error)
}

type SectionService interface {
	Create(ctx context.Context, input domain.SectionCreateInput) (domain.Section, error)
	Update(ctx context.Context, input domain.SectionUpdateInput) (domain.Section, error)
	Delete(ctx context.Context, id uuid.UUID, actorID string) error
	FindByID(ctx context.Context, id uuid.UUID) (domain.Section, error)
	FindAll(ctx context.Context, query domain.SectionQuery) (domain.PageResult[domain.Section], error)
}

type LessonService interface {
	Create(ctx context.Context, input domain.LessonCreateInput) (domain.Lesson, error)
	Update(ctx context.Context, input domain.LessonUpdateInput) (domain.Lesson, error)
	Delete(ctx context.Context, id uuid.UUID, actorID string) error
	FindByID(ctx context.Context, id uuid.UUID) (domain.Lesson, error)
	FindAll(ctx context.Context, query domain.LessonQuery) (domain.PageResult[domain.Lesson], error)
}

type UnitService interface {
	Create(ctx context.Context, input domain.UnitCreateInput) (domain.Unit, error)
	Update(ctx context.Context, input domain.UnitUpdateInput) (domain.Unit, error)
	Delete(ctx context.Context, id uuid.UUID, actorID string) error
	FindByID(ctx context.Context, id uuid.UUID) (domain.Unit, error)
	FindAll(ctx context.Context, query domain.UnitQuery) (domain.PageResult[domain.Unit], error)
}

type QuestionSetService interface {
	Create(ctx context.Context, input domain.QuestionSetCreateInput) (domain.QuestionSet, error)
	Update(ctx context.Context, input domain.QuestionSetUpdateInput) (domain.QuestionSet, error)
	Delete(ctx context.Context, id uuid.UUID, actorID string) error
	FindByID(ctx context.Context, id uuid.UUID) (domain.QuestionSet, error)
	FindAll(ctx context.Context, query domain.QuestionSetQuery) (domain.PageResult[domain.QuestionSet], error)
}

type QuestionService interface {
	Create(ctx context.Context, input domain.QuestionCreateInput) (domain.Question, error)
	Update(ctx context.Context, input domain.QuestionUpdateInput) (domain.Question, error)
	Delete(ctx context.Context, id uuid.UUID, actorID string) error
	FindByID(ctx context.Context, id uuid.UUID) (domain.Question, error)
	FindAll(ctx context.Context, query domain.QuestionQuery) (domain.PageResult[domain.Question], error)
}

type QuestionChoiceService interface {
	Create(ctx context.Context, input domain.QuestionChoiceCreateInput) (domain.QuestionChoice, error)
	Update(ctx context.Context, input domain.QuestionChoiceUpdateInput) (domain.QuestionChoice, error)
	Delete(ctx context.Context, id uuid.UUID, actorID string) error
	FindByID(ctx context.Context, id uuid.UUID) (domain.QuestionChoice, error)
	FindAll(ctx context.Context, query domain.QuestionChoiceQuery) (domain.PageResult[domain.QuestionChoice], error)
}

type TagService interface {
	Create(ctx context.Context, input domain.TagCreateInput) (domain.Tag, error)
	Update(ctx context.Context, input domain.TagUpdateInput) (domain.Tag, error)
	Delete(ctx context.Context, id uuid.UUID, actorID string) error
	FindByID(ctx context.Context, id uuid.UUID) (domain.Tag, error)
	FindAll(ctx context.Context, query domain.TagQuery) (domain.PageResult[domain.Tag], error)
}

type MediaAssetService interface {
	Create(ctx context.Context, input domain.MediaAssetCreateInput) (domain.MediaAsset, error)
	Update(ctx context.Context, input domain.MediaAssetUpdateInput) (domain.MediaAsset, error)
	Delete(ctx context.Context, id uuid.UUID, actorID string) error
	FindByID(ctx context.Context, id uuid.UUID) (domain.MediaAsset, error)
	FindAll(ctx context.Context, query domain.MediaAssetQuery) (domain.PageResult[domain.MediaAsset], error)
}
