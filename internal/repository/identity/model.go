package identity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserModel struct {
	ID          uuid.UUID      `gorm:"column:id;type:uuid;default:gen_random_uuid();primaryKey"`
	Email       string         `gorm:"column:email;size:255;uniqueIndex"`
	Username    string         `gorm:"column:username;size:255;uniqueIndex"`
	DisplayName string         `gorm:"column:display_name;size:255;not null"`
	AvatarID    *uuid.UUID     `gorm:"column:avatar_id;type:uuid"`
	Status      string         `gorm:"column:status;size:32;not null;index"`
	CreatedBy   string         `gorm:"column:created_by;size:255"`
	UpdatedBy   string         `gorm:"column:updated_by;size:255"`
	CreatedAt   time.Time      `gorm:"column:created_at;not null"`
	UpdatedAt   time.Time      `gorm:"column:updated_at;not null"`
	DeletedAt   gorm.DeletedAt `gorm:"column:deleted_at;index"`
}

func (UserModel) TableName() string {
	return "tbl_users"
}

type RoleModel struct {
	ID          uuid.UUID      `gorm:"column:id;type:uuid;default:gen_random_uuid();primaryKey"`
	Name        string         `gorm:"column:name;size:64;uniqueIndex;not null"`
	Description string         `gorm:"column:description"`
	CreatedBy   string         `gorm:"column:created_by;size:255"`
	UpdatedBy   string         `gorm:"column:updated_by;size:255"`
	CreatedAt   time.Time      `gorm:"column:created_at;not null"`
	UpdatedAt   time.Time      `gorm:"column:updated_at;not null"`
	DeletedAt   gorm.DeletedAt `gorm:"column:deleted_at;index"`
}

func (RoleModel) TableName() string {
	return "tbl_roles"
}

type UserRoleModel struct {
	UserID    uuid.UUID      `gorm:"column:user_id;type:uuid;primaryKey"`
	RoleID    uuid.UUID      `gorm:"column:role_id;type:uuid;primaryKey"`
	CreatedBy string         `gorm:"column:created_by;size:255"`
	UpdatedBy string         `gorm:"column:updated_by;size:255"`
	CreatedAt time.Time      `gorm:"column:created_at;not null"`
	UpdatedAt time.Time      `gorm:"column:updated_at;not null"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at;index"`
}

func (UserRoleModel) TableName() string {
	return "tbl_user_roles"
}

type RolePermissionModel struct {
	RoleID         uuid.UUID      `gorm:"column:role_id;type:uuid;primaryKey"`
	PermissionCode string         `gorm:"column:permission_code;size:64;primaryKey"`
	Scope          string         `gorm:"column:scope"`
	CreatedBy      string         `gorm:"column:created_by;size:255"`
	UpdatedBy      string         `gorm:"column:updated_by;size:255"`
	CreatedAt      time.Time      `gorm:"column:created_at;not null"`
	UpdatedAt      time.Time      `gorm:"column:updated_at;not null"`
	DeletedAt      gorm.DeletedAt `gorm:"column:deleted_at;index"`
}

func (RolePermissionModel) TableName() string {
	return "tbl_role_permissions"
}

type AuthIdentityModel struct {
	ID             uuid.UUID      `gorm:"column:id;type:uuid;default:gen_random_uuid();primaryKey"`
	UserID         uuid.UUID      `gorm:"column:user_id;type:uuid;index;not null"`
	Provider       string         `gorm:"column:provider;size:64;not null"`
	ProviderUserID string         `gorm:"column:provider_user_id"`
	PasswordHash   string         `gorm:"column:password_hash"`
	CreatedBy      string         `gorm:"column:created_by;size:255"`
	UpdatedBy      string         `gorm:"column:updated_by;size:255"`
	CreatedAt      time.Time      `gorm:"column:created_at;not null"`
	UpdatedAt      time.Time      `gorm:"column:updated_at;not null"`
	DeletedAt      gorm.DeletedAt `gorm:"column:deleted_at;index"`
}

func (AuthIdentityModel) TableName() string {
	return "tbl_auth_identities"
}
