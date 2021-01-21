# Tool to explore chess openings 

## Installation
  * Install [Golang](https://golang.org/doc/install) 
  * Install [MongoDB Community Server](https://www.mongodb.com/try/download/community)
  * Download the source code from https://github.com/yafred/chess-explorer/releases and unzip it (or `git clone https://github.com/yafred/chess-explorer.git`)
  * Open a cmd console and go to the root of the source code directory (where you can see LICENSE, README.md, main.go)
  * Feed your database with games:
    * `go run main.go chesscom {username}` to download games from https://www.chess.com
    * `go run main.go lichess {username}` to download games from https://lichess.org
    * `go run main.go lichess {username} --token {your lichess.org personal API access token}` to download games from https://lichess.org at a higher speed
    * `go run main.go pgntodb {path to your file or folder of PGN files} --username {username}` to import data from PGN files you already have
  * Run the command `go run main.go server` 
  * Browse your games on http://localhost:52825

