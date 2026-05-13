package content

import (
	"time"

	"github.com/google/uuid"
	"github.com/namchokGithub/vocabunny-core-api/internal/core/domain"
)

func toMediaAssetResponse(item domain.MediaAsset) MediaAssetResponse {
	return MediaAssetResponse{
		ID:              item.ID.String(),
		OwnerActorID:    uuidPointerToString(item.OwnerActorID),
		OwnerUserID:     uuidPointerToString(item.OwnerUserID),
		AssetType:       string(item.AssetType),
		Purpose:         string(item.Purpose),
		StorageMode:     string(item.StorageMode),
		StorageProvider: storageProviderPointerToString(item.StorageProvider),
		Bucket:          item.Bucket,
		ObjectKey:       item.ObjectKey,
		BinaryData:      item.BinaryData,
		URL:             item.URL,
		ContentType:     item.ContentType,
		MimeType:        item.MimeType,
		FileSizeBytes:   item.FileSizeBytes,
		IsPublic:        item.IsPublic,
		CreatedAt:       item.CreatedAt.Format(time.RFC3339),
		UpdatedAt:       item.UpdatedAt.Format(time.RFC3339),
		CreatedBy:       item.CreatedBy,
		UpdatedBy:       item.UpdatedBy,
	}
}

func uuidPointerToString(value *uuid.UUID) *string {
	if value == nil {
		return nil
	}
	result := value.String()
	return &result
}

func storageProviderPointerToString(value *domain.StorageProvider) *string {
	if value == nil {
		return nil
	}
	result := string(*value)
	return &result
}
