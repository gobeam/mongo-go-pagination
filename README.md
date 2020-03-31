# Golang Mongo Pagination For Package mongo-go-driver
[![Build][Build-Status-Image]][Build-Status-Url] [![Go Report Card](https://goreportcard.com/badge/github.com/gobeam/mongo-go-pagination?branch=master&kill_cache=1)](https://goreportcard.com/report/github.com/gobeam/mongo-go-pagination) [![GoDoc][godoc-image]][godoc-url]

Simple and easy to use Pagination with information like Total, Page, PerPage, Prev, Next and TotalPage. 


## Install

``` bash
$ go get -u -v github.com/gobeam/mongo-go-pagination
```

or with dep

``` bash
$ dep ensure -add github.com/gobeam/mongo-go-pagination
```


## Demo

``` go
package main

import (
	"context"
	"fmt"
	. "github.com/gobeam/mongo-go-pagination"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type Product struct {
	Id       primitive.ObjectID `json:"_id" bson:"_id"`
	Name     string             `json:"name" bson:"name"`
	Quantity float64            `json:"qty" bson:"qty"`
	Price    float64            `json:"price" bson:"price"`
}

func main() {
	// Establishing mongo db connection
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		panic(err)
	}

	// Example for Normal Find query
	filter := bson.M{}
	var limit int64 = 10
	var page int64 = 1
	collection := client.Database("myaggregate").Collection("stocks")
	projection := bson.D{
		{"name", 1},
		{"qty", 1},
	}
	// Querying paginated data
	paginatedData, err := New(collection).Limit(limit).Page(page).Sort("price", -1).Select(projection).Filter(filter).Find()
	if err != nil {
		panic(err)
	}

	// paginated data is in paginatedData.Data
	// pagination info can be accessed in  paginatedData.Pagination
	// if you want to marshal data to your defined struct

	var lists []Product
	for _, raw := range paginatedData.Data {
		var product *Product
		if marshallErr := bson.Unmarshal(raw, &product); marshallErr == nil {
			lists = append(lists, *product)
		}

	}
	// print ProductList
	fmt.Printf("Norm Find Data: %+v\n", lists)

	// print pagination data
	fmt.Printf("Normal find pagination info: %+v\n", paginatedData.Pagination)




	//Example for Aggregation

	//match query
	match := bson.M{"$match": bson.M{"qty": bson.M{"$gt": 10}}}

	//group query
	projectQuery := bson.M{"$project": bson.M{"_id": 1, "qty": 1}}

	// you can easily chain function and pass multiple query like here we are passing match
	// query and projection query as params in Aggregate function you cannot use filter with Aggregate
	// because you can pass filters directly through Aggregate param
	aggPaginatedData, err := New(collection).Limit(limit).Page(page).Aggregate(match, projectQuery)
	if err != nil {
		panic(err)
	}

	var aggProductList []Product
	for _, raw := range aggPaginatedData.Data {
		var product *Product
		if marshallErr := bson.Unmarshal(raw, &product); marshallErr == nil {
			aggProductList = append(aggProductList, *product)
		}

	}

	// print ProductList
	fmt.Printf("Aggregate Product List: %+v\n", aggProductList)

	// print pagination data
	fmt.Printf("Aggregate Pagination Data: %+v\n", aggPaginatedData.Data)
}
    
```

## Running the tests

``` bash
$ go test
```

## Contributing
Pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change.

Please make sure to update tests as appropriate.



## MIT License

```
Copyright (c) 2020
```

[Build-Status-Url]: https://travis-ci.org/gobeam/mongo-go-pagination
[Build-Status-Image]: https://travis-ci.org/gobeam/mongo-go-pagination.svg?branch=master
[godoc-url]: https://pkg.go.dev/github.com/gobeam/mongo-go-pagination?tab=doc
[godoc-image]: https://godoc.org/github.com/gobeam/mongo-go-pagination?status.svg
