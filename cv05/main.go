package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/SaimonWoidig/DPB/cv05/pkg/utils"
	_ "github.com/joho/godotenv/autoload"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

const (
	MongoDBConnectionStringKey = "MONGODB_URI"
	MongoDBDatabaseNameKey     = "MONGODB_DB"
	MongoDBCollectionNameKey   = "MONGODB_COLLECTION"
)

var DefaultFindOptions = options.Find().SetLimit(10)

func main() {
	// get mongo URI, database name and collection name
	mongoURI := os.Getenv(MongoDBConnectionStringKey)
	if mongoURI == "" {
		panic("MONGODB_URI is not set")
	}
	mongoDatabaseName := os.Getenv(MongoDBDatabaseNameKey)
	if mongoDatabaseName == "" {
		panic("MONGODB_DB is not set")
	}
	mongoCollectionName := os.Getenv(MongoDBCollectionNameKey)
	if mongoCollectionName == "" {
		panic("MONGODB_COLLECTION is not set")
	}

	// connect to mongo
	slog.Debug("Connecting to MongoDB")
	mdb, err := connectToMongo(context.Background(), mongoURI, 3*time.Second)
	if err != nil {
		slog.Error("Failed to connect to MongoDB", "error", err)
		panic(err.Error())
	}
	slog.Debug("Connected to MongoDB")

	// defer disconnect
	defer func(mdb *mongo.Client) {
		slog.Debug("Disconnecting from MongoDB")
		if err := disconnectMongo(mdb, 3*time.Second); err != nil {
			slog.Error("Failed to disconnect from MongoDB", "error", err)
			panic(err.Error())
		}
		slog.Debug("Disconnected from MongoDB")
	}(mdb)

	// get restaurant collection
	restaurantsCollection := mdb.Database(mongoDatabaseName).Collection(mongoCollectionName)

	// run tasks
	assignmentPart01(restaurantsCollection)
	assignmentPart02(restaurantsCollection)
	assignmentPart03(restaurantsCollection)
}

