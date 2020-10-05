package stat

import (
	"encoding/json"
	"net/http"
	"time"
)

/*
Cache: keep responses from chess.com to reduce time and network load for subsequent stats

Notes:
- Games: only most recent month may have changed



Functions needed:

- response to file
func main() {
  resp, err := http.Get("...")
  check(err)
  defer resp.Body.Close()
  out, err := os.Create("filename.ext")
  if err != nil {
    // panic?
  }
  defer out.Close()
  io.Copy(out, resp.Body)
}

- response to struct
var chessClient = &http.Client{Timeout: 10 * time.Second}
archivesURL := "https://api.chess.com/pub/player/" + player + "/games/archives"

	archiveResponse := ArchivesResponse{}
	r, err := chessClient.Get(archivesURL)
	if err != nil {
	}

	json.NewDecoder(r.Body).Decode(&archiveResponse)
	r.Body.Close()


- file to struct
https://www.golangprograms.com/golang-read-json-file-into-struct.html

*/

var chessClient = &http.Client{Timeout: 10 * time.Second}

func getArchives(player string, archivesContainer *ArchivesContainer, cachePath string) {
	archivesURL := "https://api.chess.com/pub/player/" + player + "/games/archives"

	r, err := chessClient.Get(archivesURL)
	if err != nil {
	}

	json.NewDecoder(r.Body).Decode(archivesContainer)
	r.Body.Close()
}

func getGames(gamesContainer *GamesContainer, archiveURL string, cachePath string) {
	r, err := chessClient.Get(archiveURL)
	if err != nil {
	}

	json.NewDecoder(r.Body).Decode(gamesContainer)
	r.Body.Close()
}
