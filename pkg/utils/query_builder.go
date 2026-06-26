package utils

import (
	"fmt"
	"math"
	"spot-sync/pkg/httpresponse"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

// QueryBuilder parses query parameters from an HTTP request and applies
// pagination, sorting, and search filtering to GORM queries.
type QueryBuilder struct {
	Page   int
	Limit  int
	Sort   string
	Order  string
	Search string
}

// NewQueryBuilder creates a QueryBuilder by parsing query parameters from an Echo context.
// Defaults: page=1, limit=10, sort=created_at, order=desc.
func NewQueryBuilder(c echo.Context) *QueryBuilder {
	qb := &QueryBuilder{
		Page:   1,
		Limit:  10,
		Sort:   "created_at",
		Order:  "desc",
		Search: "",
	}

	if page, err := strconv.Atoi(c.QueryParam("page")); err == nil && page > 0 {
		qb.Page = page
	}

	if limit, err := strconv.Atoi(c.QueryParam("limit")); err == nil && limit > 0 {
		if limit > 100 {
			limit = 100 // cap at 100 to prevent abuse
		}
		qb.Limit = limit
	}

	if sort := c.QueryParam("sort"); sort != "" {
		qb.Sort = sort
	}

	if order := strings.ToLower(c.QueryParam("order")); order == "asc" || order == "desc" {
		qb.Order = order
	}

	if search := strings.TrimSpace(c.QueryParam("search")); search != "" {
		qb.Search = search
	}

	return qb
}

// Offset returns the calculated offset for pagination.
func (qb *QueryBuilder) Offset() int {
	return (qb.Page - 1) * qb.Limit
}

// ApplySearch applies an ILIKE search filter across the given fields.
// Returns the modified query. If no search term or no fields are provided,
// returns the query unchanged.
func (qb *QueryBuilder) ApplySearch(db *gorm.DB, searchFields []string) *gorm.DB {
	if qb.Search == "" || len(searchFields) == 0 {
		return db
	}

	searchTerm := "%" + qb.Search + "%"
	conditions := make([]string, 0, len(searchFields))
	args := make([]interface{}, 0, len(searchFields))

	for _, field := range searchFields {
		conditions = append(conditions, fmt.Sprintf("%s ILIKE ?", field))
		args = append(args, searchTerm)
	}

	return db.Where(strings.Join(conditions, " OR "), args...)
}

// ApplyPaginationAndSort applies ordering, limit, and offset to the query.
func (qb *QueryBuilder) ApplyPaginationAndSort(db *gorm.DB) *gorm.DB {
	orderClause := fmt.Sprintf("%s %s", qb.Sort, qb.Order)
	return db.Order(orderClause).Limit(qb.Limit).Offset(qb.Offset())
}

// Apply is a convenience method that applies search, sort, and pagination in one call.
func (qb *QueryBuilder) Apply(db *gorm.DB, searchFields []string) *gorm.DB {
	query := qb.ApplySearch(db, searchFields)
	return qb.ApplyPaginationAndSort(query)
}

// GetMeta returns the pagination metadata for a response.
func (qb *QueryBuilder) GetMeta(total int64) *httpresponse.Meta {
	totalPage := int64(math.Ceil(float64(total) / float64(qb.Limit)))
	if totalPage < 1 {
		totalPage = 1
	}

	return &httpresponse.Meta{
		Page:      qb.Page,
		Limit:     qb.Limit,
		Total:     total,
		TotalPage: totalPage,
	}
}
