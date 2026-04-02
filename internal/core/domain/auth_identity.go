package domain

import "github.com/google/uuid"

type AuthProvider string

const (
	AuthProviderPassword AuthProvider = "PASSWORD"
	AuthProviderGoogle   AuthProvider = "GOOGLE"
	AuthProviderApple    AuthProvider = "APPLE"
)

type AuthIdentity struct {
	ID             uuid.UUID
	UserID         uuid.UUID
	Provider       AuthProvider
	ProviderUserID string
	PasswordHash   string
	AuditFields
}

type AuthIdentityCreateInput struct {
	UserID         uuid.UUID
	Provider       AuthProvider
	ProviderUserID string
	Password       string
	ActorID        string
}

type AuthIdentityUpdateInput struct {
	ID             uuid.UUID
	ProviderUserID EntityField[string]
	Password       EntityField[string]
	ActorID        string
}

type AuthIdentityQuery struct {
	Paging    Paging
	UserID    *uuid.UUID
	Provider  *AuthProvider
	SortBy    string
	SortOrder string
}

type PasswordLoginInput struct {
	EmailOrUsername string
	Password        string
}

type TokenType string

const (
	TokenTypeBearer  TokenType = "Bearer"
	TokenUseAccess   string    = "access"
	TokenUseRefresh  string    = "refresh"
	TokenScopeAccess string    = "api"
)

type AuthToken struct {
	AccessToken      string
	RefreshToken     string
	TokenType        TokenType
	ExpiresIn        int64
	RefreshExpiresIn int64
	User             User
}
