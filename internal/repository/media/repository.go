package media

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/namchokGithub/vocabunny-core-api/internal/core/domain"
	"github.com/namchokGithub/vocabunny-core-api/internal/core/helper"
	"github.com/namchokGithub/vocabunny-core-api/internal/core/port"
	"gorm.io/gorm"
)

type Dependencies struct {
	DB *gorm.DB
}

type Repository struct {
	MediaAsset port.MediaAssetRepository
}

func NewRepository(deps Dependencies) *Repository {
	base := &baseRepository{db: deps.DB}
	return &Repository{
		MediaAsset: &mediaAssetRepository{baseRepository: base},
	}
}

type baseRepository struct {
	db *gorm.DB
}

// dbWithContext keeps repository methods transaction-safe.
// If an upper layer already started a transaction, repository calls inside the
// same request automatically reuse that transaction.
func (r *baseRepository) dbWithContext(ctx context.Context) *gorm.DB {
	tx, ok := helper.TxFromContext(ctx).(*gorm.DB)
	if ok && tx != nil {
		return tx.WithContext(ctx)
	}
	return r.db.WithContext(ctx)
}

type mediaAssetRepository struct {
	*baseRepository
}

func (r *mediaAssetRepository) Create(ctx context.Context, asset domain.MediaAsset) (domain.MediaAsset, error) {
	now := time.Now()
	model := MediaAssetModel{
		ID:              asset.ID,
		OwnerActorID:    asset.OwnerActorID,
		OwnerUserID:     asset.OwnerUserID,
		AssetType:       string(asset.AssetType),
		Purpose:         string(asset.Purpose),
		StorageMode:     string(asset.StorageMode),
		StorageProvider: storageProviderToModel(asset.StorageProvider),
		Bucket:          asset.Bucket,
		ObjectKey:       asset.ObjectKey,
		BinaryData:      asset.BinaryData,
		URL:             asset.URL,
		ContentType:     asset.ContentType,
		MimeType:        asset.MimeType,
		FileSizeBytes:   asset.FileSizeBytes,
		IsPublic:        asset.IsPublic,
		CreatedBy:       asset.CreatedBy,
		UpdatedBy:       asset.UpdatedBy,
		CreatedAt:       now,
		UpdatedAt:       now,
	}
	if model.ID == uuid.Nil {
		model.ID = uuid.New()
	}

	if err := r.dbWithContext(ctx).Create(&model).Error; err != nil {
		return domain.MediaAsset{}, helper.Internal("create_media_asset_failed", "failed to create media asset", err)
	}
	return r.FindByID(ctx, model.ID)
}

func (r *mediaAssetRepository) Update(ctx context.Context, input domain.MediaAssetUpdateInput) (domain.MediaAsset, error) {
	var model MediaAssetModel
	if err := r.dbWithContext(ctx).First(&model, "id = ?", input.ID).Error; err != nil {
		return domain.MediaAsset{}, mapGormNotFound(err, "media_asset_not_found", "media asset not found", "find_media_asset_failed", "failed to load media asset")
	}

	// Each field is updated only when Set=true.
	// This lets callers distinguish between "do not touch this column" and
	// "intentionally clear this nullable column".
	if input.OwnerActorID.Set {
		model.OwnerActorID = input.OwnerActorID.Value
	}
	if input.OwnerUserID.Set {
		model.OwnerUserID = input.OwnerUserID.Value
	}
	if input.AssetType.Set {
		model.AssetType = string(input.AssetType.Value)
	}
	if input.Purpose.Set {
		model.Purpose = string(input.Purpose.Value)
	}
	if input.StorageMode.Set {
		model.StorageMode = string(input.StorageMode.Value)
	}
	if input.StorageProvider.Set {
		model.StorageProvider = storageProviderToModel(input.StorageProvider.Value)
	}
	if input.Bucket.Set {
		model.Bucket = input.Bucket.Value
	}
	if input.ObjectKey.Set {
		model.ObjectKey = input.ObjectKey.Value
	}
	if input.BinaryData.Set {
		model.BinaryData = input.BinaryData.Value
	}
	if input.URL.Set {
		model.URL = input.URL.Value
	}
	if input.ContentType.Set {
		model.ContentType = input.ContentType.Value
	}
	if input.MimeType.Set {
		model.MimeType = input.MimeType.Value
	}
	if input.FileSizeBytes.Set {
		model.FileSizeBytes = input.FileSizeBytes.Value
	}
	if input.IsPublic.Set {
		model.IsPublic = input.IsPublic.Value
	}
	model.UpdatedBy = input.ActorID
	model.UpdatedAt = time.Now()

	if err := r.dbWithContext(ctx).Save(&model).Error; err != nil {
		return domain.MediaAsset{}, helper.Internal("update_media_asset_failed", "failed to update media asset", err)
	}
	return r.FindByID(ctx, model.ID)
}

func (r *mediaAssetRepository) Delete(ctx context.Context, id uuid.UUID, actorID string) error {
	if err := r.dbWithContext(ctx).Model(&MediaAssetModel{}).Where("id = ?", id).Updates(map[string]any{
		"updated_by": actorID,
		"updated_at": time.Now(),
	}).Delete(&MediaAssetModel{}, "id = ?", id).Error; err != nil {
		return helper.Internal("delete_media_asset_failed", "failed to delete media asset", err)
	}
	return nil
}

