package main

import (
	"context"
	. "github.com/gobeam/mongo-go-pagination"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type ToDo struct {
	TaskName   string `json:"task_name"`
	TaskStatus string `json:"task_status"`
}

func main() {
	// Establishing mongo db connection
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		panic(err)
	}

	// making example filter
	filter := bson.M{}
	filter["status"] = 1

	var limit int64 = 10
	var page int64 = 1
	collection := client.Database("db_name").Collection("collection_name")
	projection := bson.D{
		{"task_name", 1},
	}
	// Querying paginated data
	paging := PagingQuery{
		Collection: collection,
		Filter:     filter,
		Limit:      limit,
		Page:       page,
		Projection: projection,
		SortField:  "createdAt",
		SortValue:  -1,
	}
	paginatedData, err := paging.Paginate()

	// paginated data is in paginatedData.Data
	// pagination info can be accessed in  paginatedData.Pagination
	// if you want to marshal data to your defined struct

	var lists []ToDo
	for _, raw := range paginatedData.Data {
		var todo ToDo
		if err := bson.Unmarshal(raw, &todo); err == nil {
			lists = append(lists, todo)
		}
	}
}