// task 1 through 8
func assignmentPart01(restaurantsCollection *mongo.Collection) {
	var cur *mongo.Cursor
	var err error

	// print all restaurants
	cur, err = restaurantsCollection.Find(context.Background(), bson.M{}, DefaultFindOptions)
	if err != nil {
		slog.Error("Failed to find restaurants", "error", err)
	}
	allRestaurants := make([]map[string]any, 10)
	cur.All(context.Background(), &allRestaurants)
	fmt.Println("#1 - All restaurants:")
	printAllRestaurants(allRestaurants)

	// print all restaurants - only names and sorted
	cur, err = restaurantsCollection.Find(context.Background(), bson.M{},
		options.Find().
			SetSort(bson.M{"name": 1}).
			SetProjection(bson.M{"_id": 0, "name": 1}),
	)
	if err != nil {
		slog.Error("Failed to find restaurants", "error", err)
	}
	allRestaurantNamesSorted := make([]map[string]any, 10)
	cur.All(context.Background(), &allRestaurantNamesSorted)
	fmt.Println("#2 - All restaurant names sorted:")
	printAllRestaurants(allRestaurantNamesSorted)

	// same as previous, only limited to 5
	cur, err = restaurantsCollection.Find(context.Background(), bson.M{},
		options.Find().
			SetLimit(5).
			SetSort(bson.M{"name": 1}).
			SetProjection(bson.M{"_id": 0, "name": 1}),
	)
	if err != nil {
		slog.Error("Failed to find restaurants", "error", err)
	}
	allRestaurantNamesSortedLimited := make([]map[string]any, 5)
	cur.All(context.Background(), &allRestaurantNamesSortedLimited)
	fmt.Println("#3 - 5 restaurant names sorted:")
	printAllRestaurants(allRestaurantNamesSortedLimited)

	// same as previous, but next ten
	cur, err = restaurantsCollection.Find(context.Background(), bson.M{},
		options.Find().
			SetSkip(5).
			SetLimit(10).
			SetSort(bson.M{"name": 1}).
			SetProjection(bson.M{"_id": 0, "name": 1}),
	)
	if err != nil {
		slog.Error("Failed to find restaurants", "error", err)
	}
	allRestaurantNamesSortedLimitedNext10 := make([]map[string]any, 10)
	cur.All(context.Background(), &allRestaurantNamesSortedLimitedNext10)
	fmt.Println("#4 - Next 10 restaurant names sorted after first 5:")
	printAllRestaurants(allRestaurantNamesSortedLimitedNext10)

	// all restaurants in borough "Bronx"
	cur, err = restaurantsCollection.Find(context.Background(), bson.M{"borough": "Bronx"}, DefaultFindOptions)
	if err != nil {
		slog.Error("Failed to find restaurants", "error", err)
	}
	allBronxRestaurants := make([]map[string]any, 10)
	cur.All(context.Background(), &allBronxRestaurants)
	fmt.Println("#5 - All restaurants in Bronx:")
	printAllRestaurants(allBronxRestaurants)

	// all restaurants with their name beginning with "M"
	cur, err = restaurantsCollection.Find(context.Background(), bson.M{"name": bson.M{"$regex": "^M"}}, DefaultFindOptions)
	if err != nil {
		slog.Error("Failed to find restaurants", "error", err)
	}
	allMRestaurants := make([]map[string]any, 10)
	cur.All(context.Background(), &allMRestaurants)
	fmt.Println("#6 - All restaurants with their name beginning with M:")
	printAllRestaurants(allMRestaurants)

	// all restaurants with a score greater than 80
	cur, err = restaurantsCollection.Find(context.Background(), bson.M{"grades.score": bson.M{"$gt": 80}}, DefaultFindOptions)
	if err != nil {
		slog.Error("Failed to find restaurants", "error", err)
	}
	allHighScoreRestaurants := make([]map[string]any, 10)
	cur.All(context.Background(), &allHighScoreRestaurants)
	fmt.Println("#7 - All restaurants with a score greater than 80:")
	printAllRestaurants(allHighScoreRestaurants)

	// all restaurants with a score between 80 and 90
	cur, err = restaurantsCollection.Find(context.Background(), bson.M{"grades.score": bson.M{"$gte": 80, "$lte": 90}}, DefaultFindOptions)
	if err != nil {
		slog.Error("Failed to find restaurants", "error", err)
	}
	all80to90ScoreRestaurants := make([]map[string]any, 10)
	cur.All(context.Background(), &all80to90ScoreRestaurants)
	fmt.Println("#8 - All restaurants with a score between 80 and 90:")
	printAllRestaurants(all80to90ScoreRestaurants)
}

