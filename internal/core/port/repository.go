package port

import (
	"context"

	"github.com/google/uuid"
	"github.com/namchokGithub/vocabunny-core-api/internal/core/domain"
)

type UserRepository interface {
	Create(ctx context.Context, user domain.User) (domain.User, error)
	Update(ctx context.Context, user domain.UserUpdateInput) (domain.User, error)
	Delete(ctx context.Context, id uuid.UUID, actorID string) error
	FindByID(ctx context.Context, id uuid.UUID) (domain.User, error)
	FindByEmail(ctx context.Context, email string) (domain.User, error)
	FindByUsername(ctx context.Context, username string) (domain.User, error)
	FindAll(ctx context.Context, query domain.UserQuery) (domain.PageResult[domain.User], error)
	ExistsByEmail(ctx context.Context, email string, excludeID *uuid.UUID) (bool, error)
	ExistsByUsername(ctx context.Context, username string, excludeID *uuid.UUID) (bool, error)
	ReplaceRoles(ctx context.Context, userID uuid.UUID, roleIDs []uuid.UUID, actorID string) error
	FindBySubject(ctx context.Context, subject string) (domain.User, error)
}

type RoleRepository interface {
	Create(ctx context.Context, role domain.Role) (domain.Role, error)
	Update(ctx context.Context, role domain.RoleUpdateInput) (domain.Role, error)
	Delete(ctx context.Context, id uuid.UUID, actorID string) error
	FindByID(ctx context.Context, id uuid.UUID) (domain.Role, error)
	FindByName(ctx context.Context, name domain.RoleName) (domain.Role, error)
	FindAll(ctx context.Context, query domain.RoleQuery) (domain.PageResult[domain.Role], error)
	ExistsByName(ctx context.Context, name domain.RoleName, excludeID *uuid.UUID) (bool, error)
	ReplacePermissions(ctx context.Context, roleID uuid.UUID, permissions []domain.RolePermissionInput, actorID string) error
}

type AuthIdentityRepository interface {
	Create(ctx context.Context, identity domain.AuthIdentity) (domain.AuthIdentity, error)
	Update(ctx context.Context, input domain.AuthIdentityUpdateInput) (domain.AuthIdentity, error)
	Delete(ctx context.Context, id uuid.UUID, actorID string) error
	FindByID(ctx context.Context, id uuid.UUID) (domain.AuthIdentity, error)
	FindAll(ctx context.Context, query domain.AuthIdentityQuery) (domain.PageResult[domain.AuthIdentity], error)
	FindPasswordIdentityByUserID(ctx context.Context, userID uuid.UUID) (domain.AuthIdentity, error)
	FindPasswordIdentityByLogin(ctx context.Context, login string) (domain.AuthIdentity, error)
	ExistsByProviderIdentity(ctx context.Context, provider domain.AuthProvider, providerUserID string, excludeID *uuid.UUID) (bool, error)
}

type SectionRepository interface {
	Create(ctx context.Context, section domain.Section) (domain.Section, error)
	Update(ctx context.Context, input domain.SectionUpdateInput) (domain.Section, error)
	Delete(ctx context.Context, id uuid.UUID, actorID string) error
	FindByID(ctx context.Context, id uuid.UUID) (domain.Section, error)
	FindAll(ctx context.Context, query domain.SectionQuery) (domain.PageResult[domain.Section], error)
	ExistsBySlug(ctx context.Context, slug string, excludeID *uuid.UUID) (bool, error)
}

type LessonRepository interface {
	Create(ctx context.Context, lesson domain.Lesson) (domain.Lesson, error)
	Update(ctx context.Context, input domain.LessonUpdateInput) (domain.Lesson, error)
	Delete(ctx context.Context, id uuid.UUID, actorID string) error
	FindByID(ctx context.Context, id uuid.UUID) (domain.Lesson, error)
	FindAll(ctx context.Context, query domain.LessonQuery) (domain.PageResult[domain.Lesson], error)
	ExistsBySlug(ctx context.Context, sectionID uuid.UUID, slug string, excludeID *uuid.UUID) (bool, error)
}

