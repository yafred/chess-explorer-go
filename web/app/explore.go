package server

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/yafred/chess-stat/internal/pgntodb"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func exploreHandler(w http.ResponseWriter, r *http.Request) {

	type Result struct {
		Result string `json:"result,omitempty"`
		Sum    uint16 `json:"sum,omitempty"`
	}
	type Exploration struct {
		Move    string `json:"move,omitempty"`
		Move01  string `json:"move01,omitempty"`
		Move02  string `json:"move02,omitempty"`
		Move03  string `json:"move03,omitempty"`
		Move04  string `json:"move04,omitempty"`
		Move05  string `json:"move05,omitempty"`
		Move06  string `json:"move06,omitempty"`
		Move07  string `json:"move07,omitempty"`
		Move08  string `json:"move08,omitempty"`
		Move09  string `json:"move09,omitempty"`
		Move10  string `json:"move10,omitempty"`
		Total   uint16 `json:"total,omitempty"`
		Link    string `json:"link,omitempty"` // when Total = 1
		Results []Result
	}

	var explorations []Exploration

	pgn := ""

	switch r.Method {
	case "POST":
		// Call ParseForm() to parse the raw query and update r.PostForm and r.Form.
		if err := r.ParseForm(); err != nil {
			fmt.Fprintf(w, "ParseForm() err: %v", err)
			return
		}
		pgn = strings.TrimSpace(r.FormValue("pgn"))
	default:
		fmt.Fprintf(w, "Sorry, only POST methods is supported.")
		return
	}

	// Process input pgn (remove "1." etc)
	var pgnMoves []string
	if len(pgn) > 0 {
		pgnMoves = strings.Split(pgn, " ")
	}

	i := 0 // output index
	for _, x := range pgnMoves {
		if !strings.HasSuffix(x, ".") {
			// copy and increment index
			pgnMoves[i] = x
			i++
		}
	}
	pgnMoves = pgnMoves[:i]

	// Our logic allows input pgn to have 0 to 9 moves
	if len(pgnMoves) > 9 {
		json.NewEncoder(w).Encode(explorations) // empty array
		return
	}

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

	var andClause []bson.M
	andClause = append(andClause, bson.M{"site": site})
	andClause = append(andClause, bson.M{"white": player})

	// filter on previous moves
	for i := 1; i < len(pgnMoves)+1; i++ {
		moveField := buildMoveFieldName(i)
		andClause = append(andClause, bson.M{moveField: pgnMoves[i-1]})
	}

	// move field for aggregate
	fieldNum := len(pgnMoves) + 1
	moveField := buildMoveFieldName(fieldNum)

	// make sure next move exists
	andClause = append(andClause, bson.M{moveField: bson.M{"$exists": true, "$ne": ""}})

	pipeline = append(pipeline, bson.M{"$match": bson.M{"$and": andClause}})

	groupStage := bson.M{
		"$group": bson.M{
			"_id":    bson.M{moveField: "$" + moveField, "result": "$result"},
			"total":  bson.M{"$sum": 1},
			"result": bson.M{"$push": "$result"},
		},
	}
	pipeline = append(pipeline, groupStage)

	subGroupStage := bson.M{
		"$group": bson.M{
			"_id":     bson.M{moveField: "$_id." + moveField},
			"results": bson.M{"$addToSet": bson.M{"result": "$_id.result", "sum": "$total"}},
		},
	}
	pipeline = append(pipeline, subGroupStage)

	projectStage := bson.M{
		"$project": bson.M{
			"_id":     false,
			"move":    "$_id." + moveField,
			"results": "$results",
		},
	}
	pipeline = append(pipeline, projectStage)

	aggregateCursor, err := games.Aggregate(ctx, pipeline)
	if err != nil {
		log.Fatal(err)
	}

	defer aggregateCursor.Close(ctx)

	if err = aggregateCursor.All(ctx, &explorations); err != nil {
		log.Fatal(err)
	}

	// add a total
	for iExploration, x := range explorations {
		var total uint16
		var move string
		for _, y := range x.Results {
			total += y.Sum
			move = x.Move
		}
		explorations[iExploration].Total = total
		if total == 1 {
			// get link for moves pgn + move
			game := getGame(ctx, games, pgnMoves, move)
			if game != nil {
				explorations[iExploration].Link = game.Link
			}
		}
		sort.Slice(explorations[iExploration].Results, func(i, j int) bool {
			if explorations[iExploration].Results[i].Result == "1-0" {
				return true
			} else if explorations[iExploration].Results[i].Result == "0-1" {
				return false
			} else if explorations[iExploration].Results[j].Result == "1-0" {
				return false
			} else {
				return true
			}
		})
	}

	// sort by counts
	sort.Slice(explorations, func(i, j int) bool {
		return explorations[i].Total > explorations[j].Total
	})

	// send the response
	json.NewEncoder(w).Encode(explorations)
}

func buildMoveFieldName(fieldNum int) (moveField string) {
	moveField = "move"
	if fieldNum < 10 {
		moveField = moveField + "0"
	}
	moveField = moveField + strconv.Itoa(fieldNum)
	return moveField
}

func getGame(ctx context.Context, games *mongo.Collection, pgnMoves []string, move string) (game *pgntodb.Game) {
	var andClause []bson.M
	andClause = append(andClause, bson.M{"site": "Chess.com"})
	andClause = append(andClause, bson.M{"white": "fredo599"})
	for i := 0; i < len(pgnMoves); i++ {
		andClause = append(andClause, bson.M{buildMoveFieldName(i + 1): pgnMoves[i]})
	}
	andClause = append(andClause, bson.M{buildMoveFieldName(len(pgnMoves) + 1): move})

	cursor, err := games.Find(ctx, bson.M{"$and": andClause})
	defer cursor.Close(ctx)
	if err != nil {
		log.Fatal(err)
	}

	var resultGame []pgntodb.Game
	err = cursor.All(ctx, &resultGame)
	if err != nil {
		log.Fatal(err)
	}

	var ret *pgntodb.Game
	if len(resultGame) != 0 {
		return &resultGame[0]
	}
	return ret
}
