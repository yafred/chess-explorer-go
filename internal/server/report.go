package server

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type result struct {
	Name  string `json:"name"`
	Count int    `json:"count"`
}

type report struct {
	TotalGames   int64 `json:"totalgames,omitempty"`
	Sites        []result
	UsersAsWhite []result
	TimeControls []result
}

func reportHandler(w http.ResponseWriter, r *http.Request) {
	report := report{}

	// Connect to DB
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://127.0.0.1:27017"))
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
		log.Fatal("Cannot connect to DB")
	}

	games := client.Database("chess-explorer").Collection("games")

	// Total games
	totalGames, error := games.CountDocuments(ctx, bson.M{})
	if error != nil {
		log.Fatal(error)
	}
	report.TotalGames = totalGames

	reportGames(ctx, games, &report)
	reportSites(ctx, games, &report)
	reportUsersAsWhite(ctx, games, &report)
	reportTimeControls(ctx, games, &report)

	// send the response
	json.NewEncoder(w).Encode(report)
}

// Games
func reportGames(ctx context.Context, games *mongo.Collection, report *report) {
	totalGames, error := games.CountDocuments(ctx, bson.M{})
	if error != nil {
		log.Fatal(error)
	}
	report.TotalGames = totalGames
}

// Distinct sites
func reportSites(ctx context.Context, games *mongo.Collection, report *report) {
	filter := bson.M{"$match": bson.M{}}
	pipeline := make([]bson.M, 0)
	pipeline = append(pipeline, filter)

	groupStage := bson.M{
		"$group": bson.M{
			"_id":   bson.M{"site": "$site"},
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
			"name":  "$_id.site",
			"count": "$count",
		},
	}
	pipeline = append(pipeline, projectStage)

	aggregateCursor, err := games.Aggregate(ctx, pipeline)
	if err != nil {
		log.Fatal(err)
	}

	defer aggregateCursor.Close(ctx)

	var siteResults []result
	if err = aggregateCursor.All(ctx, &siteResults); err != nil {
		log.Fatal(err)
	}

	report.Sites = siteResults
}

// Distinct users
func reportUsersAsWhite(ctx context.Context, games *mongo.Collection, report *report) {
	filter := bson.M{"$match": bson.M{}}
	pipeline := make([]bson.M, 0)
	pipeline = append(pipeline, filter)

	groupStage := bson.M{
		"$group": bson.M{
			"_id":   bson.M{"white": "$white"},
			"count": bson.M{"$sum": 1},
		},
	}
	pipeline = append(pipeline, groupStage)

	reduceStage := bson.M{"$match": bson.M{
		"count": bson.M{"$gte": 10},
	}}
	pipeline = append(pipeline, reduceStage)

	sortStage := bson.M{
		"$sort": bson.M{
			"count": -1,
		},
	}
	pipeline = append(pipeline, sortStage)

	projectStage := bson.M{
		"$project": bson.M{
			"_id":   false,
			"name":  "$_id.white",
			"count": "$count",
		},
	}
	pipeline = append(pipeline, projectStage)

	aggregateCursor, err := games.Aggregate(ctx, pipeline)
	if err != nil {
		log.Fatal(err)
	}

	defer aggregateCursor.Close(ctx)

	var usersAsWhiteResult []result
	if err = aggregateCursor.All(ctx, &usersAsWhiteResult); err != nil {
		log.Fatal(err)
	}

	report.UsersAsWhite = usersAsWhiteResult
}

// Distinct users
func reportTimeControls(ctx context.Context, games *mongo.Collection, report *report) {
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
			"count": "$count",
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

	report.TimeControls = timeControlResults
}
