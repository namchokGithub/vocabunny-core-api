package content

import (
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/namchokGithub/vocabunny-core-api/internal/core/domain"
	"github.com/namchokGithub/vocabunny-core-api/internal/core/helper"
	"github.com/namchokGithub/vocabunny-core-api/internal/core/port"
)

type MediaAssetHandler struct{ service port.MediaAssetService }

func NewMediaAssetHandler(service port.MediaAssetService) *MediaAssetHandler {
	return &MediaAssetHandler{service: service}
}

func (h *MediaAssetHandler) Create(c echo.Context) error {
	var req CreateMediaAssetRequest
	if err := helper.BindAndValidate(c, &req); err != nil {
		return helper.RespondError(c, err)
	}

	ownerActorID, err := parseOptionalUUIDField(req.OwnerActorID, "owner_actor_id")
	if err != nil {
		return helper.RespondError(c, err)
	}
	ownerUserID, err := parseOptionalUUIDField(req.OwnerUserID, "owner_user_id")
	if err != nil {
		return helper.RespondError(c, err)
	}
	storageProvider, err := parseOptionalStorageProvider(req.StorageProvider)
	if err != nil {
		return helper.RespondError(c, err)
	}

	item, err := h.service.Create(c.Request().Context(), domain.MediaAssetCreateInput{
		OwnerActorID:    ownerActorID,
		OwnerUserID:     ownerUserID,
		AssetType:       domain.MediaAssetType(strings.TrimSpace(req.AssetType)),
		Purpose:         domain.MediaPurposeType(strings.TrimSpace(req.Purpose)),
		StorageMode:     domain.StorageMode(strings.TrimSpace(req.StorageMode)),
		StorageProvider: storageProvider,
		Bucket:          req.Bucket,
		ObjectKey:       req.ObjectKey,
		BinaryData:      req.BinaryData,
		URL:             req.URL,
		ContentType:     req.ContentType,
		MimeType:        req.MimeType,
		FileSizeBytes:   req.FileSizeBytes,
		IsPublic:        req.IsPublic,
		ActorID:         helper.ActorIDFromContext(c),
	})
	if err != nil {
		return helper.RespondError(c, err)
	}

	return helper.RespondSuccess(c, http.StatusCreated, toMediaAssetResponse(item))
}

func (h *MediaAssetHandler) Update(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return helper.RespondError(c, helper.BadRequest("invalid_media_asset_id", "media asset id must be a valid uuid", err))
	}

	var req UpdateMediaAssetRequest
	if err := helper.BindAndValidate(c, &req); err != nil {
		return helper.RespondError(c, err)
	}

	input := domain.MediaAssetUpdateInput{
		ID:      id,
		ActorID: helper.ActorIDFromContext(c),
	}

	if req.OwnerActorID != nil {
		value, err := parseOptionalUUIDField(req.OwnerActorID, "owner_actor_id")
		if err != nil {
			return helper.RespondError(c, err)
		}
		input.OwnerActorID = domain.NewEntityField(value)
	}
	if req.OwnerUserID != nil {
		value, err := parseOptionalUUIDField(req.OwnerUserID, "owner_user_id")
		if err != nil {
			return helper.RespondError(c, err)
		}
		input.OwnerUserID = domain.NewEntityField(value)
	}
	if req.AssetType != nil {
		input.AssetType = domain.NewEntityField(domain.MediaAssetType(strings.TrimSpace(*req.AssetType)))
	}
	if req.Purpose != nil {
		input.Purpose = domain.NewEntityField(domain.MediaPurposeType(strings.TrimSpace(*req.Purpose)))
	}
	if req.StorageMode != nil {
		input.StorageMode = domain.NewEntityField(domain.StorageMode(strings.TrimSpace(*req.StorageMode)))
	}
	if req.StorageProvider != nil {
		value, err := parseOptionalStorageProvider(req.StorageProvider)
		if err != nil {
			return helper.RespondError(c, err)
		}
		input.StorageProvider = domain.NewEntityField(value)
	}
	if req.Bucket != nil {
		input.Bucket = domain.NewEntityField(req.Bucket)
	}
	if req.ObjectKey != nil {
		input.ObjectKey = domain.NewEntityField(req.ObjectKey)
	}
	if req.BinaryData != nil {
		input.BinaryData = domain.NewEntityField(*req.BinaryData)
	}
	if req.URL != nil {
		input.URL = domain.NewEntityField(req.URL)
	}
	if req.ContentType != nil {
		input.ContentType = domain.NewEntityField(req.ContentType)
	}
	if req.MimeType != nil {
		input.MimeType = domain.NewEntityField(strings.TrimSpace(*req.MimeType))
	}
	if req.FileSizeBytes != nil {
		input.FileSizeBytes = domain.NewEntityField(req.FileSizeBytes)
	}
	if req.IsPublic != nil {
		input.IsPublic = domain.NewEntityField(*req.IsPublic)
	}

	item, err := h.service.Update(c.Request().Context(), input)
	if err != nil {
		return helper.RespondError(c, err)
	}

	return helper.RespondSuccess(c, http.StatusOK, toMediaAssetResponse(item))
}

