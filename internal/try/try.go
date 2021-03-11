package try

import (
	"context"
	"log"
	"strings"
	"time"

	"github.com/notnil/chess"
	"github.com/spf13/viper"
	"github.com/yafred/chess-explorer/internal/pgntodb"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// Something ... test for me
func Something() {
	log.Println("starting ...")

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

	cur, error := gamesCollection.Find(ctx, bson.D{{}})
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
		go replay(gameHolder, sem)

		if err != nil {
			log.Fatal(err)
		}
		count++
		if count%1000 == 0 {
			log.Println(count)
		}
		if count == 10000 {
			break
		}
	}

	// wait for everything to be finished
	for i := 0; i < cap(sem); i++ {
		sem <- true
	}

	log.Println("read ", count, " records")
}

func replay(game pgntodb.Game, sem chan bool) {

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
		petrovFEN := "rnbqkb1r/pppp1ppp/5n2/4p3/4P3/5N2/PPPP1PPP/RNBQKB1R w KQkq - 2 3 petrov"
		if chessGame.Position().String() == petrovFEN {
			log.Println(game.Link)
			break
		}

		iMove++
		if iMove == 40 {
			break
		}
	}
}
