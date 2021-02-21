
// NOTE: this uses chessboardjs and chess.js libraries:
// https://github.com/oakmac/chessboardjs
// https://github.com/jhlywa/chess.js

var apiPort = '52825'
var board = null
var game = new Chess()

var $white = $('#white')
var $black = $('#black')
var $timecontrol = $('#timecontrol')
var $fromDate = $('#from')
var $toDate = $('#to')
var $minelo = $('#minelo')
var $maxelo = $('#maxelo')
var $site = $('#site')
var mostPopularMove = ''
var uiMode = 'opening' // opening, replay

var gameReplaying

var nextMovesTpl = document.getElementById('nextMovesTpl').innerHTML;
var usernameListTpl = document.getElementById('usernameListTpl').innerHTML;
var timecontrolListTpl = document.getElementById('timecontrolListTpl').innerHTML;
var nameListTpl = document.getElementById('nameListTpl').innerHTML;
var openingBreadcrumbsTpl = document.getElementById('openingBreadcrumbsTpl').innerHTML;
var replayBreadcrumbsTpl = document.getElementById('replayBreadcrumbsTpl').innerHTML;
var gameDetailsTpl = document.getElementById('gameDetailsTpl').innerHTML;



$fromDate.change(function () {
    getNextMoves()
    updateReport()
});

$toDate.change(function () {
    getNextMoves()
    updateReport()
});

$white.change(function () {
    getNextMoves()
});

$black.change(function () {
    getNextMoves()
});

$timecontrol.change(function () {
    getNextMoves()
});

$minelo.change(function () {
    getNextMoves()
});

$maxelo.change(function () {
    getNextMoves()
});

$site.change(function () {
    getNextMoves()
});


$('#swap').click(function (e) {
    e.preventDefault();
    var black = $black.val()
    $black.val($white.val())
    $white.val(black)
    getNextMoves()
    board.orientation('flip')
});

$('#undo').click(function (e) {
    e.preventDefault();
    game.undo()
    board.position(game.fen())
    if (uiMode == 'opening') {
        openingUpdated()
    }
    if (uiMode == 'replay') {
        highlightMove()
    }
});

$('#next').click(function (e) {
    e.preventDefault();
    if (uiMode == 'opening') {
        if (mostPopularMove != '') {
            move(mostPopularMove)
        }
    }
    else {
        replayNext()
    }
});

$('#reset-usernames').click(function (e) {
    e.preventDefault();
    $white.val('')
    $black.val('')
    getNextMoves()
    updateReport()
    board.orientation('white')
});

$('#reset-timecontrols').click(function (e) {
    e.preventDefault();
    $timecontrol.val('')
    getNextMoves()
});

$('#reset-dates').click(function (e) {
    e.preventDefault();
    $fromDate.val('')
    $toDate.val('')
    getNextMoves()
    updateReport()
});

$('#reset-sites').click(function (e) {
    e.preventDefault();
    $site.val('')
    getNextMoves()
});

$('#reset-elos').click(function (e) {
    e.preventDefault();
    $minelo.val('')
    $maxelo.val('')
    getNextMoves()
});

$('#reset').click(function (e) {
    e.preventDefault();
    resetBoard()
});

$('#opening-link').click(function (e) {
    e.preventDefault();
    if (uiMode == 'opening') {
        if ($('#edit-pgn').css('display') == 'none') {
            $('#edit-pgn').val(game.pgn())
            $('#edit-pgn').show()
            $('#edit-pgn').change(function () {
                $(this).hide()
                game.load_pgn($('#edit-pgn').val())
                board.position(game.fen())
                updateOpeningBreadcrumbs()
                getNextMoves()
            });
        }
        else {
            $('#edit-pgn').hide()
        }
    }
    if (uiMode == 'replay') {
        setOpeningMode()
        resetBoard()
    }
});

