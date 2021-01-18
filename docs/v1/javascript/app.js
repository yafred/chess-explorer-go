
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
var $timecontrol = $('#timecontrol')
var $from = $('#from')
var $to = $('#to')
var $minelo = $('#minelo')
var $maxelo = $('#maxelo')
var $site = $('#site')
var browsingGame = ""

var nextMoveTpl = document.getElementById('nextMoveTpl').innerHTML;
var nameListTpl = document.getElementById('nameListTpl').innerHTML;



function userNameClicked(name) {
    console.log(`user ${name}`)
}

function siteNameClicked(name) {
    console.log(`site ${name}`)
}

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

$timecontrol.change(function () {
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

function flipClicked() {
    board.flip()
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

function clearClicked(type) {
    switch (type) {
        case "site":
            $site.val("")
            break;
        case "username":
            $white.val("")
            $black.val("")
            break;
        case "timecontrol":
            $timecontrol.val("")
            break;
        default:
            break;
    }
    resetClicked()
}

function handleNameClicked(event, control, name) {
    if (control.val().trim() == "" || !event.ctrlKey) {
        control.val(name)
        resetClicked()
    }
    else {
        values = control.val().trim().split(",")
        if (values.indexOf(name) == -1) {
            values.push(name)
            control.val(values.join(","))
            resetClicked()
        }
    }
}


function getNextMove() {
    $("#result").html("");
    $.post(`http://127.0.0.1:${apiPort}/nextmove`, {
        pgn: game.pgn(),
        white: $white.val(),
        black: $black.val(),
        timecontrol: $timecontrol.val(),
        from: $from.val(),
        to: $to.val(),
        minelo: $minelo.val(),
        maxelo: $maxelo.val(),
        site: $site.val()
    }, function (data) {
        nextMoveToHtml(JSON.parse(data));
    });
}

function updateReport() {
    $.get(`http://127.0.0.1:${apiPort}/report`, function (data) {
        ret = JSON.parse(data);
        if (Array.isArray(ret.Sites) != false) {
            $("#siteNames").html(Mustache.render(nameListTpl, ret.Sites))
            $("#siteNames a").bind("click", function (e) {
                e.preventDefault();
                handleNameClicked(e, $site, $(this).html())
            });
        }
        if (Array.isArray(ret.UsersAsWhite) != false) {
            $("#userNames").html(Mustache.render(nameListTpl, ret.UsersAsWhite))
            $("#userNames a").bind("click", function (e) {
                e.preventDefault();
                handleNameClicked(e, $white, $(this).html())
            });
        }
        if (Array.isArray(ret.TimeControls) != false) {
            $("#timeControlNames").html(Mustache.render(nameListTpl, ret.TimeControls))
            $("#timeControlNames a").bind("click", function (e) {
                e.preventDefault();
                handleNameClicked(e, $timecontrol, $(this).html())
            });
        }
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
        losePercent = Math.round(100 * element.lose / element.total)
        drawPercent = 100 - winPercent - losePercent
        drawPercentText = ""
        if (drawPercent > 12) {
            drawPercentText = "" + drawPercent + "%"
        }

        internalLink = false
        externalLink = false
        if (element.total == 1) {
            externalLink = true
            element.game.userlink = "https://www.chess.com/member/"
            if (element.game.site == "lichess.org") {
                element.game.userlink = "https://lichess.org/@/"
            }
            // win,draw,lose
            win = false
            lose = false
            draw = false
            if (element.game.result == "1-0") {
                win = true
            } else if (element.game.result == "0-1") {
                lose = true
            } else {
                element.game.result = "1/2"
                draw = true
            }
            // date
            element.game.date = new Date(Date.parse(element.game.datetime)).toLocaleDateString()
            moves.push({
                internalLink: internalLink,
                externalLink: externalLink,
                win: win,
                lose: lose,
                draw: draw,
                game: element.game,
                move: element.move,
            })
        }
        else {
            internalLink = true
            moves.push({
                internalLink: internalLink,
                externalLink: externalLink,
                move: element.move,
                total: element.total,
                winPercent: winPercent,
                drawPercent: drawPercent,
                drawPercentText: drawPercentText,
                losePercent: losePercent,
            })
        }

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

updateReport()