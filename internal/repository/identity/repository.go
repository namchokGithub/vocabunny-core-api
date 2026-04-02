package identity

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/namchokGithub/vocabunny-core-api/internal/core/domain"
	"github.com/namchokGithub/vocabunny-core-api/internal/core/helper"
	"github.com/namchokGithub/vocabunny-core-api/internal/core/port"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type Dependencies struct {
	DB *gorm.DB
}

type Repository struct {
	User         port.UserRepository
	Role         port.RoleRepository
	AuthIdentity port.AuthIdentityRepository
}

func NewRepository(deps Dependencies) *Repository {
	base := &baseRepository{db: deps.DB}
	return &Repository{
		User:         &userRepository{baseRepository: base},
		Role:         &roleRepository{baseRepository: base},
		AuthIdentity: &authIdentityRepository{baseRepository: base},
	}
}

type baseRepository struct {
	db *gorm.DB
}

func (r *baseRepository) dbWithContext(ctx context.Context) *gorm.DB {
	tx, ok := helper.TxFromContext(ctx).(*gorm.DB)
	if ok && tx != nil {
		return tx.WithContext(ctx)
	}
	return r.db.WithContext(ctx)
}

type userRepository struct {
	*baseRepository
}

func (r *userRepository) Create(ctx context.Context, user domain.User) (domain.User, error) {
	now := time.Now()
	model := UserModel{
		ID:          user.ID,
		Email:       user.Email,
		Username:    user.Username,
		DisplayName: user.DisplayName,
		AvatarID:    user.AvatarID,
		Status:      string(user.Status),
		CreatedBy:   user.CreatedBy,
		UpdatedBy:   user.UpdatedBy,
		CreatedAt:   now,
		UpdatedAt:   now,
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
	if err := r.dbWithContext(ctx).First(&model, "id = ?", input.ID).Error; err != nil {
		return domain.User{}, mapGormNotFound(err, "user_not_found", "user not found", "find_user_failed", "failed to load user")
	}

	if input.Email.Set {
		model.Email = input.Email.Value
	}
	if input.Username.Set {
		model.Username = input.Username.Value
	}
	if input.DisplayName.Set {
		model.DisplayName = input.DisplayName.Value
	}
	if input.AvatarID.Set {
		model.AvatarID = input.AvatarID.Value
	}
	if input.Status.Set {
		model.Status = string(input.Status.Value)
	}
	model.UpdatedBy = input.ActorID
	model.UpdatedAt = time.Now()

	if err := r.dbWithContext(ctx).Save(&model).Error; err != nil {
		return domain.User{}, helper.Internal("update_user_failed", "failed to update user", err)
	}

	return r.FindByID(ctx, model.ID)
}

func (r *userRepository) Delete(ctx context.Context, id uuid.UUID, actorID string) error {
	now := time.Now()
	if err := r.dbWithContext(ctx).Model(&UserModel{}).Where("id = ?", id).Updates(map[string]any{
		"updated_by": actorID,
		"updated_at": now,
		"status":     string(domain.UserStatusDeleted),
	}).Delete(&UserModel{}, "id = ?", id).Error; err != nil {
		return helper.Internal("delete_user_failed", "failed to delete user", err)
	}
	return nil
}

func (r *userRepository) FindByID(ctx context.Context, id uuid.UUID) (domain.User, error) {
	var model UserModel
	if err := r.dbWithContext(ctx).First(&model, "id = ?", id).Error; err != nil {
		return domain.User{}, mapGormNotFound(err, "user_not_found", "user not found", "find_user_failed", "failed to find user")
	}
	return r.loadUserRelations(ctx, model)
}

func (r *userRepository) FindByEmail(ctx context.Context, email string) (domain.User, error) {
	var model UserModel
	if err := r.dbWithContext(ctx).First(&model, "LOWER(email) = ?", strings.ToLower(email)).Error; err != nil {
		return domain.User{}, mapGormNotFound(err, "user_not_found", "user not found", "find_user_failed", "failed to find user")
	}
	return r.loadUserRelations(ctx, model)
}

func (r *userRepository) FindByUsername(ctx context.Context, username string) (domain.User, error) {
	var model UserModel
	if err := r.dbWithContext(ctx).First(&model, "LOWER(username) = ?", strings.ToLower(username)).Error; err != nil {
		return domain.User{}, mapGormNotFound(err, "user_not_found", "user not found", "find_user_failed", "failed to find user")
	}
	return r.loadUserRelations(ctx, model)
}

func (r *userRepository) FindBySubject(ctx context.Context, subject string) (domain.User, error) {
	parsed, err := uuid.Parse(subject)
	if err == nil {
		return r.FindByID(ctx, parsed)
	}

	user, err := r.FindByEmail(ctx, subject)
	if err == nil {
		return user, nil
	}

	return r.FindByUsername(ctx, subject)
}

func (r *userRepository) FindAll(ctx context.Context, query domain.UserQuery) (domain.PageResult[domain.User], error) {
	db := r.dbWithContext(ctx).Model(&UserModel{})
	if query.Search != "" {
		search := "%" + strings.ToLower(query.Search) + "%"
		db = db.Where("LOWER(email) LIKE ? OR LOWER(username) LIKE ? OR LOWER(display_name) LIKE ?", search, search, search)
	}
	if query.Status != nil {
		db = db.Where("status = ?", string(*query.Status))
	}
	if query.RoleID != nil {
		db = db.Joins("JOIN tbl_user_roles ON tbl_user_roles.user_id = tbl_users.id AND tbl_user_roles.deleted_at IS NULL").
			Where("tbl_user_roles.role_id = ?", *query.RoleID)
	}

	var total int64
	if err := db.Count(&total).Error; err != nil {
		return domain.PageResult[domain.User]{}, helper.Internal("count_users_failed", "failed to count users", err)
	}

	sortBy := safeSort(query.SortBy, []string{"email", "username", "display_name", "status", "created_at", "updated_at"}, "created_at")
	sortOrder := safeSortOrder(query.SortOrder)

	var models []UserModel
	if err := db.Order("tbl_users." + sortBy + " " + sortOrder).
		Offset(query.Paging.Offset()).
		Limit(query.Paging.Limit).
		Find(&models).Error; err != nil {
		return domain.PageResult[domain.User]{}, helper.Internal("list_users_failed", "failed to list users", err)
	}

	items := make([]domain.User, 0, len(models))
	for _, model := range models {
		user, err := r.loadUserRelations(ctx, model)
		if err != nil {
			return domain.PageResult[domain.User]{}, err
		}
		if !query.IncludeAuth {
			user.Identities = nil
		}
		items = append(items, user)
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

func (r *userRepository) ExistsByUsername(ctx context.Context, username string, excludeID *uuid.UUID) (bool, error) {
	return r.exists(ctx, "LOWER(username) = ?", strings.ToLower(username), excludeID)
}

func (r *userRepository) ReplaceRoles(ctx context.Context, userID uuid.UUID, roleIDs []uuid.UUID, actorID string) error {
	db := r.dbWithContext(ctx)
	if err := db.Where("user_id = ?", userID).Delete(&UserRoleModel{}).Error; err != nil {
		return helper.Internal("replace_user_roles_failed", "failed to replace user roles", err)
	}

	now := time.Now()
	for _, roleID := range dedupeUUIDs(roleIDs) {
		model := UserRoleModel{
			UserID:    userID,
			RoleID:    roleID,
			CreatedBy: actorID,
			UpdatedBy: actorID,
			CreatedAt: now,
			UpdatedAt: now,
		}
		if err := db.Create(&model).Error; err != nil {
			return helper.Internal("create_user_role_failed", "failed to assign role to user", err)
		}
	}

	return nil
}

func (r *userRepository) exists(ctx context.Context, clause string, value any, excludeID *uuid.UUID) (bool, error) {
	db := r.dbWithContext(ctx).Model(&UserModel{}).Where(clause, value)
	if excludeID != nil {
		db = db.Where("id <> ?", *excludeID)
	}
	var count int64
	if err := db.Count(&count).Error; err != nil {
		return false, helper.Internal("exists_user_failed", "failed to check user uniqueness", err)
	}
	return count > 0, nil
}

func (r *userRepository) loadUserRelations(ctx context.Context, model UserModel) (domain.User, error) {
	roles, err := r.loadRolesByUserID(ctx, model.ID)
	if err != nil {
		return domain.User{}, err
	}
	identities, err := r.loadAuthIdentitiesByUserID(ctx, model.ID)
	if err != nil {
		return domain.User{}, err
	}
	return toDomainUser(model, roles, identities), nil
}

func (r *userRepository) loadRolesByUserID(ctx context.Context, userID uuid.UUID) ([]domain.Role, error) {
	var models []RoleModel
	if err := r.dbWithContext(ctx).
		Table("tbl_roles").
		Joins("JOIN tbl_user_roles ON tbl_user_roles.role_id = tbl_roles.id AND tbl_user_roles.deleted_at IS NULL").
		Where("tbl_user_roles.user_id = ?", userID).
		Order("tbl_roles.name ASC").
		Find(&models).Error; err != nil {
		return nil, helper.Internal("load_user_roles_failed", "failed to load user roles", err)
	}

	items := make([]domain.Role, 0, len(models))
	for _, model := range models {
		role, err := (&roleRepository{baseRepository: r.baseRepository}).loadRoleRelations(ctx, model)
		if err != nil {
			return nil, err
		}
		items = append(items, role)
	}
	return items, nil
}

func (r *userRepository) loadAuthIdentitiesByUserID(ctx context.Context, userID uuid.UUID) ([]domain.AuthIdentity, error) {
	var models []AuthIdentityModel
	if err := r.dbWithContext(ctx).Where("user_id = ?", userID).Order("created_at ASC").Find(&models).Error; err != nil {
		return nil, helper.Internal("load_auth_identities_failed", "failed to load auth identities", err)
	}
	items := make([]domain.AuthIdentity, 0, len(models))
	for _, model := range models {
		items = append(items, toDomainAuthIdentity(model))
	}
	return items, nil
}

type roleRepository struct {
	*baseRepository
}

func (r *roleRepository) Create(ctx context.Context, role domain.Role) (domain.Role, error) {
	now := time.Now()
	model := RoleModel{
		ID:          role.ID,
		Name:        string(role.Name),
		Description: role.Description,
		CreatedBy:   role.CreatedBy,
		UpdatedBy:   role.UpdatedBy,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	if model.ID == uuid.Nil {
		model.ID = uuid.New()
	}
	if err := r.dbWithContext(ctx).Create(&model).Error; err != nil {
		return domain.Role{}, helper.Internal("create_role_failed", "failed to create role", err)
	}
	return r.FindByID(ctx, model.ID)
}

func (r *roleRepository) Update(ctx context.Context, input domain.RoleUpdateInput) (domain.Role, error) {
	var model RoleModel
	if err := r.dbWithContext(ctx).First(&model, "id = ?", input.ID).Error; err != nil {
		return domain.Role{}, mapGormNotFound(err, "role_not_found", "role not found", "find_role_failed", "failed to load role")
	}
	if input.Name.Set {
		model.Name = string(input.Name.Value)
	}
	if input.Description.Set {
		model.Description = input.Description.Value
	}
	model.UpdatedBy = input.ActorID
	model.UpdatedAt = time.Now()
	if err := r.dbWithContext(ctx).Save(&model).Error; err != nil {
		return domain.Role{}, helper.Internal("update_role_failed", "failed to update role", err)
	}
	return r.FindByID(ctx, model.ID)
}

func (r *roleRepository) Delete(ctx context.Context, id uuid.UUID, actorID string) error {
	if err := r.dbWithContext(ctx).Model(&RoleModel{}).Where("id = ?", id).Updates(map[string]any{
		"updated_by": actorID,
		"updated_at": time.Now(),
	}).Delete(&RoleModel{}, "id = ?", id).Error; err != nil {
		return helper.Internal("delete_role_failed", "failed to delete role", err)
	}
	return nil
}

func (r *roleRepository) FindByID(ctx context.Context, id uuid.UUID) (domain.Role, error) {
	var model RoleModel
	if err := r.dbWithContext(ctx).First(&model, "id = ?", id).Error; err != nil {
		return domain.Role{}, mapGormNotFound(err, "role_not_found", "role not found", "find_role_failed", "failed to find role")
	}
	return r.loadRoleRelations(ctx, model)
}

func (r *roleRepository) FindByName(ctx context.Context, name domain.RoleName) (domain.Role, error) {
	var model RoleModel
	if err := r.dbWithContext(ctx).First(&model, "name = ?", string(name)).Error; err != nil {
		return domain.Role{}, mapGormNotFound(err, "role_not_found", "role not found", "find_role_failed", "failed to find role")
	}
	return r.loadRoleRelations(ctx, model)
}

func (r *roleRepository) FindAll(ctx context.Context, query domain.RoleQuery) (domain.PageResult[domain.Role], error) {
	db := r.dbWithContext(ctx).Model(&RoleModel{})
	if query.Search != "" {
		search := "%" + strings.ToLower(query.Search) + "%"
		db = db.Where("LOWER(name) LIKE ? OR LOWER(description) LIKE ?", search, search)
	}
	if query.Permission != nil {
		db = db.Joins("JOIN tbl_role_permissions ON tbl_role_permissions.role_id = tbl_roles.id AND tbl_role_permissions.deleted_at IS NULL").
			Where("tbl_role_permissions.permission_code = ?", string(*query.Permission))
	}
	var total int64
	if err := db.Count(&total).Error; err != nil {
		return domain.PageResult[domain.Role]{}, helper.Internal("count_roles_failed", "failed to count roles", err)
	}

	sortBy := safeSort(query.SortBy, []string{"name", "created_at", "updated_at"}, "created_at")
	sortOrder := safeSortOrder(query.SortOrder)

	var models []RoleModel
	if err := db.Order("tbl_roles." + sortBy + " " + sortOrder).
		Offset(query.Paging.Offset()).
		Limit(query.Paging.Limit).
		Find(&models).Error; err != nil {
		return domain.PageResult[domain.Role]{}, helper.Internal("list_roles_failed", "failed to list roles", err)
	}

	items := make([]domain.Role, 0, len(models))
	for _, model := range models {
		role, err := r.loadRoleRelations(ctx, model)
		if err != nil {
			return domain.PageResult[domain.Role]{}, err
		}
		items = append(items, role)
	}

	return domain.PageResult[domain.Role]{
		Items: items,
		Paging: domain.Paging{
			Page:  query.Paging.Page,
			Limit: query.Paging.Limit,
			Total: total,
		},
	}, nil
}

func (r *roleRepository) ExistsByName(ctx context.Context, name domain.RoleName, excludeID *uuid.UUID) (bool, error) {
	db := r.dbWithContext(ctx).Model(&RoleModel{}).Where("name = ?", string(name))
	if excludeID != nil {
		db = db.Where("id <> ?", *excludeID)
	}
	var count int64
	if err := db.Count(&count).Error; err != nil {
		return false, helper.Internal("exists_role_failed", "failed to check role uniqueness", err)
	}
	return count > 0, nil
}

func (r *roleRepository) ReplacePermissions(ctx context.Context, roleID uuid.UUID, permissions []domain.RolePermissionInput, actorID string) error {
	db := r.dbWithContext(ctx)
	if err := db.Where("role_id = ?", roleID).Delete(&RolePermissionModel{}).Error; err != nil {
		return helper.Internal("replace_role_permissions_failed", "failed to replace role permissions", err)
	}

	now := time.Now()
	seen := map[string]struct{}{}
	for _, permission := range permissions {
		key := string(permission.PermissionCode) + "::" + permission.Scope
		if _, ok := seen[key]; ok {
			continue
		}
		seen[key] = struct{}{}

		model := RolePermissionModel{
			RoleID:         roleID,
			PermissionCode: string(permission.PermissionCode),
			Scope:          strings.TrimSpace(permission.Scope),
			CreatedBy:      actorID,
			UpdatedBy:      actorID,
			CreatedAt:      now,
			UpdatedAt:      now,
		}
		if err := db.Create(&model).Error; err != nil {
			return helper.Internal("create_role_permission_failed", "failed to assign permission to role", err)
		}
	}
	return nil
}

func (r *roleRepository) loadRoleRelations(ctx context.Context, model RoleModel) (domain.Role, error) {
	var permissions []RolePermissionModel
	if err := r.dbWithContext(ctx).Where("role_id = ?", model.ID).Order("permission_code ASC").Find(&permissions).Error; err != nil {
		return domain.Role{}, helper.Internal("load_role_permissions_failed", "failed to load role permissions", err)
	}
	return toDomainRole(model, permissions), nil
}

type authIdentityRepository struct {
	*baseRepository
}

func (r *authIdentityRepository) Create(ctx context.Context, identity domain.AuthIdentity) (domain.AuthIdentity, error) {
	now := time.Now()
	model := AuthIdentityModel{
		ID:             identity.ID,
		UserID:         identity.UserID,
		Provider:       string(identity.Provider),
		ProviderUserID: identity.ProviderUserID,
		PasswordHash:   identity.PasswordHash,
		CreatedBy:      identity.CreatedBy,
		UpdatedBy:      identity.UpdatedBy,
		CreatedAt:      now,
		UpdatedAt:      now,
	}
	if model.ID == uuid.Nil {
		model.ID = uuid.New()
	}
	if err := r.dbWithContext(ctx).Create(&model).Error; err != nil {
		return domain.AuthIdentity{}, helper.Internal("create_auth_identity_failed", "failed to create auth identity", err)
	}
	return toDomainAuthIdentity(model), nil
}

func (r *authIdentityRepository) Update(ctx context.Context, input domain.AuthIdentityUpdateInput) (domain.AuthIdentity, error) {
	var model AuthIdentityModel
	if err := r.dbWithContext(ctx).First(&model, "id = ?", input.ID).Error; err != nil {
		return domain.AuthIdentity{}, mapGormNotFound(err, "auth_identity_not_found", "auth identity not found", "find_auth_identity_failed", "failed to load auth identity")
	}
	if input.ProviderUserID.Set {
		model.ProviderUserID = input.ProviderUserID.Value
	}
	if input.Password.Set {
		hash, err := bcrypt.GenerateFromPassword([]byte(input.Password.Value), bcrypt.DefaultCost)
		if err != nil {
			return domain.AuthIdentity{}, helper.Internal("hash_password_failed", "failed to hash password", err)
		}
		model.PasswordHash = string(hash)
	}
	model.UpdatedBy = input.ActorID
	model.UpdatedAt = time.Now()
	if err := r.dbWithContext(ctx).Save(&model).Error; err != nil {
		return domain.AuthIdentity{}, helper.Internal("update_auth_identity_failed", "failed to update auth identity", err)
	}
	return toDomainAuthIdentity(model), nil
}

func (r *authIdentityRepository) Delete(ctx context.Context, id uuid.UUID, actorID string) error {
	if err := r.dbWithContext(ctx).Model(&AuthIdentityModel{}).Where("id = ?", id).Updates(map[string]any{
		"updated_by": actorID,
		"updated_at": time.Now(),
	}).Delete(&AuthIdentityModel{}, "id = ?", id).Error; err != nil {
		return helper.Internal("delete_auth_identity_failed", "failed to delete auth identity", err)
	}
	return nil
}

func (r *authIdentityRepository) FindByID(ctx context.Context, id uuid.UUID) (domain.AuthIdentity, error) {
	var model AuthIdentityModel
	if err := r.dbWithContext(ctx).First(&model, "id = ?", id).Error; err != nil {
		return domain.AuthIdentity{}, mapGormNotFound(err, "auth_identity_not_found", "auth identity not found", "find_auth_identity_failed", "failed to find auth identity")
	}
	return toDomainAuthIdentity(model), nil
}

func (r *authIdentityRepository) FindAll(ctx context.Context, query domain.AuthIdentityQuery) (domain.PageResult[domain.AuthIdentity], error) {
	db := r.dbWithContext(ctx).Model(&AuthIdentityModel{})
	if query.UserID != nil {
		db = db.Where("user_id = ?", *query.UserID)
	}
	if query.Provider != nil {
		db = db.Where("provider = ?", string(*query.Provider))
	}
	var total int64
	if err := db.Count(&total).Error; err != nil {
		return domain.PageResult[domain.AuthIdentity]{}, helper.Internal("count_auth_identities_failed", "failed to count auth identities", err)
	}

	sortBy := safeSort(query.SortBy, []string{"provider", "created_at", "updated_at"}, "created_at")
	sortOrder := safeSortOrder(query.SortOrder)

	var models []AuthIdentityModel
	if err := db.Order("tbl_auth_identities." + sortBy + " " + sortOrder).
		Offset(query.Paging.Offset()).
		Limit(query.Paging.Limit).
		Find(&models).Error; err != nil {
		return domain.PageResult[domain.AuthIdentity]{}, helper.Internal("list_auth_identities_failed", "failed to list auth identities", err)
	}

	items := make([]domain.AuthIdentity, 0, len(models))
	for _, model := range models {
		items = append(items, toDomainAuthIdentity(model))
	}

	return domain.PageResult[domain.AuthIdentity]{
		Items: items,
		Paging: domain.Paging{
			Page:  query.Paging.Page,
			Limit: query.Paging.Limit,
			Total: total,
		},
	}, nil
}

func (r *authIdentityRepository) FindPasswordIdentityByUserID(ctx context.Context, userID uuid.UUID) (domain.AuthIdentity, error) {
	var model AuthIdentityModel
	if err := r.dbWithContext(ctx).First(&model, "user_id = ? AND provider = ?", userID, string(domain.AuthProviderPassword)).Error; err != nil {
		return domain.AuthIdentity{}, mapGormNotFound(err, "auth_identity_not_found", "password identity not found", "find_auth_identity_failed", "failed to find password identity")
	}
	return toDomainAuthIdentity(model), nil
}

func (r *authIdentityRepository) FindPasswordIdentityByLogin(ctx context.Context, login string) (domain.AuthIdentity, error) {
	login = strings.TrimSpace(strings.ToLower(login))

	var model AuthIdentityModel
	if err := r.dbWithContext(ctx).
		Table("tbl_auth_identities").
		Joins("JOIN tbl_users ON tbl_users.id = tbl_auth_identities.user_id AND tbl_users.deleted_at IS NULL").
		Where("tbl_auth_identities.provider = ?", string(domain.AuthProviderPassword)).
		Where("LOWER(tbl_users.email) = ? OR LOWER(tbl_users.username) = ?", login, login).
		First(&model).Error; err != nil {
		return domain.AuthIdentity{}, mapGormNotFound(err, "auth_identity_not_found", "password identity not found", "find_auth_identity_failed", "failed to find password identity")
	}
	return toDomainAuthIdentity(model), nil
}

func (r *authIdentityRepository) ExistsByProviderIdentity(ctx context.Context, provider domain.AuthProvider, providerUserID string, excludeID *uuid.UUID) (bool, error) {
	db := r.dbWithContext(ctx).Model(&AuthIdentityModel{}).Where("provider = ? AND provider_user_id = ?", string(provider), providerUserID)
	if excludeID != nil {
		db = db.Where("id <> ?", *excludeID)
	}
	var count int64
	if err := db.Count(&count).Error; err != nil {
		return false, helper.Internal("exists_auth_identity_failed", "failed to check auth identity uniqueness", err)
	}
	return count > 0, nil
}

func toDomainUser(model UserModel, roles []domain.Role, identities []domain.AuthIdentity) domain.User {
	user := domain.User{
		ID:          model.ID,
		Email:       model.Email,
		Username:    model.Username,
		DisplayName: model.DisplayName,
		AvatarID:    model.AvatarID,
		Status:      domain.UserStatus(model.Status),
		Roles:       roles,
		Identities:  identities,
		AuditFields: domain.AuditFields{
			CreatedAt: model.CreatedAt,
			UpdatedAt: model.UpdatedAt,
			CreatedBy: model.CreatedBy,
			UpdatedBy: model.UpdatedBy,
		},
	}
	if model.DeletedAt.Valid {
		deletedAt := model.DeletedAt.Time
		user.DeletedAt = &deletedAt
	}
	return user
}

func toDomainRole(model RoleModel, permissions []RolePermissionModel) domain.Role {
	role := domain.Role{
		ID:          model.ID,
		Name:        domain.RoleName(model.Name),
		Description: model.Description,
		AuditFields: domain.AuditFields{
			CreatedAt: model.CreatedAt,
			UpdatedAt: model.UpdatedAt,
			CreatedBy: model.CreatedBy,
			UpdatedBy: model.UpdatedBy,
		},
	}
	if model.DeletedAt.Valid {
		deletedAt := model.DeletedAt.Time
		role.DeletedAt = &deletedAt
	}
	for _, permission := range permissions {
		item := domain.RolePermission{
			RoleID:         permission.RoleID,
			PermissionCode: domain.PermissionCode(permission.PermissionCode),
			Scope:          permission.Scope,
			AuditFields: domain.AuditFields{
				CreatedAt: permission.CreatedAt,
				UpdatedAt: permission.UpdatedAt,
				CreatedBy: permission.CreatedBy,
				UpdatedBy: permission.UpdatedBy,
			},
		}
		if permission.DeletedAt.Valid {
			deletedAt := permission.DeletedAt.Time
			item.DeletedAt = &deletedAt
		}
		role.Permissions = append(role.Permissions, item)
	}
	return role
}

func toDomainAuthIdentity(model AuthIdentityModel) domain.AuthIdentity {
	identity := domain.AuthIdentity{
		ID:             model.ID,
		UserID:         model.UserID,
		Provider:       domain.AuthProvider(model.Provider),
		ProviderUserID: model.ProviderUserID,
		PasswordHash:   model.PasswordHash,
		AuditFields: domain.AuditFields{
			CreatedAt: model.CreatedAt,
			UpdatedAt: model.UpdatedAt,
			CreatedBy: model.CreatedBy,
			UpdatedBy: model.UpdatedBy,
		},
	}
	if model.DeletedAt.Valid {
		deletedAt := model.DeletedAt.Time
		identity.DeletedAt = &deletedAt
	}
	return identity
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

func dedupeUUIDs(values []uuid.UUID) []uuid.UUID {
	seen := make(map[uuid.UUID]struct{}, len(values))
	items := make([]uuid.UUID, 0, len(values))
	for _, value := range values {
		if value == uuid.Nil {
			continue
		}
		if _, ok := seen[value]; ok {
			continue
		}
		seen[value] = struct{}{}
		items = append(items, value)
	}
	return items
}
