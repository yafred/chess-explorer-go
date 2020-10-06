package stat

import (
	"fmt"
	"log"
	"os"
	"strconv"
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

// CsvToConsole ... compiles games from chess.com into a cvs file for spreadsheets
func CreateCsvFile(player string, cachePath string, cacheRefresh bool, filepath string) {

	file, err := os.Create(filepath)
	if err != nil {
		fmt.Println("Cannot open", filepath)
		return
	}
	defer file.Close()

	archivesContainer := ArchivesContainer{}
	getArchives(player, &archivesContainer, cachePath, cacheRefresh)

	columns := []string{"EndTime", "Color", "Against", "Outcome", "Result", "TimeClass", "TimeControl", "Rating", "URL"}
	fmt.Fprintln(file, strings.Join(columns, ","))

	// Get games
	for _, archiveURL := range archivesContainer.Archives {
		gamesContainer := GamesContainer{}
		getGames(player, &gamesContainer, archiveURL, cachePath, cacheRefresh)

		outcome := ""
		result := ""
		color := ""
		against := ""
		rating := 0
		for _, game := range gamesContainer.Games {
			if strings.EqualFold(game.White.Username, player) {
				color = "white"
				against = game.Black.Username
				rating = game.White.Rating
			} else {
				color = "black"
				against = game.White.Username
				rating = game.Black.Rating
			}
			if game.White.Result != "win" && game.Black.Result != "win" {
				outcome = "draw"
				result = game.White.Result
			} else if (game.White.Result == "win" && strings.EqualFold(game.White.Username, player)) || (game.Black.Result == "win" && strings.EqualFold(game.Black.Username, player)) {
				outcome = "win"
				if color == "white" {
					result = game.Black.Result
				} else {
					result = game.White.Result
				}
			} else {
				outcome = "lose"
				if color == "white" {
					result = game.White.Result
				} else {
					result = game.Black.Result
				}
			}
			/*
				{
					"url": "https://www.chess.com/live/game/4963557833",
					"pgn": "[Event \"Live Chess\"]\n[Site \"Chess.com\"]\n[Date \"2020.06.07\"]\n[Round \"-\"]\n[White \"fredo599\"]\n[Black \"Komodo2\"]\n[Result \"0-1\"]\n[ECO \"B00\"]\n[ECOUrl \"https://www.chess.com/openings/Kings-Pawn-Opening-Carr-Defense\"]\n[CurrentPosition \"2bk1bn1/1rppq2r/p3p1Qp/n3P1p1/1p2N3/1B3N2/PPP2PPP/R1B1KR2 w Q -\"]\n[Timezone \"UTC\"]\n[UTCDate \"2020.06.07\"]\n[UTCTime \"00:16:28\"]\n[WhiteElo \"233\"]\n[BlackElo \"425\"]\n[TimeControl \"600\"]\n[Termination \"Komodo2 won on time\"]\n[StartTime \"00:16:28\"]\n[EndDate \"2020.06.07\"]\n[EndTime \"00:26:39\"]\n[Link \"https://www.chess.com/live/game/4963557833\"]\n\n1. e4 {[%clk 0:09:42.2]} 1... h6 {[%clk 0:09:59.9]} 2. Nc3 {[%clk 0:08:37.9]} 2... b6 {[%clk 0:09:59.8]} 3. d4 {[%clk 0:08:28.5]} 3... a6 {[%clk 0:09:59.7]} 4. Bc4 {[%clk 0:07:53.9]} 4... b5 {[%clk 0:09:59.6]} 5. Bd5 {[%clk 0:07:41.2]} 5... Ra7 {[%clk 0:09:59.5]} 6. Qh5 {[%clk 0:06:58.4]} 6... e6 {[%clk 0:09:59.4]} 7. Bb3 {[%clk 0:06:49.3]} 7... Rh7 {[%clk 0:09:59.3]} 8. Nf3 {[%clk 0:05:39.5]} 8... g6 {[%clk 0:09:59.2]} 9. Qe5 {[%clk 0:05:27.8]} 9... f6 {[%clk 0:09:59.1]} 10. Qf4 {[%clk 0:05:11.5]} 10... Qe7 {[%clk 0:09:59]} 11. e5 {[%clk 0:04:47.7]} 11... g5 {[%clk 0:09:58.9]} 12. Qg4 {[%clk 0:04:33.3]} 12... b4 {[%clk 0:09:58.8]} 13. Ne4 {[%clk 0:04:05.5]} 13... fxe5 {[%clk 0:09:58.7]} 14. dxe5 {[%clk 0:03:46.9]} 14... Rb7 {[%clk 0:09:58.6]} 15. Rf1 {[%clk 0:03:38]} 15... Nc6 {[%clk 0:09:58.5]} 16. Qh5+ {[%clk 0:00:56.3]} 16... Kd8 {[%clk 0:09:58.4]} 17. Qg6 {[%clk 0:00:42.3]} 17... Na5 {[%clk 0:09:58.3]} 0-1",
					"time_control": "600",
					"end_time": 1591489599,
					"rated": true,
					"fen": "2bk1bn1/1rppq2r/p3p1Qp/n3P1p1/1p2N3/1B3N2/PPP2PPP/R1B1KR2 w Q -",
					"time_class": "blitz",
					"rules": "chess",
					"white": {
						"rating": 233,
						"result": "timeout",
						"@id": "https://api.chess.com/pub/player/fredo599",
						"username": "fredo599"
					},
					"black": {
						"rating": 425,
						"result": "win",
						"@id": "https://api.chess.com/pub/player/komodo2",
						"username": "Komodo2"
					}
				}
			*/
			values := []string{strconv.Itoa(game.EndTime), color, against, outcome, result, game.TimeClass, game.TimeControl, strconv.Itoa(rating), game.URL}
			fmt.Fprintln(file, strings.Join(values, ","))
		}
	}
}
