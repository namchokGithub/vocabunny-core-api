package repository

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

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) port.UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(ctx context.Context, user domain.User) (domain.User, error) {
	model := UserModel{
		ID:        user.ID,
		Code:      user.Code,
		Name:      user.Name,
		Email:     user.Email,
		Status:    string(user.Status),
		RoleID:    user.RoleID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		CreatedBy: user.CreatedBy,
		UpdatedBy: user.UpdatedBy,
	}

	if model.ID == uuid.Nil {
		model.ID = uuid.New()
	}

	if err := r.dbWithContext(ctx).Create(&model).Error; err != nil {
		return domain.User{}, helper.Internal("create_user_failed", "failed to create user", err)
	}

	return r.FindByID(ctx, model.ID)
}

func (r *userRepository) Update(ctx context.Context, input domain.UserUpdateInput) (domain.User, error) {
	var model UserModel
	if err := r.dbWithContext(ctx).Preload("Role").First(&model, "id = ?", input.ID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return domain.User{}, helper.NotFound("user_not_found", "user not found", err)
		}

		return domain.User{}, helper.Internal("find_user_failed", "failed to load user", err)
	}

	applyUserPatch(&model, input)
	if err := r.dbWithContext(ctx).Save(&model).Error; err != nil {
		return domain.User{}, helper.Internal("update_user_failed", "failed to update user", err)
	}

	return r.FindByID(ctx, model.ID)
}

func (r *userRepository) Delete(ctx context.Context, id uuid.UUID, actorID string) error {
	updates := map[string]interface{}{
		"updated_by": actorID,
		"updated_at": time.Now(),
	}

	if err := r.dbWithContext(ctx).Model(&UserModel{}).Where("id = ?", id).Updates(updates).Delete(&UserModel{}, "id = ?", id).Error; err != nil {
		return helper.Internal("delete_user_failed", "failed to delete user", err)
	}

	return nil
}

func (r *userRepository) FindByID(ctx context.Context, id uuid.UUID) (domain.User, error) {
	var model UserModel
	if err := r.dbWithContext(ctx).Preload("Role").First(&model, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return domain.User{}, helper.NotFound("user_not_found", "user not found", err)
		}

		return domain.User{}, helper.Internal("find_user_failed", "failed to find user", err)
	}

	return toDomainUser(model), nil
}

func (r *userRepository) FindByEmail(ctx context.Context, email string) (domain.User, error) {
	var model UserModel
	if err := r.dbWithContext(ctx).Preload("Role").First(&model, "LOWER(email) = ?", strings.ToLower(email)).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return domain.User{}, helper.NotFound("user_not_found", "user not found", err)
		}

		return domain.User{}, helper.Internal("find_user_failed", "failed to find user", err)
	}

	return toDomainUser(model), nil
}

func (r *userRepository) FindByCode(ctx context.Context, code string) (domain.User, error) {
	var model UserModel
	if err := r.dbWithContext(ctx).Preload("Role").First(&model, "UPPER(code) = ?", strings.ToUpper(code)).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return domain.User{}, helper.NotFound("user_not_found", "user not found", err)
		}

		return domain.User{}, helper.Internal("find_user_failed", "failed to find user", err)
	}

	return toDomainUser(model), nil
}

func (r *userRepository) FindAll(ctx context.Context, query domain.UserQuery) (domain.PageResult[domain.User], error) {
	db := r.dbWithContext(ctx).Model(&UserModel{}).Preload("Role")

	if query.Search != "" {
		search := "%" + strings.ToLower(query.Search) + "%"
		db = db.Where("LOWER(users.name) LIKE ? OR LOWER(users.email) LIKE ? OR LOWER(users.code) LIKE ?", search, search, search)
	}

	if query.RoleID != nil {
		db = db.Where("users.role_id = ?", *query.RoleID)
	}

	if query.Status != nil {
		db = db.Where("users.status = ?", string(*query.Status))
	}

	var total int64
	if err := db.Count(&total).Error; err != nil {
		return domain.PageResult[domain.User]{}, helper.Internal("count_users_failed", "failed to count users", err)
	}

	sortBy := query.SortBy
	switch sortBy {
	case "name", "email", "code", "status", "created_at", "updated_at":
	default:
		sortBy = "created_at"
	}

	sortOrder := "asc"
	if strings.ToLower(query.SortOrder) == "desc" {
		sortOrder = "desc"
	}

	var models []UserModel
	if err := db.Order("users." + sortBy + " " + sortOrder).
		Offset(query.Paging.Offset()).
		Limit(query.Paging.Limit).
		Find(&models).Error; err != nil {
		return domain.PageResult[domain.User]{}, helper.Internal("list_users_failed", "failed to list users", err)
	}

	items := make([]domain.User, 0, len(models))
	for _, model := range models {
		items = append(items, toDomainUser(model))
	}

	return domain.PageResult[domain.User]{
		Items: items,
		Paging: domain.Paging{
			Page:  query.Paging.Page,
			Limit: query.Paging.Limit,
			Total: total,
		},
	}, nil
}

func (r *userRepository) ExistsByEmail(ctx context.Context, email string, excludeID *uuid.UUID) (bool, error) {
	return r.exists(ctx, "LOWER(email) = ?", strings.ToLower(email), excludeID)
}

func (r *userRepository) ExistsByCode(ctx context.Context, code string, excludeID *uuid.UUID) (bool, error) {
	return r.exists(ctx, "UPPER(code) = ?", strings.ToUpper(code), excludeID)
}

func (r *userRepository) FindRoleByID(ctx context.Context, id uuid.UUID) (domain.Role, error) {
	var model RoleModel
	if err := r.dbWithContext(ctx).First(&model, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return domain.Role{}, helper.NotFound("role_not_found", "role not found", err)
		}

		return domain.Role{}, helper.Internal("find_role_failed", "failed to find role", err)
	}

	return toDomainRole(model), nil
}

func (r *userRepository) exists(ctx context.Context, query string, value interface{}, excludeID *uuid.UUID) (bool, error) {
	db := r.dbWithContext(ctx).Model(&UserModel{}).Where(query, value)
	if excludeID != nil {
		db = db.Where("id <> ?", *excludeID)
	}

	var count int64
	if err := db.Count(&count).Error; err != nil {
		return false, helper.Internal("exists_user_failed", "failed to check user uniqueness", err)
	}

	return count > 0, nil
}

func (r *userRepository) dbWithContext(ctx context.Context) *gorm.DB {
	tx, ok := helper.TxFromContext(ctx).(*gorm.DB)
	if ok && tx != nil {
		return tx.WithContext(ctx)
	}

	return r.db.WithContext(ctx)
}
