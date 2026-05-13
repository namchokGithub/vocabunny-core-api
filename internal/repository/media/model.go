package media

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type MediaAssetModel struct {
	ID              uuid.UUID      `gorm:"column:id;type:uuid;default:gen_random_uuid();primaryKey"`
	OwnerActorID    *uuid.UUID     `gorm:"column:owner_actor_id;type:uuid;index"`
	OwnerUserID     *uuid.UUID     `gorm:"column:owner_user_id;type:uuid;index"`
	AssetType       string         `gorm:"column:asset_type;size:64;not null;index"`
	Purpose         string         `gorm:"column:purpose;size:64;not null;index"`
	StorageMode     string         `gorm:"column:storage_mode;size:32;not null;index"`
	StorageProvider *string        `gorm:"column:storage_provider;size:32;index"`
	Bucket          *string        `gorm:"column:bucket"`
	ObjectKey       *string        `gorm:"column:object_key"`
	BinaryData      []byte         `gorm:"column:binary_data"`
	URL             *string        `gorm:"column:url"`
	ContentType     *string        `gorm:"column:content_type"`
	MimeType        string         `gorm:"column:mime_type;size:255;not null;index"`
	FileSizeBytes   *int64         `gorm:"column:file_size_bytes"`
	IsPublic        bool           `gorm:"column:is_public;not null;default:false;index"`
	CreatedBy       string         `gorm:"column:created_by;size:255"`
	UpdatedBy       string         `gorm:"column:updated_by;size:255"`
	CreatedAt       time.Time      `gorm:"column:created_at;not null"`
	UpdatedAt       time.Time      `gorm:"column:updated_at;not null"`
	DeletedAt       gorm.DeletedAt `gorm:"column:deleted_at;index"`
}

func (MediaAssetModel) TableName() string {
	return "tbl_media_assets"
}
