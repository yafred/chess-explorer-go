package server

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func gamesHandler(w http.ResponseWriter, r *http.Request) {

	type Game struct {
		Pgn string `json:"pgn,omitempty"`
	}

	// allow cross origin
	w.Header().Set("Access-Control-Allow-Origin", "*")

	link := ""
	switch r.Method {
	case "POST":
		// Call ParseForm() to parse the raw query and update r.PostForm and r.Form.
		if err := r.ParseForm(); err != nil {
			fmt.Fprintf(w, "ParseForm() err: %v", err)
			return
		}
		link = strings.TrimSpace(r.FormValue("link"))
	default:
		fmt.Fprintf(w, "Sorry, only POST methods is supported.")
		return
	}

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
		log.Fatal("Cannot connect to DB")
	}

	games := client.Database("chess-explorer").Collection("games")

	cursor, err := games.Find(ctx, bson.M{"link": link})
	defer cursor.Close(ctx)
	if err != nil {
		log.Fatal(err)
	}

	var resultGames []Game
	err = cursor.All(ctx, &resultGames)
	if err != nil {
		log.Fatal(err)
	}

	// send the response
	json.NewEncoder(w).Encode(resultGames)
}
