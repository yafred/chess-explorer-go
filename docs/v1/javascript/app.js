
// NOTE: this uses chessboardjs and chess.js libraries:
// https://github.com/oakmac/chessboardjs
// https://github.com/jhlywa/chess.js

var board = null
var game = new Chess()
var $status = $('#status')
var $fen = $('#fen')
var $pgn = $('#pgn')
var $white = $('#white')
var $black = $('#black')
var browsingGame = ""


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
    $fen.html(game.fen())
    displayPgn(game.pgn())
}

function resetClicked(e) {
    browsingGame = ""
    game.reset()
    board.position(game.fen())
    updateStatus()
    $fen.html(game.fen())
    displayPgn(game.pgn())
}

function nextmove() {
    $("#result").html("");
    $.post("http://127.0.0.1:52825/nextmove", { pgn: game.pgn(), white: $white.val(), black: $black.val() }, function (data) {
        nextMoveToHtml(JSON.parse(data));
    });
}

function loadGame(link, aMove) {
    // set tool in browsing game mode
    $("#result").html("");
    browsingGame = getPgnPlusMove(aMove)
    move(aMove)
    $.post("http://127.0.0.1:52825/games", { link: link }, function (data) {
        ret = JSON.parse(data);
        displayPgn(ret[0].pgn)
    });
}

function nextMoveToHtml(dataObject) {
    if (Array.isArray(dataObject) == false) {
        console.log("not an array")
        return
    }


    var htmlAsArray = []

    dataObject.forEach(element => {

        htmlAsArray.push('<div style="display: flex">')

        htmlAsArray.push('<div style="width:20%;">')
        moveLink = `<a href="javascript:move('${element.move}');">${element.move}</a>`
        if (element.link) {
            moveLink = `<a target="_blank" href="${element.link}">${element.move}</a>`
        }
        htmlAsArray.push(moveLink)
        htmlAsArray.push('</div>')

        htmlAsArray.push('<div style="width:20%;">')
        htmlAsArray.push(element.total)
        htmlAsArray.push('</div>')

        htmlAsArray.push('<div style="width:60%; display:flex; border: 1px solid #aaa; margin-bottom: .1rem">')
        percentage = Math.round(100 * element.win / element.total)
        htmlAsArray.push(`<div style="background-color: white; width:${percentage}%">${percentage}%</div>`)
        percentage = Math.round(100 * element.draw / element.total)
        htmlAsArray.push(`<div style="background-color: #aaa; width:${percentage}%"></div>`)
        percentage = Math.round(100 * element.lose / element.total)
        htmlAsArray.push(`<div style="text-align: right; color: white; background-color: #595959; width:${percentage}%">${percentage}%</div>`)
        htmlAsArray.push('</div>')

        htmlAsArray.push('</div>')
    });

    $("#result").append(htmlAsArray.join('\n'))
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
    splitPgn = pgn.split(" ")

    organizedPgn = [] // array of array of 3 strings ("1.", "move1", "move2")

    pgnMove = []
    splitPgn.forEach(function (item, index, array) {
        if (index % 3 == 0) {
            pgnMove = []
            organizedPgn.push(pgnMove)
        }
        pgnMove.push(item)
    })

    resultString = ""
    organizedPgn.forEach(function (item, index, array) {
        resultString = resultString + item.join(' ') + '<br/>'
    })

    $pgn.html(resultString)
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
    var status = ''

    if (browsingGame == "") {
        status = "Opening mode"
    }
    else {
        status = "Browsing game " + browsingGame
    }

    var moveColor = 'White'
    if (game.turn() === 'b') {
        moveColor = 'Black'
    }

    // checkmate?
    if (game.in_checkmate()) {
        status += ', Game over, ' + moveColor + ' is in checkmate.'
    }

    // draw?
    else if (game.in_draw()) {
        status += ', Game over, drawn position'
    }

    // game still on
    else {
        status += ", " + moveColor + ' to move'

        // check?
        if (game.in_check()) {
            status += ', ' + moveColor + ' is in check'
        }
    }

    $status.html(status)
    $fen.html(game.fen())
    displayPgn(game.pgn())
    if (browsingGame == "") {
        nextmove()
    }
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