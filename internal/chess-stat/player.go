package stat

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

// Player ... a profile from Chess.com
type Player struct {
	PlayerID   int    `json:"player_id"`
	ID         string `json:"@id"`
	URL        string `json:"url"`
	Username   string `json:"username"`
	Followers  int    `json:"followers"`
	Country    string `json:"country"`
	LastOnline int    `json:"last_online"`
	Joined     int    `json:"joined"`
	Status     string `json:"status"`
	IsStreamer bool   `json:"is_streamer"`
}

// GamesResponse ... a list of Games from Chess.com
type GamesResponse struct {
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

// ArchivesResponse ... a list of available archives from Chess.com
type ArchivesResponse struct {
	Archives []string `json:"archives"`
}

// StatsToConsole ... does everything
func StatsToConsole(player string, cachePath string, refreshCache bool) {

	fmt.Println("cache", cachePath)

	var chessClient = &http.Client{Timeout: 10 * time.Second}

	// Get available archives
	archivesURL := "https://api.chess.com/pub/player/" + player + "/games/archives"

	archiveResponse := ArchivesResponse{}
	r, err := chessClient.Get(archivesURL)
	if err != nil {
	}

	json.NewDecoder(r.Body).Decode(&archiveResponse)
	r.Body.Close()

	// Get games
	var totalGames int
	var loseResults = make(map[string]int)
	var winResults = make(map[string]int)
	var drawResults = make(map[string]int)

	for _, archiveURL := range archiveResponse.Archives {
		gamesResponse := GamesResponse{}
		r, err := chessClient.Get(archiveURL)
		if err != nil {
		}

		json.NewDecoder(r.Body).Decode(&gamesResponse)
		r.Body.Close()

		totalGames += len(gamesResponse.Games)

		for _, game := range gamesResponse.Games {
			if game.White.Result != "win" && game.Black.Result != "win" { // Draw
				if game.White.Result != game.Black.Result {
					fmt.Println("Results should be the same for black and white: ", game.White, game.Black)
				}
				_, ok := drawResults[game.White.Result]
				if ok {
					drawResults[game.White.Result]++
				} else {
					drawResults[game.White.Result] = 1
				}
			} else if (game.White.Result == "win" && strings.EqualFold(game.White.Username, player)) || (game.Black.Result == "win" && strings.EqualFold(game.Black.Username, player)) { // Win
				var result string
				if game.White.Result != "win" {
					result = game.White.Result
				} else {
					result = game.Black.Result
				}
				_, ok := winResults[result]
				if ok {
					winResults[result]++
				} else {
					winResults[result] = 1
				}
			} else { // Lose
				var result string
				if game.White.Result != "win" {
					result = game.White.Result
				} else {
					result = game.Black.Result
				}
				_, ok := loseResults[result]
				if ok {
					loseResults[result]++
				} else {
					loseResults[result] = 1
				}
			}
		}
	}

	fmt.Println(">>>> Total games: ", totalGames)
	fmt.Println(">>>> Draw:")
	for key, value := range drawResults {
		fmt.Println(key, ":", value)
	}
	fmt.Println(">>>> Lose:")
	for key, value := range loseResults {
		fmt.Println(key, ":", value)
	}
	fmt.Println(">>>> Win:")
	for key, value := range winResults {
		fmt.Println(key, ":", value)
	}
}
