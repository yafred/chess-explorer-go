package server

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/yafred/chess-stat/internal/pgntodb"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func findHandler(w http.ResponseWriter, r *http.Request) {

	pgn := ""

	switch r.Method {
	case "POST":
		// Call ParseForm() to parse the raw query and update r.PostForm and r.Form.
		if err := r.ParseForm(); err != nil {
			fmt.Fprintf(w, "ParseForm() err: %v", err)
			return
		}
		pgn = strings.TrimSpace(r.FormValue("pgn"))
	default:
		fmt.Fprintf(w, "Sorry, only POST methods is supported.")
		return
	}

	log.Println(pgn)

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

	var andClause []bson.M
	andClause = append(andClause, bson.M{"site": "Chess.com"})
	andClause = append(andClause, bson.M{"white": "fredo599"})
	andClause = append(andClause, bson.M{"pgn": bson.M{"$regex": pgn + ".*"}})

	cursor, err := games.Find(ctx, bson.M{"$and": andClause})
	defer cursor.Close(ctx)
	if err != nil {
		log.Fatal(err)
	}

	var result []pgntodb.Game
	err = cursor.All(ctx, &result)
	if err != nil {
		log.Fatal(err)
	}

	// send the response
	json.NewEncoder(w).Encode(result[0])
}
