package mongo_go_pagination

import (
	"context"
	"math"
)

// Paginator struct for holding pre pagination
// stats
type Paginator struct {
	TotalRecord int64 `json:"total_record"`
	TotalPage   int64 `json:"total_page"`
	Offset      int64 `json:"offset"`
	Limit       int64 `json:"limit"`
	Page        int64 `json:"page"`
	PrevPage    int64 `json:"prev_page"`
	NextPage    int64 `json:"next_page"`
}

// Pagination struct for returning pagination stat
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
	if p.page < 1 {
		p.page = 1
	}
	if p.limit == 0 {
		p.limit = 10
	}
	var paginator Paginator
	var count int64
	var offset int64
	total, _ := p.collection.CountDocuments(context.Background(), p.filter)
	count = int64(total)

	if p.page == 1 {
		offset = 0
	} else {
		offset = (p.page - 1) * p.limit
	}
	paginator.TotalRecord = count
	paginator.Page = p.page
	paginator.Offset = offset
	paginator.Limit = p.limit
	paginator.TotalPage = int64(math.Ceil(float64(count) / float64(p.limit)))
	if p.page > 1 {
		paginator.PrevPage = p.page - 1
	} else {
		paginator.PrevPage = p.page
	}
	if p.page == paginator.TotalPage {
		paginator.NextPage = p.page
	} else {
		paginator.NextPage = p.page + 1
	}
	return &paginator
}
