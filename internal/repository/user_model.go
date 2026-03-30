package repository

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type RoleModel struct {
	ID        uuid.UUID      `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Code      string         `gorm:"size:64;uniqueIndex;not null"`
	Name      string         `gorm:"size:255;not null"`
	CreatedAt time.Time      `gorm:"not null"`
	UpdatedAt time.Time      `gorm:"not null"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
	CreatedBy string         `gorm:"size:255;not null"`
	UpdatedBy string         `gorm:"size:255;not null"`
}

func (RoleModel) TableName() string {
	return "roles"
}

type UserModel struct {
	ID        uuid.UUID      `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Code      string         `gorm:"size:64;uniqueIndex;not null"`
	Name      string         `gorm:"size:255;not null"`
	Email     string         `gorm:"size:255;uniqueIndex;not null"`
	Status    string         `gorm:"size:32;index;not null"`
	RoleID    *uuid.UUID     `gorm:"type:uuid;index"`
	Role      *RoleModel     `gorm:"foreignKey:RoleID"`
	CreatedAt time.Time      `gorm:"not null"`
	UpdatedAt time.Time      `gorm:"not null"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
	CreatedBy string         `gorm:"size:255;not null"`
	UpdatedBy string         `gorm:"size:255;not null"`
}

func (UserModel) TableName() string {
	return "users"
}
