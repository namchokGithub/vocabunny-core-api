package content

import "github.com/namchokGithub/vocabunny-core-api/internal/core/helper"

type PagingResponse = helper.PagingResponse

type ListResponse[T any, Q any] struct {
	Items  []T            `json:"items"`
	Paging PagingResponse `json:"paging"`
	Query  Q              `json:"query"`
}
