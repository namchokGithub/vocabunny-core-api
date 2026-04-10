package content

type CreateQuestionSetRequest struct {
	UnitID           string  `json:"unit_id" validate:"required,uuid"`
	Slug             string  `json:"slug" validate:"required,max=255"`
	Title            string  `json:"title" validate:"required,max=255"`
	Description      *string `json:"description"`
	OrderNo          int     `json:"order_no"`
	EstimatedSeconds *int    `json:"estimated_seconds"`
	IsPublished      bool    `json:"is_published"`
	Version          int     `json:"version"`
}

type UpdateQuestionSetRequest struct {
	UnitID           *string `json:"unit_id" validate:"omitempty,uuid"`
	Slug             *string `json:"slug" validate:"omitempty,max=255"`
	Title            *string `json:"title" validate:"omitempty,max=255"`
	Description      *string `json:"description"`
	OrderNo          *int    `json:"order_no"`
	EstimatedSeconds *int    `json:"estimated_seconds"`
	IsPublished      *bool   `json:"is_published"`
	Version          *int    `json:"version"`
}

type QuestionSetResponse struct {
	ID               string  `json:"id"`
	UnitID           string  `json:"unit_id"`
	Slug             string  `json:"slug"`
	Title            string  `json:"title"`
	Description      *string `json:"description,omitempty"`
	OrderNo          int     `json:"order_no"`
	EstimatedSeconds *int    `json:"estimated_seconds,omitempty"`
	IsPublished      bool    `json:"is_published"`
	Version          int     `json:"version"`
	CreatedAt        string  `json:"created_at"`
	UpdatedAt        string  `json:"updated_at"`
	CreatedBy        string  `json:"created_by"`
	UpdatedBy        string  `json:"updated_by"`
}
