package content

import (
	"context"
	"strings"

	"github.com/google/uuid"
	"github.com/namchokGithub/vocabunny-core-api/internal/core/domain"
	"github.com/namchokGithub/vocabunny-core-api/internal/core/helper"
	"github.com/namchokGithub/vocabunny-core-api/internal/core/port"
)

type MediaAssetServiceDependencies struct {
	MediaAssetRepository port.MediaAssetRepository
	TxManager            port.TransactionManager
}

type mediaAssetService struct {
	mediaAssetRepository port.MediaAssetRepository
	txManager            port.TransactionManager
}

func NewMediaAssetService(deps MediaAssetServiceDependencies) port.MediaAssetService {
	return &mediaAssetService{
		mediaAssetRepository: deps.MediaAssetRepository,
		txManager:            deps.TxManager,
	}
}

func (s *mediaAssetService) Create(ctx context.Context, input domain.MediaAssetCreateInput) (domain.MediaAsset, error) {
	normalized, err := normalizeMediaAssetCreateInput(input)
	if err != nil {
		return domain.MediaAsset{}, err
	}
	if err := validateMediaAssetState(normalized.StorageMode, normalized.StorageProvider, normalized.Bucket, normalized.ObjectKey, normalized.BinaryData, normalized.MimeType, normalized.FileSizeBytes); err != nil {
		return domain.MediaAsset{}, err
	}

	var created domain.MediaAsset
	err = s.txManager.RunInTx(ctx, func(txCtx context.Context) error {
		created, err = s.mediaAssetRepository.Create(txCtx, domain.MediaAsset{
			OwnerActorID:    normalized.OwnerActorID,
			OwnerUserID:     normalized.OwnerUserID,
			AssetType:       normalized.AssetType,
			Purpose:         normalized.Purpose,
			StorageMode:     normalized.StorageMode,
			StorageProvider: normalized.StorageProvider,
			Bucket:          normalized.Bucket,
			ObjectKey:       normalized.ObjectKey,
			BinaryData:      normalized.BinaryData,
			URL:             normalized.URL,
			ContentType:     normalized.ContentType,
			MimeType:        normalized.MimeType,
			FileSizeBytes:   normalized.FileSizeBytes,
			IsPublic:        normalized.IsPublic,
			AuditFields: domain.AuditFields{
				CreatedBy: normalized.ActorID,
				UpdatedBy: normalized.ActorID,
			},
		})
		return err
	})
	return created, err
}

func (s *mediaAssetService) Update(ctx context.Context, input domain.MediaAssetUpdateInput) (domain.MediaAsset, error) {
	var updated domain.MediaAsset
	err := s.txManager.RunInTx(ctx, func(txCtx context.Context) error {
		current, err := s.mediaAssetRepository.FindByID(txCtx, input.ID)
		if err != nil {
			return err
		}

		normalized, err := normalizeMediaAssetUpdateInput(input)
		if err != nil {
			return err
		}

		effective := mergeMediaAsset(current, normalized)
		if err := validateMediaAssetState(effective.StorageMode, effective.StorageProvider, effective.Bucket, effective.ObjectKey, effective.BinaryData, effective.MimeType, effective.FileSizeBytes); err != nil {
			return err
		}

		// Keep storage fields consistent with the selected mode.
		// The repository should persist only the effective state, not half-updated data.
		input = alignMediaAssetUpdateWithMode(normalized, effective)

		updated, err = s.mediaAssetRepository.Update(txCtx, input)
		return err
	})
	return updated, err
}

func (s *mediaAssetService) Delete(ctx context.Context, id uuid.UUID, actorID string) error {
	return s.txManager.RunInTx(ctx, func(txCtx context.Context) error {
		if _, err := s.mediaAssetRepository.FindByID(txCtx, id); err != nil {
			return err
		}
		return s.mediaAssetRepository.Delete(txCtx, id, actorID)
	})
}

func (s *mediaAssetService) FindByID(ctx context.Context, id uuid.UUID) (domain.MediaAsset, error) {
	return s.mediaAssetRepository.FindByID(ctx, id)
}

func (s *mediaAssetService) FindAll(ctx context.Context, query domain.MediaAssetQuery) (domain.PageResult[domain.MediaAsset], error) {
	return s.mediaAssetRepository.FindAll(ctx, query)
}

