package domain

import "github.com/google/uuid"

type UserStatus string

const (
	UserStatusInactive UserStatus = "INACTIVE"
	UserStatusActive   UserStatus = "ACTIVE"
	UserStatusBanned   UserStatus = "BANNED"
	UserStatusDeleted  UserStatus = "DELETED"
)

type User struct {
	ID          uuid.UUID
	Email       string
	Username    string
	DisplayName string
	AvatarID    *uuid.UUID
	Status      UserStatus
	Roles       []Role
	Identities  []AuthIdentity
	AuditFields
}

type UserCreateInput struct {
	Email       string
	Username    string
	DisplayName string
	AvatarID    *uuid.UUID
	Status      UserStatus
	RoleIDs     []uuid.UUID
	ActorID     string
}

type UserUpdateInput struct {
	ID          uuid.UUID
	Email       EntityField[string]
	Username    EntityField[string]
	DisplayName EntityField[string]
	AvatarID    EntityField[*uuid.UUID]
	Status      EntityField[UserStatus]
	RoleIDs     EntityField[[]uuid.UUID]
	ActorID     string
}

type UserQuery struct {
	Paging      Paging
	Search      string
	RoleID      *uuid.UUID
	Status      *UserStatus
	SortBy      string
	SortOrder   string
	IncludeAuth bool
}