$('#reset-all').click(function (e) {
    e.preventDefault();
    board.orientation('white')
    setOpeningMode()
    $white.val('')
    $black.val('')
    $timecontrol.val('')
    $fromDate.val('')
    $toDate.val('')
    $site.val('')
    $minelo.val('')
    $maxelo.val('')
    resetBoard()
});


function resetBoard() {
    game.reset()
    board.position(game.fen())
    if (uiMode == 'opening') {
        openingUpdated()
        updateReport()
    }
    if (uiMode == 'replay') {
        highlightMove()
    }
}


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
    if (control.val().trim() == '' || !event.ctrlKey) {
        control.val(name)
        getNextMoves()
    }
    else {
        values = control.val().trim().split(',')
        if (values.indexOf(name) == -1) {
            values.push(name)
            control.val(values.join(','))
            getNextMoves()
        }
        else {
            values.remove(name)
            control.val(values.join(','))
            getNextMoves()
        }
    }
}


function getNextMoves() {
    $('#next-moves').html('');
    $.post(`http://127.0.0.1:${apiPort}/nextmoves`, {
        pgn: game.pgn(),
        white: $white.val(),
        black: $black.val(),
        timecontrol: $timecontrol.val(),
        from: $fromDate.val(),
        to: $toDate.val(),
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
        black: $black.val(),
        from: $fromDate.val(),
        to: $toDate.val()
    }, function (data) {
        ret = JSON.parse(data);
        if (Array.isArray(ret.Sites) != false) {
            $('#siteNames').html(Mustache.render(nameListTpl, ret.Sites))
            $('#siteNames a').bind('click', function (e) {
                e.preventDefault();
                handleNameClicked(e, $site, $(this).html())
            });
        }
        if (Array.isArray(ret.Users) != false) {
            ret.Users.forEach((element) => {
                if (element.sitename == 'lichess.org') {
                    element.imgpath = '/img/logos/lichessorg-48.png'
                }
                if (element.sitename == 'chess.com') {
                    element.imgpath = '/img/logos/chesscom-48.png'
                }
            })
            $('#userNames').html(Mustache.render(usernameListTpl, ret.Users))
            $('#userNames a').bind('click', function (e) {
                e.preventDefault();
                username = $(this).html()
                if ($(this).data('sitename') == 'chess.com') {
                    username = 'c:' + username
                }
                if ($(this).data('sitename') == 'lichess.org') {
                    username = 'l:' + username
                }
                if ($black.val() != '' && $white.val() == '') {
                    handleNameClicked(e, $black, username)
                }
                else {
                    handleNameClicked(e, $white, username)
                }
                updateReport()
            });
        }
        if (Array.isArray(ret.TimeControls) != false) {
            ret.TimeControls.sort(compareTimecontrolsByName)
            $('#ultra-bullet-timeControlNames').html('')
            $('#bullet-timeControlNames').html('')
            $('#blitz-timeControlNames').html('')
            $('#rapid-timeControlNames').html('')
            $('#classic-timeControlNames').html('')
            ret.TimeControls = groupTimecontrols(ret.TimeControls)
            // groups
            for (key in ret.TimeControls.grouped) {
                $('#' + key + '-timeControlNames').html(Mustache.render(timecontrolListTpl, ret.TimeControls.grouped[key]))
                $('#' + key + '-timeControlNames a').bind('click', function (e) {
                    e.preventDefault();
                    handleNameClicked(e, $timecontrol, $(this).html())
                });
            }
            $('.timeControlLabel').show()
        }
    });
}

