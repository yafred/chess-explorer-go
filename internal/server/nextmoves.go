package server

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/viper"
	"github.com/yafred/chess-explorer/internal/pgntodb"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// GameFilter ... represents the filter form from the UI
type GameFilter struct {
	pgn                 string
	white               string
	black               string
	timecontrol         string
	useLooseTimecontrol string
	from                string
	to                  string
	minelo              string
	maxelo              string
	site                string
}

func nextMovesHandler(w http.ResponseWriter, r *http.Request) {

	// allow cross origin
	w.Header().Set("Access-Control-Allow-Origin", "*")

	type Result struct {
		Result string `json:"result,omitempty"`
		Sum    uint32 `json:"sum,omitempty"`
	}
	type NextMove struct {
		move01  string `bson:"m01,omitempty"`
		move02  string `bson:"m02,omitempty"`
		move03  string `bson:"m03,omitempty"`
		move04  string `bson:"m04,omitempty"`
		move05  string `bson:"m05,omitempty"`
		move06  string `bson:"m06,omitempty"`
		move07  string `bson:"m07,omitempty"`
		move08  string `bson:"m08,omitempty"`
		move09  string `bson:"m09,omitempty"`
		move10  string `bson:"m10,omitempty"`
		move11  string `bson:"m11,omitempty"`
		move12  string `bson:"m12,omitempty"`
		move13  string `bson:"m13,omitempty"`
		move14  string `bson:"m14,omitempty"`
		move15  string `bson:"m15,omitempty"`
		move16  string `bson:"m16,omitempty"`
		move17  string `bson:"m17,omitempty"`
		move18  string `bson:"m18,omitempty"`
		move19  string `bson:"m19,omitempty"`
		move20  string `bson:"m20,omitempty"`
		tmpGame pgntodb.Game
		// Only the fields below go in the response
		Results []Result     `json:"results"`
		Move    string       `json:"move"`
		Win     uint32       `json:"win"`
		Draw    uint32       `json:"draw"`
		Lose    uint32       `json:"lose"`
		Total   uint32       `json:"total"`
		Game    pgntodb.Game `json:"game,omitempty"` // when Total = 1
	}

	var nextmoves []NextMove
	var filter GameFilter
	mongoAggregation := true

	switch r.Method {
	case "POST":
		// Call ParseForm() to parse the raw query and update r.PostForm and r.Form.
		if err := r.ParseForm(); err != nil {
			fmt.Fprintf(w, "ParseForm() err: %v", err)
			return
		}
		filter.pgn = strings.TrimSpace(r.FormValue("pgn"))
		filter.white = strings.TrimSpace(r.FormValue("white"))
		filter.black = strings.TrimSpace(r.FormValue("black"))
		filter.timecontrol = strings.TrimSpace(r.FormValue("timecontrol"))
		filter.useLooseTimecontrol = strings.TrimSpace(r.FormValue("useLooseTimecontrol"))
		filter.from = strings.TrimSpace(r.FormValue("from"))
		filter.to = strings.TrimSpace(r.FormValue("to"))
		filter.minelo = strings.TrimSpace(r.FormValue("minelo"))
		filter.maxelo = strings.TrimSpace(r.FormValue("maxelo"))
		filter.site = strings.ToLower(strings.TrimSpace(r.FormValue("site")))

	default:
		fmt.Fprintf(w, "Sorry, only POST methods is supported.")
		return
	}

	// Process input pgn (remove "1." etc)
	var pgnMoves []string
	if len(filter.pgn) > 0 {
		pgnMoves = strings.Split(filter.pgn, " ")
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

	if len(pgnMoves) < 20 {
		mongoAggregation = true
	} else {
		mongoAggregation = false
	}

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
		log.Fatal("Cannot connect to DB")
	}

	games := client.Database(viper.GetString("mongo-db-name")).Collection("games")

	// Distinct moves with counts
	var andClause []bson.M

	// create game filter
	gameFilterBson := processGameFilter(filter)
	andClause = append(andClause, gameFilterBson)

	if mongoAggregation {
		// Our logic allows input pgn to have 0 to 19 moves
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

		pipeline := make([]bson.M, 0)
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

		if err = aggregateCursor.All(ctx, &nextmoves); err != nil {
			log.Fatal(err)
		}
	} else {
		// algorythmic aggregation
		quotedPgn := regexp.QuoteMeta(filter.pgn)
		andClause = append(andClause, bson.M{"pgn": bson.M{"$regex": quotedPgn}})

		cursor, err := games.Find(ctx, bson.M{"$and": andClause})
		defer cursor.Close(ctx)
		if err != nil {
			log.Fatal(err)
		}

		var resultGames []pgntodb.Game
		err = cursor.All(ctx, &resultGames)
		if err != nil {
			log.Fatal(err)
		}

		filterPgn := strings.Split(filter.pgn, " ")
		for _, game := range resultGames {
			gamePgn := strings.Split(game.PGN, " ")
			gamePgn = gamePgn[0 : len(gamePgn)-1] // remove last bit which is the result
			nextmove := ""
			if len(gamePgn) > len(filterPgn) {
				if strings.HasSuffix(gamePgn[len(filterPgn)], ".") {
					nextmove = gamePgn[len(filterPgn)+1]
				} else {
					nextmove = gamePgn[len(filterPgn)]
				}
			}
			if nextmove != "" {
				foundNextMove := -1
				for iNextMove := range nextmoves {
					if nextmoves[iNextMove].Move == nextmove {
						foundNextMove = iNextMove
						break
					}
				}
				if foundNextMove == -1 {
					nextmoves = append(nextmoves, NextMove{Move: nextmove, Results: make([]Result, 0), tmpGame: game})
					foundNextMove = len(nextmoves) - 1
				}
				foundResult := -1
				for iResult := range nextmoves[foundNextMove].Results {
					if nextmoves[foundNextMove].Results[iResult].Result == game.Result {
						foundResult = iResult
						nextmoves[foundNextMove].Results[iResult].Sum = nextmoves[foundNextMove].Results[iResult].Sum + 1
						break
					}
				}
				if foundResult == -1 {
					nextmoves[foundNextMove].Results = append(nextmoves[foundNextMove].Results, Result{Result: game.Result, Sum: 1})
				}
			}
		}

	}

	// add a total
	for iNextMove := range nextmoves {
		for _, y := range nextmoves[iNextMove].Results {
			if y.Result == "1-0" {
				nextmoves[iNextMove].Win = y.Sum
			} else if y.Result == "0-1" {
				nextmoves[iNextMove].Lose = y.Sum
			} else {
				nextmoves[iNextMove].Draw = y.Sum
			}
		}

		nextmoves[iNextMove].Total = nextmoves[iNextMove].Win + nextmoves[iNextMove].Draw + nextmoves[iNextMove].Lose

		if nextmoves[iNextMove].Total == 1 {
			if mongoAggregation {
				// get link for moves pgn + move
				// Note: this slows down the results if there are a lot of single games (eg: EricRosen)
				game := getGame(ctx, games, pgnMoves, nextmoves[iNextMove].Move, gameFilterBson)
				if game != nil {
					nextmoves[iNextMove].Game = *game
				}
			} else {
				nextmoves[iNextMove].Game = nextmoves[iNextMove].tmpGame
			}
		}
	}

	// sort by counts
	sort.Slice(nextmoves, func(i, j int) bool {
		return nextmoves[i].Total > nextmoves[j].Total
	})

	// send the response
	json.NewEncoder(w).Encode(nextmoves)
}