type UnitRepository interface {
	Create(ctx context.Context, unit domain.Unit) (domain.Unit, error)
	Update(ctx context.Context, input domain.UnitUpdateInput) (domain.Unit, error)
	Delete(ctx context.Context, id uuid.UUID, actorID string) error
	FindByID(ctx context.Context, id uuid.UUID) (domain.Unit, error)
	FindAll(ctx context.Context, query domain.UnitQuery) (domain.PageResult[domain.Unit], error)
	ExistsBySlug(ctx context.Context, lessonID uuid.UUID, slug string, excludeID *uuid.UUID) (bool, error)
}

type QuestionSetRepository interface {
	Create(ctx context.Context, questionSet domain.QuestionSet) (domain.QuestionSet, error)
	Update(ctx context.Context, input domain.QuestionSetUpdateInput) (domain.QuestionSet, error)
	Delete(ctx context.Context, id uuid.UUID, actorID string) error
	FindByID(ctx context.Context, id uuid.UUID) (domain.QuestionSet, error)
	FindAll(ctx context.Context, query domain.QuestionSetQuery) (domain.PageResult[domain.QuestionSet], error)
	ExistsBySlugVersion(ctx context.Context, unitID uuid.UUID, slug string, version int, excludeID *uuid.UUID) (bool, error)
}

type QuestionRepository interface {
	Create(ctx context.Context, question domain.Question) (domain.Question, error)
	Update(ctx context.Context, input domain.QuestionUpdateInput) (domain.Question, error)
	Delete(ctx context.Context, id uuid.UUID, actorID string) error
	FindByID(ctx context.Context, id uuid.UUID) (domain.Question, error)
	FindAll(ctx context.Context, query domain.QuestionQuery) (domain.PageResult[domain.Question], error)
	ReplaceChoices(ctx context.Context, questionID uuid.UUID, choices []domain.QuestionChoiceInput, actorID string) error
	ReplaceTags(ctx context.Context, questionID uuid.UUID, tagIDs []uuid.UUID, actorID string) error
}

type QuestionChoiceRepository interface {
	Create(ctx context.Context, choice domain.QuestionChoice) (domain.QuestionChoice, error)
	Update(ctx context.Context, input domain.QuestionChoiceUpdateInput) (domain.QuestionChoice, error)
	Delete(ctx context.Context, id uuid.UUID, actorID string) error
	FindByID(ctx context.Context, id uuid.UUID) (domain.QuestionChoice, error)
	FindAll(ctx context.Context, query domain.QuestionChoiceQuery) (domain.PageResult[domain.QuestionChoice], error)
}

type TagRepository interface {
	Create(ctx context.Context, tag domain.Tag) (domain.Tag, error)
	Update(ctx context.Context, input domain.TagUpdateInput) (domain.Tag, error)
	Delete(ctx context.Context, id uuid.UUID, actorID string) error
	FindByID(ctx context.Context, id uuid.UUID) (domain.Tag, error)
	FindAll(ctx context.Context, query domain.TagQuery) (domain.PageResult[domain.Tag], error)
	ExistsByName(ctx context.Context, name string, excludeID *uuid.UUID) (bool, error)
}

type MediaAssetRepository interface {
	Create(ctx context.Context, asset domain.MediaAsset) (domain.MediaAsset, error)
	Update(ctx context.Context, input domain.MediaAssetUpdateInput) (domain.MediaAsset, error)
	Delete(ctx context.Context, id uuid.UUID, actorID string) error
	FindByID(ctx context.Context, id uuid.UUID) (domain.MediaAsset, error)
	FindAll(ctx context.Context, query domain.MediaAssetQuery) (domain.PageResult[domain.MediaAsset], error)
}

type TransactionManager interface {
	RunInTx(ctx context.Context, fn func(txCtx context.Context) error) error
}

type FileStorage interface {
	Save(ctx context.Context, path string, payload []byte) (string, error)
	Delete(ctx context.Context, path string) error
}

type TokenManager interface {
	GenerateAccessToken(subject string, scope string) (string, error)
	GenerateRefreshToken(subject string, scope string) (string, error)
	AccessTokenTTLSeconds() int64
	RefreshTokenTTLSeconds() int64
}
