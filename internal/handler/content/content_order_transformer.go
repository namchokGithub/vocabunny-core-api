package content

import "github.com/namchokGithub/vocabunny-core-api/internal/core/domain"

func toContentOrderNoResponse(item domain.ContentOrderNoSummary) ContentOrderNoResponse {
	return ContentOrderNoResponse{
		Sections:     item.Sections,
		Lessons:      item.Lessons,
		Units:        item.Units,
		QuestionSets: item.QuestionSets,
		Questions:    item.Questions,
	}
}
