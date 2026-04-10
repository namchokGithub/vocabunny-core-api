package content

type QuestionChoicePayload struct {
	ID          *string `json:"id,omitempty" validate:"omitempty,uuid"`
	ChoiceText  string  `json:"choice_text" validate:"required"`
	ChoiceOrder int     `json:"choice_order"`
	IsCorrect   bool    `json:"is_correct"`
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
