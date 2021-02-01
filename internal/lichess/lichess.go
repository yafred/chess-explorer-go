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

// DownloadGames ... Downloads games from lichess.org for user {user}
// https://lichess.org/api#operation/apiGamesUser
func DownloadGames(username string, keepPgn string) {

	url := "https://lichess.org/api/games/user/" + username

	client := &http.Client{}
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
	lastGame := pgntodb.FindLastGame(username, "lichess.org")

	if lastGame.DateTime.IsZero() {
		log.Println("New user")
	} else {
		log.Println("Last game in database: " + lastGame.GameID)
		since := lastGame.DateTime.UnixNano() / int64(time.Millisecond)
		since += 1000 // add 1 sec to avoid downloading the last game we have
		q.Add("since", strconv.FormatInt(since, 10))
	}

	req.URL.RawQuery = q.Encode()

	fmt.Println("GET " + req.URL.String())

	// Get data
	resp, err := client.Do(req)

	if err != nil {
		log.Fatal(err)
	}

	fileName := keepPgn

	if fileName == "" {
		// Create a temp file
		tmpfile, err := ioutil.TempFile("", "lichess")
		if err != nil {
			log.Fatal(err)
		}
		fileName = tmpfile.Name()
		defer os.Remove(tmpfile.Name()) // clean up
	}

	f, err := os.OpenFile(fileName, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

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

		if err != nil {
			log.Fatal("Error reading HTTP response: ", err.Error())
		}
	}

	fmt.Println()

	log.Println(numBytesRead, " bytes read")
	pgntodb.Process(fileName, lastGame)
}