func normalizeMediaAssetCreateInput(input domain.MediaAssetCreateInput) (domain.MediaAssetCreateInput, error) {
	input.MimeType = normalizeText(input.MimeType)
	input.Bucket = normalizeOptionalString(input.Bucket)
	input.ObjectKey = normalizeOptionalString(input.ObjectKey)
	input.URL = normalizeOptionalString(input.URL)
	input.ContentType = normalizeOptionalString(input.ContentType)

	if err := validateRequired(string(input.AssetType), "invalid_asset_type", "asset_type is required"); err != nil {
		return domain.MediaAssetCreateInput{}, err
	}
	if err := validateMediaAssetType(input.AssetType); err != nil {
		return domain.MediaAssetCreateInput{}, err
	}
	if err := validateRequired(string(input.Purpose), "invalid_purpose", "purpose is required"); err != nil {
		return domain.MediaAssetCreateInput{}, err
	}
	if err := validateMediaPurposeType(input.Purpose); err != nil {
		return domain.MediaAssetCreateInput{}, err
	}
	if err := validateRequired(string(input.StorageMode), "invalid_storage_mode", "storage_mode is required"); err != nil {
		return domain.MediaAssetCreateInput{}, err
	}
	if err := validateRequired(input.MimeType, "invalid_mime_type", "mime_type is required"); err != nil {
		return domain.MediaAssetCreateInput{}, err
	}

	return input, nil
}

func normalizeMediaAssetUpdateInput(input domain.MediaAssetUpdateInput) (domain.MediaAssetUpdateInput, error) {
	if input.AssetType.Set {
		if err := validateRequired(string(input.AssetType.Value), "invalid_asset_type", "asset_type cannot be empty"); err != nil {
			return input, err
		}
		if err := validateMediaAssetType(input.AssetType.Value); err != nil {
			return input, err
		}
	}
	if input.Purpose.Set {
		if err := validateRequired(string(input.Purpose.Value), "invalid_purpose", "purpose cannot be empty"); err != nil {
			return input, err
		}
		if err := validateMediaPurposeType(input.Purpose.Value); err != nil {
			return input, err
		}
	}
	if input.StorageMode.Set {
		if err := validateRequired(string(input.StorageMode.Value), "invalid_storage_mode", "storage_mode cannot be empty"); err != nil {
			return input, err
		}
	}
	if input.Bucket.Set {
		input.Bucket.Value = normalizeOptionalString(input.Bucket.Value)
	}
	if input.ObjectKey.Set {
		input.ObjectKey.Value = normalizeOptionalString(input.ObjectKey.Value)
	}
	if input.URL.Set {
		input.URL.Value = normalizeOptionalString(input.URL.Value)
	}
	if input.ContentType.Set {
		input.ContentType.Value = normalizeOptionalString(input.ContentType.Value)
	}
	if input.MimeType.Set {
		input.MimeType.Value = normalizeText(input.MimeType.Value)
		if err := validateRequired(input.MimeType.Value, "invalid_mime_type", "mime_type cannot be empty"); err != nil {
			return input, err
		}
	}
	if input.FileSizeBytes.Set && input.FileSizeBytes.Value != nil && *input.FileSizeBytes.Value < 0 {
		return input, helper.BadRequest("invalid_file_size_bytes", "file_size_bytes must be zero or greater", nil)
	}

	return input, nil
}

func mergeMediaAsset(current domain.MediaAsset, input domain.MediaAssetUpdateInput) domain.MediaAsset {
	effective := current
	if input.OwnerActorID.Set {
		effective.OwnerActorID = input.OwnerActorID.Value
	}
	if input.OwnerUserID.Set {
		effective.OwnerUserID = input.OwnerUserID.Value
	}
	if input.AssetType.Set {
		effective.AssetType = input.AssetType.Value
	}
	if input.Purpose.Set {
		effective.Purpose = input.Purpose.Value
	}
	if input.StorageMode.Set {
		effective.StorageMode = input.StorageMode.Value
	}
	if input.StorageProvider.Set {
		effective.StorageProvider = input.StorageProvider.Value
	}
	if input.Bucket.Set {
		effective.Bucket = input.Bucket.Value
	}
	if input.ObjectKey.Set {
		effective.ObjectKey = input.ObjectKey.Value
	}
	if input.BinaryData.Set {
		effective.BinaryData = input.BinaryData.Value
	}
	if input.URL.Set {
		effective.URL = input.URL.Value
	}
	if input.ContentType.Set {
		effective.ContentType = input.ContentType.Value
	}
	if input.MimeType.Set {
		effective.MimeType = input.MimeType.Value
	}
	if input.FileSizeBytes.Set {
		effective.FileSizeBytes = input.FileSizeBytes.Value
	}
	if input.IsPublic.Set {
		effective.IsPublic = input.IsPublic.Value
	}
	return effective
}

