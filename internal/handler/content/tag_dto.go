package content

type CreateTagRequest struct {
	Name string `json:"name" validate:"required,max=255"`
}

type UpdateTagRequest struct {
	Name *string `json:"name" validate:"omitempty,max=255"`
}

type TagResponse struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
	CreatedBy string `json:"created_by"`
	UpdatedBy string `json:"updated_by"`
}
