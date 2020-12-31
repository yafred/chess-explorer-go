package server

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func exploreHandler(w http.ResponseWriter, r *http.Request) {

	pgn := ""

	switch r.Method {
	case "POST":
		// Call ParseForm() to parse the raw query and update r.PostForm and r.Form.
		if err := r.ParseForm(); err != nil {
			fmt.Fprintf(w, "ParseForm() err: %v", err)
			return
		}
		pgn = r.FormValue("pgn")
		log.Println(pgn)
	default:
		fmt.Fprintf(w, "Sorry, only POST methods is supported.")
		return
	}

	// Process input pgn (remove "1." etc)
	pgnMoves := strings.Split(pgn, " ")

	i := 0 // output index
	for _, x := range pgnMoves {
		if !strings.HasSuffix(x, ".") {
			// copy and increment index
			pgnMoves[i] = x
			i++
		}
	}
	pgnMoves = pgnMoves[:i]
	log.Println(pgnMoves)

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
	player := "fredo599"
	site := "Chess.com"
	pipeline := make([]bson.M, 0)

	matchStage := bson.M{
		"$match": bson.M{
			"white": player,
		},
	}
	pipeline = append(pipeline, matchStage)

	matchStage2 := bson.M{
		"$match": bson.M{
			"site": site,
		},
	}
	pipeline = append(pipeline, matchStage2)

	groupStage := bson.M{
		"$group": bson.M{
			"_id":    bson.M{"movew01": "$movew01", "result": "$result"},
			"total":  bson.M{"$sum": 1},
			"result": bson.M{"$push": "$result"},
		},
	}
	pipeline = append(pipeline, groupStage)

	subGroupStage := bson.M{
		"$group": bson.M{
			"_id":     bson.M{"movew01": "$_id.movew01"},
			"results": bson.M{"$addToSet": bson.M{"result": "$_id.result", "sum": "$total"}},
		},
	}
	pipeline = append(pipeline, subGroupStage)

	projectStage := bson.M{
		"$project": bson.M{
			"_id":     false,
			"movew01": "$_id.movew01",
			"results": "$results",
		},
	}
	pipeline = append(pipeline, projectStage)

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
