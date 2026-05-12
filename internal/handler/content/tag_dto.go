package content

type CreateTagRequest struct {
	Name  string `json:"name" validate:"required,max=255"`
	Color string `json:"color" validate:"omitempty,len=7"`
}

type UpdateTagRequest struct {
	Name  *string `json:"name" validate:"omitempty,max=255"`
	Color *string `json:"color" validate:"omitempty,len=7"`
}

type TagResponse struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Color     string `json:"color,omitempty"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
	CreatedBy string `json:"created_by"`
	UpdatedBy string `json:"updated_by"`
}
