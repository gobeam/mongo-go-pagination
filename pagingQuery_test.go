package mongopagination

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"strings"
	"sync"
	"testing"
	"time"
)

type TodoTest struct {
	Title     string    `json:"title" bson:"title"`
	Status    string    `json:"status" bson:"status"`
	CreatedAt time.Time `json:"createdAt" bson:"createdAt"`
}

const (
	DatabaseHost       string = "mongodb://localhost:27017"
	DatabaseName       string = "todo"
	DatabaseCollection string = "TodoTest"
)

func cleanup(db *mongo.Database) (err error) {
	err = db.Collection(DatabaseCollection).Drop(context.Background())
	return
}

func insertExamples(db *mongo.Database) (insertedIds []interface{}, err error) {
	var data []interface{}
	for i := 0; i < 20; i++ {
		data = append(data, bson.M{
			"title":     fmt.Sprintf("todo-%d", i),
			"status":    "active",
			"createdAt": time.Now(),
		})
	}
	result, err := db.Collection(DatabaseCollection).InsertMany(
		context.Background(), data)
	if err != nil {
		return nil, err
	}
	return result.InsertedIDs, nil
}

func insertExamplesForSortingTest(db *mongo.Database) (insertedIds []interface{}, err error) {
	var data []interface{}
	data = append(data, bson.M{"title": "a"})
	data = append(data, bson.M{"title": "c"})
	data = append(data, bson.M{"title": "b"})
	data = append(data, bson.M{"title": "A"})
	data = append(data, bson.M{"title": "C"})
	data = append(data, bson.M{"title": "B"})
	result, err := db.Collection(DatabaseCollection).InsertMany(
		context.Background(), data)
	if err != nil {
		return nil, err
	}
	return result.InsertedIDs, nil
}

func TestPagingQuery_FindWithCollation(t *testing.T) {
	_, session := NewConnection()
	db := session.Database(DatabaseName)
	defer cleanup(db)
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	insertedIds, err := insertExamplesForSortingTest(db)
	if len(insertedIds) < 1 {
		t.Errorf("Empty insert ids")
	}
	if err != nil {
		t.Errorf("Data insert error. Error: %s", err.Error())
	}
	filter := bson.M{}
	var limit int64 = 10
	var page int64
	collection := db.Collection(DatabaseCollection)
	var todos []TodoTest

	collation := options.Collation{
		Locale:    "en",
		CaseLevel: true,
	}

	paginatedData, err := New(collection).SetCollation(&collation).Context(ctx).Limit(limit).Page(page).Sort("title", 1).Filter(filter).Decode(&todos).Find()

	if err != nil {
		t.Errorf("Error while pagination. Error: %s", err.Error())
	}
	if paginatedData == nil {
		t.Errorf("Empty Pagination data error")
		return
	}

	if len(todos) < 1 {
		t.Errorf("Error fetching data")
	}

	if paginatedData.Pagination.Total != 6 || paginatedData.Pagination.Page != 1 {
		t.Errorf("False Pagination data should be 6 but got: %d", paginatedData.Pagination.Total)
	}

	//Check if all title are in the right order
	if !strings.EqualFold(todos[0].Title, "a") || !strings.EqualFold(todos[1].Title, "a") {
		t.Errorf("Index 0 and 1 should be a and A. But index 0 was %s and index 1 was %s", todos[0].Title, todos[1].Title)
	}
	if !strings.EqualFold(todos[2].Title, "b") || !strings.EqualFold(todos[3].Title, "b") {
		t.Errorf("Index 2 and 3 should be b and B. But index 2 was %s and index 3 was %s", todos[2].Title, todos[3].Title)
	}
	if !strings.EqualFold(todos[4].Title, "c") || !strings.EqualFold(todos[5].Title, "c") {
		t.Errorf("Index 4 and 5 should be c and C. But index 4 was %s and index 5 was %s", todos[4].Title, todos[5].Title)
	}
}

