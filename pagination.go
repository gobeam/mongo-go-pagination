package mongopagination

import (
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
func Paging(p *pagingQuery, paginationInfo chan<- *Paginator, aggregate bool, aggCount int64) {
	var paginator Paginator
	var offset int64
	var count int64
	ctx := p.getContext()
	if !aggregate {
		count, _ = p.Collection.CountDocuments(ctx, p.FilterQuery)
	} else {
		count = aggCount
	}

	if p.PageCount > 0 {
		offset = ((p.PageCount - 1) * p.LimitCount) + 1
	} else {
		offset = 0
	}
	paginator.TotalRecord = count
	paginator.Page = p.PageCount
	paginator.Offset = offset
	paginator.Limit = p.LimitCount
	paginator.TotalPage = int64(math.Ceil(float64(count) / float64(p.LimitCount)))
	if p.PageCount > 1 {
		paginator.PrevPage = p.PageCount - 1
	} else {
		paginator.PrevPage = p.PageCount
	}
	if p.PageCount == paginator.TotalPage {
		paginator.NextPage = p.PageCount
	} else {
		paginator.NextPage = p.PageCount + 1
	}
	paginationInfo <- &paginator
}