// task 9 through 15
func assignmentPart02(restaurantsCollection *mongo.Collection) {
	var cur *mongo.Cursor
	var err error

	// all restaurants with a score between 80 and 90 with non-American cuisine
	cur, err = restaurantsCollection.Find(context.Background(),
		bson.M{"grades.score": bson.M{"$gte": 80, "$lte": 90}, "cuisine": bson.M{"$ne": "American"}},
		DefaultFindOptions,
	)
	if err != nil {
		slog.Error("Failed to find restaurants", "error", err)
	}
	all80to90ScoreNonAmericanCuisineRestaurants := make([]map[string]any, 10)
	cur.All(context.Background(), &all80to90ScoreNonAmericanCuisineRestaurants)
	fmt.Println("#9 - All restaurants with a score between 80 and 90 with non-American cuisine:")
	printAllRestaurants(all80to90ScoreNonAmericanCuisineRestaurants)

	// all restaurants with atleast 8 grades
	cur, err = restaurantsCollection.Find(context.Background(),
		// if an element in grades at index 7 exists (index is zero-based)
		bson.M{"grades.7": bson.M{"$exists": "true"}},
		DefaultFindOptions,
	)
	if err != nil {
		slog.Error("Failed to find restaurants", "error", err)
	}
	allRestaurantsWithAtleast8Grades := make([]map[string]any, 10)
	err = cur.All(context.Background(), &allRestaurantsWithAtleast8Grades)
	if err != nil {
		slog.Error("Failed to find restaurants", "error", err.Error())
	}
	fmt.Println("#10 - All restaurants with atleast 8 grades:")
	printAllRestaurants(allRestaurantsWithAtleast8Grades)

	// all restaurants with atleast one grade from the year 2014
	cur, err = restaurantsCollection.Find(context.Background(),
		bson.M{"grades.date": bson.M{"$gte": "2014-01-01", "$lte": "2015-01-01"}},
		DefaultFindOptions,
	)
	if err != nil {
		slog.Error("Failed to find restaurants", "error", err)
	}
	allRestaurantsWithAtleastOneGradeFromTheYear2014 := make([]map[string]any, 10)
	cur.All(context.Background(), &allRestaurantsWithAtleastOneGradeFromTheYear2014)
	fmt.Println("#11 - All restaurants with atleast one grade from the year 2014:")
	printAllRestaurants(allRestaurantsWithAtleastOneGradeFromTheYear2014)

	// create a restaurant with any name and address
	baseRestaurantName := "Barvírna Craft Beer & Burgers"
	newRestaurant := map[string]any{
		"name":          baseRestaurantName,
		"restaurant_id": "123456789",
		"cuisine":       "Burgers",
		"address": map[string]any{
			"coord":    []float64{50.7675691, 15.0519803},
			"street":   "Barvířská",
			"building": "33/4",
			"zipcode":  "460 07",
			"city":     "Liberec",
		},
		"borough": "Jeřáb",
	}
	insertResult, err := restaurantsCollection.InsertOne(context.Background(), newRestaurant)
	if err != nil {
		slog.Error("Failed to insert restaurant", "error", err)
	}
	fmt.Printf("#12 - Inserted restaurant %s\n", insertResult.InsertedID)

	// print the newly created restaurant
	createdRestaurantRes := restaurantsCollection.FindOne(context.Background(), bson.M{"restaurant_id": "123456789"})
	if createdRestaurantRes.Err() != nil {
		slog.Error("Failed to find restaurant", "error", createdRestaurantRes.Err())
	}
	createdRestaurant := make(map[string]any)
	createdRestaurantRes.Decode(&createdRestaurant)
	fmt.Println("#13 - Newly created restaurant:")
	utils.PrettyPrintJSON(createdRestaurant)

	// change the name of the newly created restaurant
	updatedRestaurantName := "Barvírna Craft Beer & Burgerzzz"
	updateResult, err := restaurantsCollection.UpdateOne(context.Background(),
		bson.M{"restaurant_id": "123456789"}, bson.M{"$set": bson.M{"name": updatedRestaurantName}},
	)
	if err != nil {
		slog.Error("Failed to update restaurant", "error", err)
	}
	fmt.Printf("#14 - Updated %d documents\n", updateResult.MatchedCount)

	// delete the newly created restaurant by restaurant ID
	deleteResult, err := restaurantsCollection.DeleteOne(context.Background(), bson.M{"restaurant_id": "123456789"})
	if err != nil {
		slog.Error("Failed to delete restaurant", "error", err)
	}
	fmt.Printf("#15.1 - Deleted %d documents\n", deleteResult.DeletedCount)

	// recreate the restaurant and drop by either old name or new name
	insertResult, err = restaurantsCollection.InsertOne(context.Background(), newRestaurant)
	if err != nil {
		slog.Error("Failed to insert restaurant", "error", err)
	}
	fmt.Printf("Recreated restaurant %s\n", insertResult.InsertedID)
	deleteResult, err = restaurantsCollection.DeleteMany(context.Background(),
		bson.M{"name": bson.M{"$in": []string{baseRestaurantName, updatedRestaurantName}}},
	)
	if err != nil {
		slog.Error("Failed to delete restaurant", "error", err)
	}
	fmt.Printf("#15.2 - Deleted %d documents\n", deleteResult.DeletedCount)
}

