# Tool to explore openings of chess games 


## Quick start

  * Install [MongoDB Community Server](https://www.mongodb.com/try/download/community)
  * Download latest release of the tool for your platform [here](https://github.com/yafred/chess-explorer/releases)
  * Gather some PGN files in a folder on your machine
    * For [Chess.com](https://chess.com), you can use the API: https://api.chess.com/pub/player/{user}/games/{year}/{month}/pgn
    * For [Lichess.org](https://lichess.org), you can use the API: http://lichess.org/api/games/user/{user} or your account page
  * Run the command chess-explorer pgntodb {path to your file or folder of PGN files}
  * Run the command chess-explorer server and start browsing your games