func (h *MediaAssetHandler) Delete(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return helper.RespondError(c, helper.BadRequest("invalid_media_asset_id", "media asset id must be a valid uuid", err))
	}
	if err := h.service.Delete(c.Request().Context(), id, helper.ActorIDFromContext(c)); err != nil {
		return helper.RespondError(c, err)
	}
	return helper.RespondSuccess(c, http.StatusOK, map[string]string{"id": id.String(), "status": "deleted"})
}

func (h *MediaAssetHandler) FindByID(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return helper.RespondError(c, helper.BadRequest("invalid_media_asset_id", "media asset id must be a valid uuid", err))
	}
	item, err := h.service.FindByID(c.Request().Context(), id)
	if err != nil {
		return helper.RespondError(c, err)
	}
	return helper.RespondSuccess(c, http.StatusOK, toMediaAssetResponse(item))
}

func (h *MediaAssetHandler) FindAll(c echo.Context) error {
	query := domain.MediaAssetQuery{
		Paging: helper.BuildPaging(c),
		Search: strings.TrimSpace(c.QueryParam("search")),
	}
	query.SortBy, query.SortOrder = helper.BuildSort(c)

	if value := strings.TrimSpace(c.QueryParam("owner_actor_id")); value != "" {
		parsed, err := parseUUID(value, "owner_actor_id")
		if err != nil {
			return helper.RespondError(c, err)
		}
		query.OwnerActorID = &parsed
	}
	if value := strings.TrimSpace(c.QueryParam("owner_user_id")); value != "" {
		parsed, err := parseUUID(value, "owner_user_id")
		if err != nil {
			return helper.RespondError(c, err)
		}
		query.OwnerUserID = &parsed
	}
	if value := strings.TrimSpace(c.QueryParam("asset_type")); value != "" {
		parsed, err := parseMediaAssetType(value)
		if err != nil {
			return helper.RespondError(c, err)
		}
		query.AssetType = &parsed
	}
	if value := strings.TrimSpace(c.QueryParam("purpose")); value != "" {
		parsed, err := parseMediaPurposeType(value)
		if err != nil {
			return helper.RespondError(c, err)
		}
		query.Purpose = &parsed
	}
	if value := strings.TrimSpace(c.QueryParam("storage_mode")); value != "" {
		parsed, err := parseStorageMode(value)
		if err != nil {
			return helper.RespondError(c, err)
		}
		query.StorageMode = &parsed
	}
	if value := strings.TrimSpace(c.QueryParam("storage_provider")); value != "" {
		parsed, err := parseStorageProvider(value)
		if err != nil {
			return helper.RespondError(c, err)
		}
		query.StorageProvider = &parsed
	}
	if value := strings.TrimSpace(c.QueryParam("is_public")); value != "" {
		parsed := strings.EqualFold(value, "true")
		query.IsPublic = &parsed
	}

	result, err := h.service.FindAll(c.Request().Context(), query)
	if err != nil {
		return helper.RespondError(c, err)
	}

	items := make([]MediaAssetResponse, 0, len(result.Items))
	for _, item := range result.Items {
		items = append(items, toMediaAssetResponse(item))
	}

	return helper.RespondSuccess(c, http.StatusOK, ListResponse[MediaAssetResponse, domain.MediaAssetQuery]{
		Items: items,
		Paging: PagingResponse{
			Page:  result.Paging.Page,
			Limit: result.Paging.Limit,
			Total: result.Paging.Total,
		},
		Query: query,
	})
}

func parseOptionalStorageProvider(value *string) (*domain.StorageProvider, error) {
	if value == nil {
		return nil, nil
	}

	trimmed := strings.TrimSpace(*value)
	if trimmed == "" {
		return nil, nil
	}

	provider, err := parseStorageProvider(trimmed)
	if err != nil {
		return nil, err
	}

	return &provider, nil
}

func parseMediaAssetType(value string) (domain.MediaAssetType, error) {
	parsed := domain.MediaAssetType(strings.TrimSpace(value))
	switch parsed {
	case domain.MediaAssetTypeImage, domain.MediaAssetTypeVideo, domain.MediaAssetTypeDocument:
		return parsed, nil
	default:
		return "", helper.BadRequest("invalid_asset_type", "asset_type is invalid", nil)
	}
}

func parseMediaPurposeType(value string) (domain.MediaPurposeType, error) {
	parsed := domain.MediaPurposeType(strings.TrimSpace(value))
	switch parsed {
	case domain.MediaPurposeAvatar, domain.MediaPurposeQuestionImage, domain.MediaPurposeBadgeIcon, domain.MediaPurposeTrophyIcon, domain.MediaPurposeBanner:
		return parsed, nil
	default:
		return "", helper.BadRequest("invalid_purpose", "purpose is invalid", nil)
	}
}

func parseStorageMode(value string) (domain.StorageMode, error) {
	parsed := domain.StorageMode(strings.TrimSpace(value))
	switch parsed {
	case domain.StorageModeExternal, domain.StorageModeDatabase:
		return parsed, nil
	default:
		return "", helper.BadRequest("invalid_storage_mode", "storage_mode is invalid", nil)
	}
}

func parseStorageProvider(value string) (domain.StorageProvider, error) {
	parsed := domain.StorageProvider(strings.TrimSpace(value))
	switch parsed {
	case domain.StorageProviderS3, domain.StorageProviderR2, domain.StorageProviderGCS, domain.StorageProviderLocal:
		return parsed, nil
	default:
		return "", helper.BadRequest("invalid_storage_provider", "storage_provider is invalid", nil)
	}
}
