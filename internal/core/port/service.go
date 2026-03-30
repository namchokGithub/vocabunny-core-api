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
	FindAll(ctx context.Context, query domain.UserQuery) (domain.PageResult[domain.User], error)
}
