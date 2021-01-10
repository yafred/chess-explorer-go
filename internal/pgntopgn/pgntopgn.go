package pgntopgn

import (
	"bufio"
	"log"
	"os"
	"strconv"
	"strings"
)

/*
PGN parse is very similar to pgntodb package
*/

// Process ... process a single file or all the files of a folder
func Process(filepath string) {
	processFile(filepath)
}

func processFile(filepath string) {
	// Open file
	file, err := os.Open(filepath)
	defer file.Close()

	if err != nil {
		log.Fatal("Cannot open file " + filepath)
	}

	// Scan file
	scanner := bufio.NewScanner(file)

	gameCounter := 0
	elo1200to1300 := 0
	inGame := false
	keyValues := make(map[string]string)
	for i := 1; scanner.Scan(); i++ {
		line := scanner.Text()
		line = strings.Trim(line, " ")
		if len(line) == 0 {
			continue
		}
		switch line[0] {
		case '[':
			if !inGame {
				inGame = true
			}
			key, value := parseKeyValue(line)
			if key != "" && value != "" {
				keyValues[key] = value
			}
			break
		case '0':
		case '1':
			// If game was abandoned, pgn will be 0-1 or 1-0 (skip it)
			if line != "0-1" && line != "1-0" {
				keyValues["PGN"] = stripPgn(line)
			}
			gameCounter++
			if gameCounter%10000 == 0 {
				log.Println("Scanned " + strconv.Itoa(gameCounter))
			}
			whiteElo, _ := strconv.Atoi(keyValues["WhiteElo"])
			blackElo, _ := strconv.Atoi(keyValues["BlackElo"])
			if whiteElo >= 1200 && whiteElo < 1300 && blackElo >= 1200 && blackElo < 1300 {
				elo1200to1300++
			}
			keyValues = make(map[string]string) // for next game
			break
		default:
			// not a valid char, skip
			break
		}
	}

	log.Println("Scanned " + strconv.Itoa(gameCounter))
	log.Println("Elo 1200-1300: " + strconv.Itoa(elo1200to1300))

}

// [Key "value"]
func parseKeyValue(line string) (key string, value string) {
	line = strings.Trim(line, "[] ")
	if len(line) == 0 {
		return "", ""
	}
	split := strings.Split(line, " ")
	if len(split) == 0 {
		return "", ""
	}
	key = split[0]

	quotedValue := strings.Join(split[1:], " ")
	value = strings.Trim(quotedValue, "\" ")

	return key, value
}

// lichess: 1. d4 Nf6 2. e3 d5
// chess.com: 1. d4 {[%clk 0:29:56.7]} 1... d5 {[%clk 0:29:52.9]} 2. Bf4 {[%clk 0:29:52.9]} 2... Nf6 {[%clk 0:29:24.1]}
func stripPgn(line string) (pgn string) {
	split := strings.Split(line, " ")
	i := 0 // output index
	skip := false
	for _, bit := range split {
		if strings.HasPrefix(bit, "{") {
			skip = true
		}
		if skip == false && !strings.HasSuffix(bit, "...") {
			// copy and increment index
			bit = strings.Replace(bit, "!", "", -1)
			bit = strings.Replace(bit, "?", "", -1)
			split[i] = bit
			i++
		}
		if strings.HasSuffix(bit, "}") {
			skip = false
		}
	}
	pgn = strings.Join(split[:i], " ")
	return pgn
}
