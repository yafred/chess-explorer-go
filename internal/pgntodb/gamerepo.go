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
	Move01      string `json:"move01,omitempty" bson:"move01,omitempty"`
	Move02      string `json:"move02,omitempty" bson:"move02,omitempty"`
	Move03      string `json:"move03,omitempty" bson:"move03,omitempty"`
	Move04      string `json:"move04,omitempty" bson:"move04,omitempty"`
	Move05      string `json:"move05,omitempty" bson:"move05,omitempty"`
	Move06      string `json:"move06,omitempty" bson:"move06,omitempty"`
	Move07      string `json:"move07,omitempty" bson:"move07,omitempty"`
	Move08      string `json:"move08,omitempty" bson:"move08,omitempty"`
	Move09      string `json:"move09,omitempty" bson:"move09,omitempty"`
	Move10      string `json:"move10,omitempty" bson:"move10,omitempty"`
	Move11      string `json:"move11,omitempty" bson:"move11,omitempty"`
	Move12      string `json:"move12,omitempty" bson:"move12,omitempty"`
	Move13      string `json:"move13,omitempty" bson:"move13,omitempty"`
	Move14      string `json:"move14,omitempty" bson:"move14,omitempty"`
	Move15      string `json:"move15,omitempty" bson:"move15,omitempty"`
	Move16      string `json:"move16,omitempty" bson:"move16,omitempty"`
	Move17      string `json:"move17,omitempty" bson:"move17,omitempty"`
	Move18      string `json:"move18,omitempty" bson:"move18,omitempty"`
	Move19      string `json:"move19,omitempty" bson:"move19,omitempty"`
	Move20      string `json:"move20,omitempty" bson:"move20,omitempty"`
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

// Remider: last item of the pgn is "0-1" or "1-0" or "1/2-1/2" (for len(pgnElements) test)
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
	if len(pgnElements) > 17 {
		game.Move11 = pgnElements[16]
	}
	if len(pgnElements) > 18 {
		game.Move12 = pgnElements[17]
	}
	if len(pgnElements) > 20 {
		game.Move13 = pgnElements[19]
	}
	if len(pgnElements) > 21 {
		game.Move14 = pgnElements[20]
	}
	if len(pgnElements) > 23 {
		game.Move15 = pgnElements[22]
	}
	if len(pgnElements) > 24 {
		game.Move16 = pgnElements[23]
	}
	if len(pgnElements) > 26 {
		game.Move17 = pgnElements[25]
	}
	if len(pgnElements) > 27 {
		game.Move18 = pgnElements[26]
	}
	if len(pgnElements) > 29 {
		game.Move19 = pgnElements[28]
	}
	if len(pgnElements) > 30 {
		game.Move20 = pgnElements[29]
	}
}
