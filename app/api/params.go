package api

import (
	"fmt"
	"net/url"
)

type ListParams struct {
	PageIndex int               `validate:"gte=0"`
	Limit     int               `validate:"gt=0,lte=100"`
	Order     string            `validate:"oneof=ASC DESC asc desc"`
	OrderBy   string            `validate:"oneof=name createdAt"`
	Filter    map[string]string `validate:"omitempty"`
}

func DefaultListParams() *ListParams {
	return &ListParams{
		PageIndex: 0,
		Limit:     10,
		Order:     "DESC",
		OrderBy:   "createdAt",
	}
}

func (p *ListParams) ToQuery() url.Values {
	query := url.Values{}
	query.Set("pageIndex", fmt.Sprintf("%d", p.PageIndex))
	query.Set("limit", fmt.Sprintf("%d", p.Limit))
	query.Set("order", p.Order)
	query.Set("orderBy", p.OrderBy)

	for key, value := range p.Filter {
		query.Set(key, value)
	}

	return query
}

func (p *ListParams) Validate() error {
	if p.PageIndex < 0 {
		return fmt.Errorf("pageIndex must be >= 0, got %d", p.PageIndex)
	}
	if p.Limit <= 0 || p.Limit > 100 {
		return fmt.Errorf("limit must be between 1 and 100, got %d", p.Limit)
	}
	if p.Order != "ASC" && p.Order != "DESC" && p.Order != "asc" && p.Order != "desc" {
		return fmt.Errorf("order must be ASC or DESC, got %s", p.Order)
	}
	return nil
}
