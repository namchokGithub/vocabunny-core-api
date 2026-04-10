package content

type CreateMediaAssetRequest struct {
	OwnerActorID    *string `json:"owner_actor_id" validate:"omitempty,uuid"`
	OwnerUserID     *string `json:"owner_user_id" validate:"omitempty,uuid"`
	AssetType       string  `json:"asset_type" validate:"required,oneof=IMAGE VIDEO DOCUMENT"`
	Purpose         string  `json:"purpose" validate:"required,oneof=AVATAR QUESTION_IMAGE BADGE_ICON TROPHY_ICON BANNER"`
	StorageMode     string  `json:"storage_mode" validate:"required,oneof=EXTERNAL DATABASE"`
	StorageProvider *string `json:"storage_provider" validate:"omitempty,oneof=S3 R2 GCS LOCAL"`
	Bucket          *string `json:"bucket"`
	ObjectKey       *string `json:"object_key"`
	BinaryData      []byte  `json:"binary_data"`
	URL             *string `json:"url" validate:"omitempty,url"`
	ContentType     *string `json:"content_type" validate:"omitempty,max=255"`
	MimeType        string  `json:"mime_type" validate:"required,max=255"`
	FileSizeBytes   *int64  `json:"file_size_bytes" validate:"omitempty,gte=0"`
	IsPublic        bool    `json:"is_public"`
}

type UpdateMediaAssetRequest struct {
	OwnerActorID    *string `json:"owner_actor_id" validate:"omitempty,uuid"`
	OwnerUserID     *string `json:"owner_user_id" validate:"omitempty,uuid"`
	AssetType       *string `json:"asset_type" validate:"omitempty,oneof=IMAGE VIDEO DOCUMENT"`
	Purpose         *string `json:"purpose" validate:"omitempty,oneof=AVATAR QUESTION_IMAGE BADGE_ICON TROPHY_ICON BANNER"`
	StorageMode     *string `json:"storage_mode" validate:"omitempty,oneof=EXTERNAL DATABASE"`
	StorageProvider *string `json:"storage_provider" validate:"omitempty,oneof=S3 R2 GCS LOCAL"`
	Bucket          *string `json:"bucket"`
	ObjectKey       *string `json:"object_key"`
	BinaryData      *[]byte `json:"binary_data"`
	URL             *string `json:"url" validate:"omitempty,url"`
	ContentType     *string `json:"content_type" validate:"omitempty,max=255"`
	MimeType        *string `json:"mime_type" validate:"omitempty,max=255"`
	FileSizeBytes   *int64  `json:"file_size_bytes" validate:"omitempty,gte=0"`
	IsPublic        *bool   `json:"is_public"`
}

type MediaAssetResponse struct {
	ID              string  `json:"id"`
	OwnerActorID    *string `json:"owner_actor_id,omitempty"`
	OwnerUserID     *string `json:"owner_user_id,omitempty"`
	AssetType       string  `json:"asset_type"`
	Purpose         string  `json:"purpose"`
	StorageMode     string  `json:"storage_mode"`
	StorageProvider *string `json:"storage_provider,omitempty"`
	Bucket          *string `json:"bucket,omitempty"`
	ObjectKey       *string `json:"object_key,omitempty"`
	BinaryData      []byte  `json:"binary_data,omitempty"`
	URL             *string `json:"url,omitempty"`
	ContentType     *string `json:"content_type,omitempty"`
	MimeType        string  `json:"mime_type"`
	FileSizeBytes   *int64  `json:"file_size_bytes,omitempty"`
	IsPublic        bool    `json:"is_public"`
	CreatedAt       string  `json:"created_at"`
	UpdatedAt       string  `json:"updated_at"`
	CreatedBy       string  `json:"created_by"`
	UpdatedBy       string  `json:"updated_by"`
}
