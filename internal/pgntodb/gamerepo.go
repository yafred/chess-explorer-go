package pgntodb

import (
	"context"
	"log"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// Game ... for the database
type Game struct {
	Site        string `json:"site,omitempty"`
	White       string `json:"white,omitempty"`
	Black       string `json:"black,omitempty"`
	UTCDate     string `json:"utcdate,omitempty"`
	UTCTime     string `json:"utctime,omitempty"`
	Result      string `json:"result,omitempty"`
	WhiteElo    string `json:"whiteelo,omitempty"`
	BlackElo    string `json:"blackelo,omitempty"`
	TimeControl string `json:"timecontrol,omitempty"`
	Link        string `json:"link,omitempty"`
	PGN         string `json:"pgn,omitempty"`
	Move01      string `json:"move01,omitempty"`
	Move02      string `json:"move02,omitempty"`
	Move03      string `json:"move03,omitempty"`
	Move04      string `json:"move04,omitempty"`
	Move05      string `json:"move05,omitempty"`
	Move06      string `json:"move06,omitempty"`
	Move07      string `json:"move07,omitempty"`
	Move08      string `json:"move08,omitempty"`
	Move09      string `json:"move09,omitempty"`
	Move10      string `json:"move10,omitempty"`
}

var client *mongo.Client

func insertGame(gameMap map[string]string, client *mongo.Client) {

	// Clean up data
	if strings.Index(gameMap["Site"], "lichess.org") != -1 {
		gameMap["Link"] = gameMap["Site"]
		gameMap["Site"] = "Lichess.org"
	}

	game := Game{
		Site:        gameMap["Site"],
		White:       gameMap["White"],
		Black:       gameMap["Black"],
		UTCDate:     gameMap["UTCDate"],
		UTCTime:     gameMap["UTCTime"],
		Result:      gameMap["Result"],
		WhiteElo:    gameMap["WhiteElo"],
		BlackElo:    gameMap["BlackElo"],
		TimeControl: gameMap["TimeControl"],
		Link:        gameMap["Link"],
		PGN:         gameMap["PGN"],
	}

	// Itemize first moves of the pgn
	itemizePgn(&game)

	// Look for a duplicate before inserting
	games := client.Database("chess-explorer").Collection("games")

	count, error := games.CountDocuments(context.TODO(), bson.M{"white": game.White, "black": game.Black, "utcdate": game.UTCDate, "utctime": game.UTCTime})

	// Insert
	if count == 0 && error == nil {
		_, err := games.InsertOne(context.TODO(), game)

		if err != nil {
			log.Fatal(err)
		}
	}
}

// Remider: last item of the pgn is "0-1" or "1-0" or "1/2-1/2"
func itemizePgn(game *Game) {
	pgn := game.PGN
	pgnElements := strings.Split(pgn, " ")
	if len(pgnElements) > 2 {
		game.Move01 = pgnElements[1]
	}
	if len(pgnElements) > 3 {
		game.Move02 = pgnElements[2]
	}
	if len(pgnElements) > 5 {
		game.Move03 = pgnElements[4]
	}
	if len(pgnElements) > 6 {
		game.Move04 = pgnElements[5]
	}
	if len(pgnElements) > 8 {
		game.Move05 = pgnElements[7]
	}
	if len(pgnElements) > 9 {
		game.Move06 = pgnElements[8]
	}
	if len(pgnElements) > 11 {
		game.Move07 = pgnElements[10]
	}
	if len(pgnElements) > 12 {
		game.Move08 = pgnElements[11]
	}
	if len(pgnElements) > 14 {
		game.Move09 = pgnElements[13]
	}
	if len(pgnElements) > 15 {
		game.Move10 = pgnElements[14]
	}
}
