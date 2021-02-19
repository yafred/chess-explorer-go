# Tool to explore chess openings 

## This tool needs a local database
  * Install [MongoDB Community Server](https://www.mongodb.com/try/download/community)

## Alternative 1: using executable
  * (You need NOT install Golang)
  * Download the executable for your platform from the [releases page](https://github.com/yafred/chess-explorer-go/releases)
  * Follow instructions below replacing `{command}` with `chess-explorer-{os}-{arch}`

## Alternative 2: Using source code
  * Install [Golang](https://golang.org/doc/install) 
  * `git clone https://github.com/yafred/chess-explorer.git`
  * Open a cmd console and go to the root of the source code directory (where you can see LICENSE, README.md, main.go)
  * Follow instructions below replacing `{command}` with `go run main.go`

## Commands
  * Help
    * `{command} help`
  * Feed your database with games:
    * `{command} chesscom {username}` to download games from https://www.chess.com
    * `{command} lichess {username}` to download games from https://lichess.org
    * `{command} lichess {username} --token {your lichess.org personal API access token}` to download games from https://lichess.org at a higher speed
    * `{command} sync` to download recent games for all users you have already downloaded games for (see commands above)
  * Run the command `{command} server` 
  * Browse your games on http://localhost:52825

  * You can keep your initial download (saves time if you need to reinitialize your database)
    * `{command} chesscom {username} --keep {path to a new file}`
    * `{command} lichess {username} --keep {path to a new file}` 
  * Reinitialize database 
    * `{command} delete {username}` 
    * `{command} delete lichess.org:{username}` 
    * `{command} delete chess.com:{username}` 
    * `{command} pgntodb {path to your PGN file} --username {username}` 

