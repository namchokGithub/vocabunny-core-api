package content

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
