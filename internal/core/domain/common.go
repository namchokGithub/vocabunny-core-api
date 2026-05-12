package domain

import (
	"strings"
	"time"
)

// Includes is a set of relation names requested by the caller via the ?include= query param.
type Includes map[string]struct{}

func ParseIncludes(raw string) Includes {
	inc := make(Includes)
	for _, part := range strings.Split(raw, ",") {
		if s := strings.TrimSpace(strings.ToLower(part)); s != "" {
			inc[s] = struct{}{}
		}
	}
	return inc
}

func (i Includes) Has(key string) bool {
	if i == nil {
		return false
	}
	_, ok := i[key]
	return ok
}

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

func (p Paging) TotalPages() int {
	if p.Total <= 0 {
		return 0
	}

	limit := p.Limit
	if limit <= 0 {
		return 1
	}

	return int((p.Total + int64(limit) - 1) / int64(limit))
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
