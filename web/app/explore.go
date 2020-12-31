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

func exploreHandler(w http.ResponseWriter, r *http.Request) {

	// Connect to DB
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
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

	// Distinct moves with counts
	matchStage := bson.M{
		"$match": bson.M{
			"white": "fredo599",
			"site":  "Chess.com",
		},
	}

	groupStage := bson.M{
		"$group": bson.M{
			"_id":    bson.M{"movew01": "$movew01", "result": "$result"},
			"total":  bson.M{"$sum": 1},
			"result": bson.M{"$push": "$result"},
		},
	}

	subGroupStage := bson.M{
		"$group": bson.M{
			"_id":     bson.M{"movew01": "$_id.movew01"},
			"results": bson.M{"$addToSet": bson.M{"result": "$_id.result", "sum": "$total"}},
		},
	}

	projectStage := bson.M{
		"$project": bson.M{
			"_id":     false,
			"movew01": "$_id.movew01",
			"results": "$results",
		},
	}

	pipeline := make([]bson.M, 0)
	pipeline = append(pipeline, matchStage, groupStage, subGroupStage, projectStage)

	aggregateCursor, err := games.Aggregate(ctx, pipeline)
	if err != nil {
		log.Fatal(err)
	}

	defer aggregateCursor.Close(ctx)

	type Result struct {
		Result string `json:"result,omitempty"`
		Sum    uint16 `json:"sum,omitempty"`
	}
	type Exploration struct {
		MoveW01 string `json:"movew01,omitempty"`
		Results []Result
	}

	var explorations []Exploration
	if err = aggregateCursor.All(ctx, &explorations); err != nil {
		log.Fatal(err)
	}

	// send the response
	json.NewEncoder(w).Encode(explorations)
}
