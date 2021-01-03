# Tool to explore openings of chess games 

## Important for windows users

If you follow the procedure described in the Quick start, Windows Defender may report a virus (Trojan:Script/Wacatac.B!ml)

Please read: https://golang.org/doc/faq#virus

You can, either authorize this program or follow the alternate procedure



## Quick start

  * Install [MongoDB Community Server](https://www.mongodb.com/try/download/community)
  * Download latest release of the tool for your platform [here](https://github.com/yafred/chess-explorer/releases)
  * Gather some PGN files in a folder on your machine
    * For [Chess.com](https://chess.com), you can use the API: https://api.chess.com/pub/player/{user}/games/{year}/{month}/pgn
    * For [Lichess.org](https://lichess.org), you can use the API: http://lichess.org/api/games/user/{user} or your account page
  * Run the command chess-explorer pgntodb {path to your file or folder of PGN files}
  * Run the command chess-explorer server and start browsing your games

  ## Alternate procedure

  * Install Golang 