function groupTimecontrols(timecontrolList) {
    timecontrolList.grouped = []
    timecontrolList.forEach((item) => {
        baseTimeStr = item.name.split('+')[0]
        if (!isNormalInteger(baseTimeStr)) {
            baseTime = Number.MAX_SAFE_INTEGER;
        }
        baseTime = parseInt(baseTimeStr)
        groupName = ''
        if (baseTime < 60) {
            groupName = 'ultra-bullet'
        }
        else if (baseTime < 180) {
            groupName = 'bullet'
        }
        else if (baseTime < 600) {
            groupName = 'blitz'
        }
        else if (baseTime < 3600) {
            groupName = 'rapid'
        }
        else {
            groupName = 'classic'
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
    int2A = 0;
    int2B = 0;

    if (isNormalInteger(a)) {
        intA = parseInt(a)
    }
    if (isNormalInteger(b)) {
        intB = parseInt(b)
    }
    if (intA == Number.MAX_SAFE_INTEGER) {
        // try the A+B form
        if (-1 != a.indexOf('+')) {
            splitA = a.split('+')
            if (isNormalInteger(splitA[0]) && isNormalInteger(splitA[1])) {
                intA = parseInt(splitA[0])
                int2A = parseInt(splitA[1])
            }
        }
    }
    if (intB == Number.MAX_SAFE_INTEGER) {
        // try the A+B form
        if (-1 != b.indexOf('+')) {
            splitB = b.split('+')
            if (isNormalInteger(splitB[0]) && isNormalInteger(splitB[1])) {
                intB = parseInt(splitB[0])
                int2B = parseInt(splitB[1])
            }
        }
    }

    if (intA == intB) {
        return int2A - int2B
    }
    else {
        return intA - intB
    }
}

function nextMovesToHtml(dataObject) {
    mostPopularMove = ''
    if (Array.isArray(dataObject) == false) {
        console.log('not an array')
        return
    }

    var moves = []

    dataObject.forEach(element => {

        winPercent = Math.round(100 * element.win / element.total)
        losePercent = Math.round(100 * element.lose / element.total)
        drawPercent = 100 - winPercent - losePercent
        winPercentText = ''
        if (winPercent > 12) {
            winPercentText = '' + winPercent + '%'
        }
        losePercentText = ''
        if (losePercent > 12) {
            losePercentText = '' + losePercent + '%'
        }
        drawPercentText = ''
        if (drawPercent > 12) {
            drawPercentText = '' + drawPercent + '%'
        }

        openingLink = false
        replayLink = false
        if (element.total == 1) {
            replayLink = true
            element.game.userlink = 'https://www.chess.com/member/'
            if (element.game.site == 'lichess.org') {
                element.game.userlink = 'https://lichess.org/@/'
            }
            // win,draw,lose
            win = false
            lose = false
            draw = false
            if (element.game.result == '1-0') {
                win = true
            } else if (element.game.result == '0-1') {
                lose = true
            } else {
                element.game.result = '1/2'
                draw = true
            }
            // date
            element.game.date = new Date(Date.parse(element.game.datetime)).toLocaleDateString()
            element.game.link = 'replay.html?gameId=' + element.game._id + '&skip=' + (game.history().length) + '&orientation=' + board.orientation()
            moves.push({
                openingLink: openingLink,
                replayLink: replayLink,
                win: win,
                lose: lose,
                draw: draw,
                game: element.game,
                move: element.move,
            })
        }
        else {
            openingLink = true
            if (mostPopularMove == '') {
                mostPopularMove = element.move
            }
            moves.push({
                openingLink: openingLink,
                replayLink: replayLink,
                move: element.move,
                total: element.total,
                winPercent: winPercent,
                losePercent: losePercent,
                drawPercent: drawPercent,
                winPercentText: winPercentText,
                losePercentText: losePercentText,
                drawPercentText: drawPercentText,
            })
        }

    });

    $('#next-moves').html(Mustache.render(nextMovesTpl, moves))
    $('.next-move').bind('click', function (e) {
        e.preventDefault();
        move($(this).html())
    });
    $('.replay-game').bind('click', function (e) {
        e.preventDefault();
        replayGame($(this).attr('data-gameid'))
    });
}

function setReplayMode() {
    uiMode = 'replay'
    $('#filter').hide()
    $('#next-moves').hide()
    $('#values').hide()
    $('#swap').hide()
    $('#replay').show()
    $('#game-details').show()
}

function setOpeningMode() {
    uiMode = 'opening'
    $('#filter').show()
    $('#next-moves').show()
    $('#values').show()
    $('#swap').show()
    $('#replay').hide()
    $('#game-details').hide()
}

function replayGame(gameId) {
    setReplayMode()
    // load data
    $.get(`http://127.0.0.1:${apiPort}/game`, { gameId: gameId }, function (jsonData) {
        data = JSON.parse(jsonData)
        splitPgn = data.pgn.split(' ')
        gameReplaying = []
        splitPgn.forEach((value, index) => {
            round = Math.floor(index / 3)
            if (index % 3 == 0) {
                gameReplaying.push({
                    index: round,
                    round: value,
                    white: '',
                    black: '',
                    isComplete: false
                })
            }
            if (index % 3 == 1) {
                gameReplaying[round].white = value
            }
            else {
                gameReplaying[round].black = value
                gameReplaying[round].isComplete = true
            }
        })
        data.dateStr = new Date(data.datetime).toGMTString()
        $('#game-details').html(Mustache.render(gameDetailsTpl, data))
        $('#replay').html(Mustache.render(replayBreadcrumbsTpl, gameReplaying))
        $('#replay a').bind('click', function (e) {
            e.preventDefault();
            round = $(this).attr('data-index')
            color = $(this).attr('data-color')
            game.reset()
            for (i = 0; i < round; i++) {
                game.move(gameReplaying[i].white)
                game.move(gameReplaying[i].black)
            }
            game.move(gameReplaying[round].white)
            if (color == 'black') {
                game.move(gameReplaying[round].black)
            }
            board.position(game.fen(), true)
            highlightMove()
        });
        // replay first move after opening
        replayNext()
    });
}

function replayNext() {
    round = Math.floor(game.history().length / 2)
    if (game.history().length % 2 == 0) {
        move(gameReplaying[round].white)
    }
    else {
        move(gameReplaying[round].black)
    }
    highlightMove()
}


function highlightMove() {
    $('#replay a').parent().removeClass('highlight')
    if (game.history().length > 0) {
        round = Math.floor((game.history().length - 1) / 2)
        color = 'white'
        if ((game.history().length - 1) % 2 == 1) {
            color = 'black'
        }
        $('#replay a[data-index="' + round + '"][data-color="' + color + '"]').parent().addClass('highlight')
    }
}

function updateOpeningBreadcrumbs() {
    splitBreadcrumbs = []
    game.history().forEach((value, index) => {
        round = Math.floor(index / 2)
        if (index % 2 == 0) {
            splitBreadcrumbs.push({
                index: round,
                round: (round + 1) + '.',
                white: value,
                black: '',
                isComplete: false
            })
        }
        else {
            splitBreadcrumbs[round].black = value
            splitBreadcrumbs[round].isComplete = true
        }
    })

    $('#opening').html(Mustache.render(openingBreadcrumbsTpl, splitBreadcrumbs))

    $('#opening a').bind('click', function (e) {
        e.preventDefault();
        indexInHistory = 2 * $(this).attr('data-index')
        if ($(this).attr('data-color') == 'black') {
            indexInHistory += 1
        }
        saveHistory = game.history()
        game.reset()
        for (i = 0; i < indexInHistory + 1; i++) {
            game.move(saveHistory[i])
        }
        setOpeningMode()
        board.position(game.fen())
        getNextMoves()
        updateOpeningBreadcrumbs()
    });
}



function move(aMove) {
    game.move(aMove)
    if (uiMode == 'opening') {
        openingUpdated()
    }
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

    openingUpdated()
}

// update the board position after the piece snap
// for castling, en passant, pawn promotion
function onSnapEnd() {
    board.position(game.fen())
}

function openingUpdated() {
    updateOpeningBreadcrumbs()
    getNextMoves()
}

var config = {
    moveSpeed: 400,
    draggable: true,
    position: 'start',
    onDragStart: onDragStart,
    onDrop: onDrop,
    onSnapEnd: onSnapEnd
}
board = Chessboard('myBoard', config)
board.resize()

resetBoard()

