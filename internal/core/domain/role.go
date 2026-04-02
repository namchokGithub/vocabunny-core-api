package domain

import "github.com/google/uuid"

type RoleName string

const (
	RoleNameAdmin        RoleName = "ADMIN"
	RoleNameContentAdmin RoleName = "CONTENT_ADMIN"
	RoleNameModerator    RoleName = "MODERATOR"
	RoleNameUser         RoleName = "USER"
)

type PermissionCode string

const (
	PermissionContentRead    PermissionCode = "CONTENT_READ"
	PermissionContentWrite   PermissionCode = "CONTENT_WRITE"
	PermissionContentPublish PermissionCode = "CONTENT_PUBLISH"
	PermissionUserRead       PermissionCode = "USER_READ"
	PermissionUserBan        PermissionCode = "USER_BAN"
	PermissionAnalyticsRead  PermissionCode = "ANALYTICS_READ"
	PermissionSystemConfig   PermissionCode = "SYSTEM_CONFIG"
)

type Role struct {
	ID          uuid.UUID
	Name        RoleName
	Description string
	Permissions []RolePermission
	AuditFields
}

type RolePermission struct {
	RoleID         uuid.UUID
	PermissionCode PermissionCode
	Scope          string
	AuditFields
}

type RoleCreateInput struct {
	Name        RoleName
	Description string
	Permissions []RolePermissionInput
	ActorID     string
}

type RolePermissionInput struct {
	PermissionCode PermissionCode
	Scope          string
}

type RoleUpdateInput struct {
	ID          uuid.UUID
	Name        EntityField[RoleName]
	Description EntityField[string]
	Permissions EntityField[[]RolePermissionInput]
	ActorID     string
}

type RoleQuery struct {
	Paging     Paging
	Search     string
	Permission *PermissionCode
	SortBy     string
	SortOrder  string
}
