package mongopagination

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
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
	result, err := db.Collection(DatabaseCollection).InsertMany(
		context.Background(),
		[]interface{}{
			bson.D{
				{"title", "todo1"},
				{"status", "active"},
				{"createdAt", time.Now()},
			},
			bson.D{
				{"title", "todo2"},
				{"status", "active"},
				{"createdAt", time.Now()},
			},
			bson.D{
				{"title", "todo3"},
				{"status", "inactive"},
				{"createdAt", time.Now()},
			},
			bson.D{
				{"title", "todo4"},
				{"status", "active"},
				{"createdAt", time.Now()},
			},
			bson.D{
				{"title", "todo5"},
				{"status", "inactive"},
				{"createdAt", time.Now()},
			},
		})
	return result.InsertedIDs, err
}

func TestPagingQuery_Find(t *testing.T) {
	_, session := NewConnection()
	db := session.Database(DatabaseName)
	insertedIds, err := insertExamples(db)
	if len(insertedIds) < 1 {
		t.Errorf("Empty insert ids")
	}
	if err != nil {
		t.Errorf("Data insert error. Error: %s", err.Error())
	}
	filter := bson.M{}
	var limit int64 = 10
	var page int64 = 1
	projection := bson.D{
		{"title", 1},
		{"status", 1},
	}
	collection := db.Collection(DatabaseCollection)
	paginatedData, err := New(collection).Limit(limit).Page(page).Sort("price", -1).Select(projection).Filter(filter).Find()

	if err != nil {
		t.Errorf("Error while pagination. Error: %s", err.Error())
	}
	if paginatedData == nil {
		t.Errorf("Empty Pagination data error")
		return
	}

	if len(paginatedData.Data) < 1 {
		t.Errorf("Error fetching data")
	}

	if paginatedData.Pagination.Total != 5 || paginatedData.Pagination.Page != 1 {
		t.Errorf("False Pagination data should be 5 but got: %d", paginatedData.Pagination.Total)
	}

	// Aggregate pipeline pagination test
	match := bson.M{"$match": bson.M{"status": "active"}}

	aggPaginatedData, err := New(collection).Limit(limit).Page(page).Sort("price", -1).Aggregate(match)
	if err != nil {
		t.Errorf("Error while Aggregation pagination. Error: %s", err.Error())
	}

	if aggPaginatedData == nil {
		t.Errorf("Empty Aggregated Pagination data error")
		return
	}

	err = cleanup(db)
	if err != nil {
		t.Errorf("Error while cleanup. Error: %s", err.Error())
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