func TestPagingQuery_Find(t *testing.T) {
	_, session := NewConnection()
	db := session.Database(DatabaseName)
	defer cleanup(db)
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	insertedIds, err := insertExamples(db)
	if len(insertedIds) < 1 {
		t.Errorf("Empty insert ids")
	}
	if err != nil {
		t.Errorf("Data insert error. Error: %s", err.Error())
	}
	filter := bson.M{}
	var limit int64 = 10
	var page int64
	projection := bson.D{
		{"title", 1},
		{"status", 1},
	}
	collection := db.Collection(DatabaseCollection)
	var todos []TodoTest
	paginatedData, err := New(collection).Context(ctx).Limit(limit).Page(page).Sort("price", -1).Select(projection).Filter(filter).Decode(&todos).Find()

	if err != nil {
		t.Errorf("Error while pagination. Error: %s", err.Error())
	}
	if paginatedData == nil {
		t.Errorf("Empty Pagination data error")
		return
	}

	if len(todos) < 1 {
		t.Errorf("Error fetching data")
	}

	if paginatedData.Pagination.Total != 20 || paginatedData.Pagination.Page != 1 {
		t.Errorf("False Pagination data should be 20 but got: %d", paginatedData.Pagination.Total)
	}

	// no limit or page provided error
	_, noLimitOrPageError := New(collection).Sort("price", -1).Select(projection).Filter(filter).Find()
	if noLimitOrPageError == nil {
		t.Errorf("Error expected but got no error")
	}

	// no filter error
	_, noFilterError := New(collection).Limit(limit).Page(page).Sort("price", -1).Select(projection).Find()
	if noFilterError == nil {
		t.Errorf("Error expected but got no error")
	}
}

func TestPagingQuery_Aggregate(t *testing.T) {
	_, session := NewConnection()
	db := session.Database(DatabaseName)
	collection := db.Collection(DatabaseCollection)
	defer cleanup(db)
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	insertedIds, err := insertExamples(db)
	if len(insertedIds) < 1 {
		t.Errorf("Empty insert ids")
	}
	if err != nil {
		t.Errorf("Data insert error. Error: %s", err.Error())
	}
	// getting page 2 data
	var limit int64 = 10
	var page int64

	// Aggregate pipeline pagination test
	match := bson.M{"$match": bson.M{"status": "active"}}
	filter := bson.M{}

	//check Aggregate Error if decoder is being used which is not supported yet
	var todos []TodoTest
	_, decodeErrorTest := New(collection).Context(ctx).Limit(limit).Page(page).Decode(todos).Aggregate(match)
	if decodeErrorTest == nil {
		t.Errorf("error expected because Decode feature is not supported")
		return
	}

	aggPaginatedData, err := New(collection).Context(ctx).Limit(limit).Page(page).Sort("price", -1).Aggregate(match)
	if err != nil {
		t.Errorf("Error while Aggregation pagination. Error: %s", err.Error())
	}

	if aggPaginatedData == nil {
		t.Errorf("Empty Aggregated Pagination data error")
		return
	}

	// Aggregation error match query test
	faultyMatch := bson.M{"$matches": bson.M{"status": "active"}}
	_, faultyMatchQuery := New(collection).Sort("price", -1).Aggregate(faultyMatch)
	if faultyMatchQuery == nil {
		t.Errorf("Error expected but got no error")
	}

	// no limit or page provided error
	_, noLimitOrPageAggError := New(collection).Sort("price", -1).Aggregate(match)
	if noLimitOrPageAggError == nil {
		t.Errorf("Error expected but got no error")
	}

	// filter in aggregate error
	_, noFilterAggError := New(collection).Limit(limit).Page(page).Filter(filter).Sort("price", -1).Aggregate(match)
	if noFilterAggError == nil {
		t.Errorf("Error expected but got no error")
	}

	// without sorting test
	_, sortProvideTest := New(collection).Aggregate(match)
	if sortProvideTest == nil {
		t.Errorf("data expected")
		return
	}
}

func NewConnection() (a *mongo.Database, b *mongo.Client) {
	var connectOnce sync.Once
	var db *mongo.Database
	var session *mongo.Client
	connectOnce.Do(func() {
		db, session = connect()
	})

	return db, session
}

func connect() (a *mongo.Database, b *mongo.Client) {
	var err error
	session, err := mongo.NewClient(options.Client().ApplyURI(DatabaseHost))
	if err != nil {
		log.Fatal(err)
	}
	err = session.Connect(context.TODO())
	if err != nil {
		log.Fatal(err)
	}

	var db = session.Database(DatabaseName)
	return db, session
}

func TestGetSkip(t *testing.T) {
	tc := []struct {
		limit    int64
		page     int64
		expected int64
	}{
		{
			limit:    10,
			page:     -1,
			expected: 0,
		},
		{
			limit:    10,
			page:     1,
			expected: 0,
		}, {
			limit:    10,
			page:     2,
			expected: 10,
		}, {
			limit:    10,
			page:     3,
			expected: 20,
		},
	}

	for _, tt := range tc {
		skip := getSkip(tt.page, tt.limit)
		if skip != tt.expected {
			t.Fatalf("expected skip to be %d, got %d", tt.expected, skip)
		}
	}
}
