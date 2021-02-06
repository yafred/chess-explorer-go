// NOTE: this uses chessboardjs and chess.js libraries:
// https://github.com/oakmac/chessboardjs
// https://github.com/jhlywa/chess.js

var apiPort = "52825"
var board = null
var game = new Chess()

var gameId = getParameterByName('gameId'); 
var boardOrientation = getParameterByName('orientation'); 
var skipMoves = getParameterByName('skip'); 


// Get query parameters
function getParameterByName(name, url = window.location.href) {
    name = name.replace(/[\[\]]/g, '\\$&');
    var regex = new RegExp('[?&]' + name + '(=([^&#]*)|&|#|$)'),
        results = regex.exec(url);
    if (!results) return null;
    if (!results[2]) return '';
    return decodeURIComponent(results[2].replace(/\+/g, ' '));
}


// Load the game we are going to replay
function loadGame(gameId) {
    $.get(`http://127.0.0.1:${apiPort}/game`, { gameId: gameId }, function (jsonData) {
        data = JSON.parse(jsonData)

        // slice pgn
        splitPgn = data.pgn.split(" ")
        splitOpening = splitPgn.slice(0, 3*Math.floor(skipMoves/2) + skipMoves%2 + 1)

        game.load_pgn(splitOpening.join(' '))

        var config = {
            moveSpeed: 400,
            draggable: false,
            position: game.fen(),
            orientation: boardOrientation
        }
        board = Chessboard('myBoard', config)  

        // animate the next move
        game.move(splitPgn[3*Math.floor(skipMoves/2) + skipMoves%2 + 1])
        board.position(game.fen(), true)
    });
}




if(gameId != null) {
    loadGame(gameId)
}

