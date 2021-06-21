package main

import (
	"context"
	"encoding/json"
	"fmt"
	paginate "github.com/gobeam/mongo-go-pagination"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"net/http"
	"strconv"
)

// Product struct
type Product struct {
	Id       primitive.ObjectID `json:"_id,omitempty" bson:"_id"`
	Name     string             `json:"name,omitempty" bson:"name"`
	Quantity float64            `json:"quantity,omitempty" bson:"quantity"`
	Price    float64            `json:"price,omitempty" bson:"price"`
}

func insertExamples(db *mongo.Database) (insertedIds []interface{}, err error) {
	var data []interface{}
	for i := 0; i < 30; i++ {
		data = append(data, bson.M{
			"name":     fmt.Sprintf("product-%d", i),
			"quantity":    float64(i),
			"price": float64(i*10+5),
		})
	}
	result, err := db.Collection("products").InsertMany(
		context.Background(), data)
	if err != nil {
		return nil, err
	}
	return result.InsertedIDs, nil
}

var dbConnection *mongo.Database

func main() {
	// Establishing mongo db connection
	ctx := context.Background()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://127.0.0.1:27017/"))
	if err != nil {
		panic(err)
	}
	dbConnection = client.Database("myaggregate")
	_, insertErr := insertExamples(client.Database("myaggregate"))
	if insertErr != nil {
		panic(insertErr)
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Testing pagination go to http://localhost:8081/normal-pagination?page=1&limit=10  to beign testing")
	})

	http.HandleFunc("/normal-pagination", func(w http.ResponseWriter, r *http.Request) {
		convertedPageInt, convertedLimitInt := getPageAndLimit(r)
		// Example for Normal Find query
		filter := bson.M{}
		limit := int64(convertedLimitInt)
		page := int64(convertedPageInt)
		collection := dbConnection.Collection("products")
		projection := bson.D{
			{"name", 1},
			{"quantity", 1},
		}
		// Querying paginated data
		// If you want to do some complex sort like sort by score(weight) for full text search fields you can do it easily
		// sortValue := bson.M{
		//		"$meta" : "textScore",
		//	}
		// paginatedData, err := paginate.New(collection).Context(ctx).Limit(limit).Page(page).Sort("score", sortValue)...
		var products []Product
		paginatedData, err := paginate.New(collection).Context(ctx).Limit(limit).Page(page).Sort("price", -1).Sort("quantity", -1).Select(projection).Filter(filter).Decode(&products).Find()
		if err != nil {
			panic(err)
		}

		payload := struct {
			Data []Product `json:"data"`
			Pagination paginate.PaginationData `json:"pagination"`
		}{
			Pagination: paginatedData.Pagination,
			Data: products,
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(payload)
	})

	http.HandleFunc("/aggregate-pagination", func(w http.ResponseWriter, r *http.Request) {
		convertedPageInt, convertedLimitInt := getPageAndLimit(r)
		collection := dbConnection.Collection("products")

		//Example for Aggregation
		limit := int64(convertedLimitInt)
		page := int64(convertedPageInt)
		//match query
		match := bson.M{"$match": bson.M{"quantity": bson.M{"$gt": 0}}}
		//
		//group query
		projectQuery := bson.M{"$project": bson.M{"_id": 1, "name": 1, "quantity": 1}}

		// you can easily chain function and pass multiple query like here we are passing match
		// query and projection query as params in Aggregate function you cannot use filter with Aggregate
		// because you can pass filters directly through Aggregate param
		aggPaginatedData, err := paginate.New(collection).Context(ctx).Limit(limit).Page(page).Sort("price", -1).Aggregate( match,projectQuery)
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

		payload := struct {
			Data []Product `json:"data"`
			Pagination paginate.PaginationData `json:"pagination"`
		}{
			Pagination: aggPaginatedData.Pagination,
			Data: aggProductList,
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(payload)
	})

	fmt.Println("Application started on port http://localhost:8081")
	log.Fatal(http.ListenAndServe(":8081", nil))
}

func getPageAndLimit(r *http.Request) (convertedPageInt int, convertedLimitInt int) {
	queryPageValue := r.FormValue("page")
	if queryPageValue != "" {
		convertedPageInt, _ = strconv.Atoi(queryPageValue)
	}

	queryLimitValue := r.FormValue("limit")
	if queryLimitValue != "" {
		convertedLimitInt, _ = strconv.Atoi(queryLimitValue)
	}

	return convertedPageInt, convertedLimitInt
}
