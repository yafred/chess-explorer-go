
// NOTE: this uses chessboardjs and chess.js libraries:
// https://github.com/oakmac/chessboardjs
// https://github.com/jhlywa/chess.js

var apiPort = "52825"
var board = null
var game = new Chess()
var $status = $('#status')
var $fen = $('#fen')
var $pgn = $('#pgn')
var $white = $('#white')
var $black = $('#black')
var $from = $('#from')
var $to = $('#to')
var $minelo = $('#minelo')
var $maxelo = $('#maxelo')
var $site = $('#site')
var browsingGame = ""

var nextMoveTpl = document.getElementById('nextMoveTpl').innerHTML;

/* we could do some completion with this ...
var keydownTimeout;
var keydownTimeoutValue = 1500; // millisecs

$white.keydown(function () {
    clearTimeout(keydownTimeout);
    keydownTimeout = setTimeout(function () {
        // do stuff when user has been idle for 1.5 second
        resetClicked()
    }, keydownTimeoutValue);
});
$black.keydown(function () {
    clearTimeout(keydownTimeout);
    keydownTimeout = setTimeout(function () {
        // do stuff when user has been idle for 1.5 second
        resetClicked()
    }, keydownTimeoutValue);
});
*/

$from.change(function () {
    resetClicked()
});

$to.change(function () {
    resetClicked()
});

$white.change(function () {
    resetClicked()
});

$black.change(function () {
    resetClicked()
});

$minelo.change(function () {
    resetClicked()
});

$maxelo.change(function () {
    resetClicked()
});

$site.change(function () {
    resetClicked()
});

function swapBlackWhiteClicked(e) {
    var black = $black.val()
    $black.val($white.val())
    $white.val(black)
    resetClicked()
}

function undoClicked(e) {
    browsingGame = ""
    game.undo()
    board.position(game.fen())
    updateStatus()
}

function resetClicked(e) {
    browsingGame = ""
    game.reset()
    board.position(game.fen())
    updateStatus()
}

function getNextMove() {
    $("#result").html("");
    $.post(`http://127.0.0.1:${apiPort}/nextmove`, {
        pgn: game.pgn(),
        white: $white.val(),
        black: $black.val(),
        from: $from.val(),
        to: $to.val(),
        minelo: $minelo.val(),
        maxelo: $maxelo.val(),
        site: $site.val()
    }, function (data) {
        nextMoveToHtml(JSON.parse(data));
    });
}

function loadGame(link, aMove) {
    // set tool in browsing game mode
    $("#result").html("");
    browsingGame = getPgnPlusMove(aMove)
    move(aMove)
    $.post(`http://127.0.0.1:${apiPort}/games`, { link: link }, function (data) {
        ret = JSON.parse(data);
        displayPgn(ret[0].pgn)
    });
}

function nextMoveToHtml(dataObject) {
    if (Array.isArray(dataObject) == false) {
        console.log("not an array")
        return
    }

    var moves = []

    dataObject.forEach(element => {

        winPercent = Math.round(100 * element.win / element.total)
        drawPercent = Math.round(100 * element.draw / element.total)
        drawPercentText = ""
        if (drawPercent > 12) {
            drawPercentText = "" + drawPercent + "%"
        }
        losePercent = Math.round(100 * element.lose / element.total)

        internalLink = false
        externalLink = false
        if (element.link) {
            externalLink = true
        }
        else {
            internalLink = true
        }

        moves.push({
            move: element.move,
            link: element.link,
            internalLink: internalLink,
            externalLink: externalLink,
            total: element.total,
            winPercent: winPercent,
            drawPercent: drawPercent,
            drawPercentText: drawPercentText,
            losePercent: losePercent,
        })
    });

    $("#result").html(Mustache.render(nextMoveTpl, moves))
}

// Not used (I use game link instead)
function getPgnPlusMove(aMove) {
    pgn = game.pgn()
    splitPgn = pgn.split(" ")
    lineCount = Math.floor((splitPgn.length / 3))
    if (splitPgn.length % 3 == 0) {
        // create a new line
        pgn = pgn + " " + (lineCount + 1) + "."
    }
    pgn = pgn + " " + aMove
    return pgn
}

function displayPgn(pgn) {
    $pgn.html(pgn)
}

function move(aMove) {
    game.move(aMove)
    updateStatus()
    board.position(game.fen(), false)
}

function onDragStart(source, piece, position, orientation) {
    // do not pick up pieces if the game is over
    if (game.game_over()) return false

    // only pick up pieces for the side to move
    if ((game.turn() === 'w' && piece.search(/^b/) !== -1) ||
        (game.turn() === 'b' && piece.search(/^w/) !== -1)) {
        return false
    }
}

function onDrop(source, target) {
    // see if the move is legal
    var move = game.move({
        from: source,
        to: target,
        promotion: 'q' // NOTE: always promote to a queen for example simplicity
    })

    // illegal move
    if (move === null) return 'snapback'

    browsingGame = "" // quit browsing mode
    updateStatus()
}

// update the board position after the piece snap
// for castling, en passant, pawn promotion
function onSnapEnd() {
    board.position(game.fen())
}

function updateStatus() {
    displayPgn(game.pgn())
    $fen.html(game.fen())
    getNextMove()
}

var config = {
    draggable: true,
    position: 'start',
    onDragStart: onDragStart,
    onDrop: onDrop,
    onSnapEnd: onSnapEnd
}
board = Chessboard('myBoard', config)

updateStatus()