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
	MoveW01     string `json:"movew01,omitempty"`
	MoveB01     string `json:"moveb01,omitempty"`
	MoveW02     string `json:"movew02,omitempty"`
	MoveB02     string `json:"moveb02,omitempty"`
	MoveW03     string `json:"movew03,omitempty"`
	MoveB03     string `json:"moveb03,omitempty"`
	MoveW04     string `json:"movew04,omitempty"`
	MoveB04     string `json:"moveb04,omitempty"`
	MoveW05     string `json:"movew05,omitempty"`
	MoveB05     string `json:"moveb05,omitempty"`
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
		game.MoveW01 = pgnElements[1]
	}
	if len(pgnElements) > 3 {
		game.MoveB01 = pgnElements[2]
	}
	if len(pgnElements) > 5 {
		game.MoveW02 = pgnElements[4]
	}
	if len(pgnElements) > 6 {
		game.MoveB02 = pgnElements[5]
	}
	if len(pgnElements) > 8 {
		game.MoveW03 = pgnElements[7]
	}
	if len(pgnElements) > 9 {
		game.MoveB03 = pgnElements[8]
	}
	if len(pgnElements) > 11 {
		game.MoveW04 = pgnElements[10]
	}
	if len(pgnElements) > 12 {
		game.MoveB04 = pgnElements[11]
	}
	if len(pgnElements) > 14 {
		game.MoveW05 = pgnElements[13]
	}
	if len(pgnElements) > 15 {
		game.MoveB05 = pgnElements[14]
	}
}
