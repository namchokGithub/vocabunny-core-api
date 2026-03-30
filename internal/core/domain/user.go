package domain

import "github.com/google/uuid"

type UserStatus string

const (
	UserStatusActive   UserStatus = "active"
	UserStatusInactive UserStatus = "inactive"
)

type User struct {
	ID     uuid.UUID
	Code   string
	Name   string
	Email  string
	Status UserStatus
	RoleID *uuid.UUID
	Role   *Role
	AuditFields
}

type UserCreateInput struct {
	Code    string
	Name    string
	Email   string
	Status  UserStatus
	RoleID  *uuid.UUID
	ActorID string
}

type UserUpdateInput struct {
	ID      uuid.UUID
	Name    EntityField[string]
	Email   EntityField[string]
	Code    EntityField[string]
	Status  EntityField[UserStatus]
	RoleID  EntityField[*uuid.UUID]
	ActorID string
}

type UserQuery struct {
	Paging    Paging
	Search    string
	RoleID    *uuid.UUID
	Status    *UserStatus
	SortBy    string
	SortOrder string
}
