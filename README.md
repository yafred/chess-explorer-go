# Tool to explore your chess openings 

![Carlsen as white](/images/carlsen_white.png)
![Carlsen as black](/images/carlsen_black.png)


## Highlights
  * You download all your games only the first time. Next time, you will only download the recent ones.
  * There is no special processing step: going from black to white, or from one player to the other is instant.
  * You can explore the openings of a group of players.
  * When you select a player, only the time controls relevant to this player are displayed.
  * You can paste your opening PGN and skip the first moves of book openings.
  * You can download the recent games of all your favourite players in one command.
  * You can scan the selected games to know if they have reached a specific position (FEN)
  * When there is only one result for the next move of the opening, you can replay the game locally or go to the site where the game was played.

## This tool needs a Mongo database to cache your data
  * Either install [MongoDB Community Server](https://www.mongodb.com/try/download/community)
  * Or create a MongoDB cluster online (there are some free plans, for example: [MongoDB Atlas](https://docs.atlas.mongodb.com/tutorial/deploy-free-tier-cluster/))

## Alternative 1: using executable
  * Download the executable for your platform from the [releases page](https://github.com/yafred/chess-explorer-go/releases)
  * Follow instructions below replacing `{command}` with `chess-explorer-{os}-{arch}`

## Alternative 2: using source code
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

