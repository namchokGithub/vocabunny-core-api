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
