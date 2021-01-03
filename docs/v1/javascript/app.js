
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


function swapBlackWhiteClicked(e) {
    var black = $black.val()
    $black.val($white.val())
    $white.val(black)
    resetClicked()
}

function undoClicked(e) {
    game.undo()
    board.position(game.fen())
    updateStatus()
    $fen.html(game.fen())
    displayPgn()
}

function resetClicked(e) {
    game.reset()
    board.position(game.fen())
    updateStatus()
    $fen.html(game.fen())
    displayPgn()
}

function nextmove() {
    $("#result").html("");
    $.post("http://127.0.0.1:52825/nextmove", { pgn: game.pgn(), white: $white.val(), black: $black.val() }, function (data) {
        nextMoveToHtml(JSON.parse(data));
    });
}

function loadGame(aMove) {
    pgn = getPgnPlusMove(aMove)
    console.log("load " + pgn)
}

function nextMoveToHtml(dataObject) {
    if (Array.isArray(dataObject) == false) {
        console.log("not an array")
        return
    }

    dataObject.forEach(element => {
        moveLink = `<a href="javascript:move('${element.move}');">${element.move}</a>`
        if(element.link) {
            moveLink = `<a href="javascript:loadGame('${element.move}');">${element.move}</a>`
        }
        var htmlAsArray = [
            '<div>',
            element.total,
            moveLink]

        element.Results.forEach(result => {
            var result = [
                '<span>(',
                result.result + ':' + result.sum,
                ')</span>'
            ]
            htmlAsArray = htmlAsArray.concat(result)
        });

        if (element.link) {
            var linkElement = [
                '<span>(',
                `<a target="_blank" href="${element.link}">Go to game</a>`,
                ')</span>'
            ]
            htmlAsArray = htmlAsArray.concat(linkElement)
        }

        var tail = '</div>'
        htmlAsArray = htmlAsArray.concat(tail)
        $("#result").append(htmlAsArray.join('\n'))
    });
}

function getPgnPlusMove(aMove) {
    pgn = game.pgn()
    splitPgn = pgn.split(" ")
    lineCount = Math.floor((splitPgn.length / 3))
    if(splitPgn.length % 3 == 0) {
        // create a new line
        pgn = pgn + " " + (lineCount + 1) + "."
    }
    pgn = pgn + " " + aMove
    return pgn
}

function displayPgn() {
    pgn = game.pgn()
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

function move(position) {
    game.move(position)
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

    updateStatus()
}

// update the board position after the piece snap
// for castling, en passant, pawn promotion
function onSnapEnd() {
    board.position(game.fen())
}

function updateStatus() {
    var status = ''

    var moveColor = 'White'
    if (game.turn() === 'b') {
        moveColor = 'Black'
    }

    // checkmate?
    if (game.in_checkmate()) {
        status = 'Game over, ' + moveColor + ' is in checkmate.'
    }

    // draw?
    else if (game.in_draw()) {
        status = 'Game over, drawn position'
    }

    // game still on
    else {
        status = moveColor + ' to move'

        // check?
        if (game.in_check()) {
            status += ', ' + moveColor + ' is in check'
        }
    }

    $status.html(status)
    $fen.html(game.fen())
    displayPgn()
    nextmove()
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