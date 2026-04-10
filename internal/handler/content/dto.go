package content

type CreateSectionRequest struct {
	Slug        string  `json:"slug" validate:"required,max=255"`
	Title       string  `json:"title" validate:"required,max=255"`
	Description *string `json:"description"`
	OrderNo     int     `json:"order_no"`
	IsPublished bool    `json:"is_published"`
}

type UpdateSectionRequest struct {
	Slug        *string `json:"slug" validate:"omitempty,max=255"`
	Title       *string `json:"title" validate:"omitempty,max=255"`
	Description *string `json:"description"`
	OrderNo     *int    `json:"order_no"`
	IsPublished *bool   `json:"is_published"`
}

type CreateLessonRequest struct {
	SectionID   string  `json:"section_id" validate:"required,uuid"`
	Slug        string  `json:"slug" validate:"required,max=255"`
	Title       string  `json:"title" validate:"required,max=255"`
	Description *string `json:"description"`
	OrderNo     int     `json:"order_no"`
	IsPublished bool    `json:"is_published"`
}

type UpdateLessonRequest struct {
	SectionID   *string `json:"section_id" validate:"omitempty,uuid"`
	Slug        *string `json:"slug" validate:"omitempty,max=255"`
	Title       *string `json:"title" validate:"omitempty,max=255"`
	Description *string `json:"description"`
	OrderNo     *int    `json:"order_no"`
	IsPublished *bool   `json:"is_published"`
}

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

type QuestionChoicePayload struct {
	ID          *string `json:"id,omitempty" validate:"omitempty,uuid"`
	ChoiceText  string  `json:"choice_text" validate:"required"`
	ChoiceOrder int     `json:"choice_order"`
	IsCorrect   bool    `json:"is_correct"`
}

type CreateQuestionRequest struct {
	QuestionSetID string                  `json:"question_set_id" validate:"required,uuid"`
	Type          string                  `json:"type" validate:"required,max=64"`
	QuestionText  string                  `json:"question_text" validate:"required"`
	BlankPosition *int                    `json:"blank_position"`
	Explanation   *string                 `json:"explanation"`
	ImageURL      *string                 `json:"image_url" validate:"omitempty,url"`
	Difficulty    int                     `json:"difficulty"`
	OrderNo       int                     `json:"order_no"`
	IsActive      bool                    `json:"is_active"`
	Choices       []QuestionChoicePayload `json:"choices"`
	TagIDs        []string                `json:"tag_ids"`
}

type UpdateQuestionRequest struct {
	QuestionSetID *string                 `json:"question_set_id" validate:"omitempty,uuid"`
	Type          *string                 `json:"type" validate:"omitempty,max=64"`
	QuestionText  *string                 `json:"question_text"`
	BlankPosition *int                    `json:"blank_position"`
	Explanation   *string                 `json:"explanation"`
	ImageURL      *string                 `json:"image_url" validate:"omitempty,url"`
	Difficulty    *int                    `json:"difficulty"`
	OrderNo       *int                    `json:"order_no"`
	IsActive      *bool                   `json:"is_active"`
	Choices       []QuestionChoicePayload `json:"choices"`
	TagIDs        []string                `json:"tag_ids"`
}

type CreateQuestionChoiceRequest struct {
	QuestionID  string `json:"question_id" validate:"required,uuid"`
	ChoiceText  string `json:"choice_text" validate:"required"`
	ChoiceOrder int    `json:"choice_order"`
	IsCorrect   bool   `json:"is_correct"`
}

type UpdateQuestionChoiceRequest struct {
	ChoiceText  *string `json:"choice_text"`
	ChoiceOrder *int    `json:"choice_order"`
	IsCorrect   *bool   `json:"is_correct"`
}

type CreateTagRequest struct {
	Name string `json:"name" validate:"required,max=255"`
}

type UpdateTagRequest struct {
	Name *string `json:"name" validate:"omitempty,max=255"`
}

type SectionResponse struct {
	ID          string  `json:"id"`
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

type LessonResponse struct {
	ID          string  `json:"id"`
	SectionID   string  `json:"section_id"`
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

type QuestionChoiceResponse struct {
	ID          string `json:"id"`
	QuestionID  string `json:"question_id"`
	ChoiceText  string `json:"choice_text"`
	ChoiceOrder int    `json:"choice_order"`
	IsCorrect   bool   `json:"is_correct"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
	CreatedBy   string `json:"created_by"`
	UpdatedBy   string `json:"updated_by"`
}

type TagResponse struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
	CreatedBy string `json:"created_by"`
	UpdatedBy string `json:"updated_by"`
}

type QuestionResponse struct {
	ID            string                   `json:"id"`
	QuestionSetID string                   `json:"question_set_id"`
	Type          string                   `json:"type"`
	QuestionText  string                   `json:"question_text"`
	BlankPosition *int                     `json:"blank_position,omitempty"`
	Explanation   *string                  `json:"explanation,omitempty"`
	ImageURL      *string                  `json:"image_url,omitempty"`
	Difficulty    int                      `json:"difficulty"`
	OrderNo       int                      `json:"order_no"`
	IsActive      bool                     `json:"is_active"`
	Choices       []QuestionChoiceResponse `json:"choices,omitempty"`
	Tags          []TagResponse            `json:"tags,omitempty"`
	CreatedAt     string                   `json:"created_at"`
	UpdatedAt     string                   `json:"updated_at"`
	CreatedBy     string                   `json:"created_by"`
	UpdatedBy     string                   `json:"updated_by"`
}

type PagingResponse struct {
	Page  int   `json:"page"`
	Limit int   `json:"limit"`
	Total int64 `json:"total"`
}

type ListResponse[T any, Q any] struct {
	Items  []T            `json:"items"`
	Paging PagingResponse `json:"paging"`
	Query  Q              `json:"query"`
}
