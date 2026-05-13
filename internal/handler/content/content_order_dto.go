package content

type ContentOrderNoResponse struct {
	Sections     int `json:"sections"`
	Lessons      int `json:"lessons"`
	Units        int `json:"units"`
	QuestionSets int `json:"question_sets"`
	Questions    int `json:"questions"`
}
