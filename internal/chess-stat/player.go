package stat

import (
	"fmt"
	"log"
	"os"
	"strings"
)

// StatsToConsole ... does everything
func StatsToConsole(player string, cachePath string, cacheRefresh bool) {
	if cachePath != "" {
		if _, err := os.Stat(cachePath); os.IsNotExist(err) {
			log.Println("Folder " + cachePath + " does not exist. No caching will be used.")
			cachePath = ""
		}
	}

	archivesContainer := ArchivesContainer{}
	getArchives(player, &archivesContainer, cachePath, cacheRefresh)

	// Get games
	var totalGames int
	var loseResults = make(map[string]int)
	var winResults = make(map[string]int)
	var drawResults = make(map[string]int)
	var timeControls = make(map[string]int)
	var rules = make(map[string]int)

	for _, archiveURL := range archivesContainer.Archives {
		gamesContainer := GamesContainer{}
		getGames(player, &gamesContainer, archiveURL, cachePath, cacheRefresh)

		totalGames += len(gamesContainer.Games)

		for _, game := range gamesContainer.Games {
			timeControls[game.TimeControl]++
			rules[game.Rules]++

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

	// Print results to console
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
	fmt.Println(">>>> Rules:")
	for key, value := range rules {
		fmt.Println(key, ":", value)
	}
	fmt.Println(">>>> Time Controls:")
	for key, value := range timeControls {
		fmt.Println(key, ":", value)
	}
}
