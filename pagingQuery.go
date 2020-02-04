package mongo_go_pagination

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

/**

 */
type PagingQuery struct {
	collection *mongo.Collection
	filter     interface{}
	sortField  string
	sortValue  int
	limit      int64
	page       int64
}

/**
Paginated data response struct
 */
type PaginatedData struct {
	Data       []bson.Raw      `json:"data"`
	Pagination PaginationData `json:"pagination"`
}


/**
Find in document
 */
func (paging *PagingQuery) Find() (paginatedData *PaginatedData, err error) {
	skip := getSkip(paging.page, paging.limit)
	opt := &options.FindOptions{
		Skip:  &skip,
		Sort:  bson.D{{paging.sortField, paging.sortValue}},
		Limit: &paging.limit,
	}
	cursor, err := paging.collection.Find(context.Background(), paging.filter, opt)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())
	var docs []bson.Raw
	for cursor.Next(context.Background()) {
		var document *bson.Raw
		if err := cursor.Decode(&document); err == nil {
			docs = append(docs, *document)
		}
	}
	paginator := Paging(&PaginationParam{
		DB:     paging.collection,
		Filter: paging.filter,
		Page:   paging.page,
		Limit:  paging.limit,
	})
	result := PaginatedData{
		Pagination: *paginator.PaginationData(),
		Data:     docs,
	}
	return &result, nil
}

// Get Skip
func getSkip(page, limit int64) (skip int64) {
	if page > 0 {
		skip = (page - 1) * limit
	} else {
		skip = page
	}
	return
}
