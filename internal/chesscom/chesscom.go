package chesscom

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	http "net/http"
	"os"
	"time"

	"github.com/yafred/chess-explorer/internal/pgntodb"
)

/*
https://www.chess.com/news/view/published-data-api

No limitation but concurrent requests forbidden
*/

// archivesContainer ... a list of available archives from Chess.com
type archivesContainer struct {
	Archives []string `json:"archives"`
}

// DownloadGames ... Downloads games from Chess.com for {username}
func DownloadGames(username string) {

	// Download archive list
	chessClient := &http.Client{Timeout: 10 * time.Second}
	archivesURL := "https://api.chess.com/pub/player/" + username + "/games/archives"

	archivesContainer := archivesContainer{}
	resp, err := chessClient.Get(archivesURL)
	if err != nil {
		log.Fatal(err)
	}
	json.NewDecoder(resp.Body).Decode(&archivesContainer)
	defer resp.Body.Close()

	// Get most recent game to set 'since' if possible
	lastGame := pgntodb.FindLastGame(username, "chess.com")
	if lastGame.DateTime.IsZero() {
		log.Println("New user")
	} else {
		log.Println("Last game in database: " + lastGame.GameID)
	}
	// Download PGN files most recent first
	// Store games in database
	// Stop on first duplicate
	for i := len(archivesContainer.Archives) - 1; i > -1; i-- {
		log.Println("GET " + archivesContainer.Archives[i] + "/pgn")
		goOn := downloadArchive(chessClient, archivesContainer.Archives[i]+"/pgn", lastGame)
		if goOn == false {
			break
		}
	}
}

func downloadArchive(chessClient *http.Client, url string, lastGame *pgntodb.LastGame) bool {

	// Get data
	resp, err := chessClient.Get(url)

	if err != nil {
		log.Fatal(err)
	}

	// Create a temp file
	tmpfile, err := ioutil.TempFile("", "chesscom")
	if err != nil {
		log.Fatal(err)
	}
	defer os.Remove(tmpfile.Name()) // clean up

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

	return pgntodb.Process(tmpfile.Name(), lastGame)
}
