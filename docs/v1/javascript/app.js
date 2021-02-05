
// NOTE: this uses chessboardjs and chess.js libraries:
// https://github.com/oakmac/chessboardjs
// https://github.com/jhlywa/chess.js

var apiPort = "52825"
var board = null
var game = new Chess()

var $white = $('#white')
var $black = $('#black')
var $timecontrol = $('#timecontrol')
var $from = $('#from')
var $to = $('#to')
var $minelo = $('#minelo')
var $maxelo = $('#maxelo')
var $site = $('#site')
var browsingGame = ""
var nextMove = ""

var nextMovesTpl = document.getElementById('nextMovesTpl').innerHTML;
var usernameListTpl = document.getElementById('usernameListTpl').innerHTML;
var timecontrolListTpl = document.getElementById('timecontrolListTpl').innerHTML;
var nameListTpl = document.getElementById('nameListTpl').innerHTML;



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

$("#swap").click(function (e) {
    e.preventDefault();
    var black = $black.val()
    $black.val($white.val())
    $white.val(black)
    resetClicked()
    board.flip()
});

$("#undo").click(function (e) {
    e.preventDefault();
    browsingGame = ""
    game.undo()
    board.position(game.fen())
    updateStatus()
});

$("#next").click(function (e) {
    e.preventDefault();
    if (nextMove != "") {
        move(nextMove)
    }
});

$("#reset").click(function (e) {
    e.preventDefault();
    resetClicked()
});

function resetClicked() {
    browsingGame = ""
    game.reset()
    board.position(game.fen())
    updateStatus()
    updateReport()
}

$("#clear-usernames").click(function (e) {
    e.preventDefault();
    $white.val("")
    $black.val("")
    resetClicked()
});

$("#clear-timecontrols").click(function (e) {
    e.preventDefault();
    $timecontrol.val("")
    resetClicked()
});

$("#clear-sites").click(function (e) {
    e.preventDefault();
    $site.val("")
    resetClicked()
});


Array.prototype.remove = function () {
    var what, a = arguments, L = a.length, ax;
    while (L && this.length) {
        what = a[--L];
        while ((ax = this.indexOf(what)) !== -1) {
            this.splice(ax, 1);
        }
    }
    return this;
};


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
        else {
            values.remove(name)
            control.val(values.join(","))
            resetClicked()
        }
    }
}


function getNextMoves() {
    $("#result").html("");
    $.post(`http://127.0.0.1:${apiPort}/nextmoves`, {
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
        nextMovesToHtml(JSON.parse(data));
    });
}

function updateReport() {
    $.get(`http://127.0.0.1:${apiPort}/report`, {
        white: $white.val(),
        black: $black.val()
    }, function (data) {
        ret = JSON.parse(data);
        if (Array.isArray(ret.Sites) != false) {
            $("#siteNames").html(Mustache.render(nameListTpl, ret.Sites))
            $("#siteNames a").bind("click", function (e) {
                e.preventDefault();
                handleNameClicked(e, $site, $(this).html())
            });
        }
        if (Array.isArray(ret.Users) != false) {
            ret.Users.forEach((element) => {
                if (element.sitename == "lichess.org") {
                    element.imgpath = "/img/logos/lichessorg-48.png"
                }
                if (element.sitename == "chess.com") {
                    element.imgpath = "/img/logos/chesscom-48.png"
                }
            })
            $("#userNames").html(Mustache.render(usernameListTpl, ret.Users))
            $("#userNames a").bind("click", function (e) {
                e.preventDefault();
                username = $(this).html()
                if ($(this).data("sitename") == "chess.com") {
                    username = "c:" + username
                }
                if ($(this).data("sitename") == "lichess.org") {
                    username = "l:" + username
                }
                handleNameClicked(e, $white, username)
            });
        }
        if (Array.isArray(ret.TimeControls) != false) {
            ret.TimeControls.sort(compareTimecontrolsByName)
            $("#ultra-bullet-timeControlNames").html("")
            $("#bullet-timeControlNames").html("")
            $("#blitz-timeControlNames").html("")
            $("#rapid-timeControlNames").html("")
            $("#classic-timeControlNames").html("")
            if (ret.TimeControls.length > 10) {
                ret.TimeControls = groupTimecontrols(ret.TimeControls)
                // groups
                for (key in ret.TimeControls.grouped) {
                    $("#" + key + "-timeControlNames").html(Mustache.render(timecontrolListTpl, ret.TimeControls.grouped[key]))
                    $("#" + key + "-timeControlNames a").bind("click", function (e) {
                        e.preventDefault();
                        handleNameClicked(e, $timecontrol, $(this).html())
                    });
                }
                $(".timeControlLabel").show()
            }
            else {
                $("#timeControlNames").html(Mustache.render(timecontrolListTpl, ret.TimeControls))
                $("#timeControlNames a").bind("click", function (e) {
                    e.preventDefault();
                    handleNameClicked(e, $timecontrol, $(this).html())
                });
                $(".timeControlLabel").hide()
            }
        }
    });
}