func (r *mediaAssetRepository) FindByID(ctx context.Context, id uuid.UUID) (domain.MediaAsset, error) {
	var model MediaAssetModel
	if err := r.dbWithContext(ctx).First(&model, "id = ?", id).Error; err != nil {
		return domain.MediaAsset{}, mapGormNotFound(err, "media_asset_not_found", "media asset not found", "find_media_asset_failed", "failed to find media asset")
	}
	return toDomainMediaAsset(model), nil
}

func (r *mediaAssetRepository) FindAll(ctx context.Context, query domain.MediaAssetQuery) (domain.PageResult[domain.MediaAsset], error) {
	db := r.dbWithContext(ctx).Model(&MediaAssetModel{})

	if query.OwnerActorID != nil {
		db = db.Where("owner_actor_id = ?", *query.OwnerActorID)
	}
	if query.OwnerUserID != nil {
		db = db.Where("owner_user_id = ?", *query.OwnerUserID)
	}
	if query.AssetType != nil {
		db = db.Where("asset_type = ?", string(*query.AssetType))
	}
	if query.Purpose != nil {
		db = db.Where("purpose = ?", string(*query.Purpose))
	}
	if query.StorageMode != nil {
		db = db.Where("storage_mode = ?", string(*query.StorageMode))
	}
	if query.StorageProvider != nil {
		db = db.Where("storage_provider = ?", string(*query.StorageProvider))
	}
	if query.IsPublic != nil {
		db = db.Where("is_public = ?", *query.IsPublic)
	}
	if query.Search != "" {
		search := "%" + strings.ToLower(strings.TrimSpace(query.Search)) + "%"
		db = db.Where(`
			LOWER(COALESCE(bucket, '')) LIKE ?
			OR LOWER(COALESCE(object_key, '')) LIKE ?
			OR LOWER(COALESCE(url, '')) LIKE ?
			OR LOWER(COALESCE(content_type, '')) LIKE ?
			OR LOWER(mime_type) LIKE ?
		`, search, search, search, search, search)
	}

	var total int64
	if err := db.Count(&total).Error; err != nil {
		return domain.PageResult[domain.MediaAsset]{}, helper.Internal("count_media_assets_failed", "failed to count media assets", err)
	}

	sortBy := safeSort(query.SortBy, []string{"asset_type", "purpose", "storage_mode", "mime_type", "file_size_bytes", "is_public", "created_at", "updated_at"}, "created_at")
	sortOrder := safeSortOrder(query.SortOrder)

	var models []MediaAssetModel
	if err := db.Order("tbl_media_assets." + sortBy + " " + sortOrder).
		Order("tbl_media_assets.id ASC").
		Offset(query.Paging.Offset()).
		Limit(query.Paging.Limit).
		Find(&models).Error; err != nil {
		return domain.PageResult[domain.MediaAsset]{}, helper.Internal("list_media_assets_failed", "failed to list media assets", err)
	}

	items := make([]domain.MediaAsset, 0, len(models))
	for _, model := range models {
		items = append(items, toDomainMediaAsset(model))
	}

	return domain.PageResult[domain.MediaAsset]{
		Items: items,
		Paging: domain.Paging{
			Page:  query.Paging.Page,
			Limit: query.Paging.Limit,
			Total: total,
		},
	}, nil
}

func toDomainMediaAsset(model MediaAssetModel) domain.MediaAsset {
	return domain.MediaAsset{
		ID:              model.ID,
		OwnerActorID:    model.OwnerActorID,
		OwnerUserID:     model.OwnerUserID,
		AssetType:       domain.MediaAssetType(model.AssetType),
		Purpose:         domain.MediaPurposeType(model.Purpose),
		StorageMode:     domain.StorageMode(model.StorageMode),
		StorageProvider: storageProviderToDomain(model.StorageProvider),
		Bucket:          model.Bucket,
		ObjectKey:       model.ObjectKey,
		BinaryData:      model.BinaryData,
		URL:             model.URL,
		ContentType:     model.ContentType,
		MimeType:        model.MimeType,
		FileSizeBytes:   model.FileSizeBytes,
		IsPublic:        model.IsPublic,
		AuditFields: domain.AuditFields{
			CreatedAt: model.CreatedAt,
			UpdatedAt: model.UpdatedAt,
			DeletedAt: deletedAtToPointer(model.DeletedAt),
			CreatedBy: model.CreatedBy,
			UpdatedBy: model.UpdatedBy,
		},
	}
}

func storageProviderToModel(provider *domain.StorageProvider) *string {
	if provider == nil {
		return nil
	}
	value := string(*provider)
	return &value
}

func storageProviderToDomain(provider *string) *domain.StorageProvider {
	if provider == nil {
		return nil
	}
	value := domain.StorageProvider(*provider)
	return &value
}

func deletedAtToPointer(deletedAt gorm.DeletedAt) *time.Time {
	if !deletedAt.Valid {
		return nil
	}
	value := deletedAt.Time
	return &value
}

func mapGormNotFound(err error, notFoundCode, notFoundMessage, internalCode, internalMessage string) error {
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return helper.NotFound(notFoundCode, notFoundMessage, err)
	}
	return helper.Internal(internalCode, internalMessage, err)
}

func safeSort(value string, allowed []string, fallback string) string {
	for _, item := range allowed {
		if value == item {
			return value
		}
	}
	return fallback
}

func safeSortOrder(value string) string {
	if strings.ToLower(value) == "desc" {
		return "desc"
	}
	return "asc"
}
