package chesscom

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	http "net/http"
	"os"

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
func DownloadGames(username string, keepPgn string) {

	// Download archive list
	client := &http.Client{}
	archivesURL := "https://api.chess.com/pub/player/" + username + "/games/archives"

	archivesContainer := archivesContainer{}
	resp, err := client.Get(archivesURL)
	if err != nil {
		log.Fatal(err)
	}
	json.NewDecoder(resp.Body).Decode(&archivesContainer)
	defer resp.Body.Close()

	// Get most recent game from database to avoid downloading duplicates
	lastGame := pgntodb.FindLastGame(username, "chess.com")
	if lastGame.DateTime.IsZero() {
		log.Println("New user")
	} else {
		log.Println("Most recent game in database: " + lastGame.GameID)
	}

	// Create the keep file if needed
	var keepPgnFile *os.File
	if keepPgn != "" {
		keepPgnFile, err = os.OpenFile(keepPgn, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
		if err != nil {
			log.Fatal(err)
		}
		defer keepPgnFile.Close()
	}

	// Download PGN files most recent first
	// Store games in database
	// Stop on first duplicate
	for i := len(archivesContainer.Archives) - 1; i > -1; i-- {
		log.Println("GET " + archivesContainer.Archives[i] + "/pgn")
		goOn := downloadArchive(client, archivesContainer.Archives[i]+"/pgn", lastGame, keepPgnFile)
		if goOn == false {
			break
		}
	}
}

func downloadArchive(client *http.Client, url string, lastGame *pgntodb.LastGame, keepPgnFile *os.File) bool {

	// Random file name
	tmpfile, err := os.CreateTemp("", "chesscom")
	if err != nil {
		log.Fatal(err)
	}
	defer os.Remove(tmpfile.Name()) // clean up

	// Create the temp file
	f, err := os.OpenFile(tmpfile.Name(), os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	// Send request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal(err)
	}

	resp, err := client.Do(req)

	if err != nil {
		log.Fatal(err)
	}

	// stream response
	buf := make([]byte, 10000)

	numBytesRead := 0
	// Read the response body
	for {
		n, err := resp.Body.Read(buf)

		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		numBytesRead += n
		fmt.Print(".")

		n, err = f.Write(buf[0:n])
		if err != nil {
			log.Fatal(err)
		}

		if keepPgnFile != nil {
			n, err = keepPgnFile.Write(buf[0:n])
			if err != nil {
				log.Fatal(err)
			}
		}

		if err != nil {
			log.Fatal("Error reading HTTP response: ", err.Error())
		}
	}

	fmt.Println()

	log.Println(numBytesRead, " bytes read")

	// parse file
	return pgntodb.Process(tmpfile.Name(), lastGame)
}