function groupTimecontrols(timecontrolList) {
    timecontrolList.grouped = []
    timecontrolList.forEach((item) => {
        baseTimeStr = item.name.split("+")[0]
        if (!isNormalInteger(baseTimeStr)) {
            baseTime = Number.MAX_SAFE_INTEGER;
        }
        baseTime = parseInt(baseTimeStr)
        groupName = ""
        if (baseTime < 60) {
            groupName = "ultra-bullet"
        }
        else if (baseTime < 180) {
            groupName = "bullet"
        }
        else if (baseTime < 600) {
            groupName = "blitz"
        }
        else if (baseTime < 3600) {
            groupName = "rapid"
        }
        else {
            groupName = "classic"
        }
        if (timecontrolList.grouped[groupName] == undefined) {
            timecontrolList.grouped[groupName] = []
        }
        timecontrolList.grouped[groupName].push(item)
    })
    return timecontrolList
}


function isNormalInteger(str) {
    var n = Math.floor(Number(str));
    return n !== Infinity && String(n) === str && n >= 0;
}

function compareTimecontrolsByName(itemA, itemB) {
    a = itemA.name;
    b = itemB.name;

    intA = Number.MAX_SAFE_INTEGER;
    intB = Number.MAX_SAFE_INTEGER;

    if (isNormalInteger(a)) {
        intA = parseInt(a)
    }
    if (isNormalInteger(b)) {
        intB = parseInt(b)
    }
    if (intA == Number.MAX_SAFE_INTEGER) {
        // try the A+B form
        if (-1 != a.indexOf("+")) {
            splitA = a.split("+")
            if (isNormalInteger(splitA[0]) && isNormalInteger(splitA[1])) {
                intA = parseInt(splitA[0]) + parseInt(splitA[1])
            }
        }
    }
    if (intB == Number.MAX_SAFE_INTEGER) {
        // try the A+B form
        if (-1 != b.indexOf("+")) {
            splitB = b.split("+")
            if (isNormalInteger(splitB[0]) && isNormalInteger(splitB[1])) {
                intB = parseInt(splitB[0]) + parseInt(splitB[1])
            }
        }
    }

    return intA - intB
}

// Not used, we follow the link to chess.com or lichess.org for now
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

function nextMovesToHtml(dataObject) {
    nextMove = ""
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
            if (nextMove == "") {
                nextMove = element.move
            }
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

    $("#result").html(Mustache.render(nextMovesTpl, moves))
    $(".move").bind("click", function (e) {
        e.preventDefault();
        move($(this).html())
    });
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
    $('#pgn').html(pgn)
}

function move(aMove) {
    game.move(aMove)
    updateStatus()
    board.position(game.fen(), true)
}

// Board events
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
    $('#fen').html(game.fen())
    getNextMoves()
}

var config = {
    draggable: true,
    position: 'start',
    onDragStart: onDragStart,
    onDrop: onDrop,
    onSnapEnd: onSnapEnd
}
board = Chessboard('myBoard', config)

resetClicked()

