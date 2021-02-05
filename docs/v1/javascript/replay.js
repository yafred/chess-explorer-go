// NOTE: this uses chessboardjs and chess.js libraries:
// https://github.com/oakmac/chessboardjs
// https://github.com/jhlywa/chess.js

var apiPort = "52825"
var board = null
var game = new Chess()

var gameId = getParameterByName('gameId'); 
var flip = getParameterByName('flip'); 

if(gameId != null) {
    loadGame(gameId)
}


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
    $.get(`http://127.0.0.1:${apiPort}/game`, { gameId: gameId }, function (data) {
        pgn = JSON.parse(data);
        console.log(data)
    });
}


var config = {
    draggable: false,
    position: 'start'
}
board = Chessboard('myBoard', config)
