package chesscom

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/yafred/chess-explorer/internal/pgntodb"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

/*
https://www.chess.com/news/view/published-data-api

No limitation but concurrent requests forbidden
*/

// archivesContainer ... a list of available archives from Chess.com
type archivesContainer struct {
	Archives []string `json:"archives"`
}

// gamesContainer ... a list of Games from Chess.com
type gamesContainer struct {
	Games []struct {
		URL         string `json:"url"`
		Pgn         string `json:"pgn"`
		TimeControl string `json:"time_control"`
		EndTime     int    `json:"end_time"`
		Rated       bool   `json:"rated"`
		Fen         string `json:"fen"`
		TimeClass   string `json:"time_class"`
		Rules       string `json:"rules"`
		White       struct {
			Rating   int    `json:"rating"`
			Result   string `json:"result"`
			ID       string `json:"@id"`
			Username string `json:"username"`
		} `json:"white"`
		Black struct {
			Rating   int    `json:"rating"`
			Result   string `json:"result"`
			ID       string `json:"@id"`
			Username string `json:"username"`
		} `json:"black"`
	} `json:"games"`
}

// DownloadGames ... Downloads games from Chess.com for user {user}
func DownloadGames(player string) {

	// Download archive list
	chessClient := &http.Client{Timeout: 10 * time.Second}
	archivesURL := "https://api.chess.com/pub/player/" + player + "/games/archives"

	archivesContainer := archivesContainer{}
	resp, err := chessClient.Get(archivesURL)
	if err != nil {
		log.Fatal(err)
	}
	json.NewDecoder(resp.Body).Decode(&archivesContainer)
	defer resp.Body.Close()

	// Download PGN files most recent first
	// Store games in database
	// Stop on first duplicate
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

	goOn := true
	for i := len(archivesContainer.Archives) - 1; i > -1; i-- {
		log.Println("Downloading " + archivesContainer.Archives[i])
		gamesContainer := gamesContainer{}
		resp, err := chessClient.Get(archivesContainer.Archives[i])
		if err != nil {
			log.Fatal(err)
		}
		json.NewDecoder(resp.Body).Decode(&gamesContainer)
		for _, game := range gamesContainer.Games {
			// Note: we should make it a real insert many
			goOn = pgntodb.PgnStringToDB(game.Pgn, client)
			if goOn == false {
				break
			}
		}
		defer resp.Body.Close()
		if goOn == false {
			break
		}
	}
}
