package content

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
