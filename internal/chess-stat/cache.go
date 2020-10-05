package stat

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

var chessClient = &http.Client{Timeout: 10 * time.Second}

func getArchives(player string, archivesContainer *ArchivesContainer, cachePath string, cacheRefresh bool) {

	if cachePath != "" {
		_ = os.Mkdir(filepath.Join(cachePath, player), 0700)
	}

	archivesURL := "https://api.chess.com/pub/player/" + player + "/games/archives"

	if cachePath != "" {
		cacheFilePath := filepath.Join(cachePath, player, player+"-archives.json")

		if _, err := os.Stat(cacheFilePath); os.IsNotExist(err) || cacheRefresh {
			resp, err := chessClient.Get(archivesURL)
			if err != nil {
			}

			bodyBytes, _ := ioutil.ReadAll(resp.Body)
			ioutil.WriteFile(cacheFilePath, bodyBytes, 0700)

			resp.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))

			json.NewDecoder(resp.Body).Decode(archivesContainer)
			defer resp.Body.Close()
		} else {
			jsonFile, err := os.Open(cacheFilePath)
			if err != nil {
			}
			fmt.Println("Successfully Opened " + cacheFilePath)
			defer jsonFile.Close()

			byteValue, _ := ioutil.ReadAll(jsonFile)

			json.Unmarshal([]byte(byteValue), archivesContainer)
		}
	} else {
		resp, err := chessClient.Get(archivesURL)
		if err != nil {
		}
		json.NewDecoder(resp.Body).Decode(archivesContainer)
		defer resp.Body.Close()
	}
}

func getGames(player string, gamesContainer *GamesContainer, archiveURL string, cachePath string, cacheRefresh bool) {

	if cachePath != "" {
		_ = os.Mkdir(filepath.Join(cachePath, player), 0700)
	}

	if cachePath != "" {
		_, month, year := bitsFromArchivesURL(archiveURL)
		cacheFilePath := filepath.Join(cachePath, player, player+"-"+year+"-"+month+".json")

		if _, err := os.Stat(cacheFilePath); os.IsNotExist(err) || cacheRefresh {
			resp, err := chessClient.Get(archiveURL)

			if err != nil {
			}

			bodyBytes, _ := ioutil.ReadAll(resp.Body)
			ioutil.WriteFile(cacheFilePath, bodyBytes, 0700)

			resp.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))

			json.NewDecoder(resp.Body).Decode(gamesContainer)
			defer resp.Body.Close()
		} else {
			jsonFile, err := os.Open(cacheFilePath)
			if err != nil {
			}
			fmt.Println("Successfully Opened " + cacheFilePath)
			defer jsonFile.Close()

			byteValue, _ := ioutil.ReadAll(jsonFile)

			json.Unmarshal([]byte(byteValue), gamesContainer)
		}
	} else {
		resp, err := chessClient.Get(archiveURL)
		if err != nil {
		}
		json.NewDecoder(resp.Body).Decode(gamesContainer)
		defer resp.Body.Close()
	}
}

func bitsFromArchivesURL(archiveURL string) (player string, month string, year string) {
	// archiveURL is https://api.chess.com/pub/player/{player}/games/{year}/{month}

	bits := strings.Split(archiveURL, "/")
	month = bits[len(bits)-1]
	year = bits[len(bits)-2]
	player = bits[len(bits)-4]
	return
}
