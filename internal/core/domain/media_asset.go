package domain

import "github.com/google/uuid"

type StorageMode string

const (
	StorageModeExternal StorageMode = "EXTERNAL"
	StorageModeDatabase StorageMode = "DATABASE"
)

type MediaAssetType string

const (
	MediaAssetTypeImage    MediaAssetType = "IMAGE"
	MediaAssetTypeVideo    MediaAssetType = "VIDEO"
	MediaAssetTypeDocument MediaAssetType = "DOCUMENT"
)

type MediaPurposeType string

const (
	MediaPurposeAvatar        MediaPurposeType = "AVATAR"
	MediaPurposeQuestionImage MediaPurposeType = "QUESTION_IMAGE"
	MediaPurposeBadgeIcon     MediaPurposeType = "BADGE_ICON"
	MediaPurposeTrophyIcon    MediaPurposeType = "TROPHY_ICON"
	MediaPurposeBanner        MediaPurposeType = "BANNER"
)

type StorageProvider string

const (
	StorageProviderS3    StorageProvider = "S3"
	StorageProviderR2    StorageProvider = "R2"
	StorageProviderGCS   StorageProvider = "GCS"
	StorageProviderLocal StorageProvider = "LOCAL"
)

type MediaAsset struct {
	ID              uuid.UUID
	OwnerActorID    *uuid.UUID
	OwnerUserID     *uuid.UUID
	AssetType       MediaAssetType
	Purpose         MediaPurposeType
	StorageMode     StorageMode
	StorageProvider *StorageProvider
	Bucket          *string
	ObjectKey       *string
	BinaryData      []byte
	URL             *string
	ContentType     *string
	MimeType        string
	FileSizeBytes   *int64
	IsPublic        bool
	AuditFields
}

type MediaAssetUpdateInput struct {
	ID              uuid.UUID
	OwnerActorID    EntityField[*uuid.UUID]
	OwnerUserID     EntityField[*uuid.UUID]
	AssetType       EntityField[MediaAssetType]
	Purpose         EntityField[MediaPurposeType]
	StorageMode     EntityField[StorageMode]
	StorageProvider EntityField[*StorageProvider]
	Bucket          EntityField[*string]
	ObjectKey       EntityField[*string]
	BinaryData      EntityField[[]byte]
	URL             EntityField[*string]
	ContentType     EntityField[*string]
	MimeType        EntityField[string]
	FileSizeBytes   EntityField[*int64]
	IsPublic        EntityField[bool]
	ActorID         string
}

type MediaAssetQuery struct {
	Paging          Paging
	OwnerActorID    *uuid.UUID
	OwnerUserID     *uuid.UUID
	AssetType       *MediaAssetType
	Purpose         *MediaPurposeType
	StorageMode     *StorageMode
	StorageProvider *StorageProvider
	IsPublic        *bool
	Search          string
	SortBy          string
	SortOrder       string
}
