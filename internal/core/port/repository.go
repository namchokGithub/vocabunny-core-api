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
	FindByCode(ctx context.Context, code string) (domain.User, error)
	FindAll(ctx context.Context, query domain.UserQuery) (domain.PageResult[domain.User], error)
	ExistsByEmail(ctx context.Context, email string, excludeID *uuid.UUID) (bool, error)
	ExistsByCode(ctx context.Context, code string, excludeID *uuid.UUID) (bool, error)
	FindRoleByID(ctx context.Context, id uuid.UUID) (domain.Role, error)
}

type TransactionManager interface {
	RunInTx(ctx context.Context, fn func(txCtx context.Context) error) error
}

type FileStorage interface {
	Save(ctx context.Context, path string, payload []byte) (string, error)
	Delete(ctx context.Context, path string) error
}
