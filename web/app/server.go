package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/yafred/chess-stat/internal/pgntodb"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// Start ... start a web server on port 8080
func Start(port int) {
	fs := http.FileServer(http.Dir("../../web/static"))
	http.Handle("/", fs)

	http.HandleFunc("/test", testHandler)

	log.Println("Server is listening on port " + strconv.Itoa(port))
	http.ListenAndServe(":"+strconv.Itoa(port), nil)
}

func testHandler(w http.ResponseWriter, r *http.Request) {
	// Connect to DB
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017"))
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

	// iterate all documents
	cursor, _ := games.Find(ctx, bson.M{})

	defer cursor.Close(ctx)

	for i := 0; cursor.Next(ctx); i++ {
		var game pgntodb.Game
		if err = cursor.Decode(&game); err != nil {
			log.Fatal(err)
		}
		if strings.HasPrefix(game.PGN, "1. d4 e5") {
			fmt.Println(game.Link)
		}
	}
}
