package helper

import "github.com/namchokGithub/vocabunny-core-api/internal/core/domain"

type PagingResponse struct {
	Page       int   `json:"page"`
	Limit      int   `json:"limit"`
	Total      int64 `json:"total"`
	TotalPages int   `json:"total_pages"`
}

func NewPagingResponse(paging domain.Paging) PagingResponse {
	return PagingResponse{
		Page:       paging.Page,
		Limit:      paging.Limit,
		Total:      paging.Total,
		TotalPages: paging.TotalPages(),
	}
}
