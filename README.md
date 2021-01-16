# Tool to explore chess openings 

## What is it ?

An opening explorer allows you to browse games move by move and find out if a position is (statistically) winning or losing

## Important for windows users

If you follow the procedure described in the Quick start, your antivirus (like Windows Defender) may report a virus (Trojan:Script/Wacatac.B!ml)

Please read: https://golang.org/doc/faq#virus

You can, either authorize this program or follow the alternate procedure

## Quick start
  * Install [MongoDB Community Server](https://www.mongodb.com/try/download/community)
  * Download latest release of the tool chess-explorer for your platform [here](https://github.com/yafred/chess-explorer/releases)
  * Feed your database with games:
    * `chess-explorer chesscom {username}` to download games from https://www.chess.com
    * `chess-explorer lichess {username}` to download games from https://lichess.org
    * `chess-explorer lichess {username} --token {your lichess.org personal API access token}` to download games from https://lichess.org at a higher speed
    * `chess-explorer pgntodb {path to your file or folder of PGN files}` to import data from PGN files you already have
  * Run the command `chess-explorer server` 
  * Browse your games on http://localhost:52825

## Alternate procedure (the one I prefer)
  * Install MongoDB and gather PGN files
  * Install [Golang](https://golang.org/doc/install) 
  * Download Source Code from https://github.com/yafred/chess-explorer/releases and unzip it
  * Open a cmd console and go to the root of the source code directory (where you can see LICENSE, README.md, main.go)
  * Replace the `chess-explorer` command from Quick start with `go run main.go` (for example: `go run main.go server`)