// task 16 to create an index
func assignmentPart03(restaurantsCollection *mongo.Collection) {
	fmt.Println("#16 - Creating index on name field")
	// first, construct the query commands as raw BSON
	// found on: https://www.mongodb.com/docs/manual/reference/command/find/#syntax
	findCmd := bson.D{{Key: "find", Value: "restaurants"}, {Key: "filter", Value: bson.D{{Key: "borough", Value: "Bronx"}}}}
	// found on: https://pkg.go.dev/go.mongodb.org/mongo-driver@v1.14.0/mongo#example-Database.RunCommand
	explainCmd := bson.D{{Key: "explain", Value: findCmd}}

	// explain find command: all restaurants in the borough "Bronx"
	var explainBSON, executionStats bson.M
	explainRes := restaurantsCollection.Database().RunCommand(context.Background(),
		explainCmd,
		options.RunCmd().SetReadPreference(readpref.Primary()),
	)
	if explainRes.Err() != nil {
		slog.Error("Failed to explain find command", "error", explainRes.Err())
	}
	if err := explainRes.Decode(&explainBSON); err != nil {
		slog.Error("Failed to decode explain result", "error", err)
	}
	executionStats = explainBSON["executionStats"].(bson.M)
	docsExaminedBeforeIndex := executionStats["totalDocsExamined"].(int32)
	fmt.Printf("Before index: Total docs examined: %d\n", docsExaminedBeforeIndex)

	// create index over borough
	indexName, err := restaurantsCollection.Indexes().CreateOne(context.Background(), mongo.IndexModel{
		Keys: bson.D{{Key: "borough", Value: 1}},
	})
	if err != nil {
		slog.Error("Failed to create index", "error", err)
	}
	fmt.Printf("Created index %q over key %q\n", indexName, "borough")

	// explain find command: all restaurants in the borough "Bronx"
	explainRes = restaurantsCollection.Database().RunCommand(context.Background(),
		explainCmd,
		options.RunCmd().SetReadPreference(readpref.Primary()),
	)
	if explainRes.Err() != nil {
		slog.Error("Failed to explain find command", "error", explainRes.Err())
	}
	if err := explainRes.Decode(&explainBSON); err != nil {
		slog.Error("Failed to decode explain result", "error", err)
	}
	executionStats = explainBSON["executionStats"].(bson.M)
	docsExaminedWithIndex := executionStats["totalDocsExamined"].(int32)
	fmt.Printf("After index: Total docs examined: %d\n", docsExaminedWithIndex)

	fmt.Printf("Creating index saved %d document examinations\n", docsExaminedBeforeIndex-docsExaminedWithIndex)

	// drop index
	if _, err := restaurantsCollection.Indexes().DropOne(context.Background(), indexName); err != nil {
		slog.Error("Failed to drop index", "error", err)
	}
	fmt.Printf("Dropped index %q\n", indexName)
}

func printAllRestaurants(restaurants []map[string]any) {
	for _, restaurant := range restaurants {
		utils.PrettyPrintJSON(restaurant)
	}
}

func connectToMongo(ctx context.Context, mongoURI string, timeout time.Duration) (*mongo.Client, error) {
	connectCtx, connectCtxCancel := context.WithTimeout(ctx, timeout)
	defer connectCtxCancel()

	// connect
	mdb, err := mongo.Connect(connectCtx, options.Client().ApplyURI(mongoURI))
	if err != nil {
		return nil, err
	}

	// ping mongo to validate connection
	if err := pingMongo(connectCtx, mdb, timeout); err != nil {
		return nil, err
	}

	return mdb, nil
}

func disconnectMongo(mdb *mongo.Client, timeout time.Duration) error {
	disconnectCtx, disconnectCtxCancel := context.WithTimeout(context.Background(), timeout)
	defer disconnectCtxCancel()
	return mdb.Disconnect(disconnectCtx)
}

func pingMongo(ctx context.Context, mdb *mongo.Client, timeout time.Duration) error {
	pingCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	return mdb.Ping(pingCtx, nil)
}