func buildMoveFieldName(fieldNum int) (moveField string) {
	moveField = "m"
	if fieldNum < 10 {
		moveField = moveField + "0"
	}
	moveField = moveField + strconv.Itoa(fieldNum)
	return moveField
}

func getGame(ctx context.Context, games *mongo.Collection, pgnMoves []string, move string, gameFilterBson bson.M) (game *pgntodb.Game) {
	var andClause []bson.M

	andClause = append(andClause, gameFilterBson)

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

func processGameFilter(filter GameFilter) bson.M {
	ret := bson.M{}

	// Time Control filter
	timeControlBson := make([]bson.M, 0)
	timeControls := strings.Split(filter.timecontrol, ",")
	for _, timeControl := range timeControls {
		if strings.TrimSpace(timeControl) != "" {
			if filter.useLooseTimecontrol == "true" {
				timecontrolParts := strings.Split(strings.TrimSpace(timeControl), "+")
				orQuery := []bson.M{}
				exactQuery := bson.M{"timecontrol": timecontrolParts[0]}
				looseQuery := bson.M{"timecontrol": bson.M{"$regex": "^" + timecontrolParts[0] + "+"}}
				orQuery = append(orQuery, exactQuery, looseQuery)
				timeControlBson = append(timeControlBson, bson.M{"$or": orQuery})
			} else {
				timeControlBson = append(timeControlBson, bson.M{"timecontrol": strings.TrimSpace(timeControl)})
			}
		}
	}

	// Site filter
	siteBson := make([]bson.M, 0)
	sites := strings.Split(filter.site, ",")
	for _, site := range sites {
		if strings.TrimSpace(site) != "" {
			siteBson = append(siteBson, bson.M{"site": strings.TrimSpace(site)})
		}
	}

	// ELO filter
	eloBson := make([]bson.M, 0)

	if filter.minelo != "" {
		minelo, _ := strconv.Atoi(filter.minelo)
		eloBson = append(eloBson, bson.M{
			"whiteelo": bson.M{"$gte": minelo},
			"blackelo": bson.M{"$gte": minelo},
		})
	}

	if filter.maxelo != "" {
		maxelo, _ := strconv.Atoi(filter.maxelo)
		eloBson = append(eloBson, bson.M{
			"whiteelo": bson.M{"$lte": maxelo},
			"blackelo": bson.M{"$lte": maxelo},
		})
	}

	// date filter
	dateBson := make([]bson.M, 0)
	if filter.from != "" {
		fromDate, error := time.Parse(time.RFC3339, filter.from+"T00:00:00+00:00")
		if error != nil {
			log.Print("datetime error " + filter.from)
		} else {
			dateBson = append(dateBson, bson.M{
				"datetime": bson.M{"$gte": fromDate},
			})
		}
	}

	if filter.to != "" {
		toDate, error := time.Parse(time.RFC3339, filter.to+"T23:59:59+00:00")
		if error != nil {
			log.Print("datetime error " + filter.to)
		} else {
			dateBson = append(dateBson, bson.M{
				"datetime": bson.M{"$lte": toDate},
			})
		}
	}

	// user filter
	whiteBson := make([]bson.M, 0)

	// example: c:fred, l:john, alfredo
	whiteUsers := strings.Split(filter.white, ",")
	for _, user := range whiteUsers {
		if strings.TrimSpace(user) == "" {
			break
		}
		splitUser := strings.Split(strings.TrimSpace(user), ":")
		if len(splitUser) > 1 {
			site := convertSite(splitUser[0])
			whiteBson = append(whiteBson, bson.M{"site": site, "white": splitUser[1]})
		} else {
			whiteBson = append(whiteBson, bson.M{"white": splitUser[0]})
		}
	}

	blackBson := make([]bson.M, 0)

	blackUsers := strings.Split(filter.black, ",")
	for _, user := range blackUsers {
		if strings.TrimSpace(user) == "" {
			break
		}
		splitUser := strings.Split(strings.TrimSpace(user), ":")
		if len(splitUser) > 1 {
			site := convertSite(splitUser[0])
			blackBson = append(blackBson, bson.M{"site": site, "black": splitUser[1]})
		} else {
			blackBson = append(blackBson, bson.M{"black": splitUser[0]})
		}
	}

	// gather all filters
	finalBson := make([]bson.M, 0)

	switch len(timeControlBson) {
	case 0:
		break
	case 1:
		finalBson = append(finalBson, timeControlBson[0])
		break
	default:
		finalBson = append(finalBson, bson.M{"$or": timeControlBson})
		break
	}

	switch len(siteBson) {
	case 0:
		break
	case 1:
		finalBson = append(finalBson, siteBson[0])
		break
	default:
		finalBson = append(finalBson, bson.M{"$or": siteBson})
		break
	}

	switch len(eloBson) {
	case 0:
		break
	case 1:
		finalBson = append(finalBson, eloBson[0])
		break
	default:
		finalBson = append(finalBson, bson.M{"$and": eloBson})
		break
	}

	switch len(dateBson) {
	case 0:
		break
	case 1:
		finalBson = append(finalBson, dateBson[0])
		break
	default:
		finalBson = append(finalBson, bson.M{"$and": dateBson})
		break
	}

	switch len(whiteBson) {
	case 0:
		break
	case 1:
		finalBson = append(finalBson, whiteBson[0])
		break
	default:
		finalBson = append(finalBson, bson.M{"$or": whiteBson})
		break
	}

	switch len(blackBson) {
	case 0:
		break
	case 1:
		finalBson = append(finalBson, blackBson[0])
		break
	default:
		finalBson = append(finalBson, bson.M{"$or": blackBson})
		break
	}

	// wrap up
	switch len(finalBson) {
	case 0:
		break
	case 1:
		ret = finalBson[0]
		break
	default:
		ret = bson.M{"$and": finalBson}
	}

	return ret
}

func convertSite(shortName string) string {
	ret := ""
	switch shortName {
	case "c":
		ret = "chess.com"
		break
	case "l":
		ret = "lichess.org"
		break
	default:
		break
	}
	return ret
}
