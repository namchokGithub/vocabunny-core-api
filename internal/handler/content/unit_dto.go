package content

type CreateUnitRequest struct {
	LessonID    string  `json:"lesson_id" validate:"required,uuid"`
	Slug        string  `json:"slug" validate:"required,max=255"`
	Title       string  `json:"title" validate:"required,max=255"`
	Description *string `json:"description"`
	OrderNo     int     `json:"order_no"`
	IsPublished bool    `json:"is_published"`
}

type UpdateUnitRequest struct {
	LessonID    *string `json:"lesson_id" validate:"omitempty,uuid"`
	Slug        *string `json:"slug" validate:"omitempty,max=255"`
	Title       *string `json:"title" validate:"omitempty,max=255"`
	Description *string `json:"description"`
	OrderNo     *int    `json:"order_no"`
	IsPublished *bool   `json:"is_published"`
}

type UnitResponse struct {
	ID          string  `json:"id"`
	LessonID    string  `json:"lesson_id"`
	Slug        string  `json:"slug"`
	Title       string  `json:"title"`
	Description *string `json:"description,omitempty"`
	OrderNo     int     `json:"order_no"`
	IsPublished bool    `json:"is_published"`
	CreatedAt   string  `json:"created_at"`
	UpdatedAt   string  `json:"updated_at"`
	CreatedBy   string  `json:"created_by"`
	UpdatedBy   string  `json:"updated_by"`
}
