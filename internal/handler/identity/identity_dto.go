package identity

type CreateUserRequest struct {
	Email       string   `json:"email" validate:"required,email,max=255"`
	Username    string   `json:"username" validate:"required,max=255"`
	DisplayName string   `json:"display_name" validate:"required,max=255"`
	AvatarID    *string  `json:"avatar_id"`
	Status      string   `json:"status" validate:"required,oneof=INACTIVE ACTIVE BANNED DELETED"`
	RoleIDs     []string `json:"role_ids"`
}

type UpdateUserRequest struct {
	Email       *string  `json:"email" validate:"omitempty,email,max=255"`
	Username    *string  `json:"username" validate:"omitempty,max=255"`
	DisplayName *string  `json:"display_name" validate:"omitempty,max=255"`
	AvatarID    *string  `json:"avatar_id"`
	Status      *string  `json:"status" validate:"omitempty,oneof=INACTIVE ACTIVE BANNED DELETED"`
	RoleIDs     []string `json:"role_ids"`
}

type UserResponse struct {
	ID          string                 `json:"id"`
	Email       string                 `json:"email,omitempty"`
	Username    string                 `json:"username,omitempty"`
	DisplayName string                 `json:"display_name"`
	AvatarID    *string                `json:"avatar_id,omitempty"`
	Status      string                 `json:"status"`
	Roles       []RoleResponse         `json:"roles"`
	Identities  []AuthIdentityResponse `json:"identities,omitempty"`
	CreatedAt   string                 `json:"created_at"`
	UpdatedAt   string                 `json:"updated_at"`
	CreatedBy   string                 `json:"created_by"`
	UpdatedBy   string                 `json:"updated_by"`
}

type UsersListResponse struct {
	Items  []UserResponse   `json:"items"`
	Paging PagingResponse   `json:"paging"`
	Query  UserQueryPayload `json:"query"`
}

type UserQueryPayload struct {
	Search      string `json:"search,omitempty"`
	RoleID      string `json:"role_id,omitempty"`
	Status      string `json:"status,omitempty"`
	SortBy      string `json:"sort_by,omitempty"`
	SortOrder   string `json:"sort_order,omitempty"`
	IncludeAuth bool   `json:"include_auth"`
}

type CreateRoleRequest struct {
	Name        string                  `json:"name" validate:"required,oneof=ADMIN CONTENT_ADMIN MODERATOR USER"`
	Description string                  `json:"description"`
	Permissions []RolePermissionRequest `json:"permissions"`
}

type UpdateRoleRequest struct {
	Name        *string                 `json:"name" validate:"omitempty,oneof=ADMIN CONTENT_ADMIN MODERATOR USER"`
	Description *string                 `json:"description"`
	Permissions []RolePermissionRequest `json:"permissions"`
}

type RolePermissionRequest struct {
	PermissionCode string `json:"permission_code" validate:"required,oneof=CONTENT_READ CONTENT_WRITE CONTENT_PUBLISH USER_READ USER_BAN ANALYTICS_READ SYSTEM_CONFIG"`
	Scope          string `json:"scope"`
}

type RoleResponse struct {
	ID          string                   `json:"id"`
	Name        string                   `json:"name"`
	Description string                   `json:"description,omitempty"`
	Permissions []RolePermissionResponse `json:"permissions"`
	CreatedAt   string                   `json:"created_at"`
	UpdatedAt   string                   `json:"updated_at"`
	CreatedBy   string                   `json:"created_by"`
	UpdatedBy   string                   `json:"updated_by"`
}

type RolePermissionResponse struct {
	PermissionCode string `json:"permission_code"`
	Scope          string `json:"scope,omitempty"`
}

type RolesListResponse struct {
	Items  []RoleResponse   `json:"items"`
	Paging PagingResponse   `json:"paging"`
	Query  RoleQueryPayload `json:"query"`
}

type RoleQueryPayload struct {
	Search     string `json:"search,omitempty"`
	Permission string `json:"permission,omitempty"`
	SortBy     string `json:"sort_by,omitempty"`
	SortOrder  string `json:"sort_order,omitempty"`
}

type CreateAuthIdentityRequest struct {
	UserID         string `json:"user_id" validate:"required,uuid"`
	Provider       string `json:"provider" validate:"required,oneof=PASSWORD GOOGLE APPLE"`
	ProviderUserID string `json:"provider_user_id"`
	Password       string `json:"password"`
}

type UpdateAuthIdentityRequest struct {
	ProviderUserID *string `json:"provider_user_id"`
	Password       *string `json:"password"`
}

type PasswordLoginRequest struct {
	EmailOrUsername string `json:"email_or_username" validate:"required"`
	Password        string `json:"password" validate:"required"`
}

type AuthIdentityResponse struct {
	ID             string `json:"id"`
	UserID         string `json:"user_id"`
	Provider       string `json:"provider"`
	ProviderUserID string `json:"provider_user_id,omitempty"`
	HasPassword    bool   `json:"has_password"`
	CreatedAt      string `json:"created_at"`
	UpdatedAt      string `json:"updated_at"`
	CreatedBy      string `json:"created_by"`
	UpdatedBy      string `json:"updated_by"`
}

type AuthIdentitiesListResponse struct {
	Items  []AuthIdentityResponse   `json:"items"`
	Paging PagingResponse           `json:"paging"`
	Query  AuthIdentityQueryPayload `json:"query"`
}

type AuthIdentityQueryPayload struct {
	UserID    string `json:"user_id,omitempty"`
	Provider  string `json:"provider,omitempty"`
	SortBy    string `json:"sort_by,omitempty"`
	SortOrder string `json:"sort_order,omitempty"`
}

type LoginResponse struct {
	AccessToken      string       `json:"access_token"`
	RefreshToken     string       `json:"refresh_token,omitempty"`
	TokenType        string       `json:"token_type"`
	ExpiresIn        int64        `json:"expires_in"`
	RefreshExpiresIn int64        `json:"refresh_expires_in,omitempty"`
	User             UserResponse `json:"user"`
}

type PagingResponse struct {
	Page  int   `json:"page"`
	Limit int   `json:"limit"`
	Total int64 `json:"total"`
}
