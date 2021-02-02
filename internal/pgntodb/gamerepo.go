package pgntodb

import (
	"context"
	"log"
	"strconv"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// LastGame ... last game (in this database) for a player
type LastGame struct {
	Username string    `json:"username" bson:"username"`
	Site     string    `json:"site" bson:"site"`
	DateTime time.Time `json:"datetime" bson:"datetime"`
	GameID   string    `json:"gameid" bson:"gameid"`
	Logged   string    `json:"logged,omitempty" bson:"logged,omitempty"` // not going to database
}

// Game ... for the database
type Game struct {
	ID          string    `json:"_id" bson:"_id"`
	Site        string    `json:"site,omitempty"`
	White       string    `json:"white,omitempty"`
	Black       string    `json:"black,omitempty"`
	DateTime    time.Time `json:"datetime,omitempty"`
	Result      string    `json:"result,omitempty"`
	WhiteElo    uint16    `json:"whiteelo,omitempty"`
	BlackElo    uint16    `json:"blackelo,omitempty"`
	TimeControl string    `json:"timecontrol,omitempty"`
	Link        string    `json:"link,omitempty"`
	PGN         string    `json:"pgn,omitempty"`
	Move01      string    `json:"m01,omitempty" bson:"m01,omitempty"`
	Move02      string    `json:"m02,omitempty" bson:"m02,omitempty"`
	Move03      string    `json:"m03,omitempty" bson:"m03,omitempty"`
	Move04      string    `json:"m04,omitempty" bson:"m04,omitempty"`
	Move05      string    `json:"m05,omitempty" bson:"m05,omitempty"`
	Move06      string    `json:"m06,omitempty" bson:"m06,omitempty"`
	Move07      string    `json:"m07,omitempty" bson:"m07,omitempty"`
	Move08      string    `json:"m08,omitempty" bson:"m08,omitempty"`
	Move09      string    `json:"m09,omitempty" bson:"m09,omitempty"`
	Move10      string    `json:"m10,omitempty" bson:"m10,omitempty"`
	Move11      string    `json:"m11,omitempty" bson:"m11,omitempty"`
	Move12      string    `json:"m12,omitempty" bson:"m12,omitempty"`
	Move13      string    `json:"m13,omitempty" bson:"m13,omitempty"`
	Move14      string    `json:"m14,omitempty" bson:"m14,omitempty"`
	Move15      string    `json:"m15,omitempty" bson:"m15,omitempty"`
	Move16      string    `json:"m16,omitempty" bson:"m16,omitempty"`
	Move17      string    `json:"m17,omitempty" bson:"m17,omitempty"`
	Move18      string    `json:"m18,omitempty" bson:"m18,omitempty"`
	Move19      string    `json:"m19,omitempty" bson:"m19,omitempty"`
	Move20      string    `json:"m20,omitempty" bson:"m20,omitempty"`
}

var client *mongo.Client

var queue []interface{} // queue for insert many

// FindLastGame ... find last game (allowing prevention of duplicates)
func findLastGame(username string, site string, client *mongo.Client) *LastGame {
	lastGame := LastGame{
		Site:     site,
		Username: username,
	}

	lastgames := client.Database("chess-explorer").Collection("lastgames")
	filter := bson.M{"site": site, "username": username}
	collation := options.Collation{Locale: "en", Strength: 2}
	findOneOptions := options.FindOneOptions{Collation: &collation} // case insensitive search

	result := lastgames.FindOne(context.TODO(), filter, &findOneOptions)

	if result != nil {
		result.Decode(&lastGame)
	}

	return &lastGame
}

func logLastGame(username string, game Game, client *mongo.Client) {
	if username != "" {
		if strings.ToLower(username) == strings.ToLower(game.White) {
			username = game.White
		} else if strings.ToLower(username) == strings.ToLower(game.Black) {
			username = game.Black
		} else {
			log.Fatal("username "+username+"is not a player of ", game)
		}

		lastGame := LastGame{
			Username: username,
			Site:     game.Site,
			DateTime: game.DateTime,
			GameID:   game.ID,
		}

		lastgames := client.Database("chess-explorer").Collection("lastgames")
		filter := bson.M{"site": game.Site, "username": username}
		updateOptions := options.Update().SetUpsert(true)
		update := bson.M{
			"$set": lastGame,
		}

		// Insert
		_, error := lastgames.UpdateOne(context.TODO(), filter, update, updateOptions)

		if error != nil {
			log.Fatal(error)
		}

		log.Println("Most recent game is now: " + lastGame.GameID)
	}
}

func pushGame(gameMap map[string]string, client *mongo.Client, lastGame *LastGame) bool {
	game := Game{}
	mapToGame(gameMap, &game)
	queue = append(queue, game)
	if len(queue) > 9999 {
		return flushGames(client, lastGame)
	}
	return true
}

func flushGames(client *mongo.Client, lastGame *LastGame) bool {
	log.Println("Flushing " + strconv.Itoa(len(queue)) + " games to DB")
	if len(queue) > 0 {
		games := client.Database("chess-explorer").Collection("games")

		insertManyOptions := options.InsertMany().SetOrdered(false) // continue if duplicates are found
		_, error := games.InsertMany(context.TODO(), queue, insertManyOptions)

		if error != nil {
			//log.Println(error)
			//log.Println("It is possible to have duplicate key errors when importing games for a user who has played again a user we already have games for).")
		}
		if lastGame.Logged == "" {
			logLastGame(lastGame.Username, queue[0].(Game), client)
			lastGame.Logged = "Done"
		}
	}

	queue = queue[:0]
	return true
}

func mapToGame(gameMap map[string]string, game *Game) {
	// Clean up data
	if strings.Index(gameMap["Site"], "lichess.org") != -1 {
		gameMap["Link"] = gameMap["Site"]
		gameMap["Site"] = "lichess.org"
	}
	gameMap["Site"] = strings.ToLower(gameMap["Site"])

	whiteelo := 0
	blackelo := 0
	var error error
	if gameMap["WhiteElo"] != "" && strings.Index(gameMap["WhiteElo"], "?") == -1 {
		whiteelo, error = strconv.Atoi(gameMap["WhiteElo"])
		if error != nil {
			log.Fatal("Not a valid ELO: " + gameMap["WhiteElo"] + " for white " + gameMap["White"])
		}
	}
	if gameMap["BlackElo"] != "" && strings.Index(gameMap["BlackElo"], "?") == -1 {
		blackelo, error = strconv.Atoi(gameMap["BlackElo"])
		if error != nil {
			log.Fatal("Not a valid ELO: " + gameMap["BlackElo"] + " for black " + gameMap["Black"])
		}
	}

	game.ID = createGameID(gameMap)
	game.Site = gameMap["Site"]
	game.White = gameMap["White"]
	game.Black = gameMap["Black"]
	game.DateTime = createDateTime(gameMap)
	game.Result = gameMap["Result"]
	game.WhiteElo = uint16(whiteelo)
	game.BlackElo = uint16(blackelo)
	game.TimeControl = gameMap["TimeControl"]
	game.Link = gameMap["Link"]
	game.PGN = gameMap["PGN"]

	// Itemize first moves of the pgn
	itemizePgn(game)
}

func createDateTime(gameMap map[string]string) time.Time {
	// Create a time.Time object
	utcDate := strings.ReplaceAll(gameMap["UTCDate"], ".", "-")
	dateTimeAsUTCString := utcDate + "T" + gameMap["UTCTime"] + "+00:00"

	dateTime, error := time.Parse(time.RFC3339, dateTimeAsUTCString)
	if error != nil {
		log.Fatal("Not a valid date: " + dateTimeAsUTCString)
	}
	return dateTime
}

func createGameID(gameMap map[string]string) string {
	return strings.ToLower(gameMap["Site"]) + ":" + gameMap["White"] + ":" + gameMap["Black"] + ":" + gameMap["UTCDate"] + ":" + gameMap["UTCTime"]
}

// Reminder: last item of the pgn is "0-1" or "1-0" or "1/2-1/2" (for len(pgnElements) test)
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
