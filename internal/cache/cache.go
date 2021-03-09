package cache

import (
	"context"
	"log"
	"time"

	"github.com/spf13/viper"
	"github.com/yafred/chess-explorer/internal/pgntodb"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type result struct {
	Name  string `json:"name"`
	Count int    `json:"count"`
}

type userResult struct {
	SiteName string `json:"sitename"`
	Name     string `json:"name"`
	Count    int    `json:"count"`
}

// UpdateInitialValues ... performs a few heavy weight aggregations and keep them in database
func UpdateInitialValues() {
	// Connect to DB
	client, err := mongo.NewClient(options.Client().ApplyURI(viper.GetString("mongo-url")))
	if err != nil {
		log.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(ctx)

	// Ping MongoDB
	if err = client.Ping(ctx, readpref.Primary()); err != nil {
		log.Fatal("Cannot connect to DB " + viper.GetString("mongo-url"))
	}

	games := client.Database(viper.GetString("mongo-db-name")).Collection("games")

	// Timecontrols
	filter := bson.M{"$match": bson.M{}}
	pipeline := make([]bson.M, 0)
	pipeline = append(pipeline, filter)

	groupStage := bson.M{
		"$group": bson.M{
			"_id":   bson.M{"timecontrol": "$timecontrol"},
			"count": bson.M{"$sum": 1},
		},
	}
	pipeline = append(pipeline, groupStage)

	sortStage := bson.M{
		"$sort": bson.M{
			"count": -1,
		},
	}
	pipeline = append(pipeline, sortStage)

	projectStage := bson.M{
		"$project": bson.M{
			"_id":   false,
			"name":  "$_id.timecontrol",
			"count": true,
		},
	}
	pipeline = append(pipeline, projectStage)

	aggregateCursor, err := games.Aggregate(ctx, pipeline)
	if err != nil {
		log.Fatal(err)
	}

	defer aggregateCursor.Close(ctx)

	var timeControlResults []result
	if err = aggregateCursor.All(ctx, &timeControlResults); err != nil {
		log.Fatal(err)
	}

	//log.Println(timeControlResults)

	// Users
	lastgames := client.Database(viper.GetString("mongo-db-name")).Collection("lastgames")

	cursor, err := lastgames.Find(ctx, bson.M{})
	if err != nil {
		log.Fatal(err)
	}
	var results []pgntodb.LastGame
	if err = cursor.All(ctx, &results); err != nil {
		log.Fatal(err)
	}

	var userResults = make([]userResult, 0)
	for _, aUser := range results {
		filter := bson.M{"site": aUser.Site, "$or": []bson.M{{"white": aUser.Username}, {"black": aUser.Username}}}
		count, err := games.CountDocuments(ctx, filter)
		if err != nil {
			log.Fatal(err)
		}
		userResults = append(userResults, userResult{SiteName: aUser.Site, Name: aUser.Username, Count: int(count)})
	}

	//log.Println(userResults)
}
