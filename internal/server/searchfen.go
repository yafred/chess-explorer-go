package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/notnil/chess"
	"github.com/spf13/viper"
	"github.com/yafred/chess-explorer/internal/pgntodb"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func searchFentHandler(w http.ResponseWriter, r *http.Request) {
	defer timeTrack(time.Now(), "searchFentHandler")

	switch r.Method {
	case "POST":
		// Call ParseForm() to parse the raw query and update r.PostForm and r.Form.
		if err := r.ParseForm(); err != nil {
			fmt.Fprintf(w, "ParseForm() err: %v", err)
			return
		}
	default:
		fmt.Fprintf(w, "Sorry, only POST methods is supported.")
		return
	}

	// allow cross origin
	w.Header().Set("Access-Control-Allow-Origin", "*")

	// create game filter
	filter := gameFilterFromRequest(r)
	gameFilterBson := bsonFromGameFilter(filter)

	fen := strings.TrimSpace(r.FormValue("fen"))
	maxMoves := 40

	go searchFEN(fen, maxMoves, gameFilterBson)
}

func searchFEN(fen string, maxMoves int, gameFilterBson primitive.M) {
	log.Println("Searching for FEN: " + fen)

	// Connect to DB
	client, err := mongo.NewClient(options.Client().ApplyURI(viper.GetString("mongo-url")))
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
		log.Fatal("Cannot connect to DB " + viper.GetString("mongo-url"))
	}

	gamesCollection := client.Database(viper.GetString("mongo-db-name")).Collection("games")

	cur, error := gamesCollection.Find(ctx, gameFilterBson)
	if error != nil {
		log.Fatal(err)
	}

	concurrency := 20
	sem := make(chan bool, concurrency)

	count := 0
	for cur.Next(context.TODO()) {
		var gameHolder pgntodb.Game
		err := cur.Decode(&gameHolder)

		sem <- true // take a slot
		go replay(gameHolder, fen, maxMoves, sem)

		if err != nil {
			log.Fatal(err)
		}
		count++
	}

	// wait for everything to be finished
	for i := 0; i < cap(sem); i++ {
		sem <- true
	}

	log.Println("replayed", count, "games")
}

func replay(game pgntodb.Game, fen string, maxMoves int, sem chan bool) {

	defer func() { <-sem }() // release the slot when finished

	// Process game.PGN (remove "1." etc)
	var pgnMoves []string
	if len(game.PGN) > 0 {
		pgnMoves = strings.Split(game.PGN, " ")
	}

	i := 0 // output index
	for _, x := range pgnMoves {
		if !strings.HasSuffix(x, ".") {
			// copy and increment index
			pgnMoves[i] = x
			i++
		}
	}
	pgnMoves = pgnMoves[:i] // strip final result

	// Replay game
	chessGame := chess.NewGame()
	iMove := 0
	for _, move := range pgnMoves {
		chessGame.MoveStr(move)

		// Compare
		if chessGame.Position().String() == fen {
			iMove++
			log.Println("move", iMove, "in game "+game.Link) // should synchronize results (and ouput number of hits in the end)
			break
		}

		iMove++
		if iMove == maxMoves { // should make that configurable
			break
		}
	}
}
