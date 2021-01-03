
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
    $pgn.html(game.pgn())
}

function resetClicked(e) {
    game.reset()
    board.position(game.fen())
    updateStatus()
    $fen.html(game.fen())
    $pgn.html(game.pgn())
}

function explore() {
    $("#result").html("");
    $.post("http://127.0.0.1:52825/explore", { pgn: game.pgn(), white: $white.val(), black: $black.val() }, function (data) {
        dataToHtml(JSON.parse(data));
    });
}

function dataToHtml(dataObject) {
    if (Array.isArray(dataObject) == false) {
        console.log("not an array")
        return
    }

    dataObject.forEach(element => {
        var htmlAsArray = [
            '<div>',
            element.total,
            `<a href="javascript:move('${element.move}');">${element.move}</a>`]

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
    $pgn.html(game.pgn())
    explore()
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