package domain

import "time"

type AuditFields struct {
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
	CreatedBy string
	UpdatedBy string
}

type Paging struct {
	Page  int
	Limit int
	Total int64
}

func (p Paging) Offset() int {
	page := p.Page
	if page <= 0 {
		page = 1
	}

	limit := p.Limit
	if limit <= 0 {
		limit = 20
	}

	return (page - 1) * limit
}

type PageResult[T any] struct {
	Items  []T    `json:"items"`
	Paging Paging `json:"paging"`
}
