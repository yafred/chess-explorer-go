package pgntodb

import (
	"bufio"
	"log"
	"os"
	"strings"

	"go.mongodb.org/mongo-driver/mongo"
)

func pgnFileToDB(f *os.File, db *mongo.Client) {

	scanner := bufio.NewScanner(f)

	inGame := false
	keyValues := make(map[string]string)
	for scanner.Scan() {
		line := scanner.Text()
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
		case '1':
			keyValues["PGN"] = stripPgn(line)
			insertGame(keyValues, db)
			keyValues = make(map[string]string) // for next game
			break
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
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
	for _, bit := range split {
		if !strings.HasPrefix(bit, "{[%") && !strings.HasSuffix(bit, "]}") && !strings.HasSuffix(bit, "...") {
			// copy and increment index
			split[i] = bit
			i++
		}
	}
	pgn = strings.Join(split[:i], " ")
	return pgn
}
