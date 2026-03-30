package handler

type CreateUserRequest struct {
	Code   string  `json:"code" validate:"required,max=64"`
	Name   string  `json:"name" validate:"required,max=255"`
	Email  string  `json:"email" validate:"required,email,max=255"`
	Status string  `json:"status" validate:"required,oneof=active inactive"`
	RoleID *string `json:"role_id"`
}

type UpdateUserRequest struct {
	Code   *string `json:"code" validate:"omitempty,max=64"`
	Name   *string `json:"name" validate:"omitempty,max=255"`
	Email  *string `json:"email" validate:"omitempty,email,max=255"`
	Status *string `json:"status" validate:"omitempty,oneof=active inactive"`
	RoleID *string `json:"role_id"`
}

type UserResponse struct {
	ID        string        `json:"id"`
	Code      string        `json:"code"`
	Name      string        `json:"name"`
	Email     string        `json:"email"`
	Status    string        `json:"status"`
	Role      *RoleResponse `json:"role,omitempty"`
	CreatedAt string        `json:"created_at"`
	UpdatedAt string        `json:"updated_at"`
	CreatedBy string        `json:"created_by"`
	UpdatedBy string        `json:"updated_by"`
}

type RoleResponse struct {
	ID   string `json:"id"`
	Code string `json:"code"`
	Name string `json:"name"`
}

type UsersListResponse struct {
	Items  []UserResponse   `json:"items"`
	Paging PagingResponse   `json:"paging"`
	Query  UserQueryPayload `json:"query"`
}

type PagingResponse struct {
	Page  int   `json:"page"`
	Limit int   `json:"limit"`
	Total int64 `json:"total"`
}

type UserQueryPayload struct {
	Search string `json:"search,omitempty"`
	RoleID string `json:"role_id,omitempty"`
	Status string `json:"status,omitempty"`
	SortBy string `json:"sort_by,omitempty"`
	Order  string `json:"sort_order,omitempty"`
}
