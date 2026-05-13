package helper

import (
	"testing"

	"github.com/namchokGithub/vocabunny-core-api/internal/core/domain"
)

func TestNewPagingResponse(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		paging     domain.Paging
		totalPages int
	}{
		{
			name:       "calculates ceil for full pages",
			paging:     domain.Paging{Page: 1, Limit: 20, Total: 100},
			totalPages: 5,
		},
		{
			name:       "calculates ceil for partial last page",
			paging:     domain.Paging{Page: 1, Limit: 20, Total: 101},
			totalPages: 6,
		},
		{
			name:       "returns zero for empty results",
			paging:     domain.Paging{Page: 1, Limit: 20, Total: 0},
			totalPages: 0,
		},
		{
			name:       "handles invalid limit safely",
			paging:     domain.Paging{Page: 1, Limit: 0, Total: 5},
			totalPages: 1,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			response := NewPagingResponse(tt.paging)
			if response.TotalPages != tt.totalPages {
				t.Fatalf("expected total_pages %d, got %d", tt.totalPages, response.TotalPages)
			}
		})
	}
}
