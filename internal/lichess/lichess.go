package lichess

import (
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

	// If there is a token in the configuration, use it
	lichessToken := viper.GetString("lichess-token")
	if lichessToken != "" {
		req.Header.Add("Authorization", "Bearer "+lichessToken)
	}

	q := req.URL.Query()

	// Get most recent game to set 'since' if possible
	latestGame := pgntodb.Game{}
	pgntodb.GetLatestGame(username, "Lichess.org", &latestGame)

	if !latestGame.DateTime.IsZero() {
		since := latestGame.DateTime.UnixNano() / int64(time.Millisecond)
		since += 1000 // add 1 sec to avoid downloading the last game we have
		q.Add("since", strconv.FormatInt(since, 10))
	}

	req.URL.RawQuery = q.Encode()

	fmt.Println("GET " + req.URL.String())

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

	pgntodb.Process(tmpfile.Name(), latestGame)
}
