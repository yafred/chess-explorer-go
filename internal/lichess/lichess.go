package lichess

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/yafred/chess-explorer/internal/pgntodb"
)

// DownloadGames ... Downloads games from Chess.com for user {user}
// https://lichess.org/api#operation/apiGamesUser
func DownloadGames(username string) {

	url := "https://lichess.org/api/games/user/" + username

	chessClient := &http.Client{Timeout: 10 * time.Second}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal(err)
	}
	//req.Header.Add("Accept", "application/x-ndjson")

	q := req.URL.Query()
	//q.Add("since", "")
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
	//defer os.Remove(tmpfile.Name()) // clean up
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
