package helper

import (
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/namchokGithub/vocabunny-core-api/internal/core/domain"
)

func BuildPaging(c echo.Context) domain.Paging {
	page, _ := strconv.Atoi(c.QueryParam("page"))
	limit, _ := strconv.Atoi(c.QueryParam("limit"))

	if page <= 0 {
		page = 1
	}

	if limit <= 0 {
		limit = 20
	}

	if limit > 100 {
		limit = 100
	}

	return domain.Paging{
		Page:  page,
		Limit: limit,
	}
}

func BuildSort(c echo.Context) (string, string) {
	sortBy := strings.TrimSpace(c.QueryParam("sort_by"))
	sortOrder := strings.TrimSpace(strings.ToLower(c.QueryParam("sort_order")))
	if sortOrder != "desc" {
		sortOrder = "asc"
	}

	return sortBy, sortOrder
}
