package lichess

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/spf13/viper"
	"github.com/yafred/chess-explorer/internal/pgntodb"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

/*
TODO:
- lichess user/password to decrease throttle
- search for most recent game in DB and use parm 'since' in request
*/

// DownloadGames ... Downloads games from Chess.com for user {user}
// https://lichess.org/api#operation/apiGamesUser
func DownloadGames(username string) {

	url := "https://lichess.org/api/games/user/" + username

	chessClient := &http.Client{Timeout: 10 * time.Second}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal(err)
	}

	// If there is a token in the configuration, use it
	lichessToken := viper.GetString("lichess-token")
	if lichessToken != "" {
		log.Println(lichessToken)
		req.Header.Add("Authorization", "Bearer "+lichessToken)
	}

	// Get most recent game to set 'since' if possible
	since := getLastUnixTime(username)

	q := req.URL.Query()

	if since != 0 {
		since += 1000 // add 1 sec to avoid downloading the last game we have
		q.Add("since", strconv.FormatInt(since, 10))
	}

	req.URL.RawQuery = q.Encode()

	fmt.Println(req.URL.String())

	// Get data
	resp, err := chessClient.Do(req)

	if err != nil {
		log.Fatal(err)
	}

	// Create a temp file
	tmpfile, err := ioutil.TempFile("", "lichess")
	if err != nil {
		log.Fatal(err)
	}
	defer os.Remove(tmpfile.Name()) // clean up
	log.Println(tmpfile.Name())
	// Create the file
	out, err := os.Create(tmpfile.Name())
	if err != nil {
		log.Fatal(err)
	}
	defer out.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	pgntodb.Process(tmpfile.Name())
}

// get the date of the most recent game of {username} as a millisec timestamp
func getLastUnixTime(username string) int64 {
	// Connect to DB
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://127.0.0.1:27017"))
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

	siteBson := bson.M{"site": "Lichess.org"}
	userBson := make([]bson.M, 0)
	userBson = append(userBson, bson.M{"white": username})
	userBson = append(userBson, bson.M{"black": username})
	finalBson := make([]bson.M, 0)
	finalBson = append(finalBson, siteBson)
	finalBson = append(finalBson, bson.M{"$or": userBson})

	queryOptions := options.FindOneOptions{}
	queryOptions.SetSort(bson.M{"datetime": -1})

	game := pgntodb.Game{}
	error := games.FindOne(context.TODO(), bson.M{"$and": finalBson}, &queryOptions).Decode(&game)

	var ret = int64(0)
	if error != nil {
		log.Println(error)
	} else {
		ret = game.DateTime.UnixNano() / int64(time.Millisecond)
	}

	return ret
}
