# Tool to explore chess openings 

## What is it ?

An opening explorer allows you to browse games move by move and find out if a position is (statistically) winning or losing

## Installation
  * Install [Golang](https://golang.org/doc/install) 
  * Install [MongoDB Community Server](https://www.mongodb.com/try/download/community)
  * Download Source Code from https://github.com/yafred/chess-explorer/releases and unzip it
  * Open a cmd console and go to the root of the source code directory (where you can see LICENSE, README.md, main.go)
  * Feed your database with games:
    * `go run main.go chesscom {username}` to download games from https://www.chess.com
    * `go run main.go lichess {username}` to download games from https://lichess.org
    * `go run main.go lichess {username} --token {your lichess.org personal API access token}` to download games from https://lichess.org at a higher speed
    * `go run main.go pgntodb {path to your file or folder of PGN files}` to import data from PGN files you already have
  * Run the command `go run main.go server` 
  * Browse your games on http://localhost:52825

## Binaries

I am trying to provide binaries [here](https://github.com/yafred/chess-explorer/releases) but they don't seem to work great