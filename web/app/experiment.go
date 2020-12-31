package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/yafred/chess-stat/internal/pgntodb"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func testHandler(w http.ResponseWriter, r *http.Request) {
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

	// iterate all documents
	cursor, err := games.Find(ctx, bson.M{})
	if err != nil {
		log.Fatal(err)
	}
	defer cursor.Close(ctx)

	for i := 0; cursor.Next(ctx); i++ {
		var game pgntodb.Game
		if err = cursor.Decode(&game); err != nil {
			log.Fatal(err)
		}
		if strings.HasPrefix(game.PGN, "1. d4 e5") {
			fmt.Println(game.Link)
		}
	}

	// Get distinct first moves in all documents
	values, err := games.Distinct(ctx, "movew01", bson.M{})

	if err != nil {
		log.Fatal(err)
	}

	for _, value := range values {
		fmt.Println(value)
	}

	// Same with aggregation (to provide counts)
	pipeline := make([]bson.M, 0)

	/*
		matchStage := bson.M{
			"$match": bson.M{
				"white": "fredo599",
				"site":  "Chess.com",
			},
		}

		groupStage := bson.M{
			"$group": bson.M{
				"_id":   bson.M{"movew01": "$movew01"},
				"total": bson.M{"$sum": 1},
			},
		}

		pipeline = append(pipeline, matchStage, groupStage)
	*/

	matchStage := bson.M{
		"$match": bson.M{
			"white": "fredo599",
			"site":  "Chess.com",
		},
	}

	groupStage := bson.M{
		"$group": bson.M{
			"_id":    bson.M{"movew01": "$movew01", "result": "$result"},
			"result": bson.M{"$push": "$result"},
			"total":  bson.M{"$sum": 1},
		},
	}

	subGroupStage := bson.M{
		"$group": bson.M{
			"_id":    bson.M{"movew01": "$_id.movew01"},
			"result": bson.M{"$addToSet": bson.M{"result": "$_id.result", "sum": "$total"}},
		},
	}

	pipeline = append(pipeline, matchStage, groupStage, subGroupStage)

	showInfoCursor, err := games.Aggregate(ctx, pipeline)
	if err != nil {
		log.Fatal(err)
	}
	var showsWithInfo []bson.M
	if err = showInfoCursor.All(ctx, &showsWithInfo); err != nil {
		log.Fatal(err)
	}
	fmt.Println(showsWithInfo)
	fmt.Println(len(showsWithInfo))
}
