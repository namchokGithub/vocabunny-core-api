package content_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/namchokGithub/vocabunny-core-api/internal/core/domain"
	content "github.com/namchokGithub/vocabunny-core-api/internal/handler/content"
)

func TestMediaAssetHandlerCreate(t *testing.T) {
	t.Parallel()

	ownerUserID := uuid.New()
	assetID := uuid.New()
	service := &mediaAssetServiceStub{
		createFn: func(ctx context.Context, input domain.MediaAssetCreateInput) (domain.MediaAsset, error) {
			if input.ActorID != "actor-123" {
				t.Fatalf("expected actor id actor-123, got %q", input.ActorID)
			}
			if input.OwnerUserID == nil || *input.OwnerUserID != ownerUserID {
				t.Fatalf("unexpected owner user id: %#v", input.OwnerUserID)
			}
			if input.AssetType != domain.MediaAssetTypeImage {
				t.Fatalf("expected asset type IMAGE, got %q", input.AssetType)
			}
			if input.Purpose != domain.MediaPurposeAvatar {
				t.Fatalf("expected purpose AVATAR, got %q", input.Purpose)
			}
			if input.StorageMode != domain.StorageModeExternal {
				t.Fatalf("expected storage mode EXTERNAL, got %q", input.StorageMode)
			}
			if input.StorageProvider == nil || *input.StorageProvider != domain.StorageProviderS3 {
				t.Fatalf("unexpected storage provider: %#v", input.StorageProvider)
			}
			if input.ObjectKey == nil || *input.ObjectKey != "avatars/user-1.png" {
				t.Fatalf("unexpected object key: %#v", input.ObjectKey)
			}
			return domain.MediaAsset{
				ID:              assetID,
				OwnerUserID:     input.OwnerUserID,
				AssetType:       input.AssetType,
				Purpose:         input.Purpose,
				StorageMode:     input.StorageMode,
				StorageProvider: input.StorageProvider,
				ObjectKey:       input.ObjectKey,
				MimeType:        input.MimeType,
				IsPublic:        input.IsPublic,
				AuditFields:     testAuditFields("actor-123"),
			}, nil
		},
	}

	handler := content.NewMediaAssetHandler(service)
	body := `{
		"owner_user_id":"` + ownerUserID.String() + `",
		"asset_type":"IMAGE",
		"purpose":"AVATAR",
		"storage_mode":"EXTERNAL",
		"storage_provider":"S3",
		"object_key":"avatars/user-1.png",
		"mime_type":"image/png",
		"is_public":true
	}`
	rec := performJSONRequest(t, http.MethodPost, "/media-assets", body, func(c echo.Context) error {
		c.Request().Header.Set("X-Actor-ID", "actor-123")
		return handler.Create(c)
	})

	if rec.Code != http.StatusCreated {
		t.Fatalf("expected status %d, got %d", http.StatusCreated, rec.Code)
	}
}

func TestMediaAssetHandlerUpdateInvalidID(t *testing.T) {
	t.Parallel()

	called := false
	service := &mediaAssetServiceStub{
		updateFn: func(ctx context.Context, input domain.MediaAssetUpdateInput) (domain.MediaAsset, error) {
			called = true
			return domain.MediaAsset{}, nil
		},
	}

	handler := content.NewMediaAssetHandler(service)
	rec := performJSONRequest(t, http.MethodPut, "/media-assets/not-a-uuid", `{"mime_type":"image/webp"}`, func(c echo.Context) error {
		c.SetParamNames("id")
		c.SetParamValues("not-a-uuid")
		return handler.Update(c)
	})

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, rec.Code)
	}
	if called {
		t.Fatal("expected service update not to be called")
	}
}

func TestMediaAssetHandlerFindAllBuildsQuery(t *testing.T) {
	t.Parallel()

	ownerActorID := uuid.New()
	var captured domain.MediaAssetQuery
	service := &mediaAssetServiceStub{
		findAllFn: func(ctx context.Context, query domain.MediaAssetQuery) (domain.PageResult[domain.MediaAsset], error) {
			captured = query
			return domain.PageResult[domain.MediaAsset]{
				Items:  []domain.MediaAsset{},
				Paging: domain.Paging{Page: query.Paging.Page, Limit: query.Paging.Limit, Total: 0},
			}, nil
		},
	}

	handler := content.NewMediaAssetHandler(service)
	target := "/media-assets?search=%20avatar%20&owner_actor_id=" + ownerActorID.String() + "&asset_type=IMAGE&purpose=AVATAR&storage_mode=EXTERNAL&storage_provider=S3&is_public=true&sort_by=created_at&sort_order=desc&page=2&limit=5"
	rec := performJSONRequest(t, http.MethodGet, target, "", func(c echo.Context) error {
		return handler.FindAll(c)
	})

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rec.Code)
	}
	if captured.Search != "avatar" {
		t.Fatalf("expected search avatar, got %q", captured.Search)
	}
	if captured.OwnerActorID == nil || *captured.OwnerActorID != ownerActorID {
		t.Fatalf("unexpected owner actor id: %#v", captured.OwnerActorID)
	}
	if captured.AssetType == nil || *captured.AssetType != domain.MediaAssetTypeImage {
		t.Fatalf("unexpected asset type: %#v", captured.AssetType)
	}
	if captured.Purpose == nil || *captured.Purpose != domain.MediaPurposeAvatar {
		t.Fatalf("unexpected purpose: %#v", captured.Purpose)
	}
	if captured.StorageMode == nil || *captured.StorageMode != domain.StorageModeExternal {
		t.Fatalf("unexpected storage mode: %#v", captured.StorageMode)
	}
	if captured.StorageProvider == nil || *captured.StorageProvider != domain.StorageProviderS3 {
		t.Fatalf("unexpected storage provider: %#v", captured.StorageProvider)
	}
	if captured.IsPublic == nil || !*captured.IsPublic {
		t.Fatalf("unexpected is_public filter: %#v", captured.IsPublic)
	}
	if captured.SortBy != "created_at" || captured.SortOrder != "desc" {
		t.Fatalf("unexpected sort: %s %s", captured.SortBy, captured.SortOrder)
	}
	if captured.Paging.Page != 2 || captured.Paging.Limit != 5 {
		t.Fatalf("unexpected paging: %#v", captured.Paging)
	}
}
