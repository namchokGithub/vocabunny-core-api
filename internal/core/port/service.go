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