func alignMediaAssetUpdateWithMode(input domain.MediaAssetUpdateInput, effective domain.MediaAsset) domain.MediaAssetUpdateInput {
	if effective.StorageMode == domain.StorageModeExternal {
		input.BinaryData = domain.NewEntityField([]byte(nil))
	}
	if effective.StorageMode == domain.StorageModeDatabase {
		input.StorageProvider = domain.NewEntityField[*domain.StorageProvider](nil)
		input.Bucket = domain.NewEntityField[*string](nil)
		input.ObjectKey = domain.NewEntityField[*string](nil)
	}
	return input
}

func validateMediaAssetState(
	storageMode domain.StorageMode,
	storageProvider *domain.StorageProvider,
	bucket *string,
	objectKey *string,
	binaryData []byte,
	mimeType string,
	fileSizeBytes *int64,
) error {
	if err := validateMediaStorageMode(storageMode); err != nil {
		return err
	}
	if err := validateMediaMimeType(mimeType); err != nil {
		return err
	}
	if fileSizeBytes != nil && *fileSizeBytes < 0 {
		return helper.BadRequest("invalid_file_size_bytes", "file_size_bytes must be zero or greater", nil)
	}

	switch storageMode {
	case domain.StorageModeExternal:
		if storageProvider == nil || strings.TrimSpace(string(*storageProvider)) == "" {
			return helper.BadRequest("invalid_storage_provider", "storage_provider is required when storage_mode is EXTERNAL", nil)
		}
		if err := validateStorageProvider(*storageProvider); err != nil {
			return err
		}
		if objectKey == nil || strings.TrimSpace(*objectKey) == "" {
			return helper.BadRequest("invalid_object_key", "object_key is required when storage_mode is EXTERNAL", nil)
		}
	case domain.StorageModeDatabase:
		if len(binaryData) == 0 {
			return helper.BadRequest("invalid_binary_data", "binary_data is required when storage_mode is DATABASE", nil)
		}
		if storageProvider != nil || bucket != nil || objectKey != nil {
			return helper.BadRequest("invalid_storage_fields", "storage_provider, bucket and object_key must be empty when storage_mode is DATABASE", nil)
		}
	default:
		return helper.BadRequest("invalid_storage_mode", "storage_mode is invalid", nil)
	}

	return nil
}

func validateMediaAssetType(value domain.MediaAssetType) error {
	switch value {
	case domain.MediaAssetTypeImage, domain.MediaAssetTypeVideo, domain.MediaAssetTypeDocument:
		return nil
	default:
		return helper.BadRequest("invalid_asset_type", "asset_type is invalid", nil)
	}
}

func validateMediaPurposeType(value domain.MediaPurposeType) error {
	switch value {
	case domain.MediaPurposeAvatar, domain.MediaPurposeQuestionImage, domain.MediaPurposeBadgeIcon, domain.MediaPurposeTrophyIcon, domain.MediaPurposeBanner:
		return nil
	default:
		return helper.BadRequest("invalid_purpose", "purpose is invalid", nil)
	}
}

func validateMediaStorageMode(value domain.StorageMode) error {
	if value != domain.StorageModeExternal && value != domain.StorageModeDatabase {
		return helper.BadRequest("invalid_storage_mode", "storage_mode is invalid", nil)
	}
	return nil
}

func validateStorageProvider(value domain.StorageProvider) error {
	switch value {
	case domain.StorageProviderS3, domain.StorageProviderR2, domain.StorageProviderGCS, domain.StorageProviderLocal:
		return nil
	default:
		return helper.BadRequest("invalid_storage_provider", "storage_provider is invalid", nil)
	}
}

func validateMediaMimeType(value string) error {
	if strings.TrimSpace(value) == "" {
		return helper.BadRequest("invalid_mime_type", "mime_type is required", nil)
	}
	return nil
}

func normalizeOptionalString(value *string) *string {
	if value == nil {
		return nil
	}
	trimmed := strings.TrimSpace(*value)
	if trimmed == "" {
		return nil
	}
	return &trimmed
}
