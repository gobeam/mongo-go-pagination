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
import (
	. "github.com/gobeam/mongo-go-pagination"
	)

    filter := bson.M{}
	var limit int64 = 10
	var page int64 = 1
	paging := PagingQuery{
		collection: db.Collection(DatabaseCollectionName),
		filter: filter,
		limit: limit,
		page: page,
		sortField: "createdAt",
		sortValue: -1,
	}
	paginatedData, err := paging.Find()
	
	// paginated data is in paginatedData.Data
	// pagination info can be accessed in  paginatedData.Pagination
	// if you want to marshal data to your defined struct
	
	var lists []TodoTest
    for _, raw := range paginatedData.Data {
        var todo TodoTest
        if err := bson.Unmarshal(raw, &todo); err == nil {
            lists = append(lists, todo)
        }
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
[godoc-url]: https://godoc.org/github.com/gobeam/mongo-go-pagination
[godoc-image]: https://godoc.org/github.com/gobeam/mongo-go-pagination?status.svg
