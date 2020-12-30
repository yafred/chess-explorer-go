package pgntodb

import (
	"context"
	"log"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// Game ... for the database
type game struct {
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
}

var client *mongo.Client

func insertGame(gameMap map[string]string, client *mongo.Client) {

	// Clean up data
	if strings.Index(gameMap["Site"], "lichess.org") != -1 {
		gameMap["Link"] = gameMap["Site"]
		gameMap["Site"] = "Lichess.org"
	}

	game := game{
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
