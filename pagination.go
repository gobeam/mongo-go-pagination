package mongopagination

import (
	"context"
	"math"
)

// Paginator struct for holding pagination info
type Paginator struct {
	TotalRecord int64 `json:"total_record"`
	TotalPage   int64 `json:"total_page"`
	Offset      int64 `json:"offset"`
	Limit       int64 `json:"limit"`
	Page        int64 `json:"page"`
	PrevPage    int64 `json:"prev_page"`
	NextPage    int64 `json:"next_page"`
}

// PaginationData struct for returning pagination stat
type PaginationData struct {
	Total     int64 `json:"total"`
	Page      int64 `json:"page"`
	PerPage   int64 `json:"perPage"`
	Prev      int64 `json:"prev"`
	Next      int64 `json:"next"`
	TotalPage int64 `json:"totalPage"`
}

// PaginationData returns PaginationData struct which
// holds information of all stats needed for pagination
func (p *Paginator) PaginationData() *PaginationData {
	data := PaginationData{
		Total:     p.TotalRecord,
		Page:      p.Page,
		PerPage:   p.Limit,
		Prev:      0,
		Next:      0,
		TotalPage: p.TotalPage,
	}
	if p.Page != p.PrevPage && p.TotalRecord > 0 {
		data.Prev = p.PrevPage
	}
	if p.Page != p.NextPage && p.TotalRecord > 0 && p.Page <= p.TotalPage {
		data.Next = p.NextPage
	}

	return &data
}

// Paging returns Paginator struct which hold pagination
// stats
func Paging(p *PagingQuery) *Paginator {
	if p.Page < 1 {
		p.Page = 1
	}
	if p.Limit == 0 {
		p.Limit = 10
	}
	var paginator Paginator
	var count int64
	var offset int64
	total, _ := p.Collection.CountDocuments(context.Background(), p.Filter)
	count = int64(total)

	if p.Page == 1 {
		offset = 0
	} else {
		offset = (p.Page - 1) * p.Limit
	}
	paginator.TotalRecord = count
	paginator.Page = p.Page
	paginator.Offset = offset
	paginator.Limit = p.Limit
	paginator.TotalPage = int64(math.Ceil(float64(count) / float64(p.Limit)))
	if p.Page > 1 {
		paginator.PrevPage = p.Page - 1
	} else {
		paginator.PrevPage = p.Page
	}
	if p.Page == paginator.TotalPage {
		paginator.NextPage = p.Page
	} else {
		paginator.NextPage = p.Page + 1
	}
	return &paginator
}
