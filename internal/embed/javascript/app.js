// 'use strict'; // turn on for test


// NOTE: this uses chessboardjs and chess.js libraries:
// https://github.com/oakmac/chessboardjs
// https://github.com/jhlywa/chess.js

var apiHost = location.protocol + '//' + location.host
var board = null
var game = new Chess()

// states
var mostPopularMove = ''
var simplifyTimecontrol = true // make m+s equivalent to m (for example: 600 will include 600+5) and 1/n equivalent to -
var uiMode = 'opening' // opening, replay
var playerInputMode = 'white' // changes when input fields are clicked
var gameReplaying

// mustache templates
var nextMovesTpl = document.getElementById('nextMovesTpl').innerHTML;
var usernameListTpl = document.getElementById('usernameListTpl').innerHTML;
var timecontrolListTpl = document.getElementById('timecontrolListTpl').innerHTML;
var nameListTpl = document.getElementById('nameListTpl').innerHTML;
var openingBreadcrumbsTpl = document.getElementById('openingBreadcrumbsTpl').innerHTML;
var replayBreadcrumbsTpl = document.getElementById('replayBreadcrumbsTpl').innerHTML;
var gameDetailsTpl = document.getElementById('gameDetailsTpl').innerHTML;


// events
$('#from').change(function() {
    getNextMoves()
    updateReport()
});

$('#to').change(function() {
    getNextMoves()
    updateReport()
});

$('#white').change(function() {
    getNextMoves()
});

$('#black').change(function() {
    getNextMoves()
});

$('#white').click(function(e) {
    playerInputMode = 'white'
});

$('#black').click(function(e) {
    playerInputMode = 'black'
});

$('#timecontrol').change(function() {
    getNextMoves()
});

$('#minelo').change(function() {
    getNextMoves()
});

$('#maxelo').change(function() {
    getNextMoves()
});

$('#site').change(function() {
    getNextMoves()
});

$('#swap').click(function(e) {
    e.preventDefault();
    var black = $('#black').val()
    $('#black').val($('#white').val())
    $('#white').val(black)
    getNextMoves()
    board.orientation('flip')
});

$('#undo').click(function(e) {
    e.preventDefault();
    game.undo()

    board.position(game.fen())
    if (uiMode == 'opening') {
        openingUpdated()
    }
    if (uiMode == 'replay') {
        highlightMoveOnReplayList()
    }
});

$('#next').click(function(e) {
    e.preventDefault();
    if (uiMode == 'opening') {
        if (mostPopularMove != '') {
            move(mostPopularMove)
        }
    } else {
        replayNext()
    }
});

$('#reset-usernames').click(function(e) {
    e.preventDefault();
    $('#white').val('')
    $('#black').val('')
    getNextMoves()
    updateReport()
    board.orientation('white')
});

$('#reset-timecontrols').click(function(e) {
    e.preventDefault();
    $('#timecontrol').val('')
    getNextMoves()
});

$('#reset-dates').click(function(e) {
    e.preventDefault();
    $('#from').val('')
    $('#to').val('')
    getNextMoves()
    updateReport()
});

$('#reset-sites').click(function(e) {
    e.preventDefault();
    $('#site').val('')
    getNextMoves()
});

$('#reset-elos').click(function(e) {
    e.preventDefault();
    $('#minelo').val('')
    $('#maxelo').val('')
    getNextMoves()
});

$('#reset').click(function(e) {
    e.preventDefault();
    resetBoard()
});

$('#show-filter').click(function(e) {
    $('#show-filter').hide()
    $('#book-moves-panel').hide()
    $('#show-book-moves').show()
    $('#filter').show()
    e.preventDefault();
});

$('#show-game-details').click(function(e) {
    $('#show-game-details').hide()
    $('#book-moves-panel').hide()
    $('#show-book-moves').show()
    $('#game-details').show()
    e.preventDefault();
});

$('#show-book-moves').click(function(e) {
    $('#show-book-moves').hide()
    $('#filter').hide()
    $('#game-details').hide()
    if (uiMode == 'opening') {
        $('#show-filter').show()
    } else {
        $('#show-game-details').show()
    }
    $('#book-moves-panel').show()
    updateBookMoves()
    e.preventDefault();
});

$('#edit-pgn-link').click(function(e) {
    e.preventDefault();
    if ($('#edit-pgn').css('display') == 'none') {
        $('#edit-pgn').val(game.pgn())
        $('#edit-pgn').show()
        $('#edit-pgn').change(function() {
            $(this).hide()
            game.load_pgn($('#edit-pgn').val())
            board.position(game.fen())
            openingUpdated()
        });
    } else {
        $('#edit-pgn').hide()
    }
});

$('#back-to-opening-link').click(function(e) {
    e.preventDefault();
    setOpeningMode()
    game.load_pgn($('#edit-pgn').val())
    board.position(game.fen())
    openingUpdated()
});

$('#opening-link').click(function(e) {
    e.preventDefault();
    setOpeningMode()
    resetBoard()
});

$('#reset-all').click(function(e) {
    e.preventDefault();
    board.orientation('white')
    setOpeningMode()
    $('#white').val('')
    $('#black').val('')
    $('#timecontrol').val('')
    $('#from').val('')
    $('#to').val('')
    $('#site').val('')
    $('#minelo').val('')
    $('#maxelo').val('')
    resetBoard()
});

$('#simplify-timecontrol-checked').click(function(e) {
    e.preventDefault();
    $(this).hide()
    simplifyTimecontrol = false
    updateReport()
    $('#simplify-timecontrol-unchecked').show()
});

$('#simplify-timecontrol-unchecked').click(function(e) {
    e.preventDefault();
    $(this).hide()
    simplifyTimecontrol = true
    updateReport()
    $('#simplify-timecontrol-checked').show()
});

$('#show-fen-checked').click(function(e) {
    e.preventDefault();
    $('#fen-container').hide()
    $(this).hide()
    $('#show-fen-unchecked').show()
});

$('#show-fen-unchecked').click(function(e) {
    e.preventDefault();
    $('#fen-container').show()
    $(this).hide()
    $('#show-fen-checked').show()
});

$('#show-search-fen-form').click(function(e) {
    e.preventDefault();
    $('#search-fen-form').show()
});

$('#cancel-search-fen-form').click(function(e) {
    e.preventDefault();
    $('#search-fen-form').hide()
});

$('#search-fen').click(function(e) {
    e.preventDefault();
    $('#search-fen-form').hide()
    $.post(`${apiHost}/searchfen`, {
        fen: $('#fen-input').val(),
        maxMoves: $('#search-fen-max-moves').val(),
        pgn: game.pgn(),
        white: $('#white').val(),
        black: $('#black').val(),
        timecontrol: $('#timecontrol').val(),
        simplifyTimecontrol: simplifyTimecontrol,
        from: $('#from').val(),
        to: $('#to').val(),
        minelo: $('#minelo').val(),
        maxelo: $('#maxelo').val(),
        site: $('#site').val()
    }, function(response) {}).fail(function() {
        showError('Error connecting to ' + apiHost)
    });

});

$('#book-moves-db-master-unchecked').click(function(e) {
    e.preventDefault();
    $(this).hide()
    $('#book-moves-db-master-checked').show()
    $('#book-moves-db-lichess-checked').hide()
    $('#book-moves-db-lichess-unchecked').show()
    updateBookMoves()
});

$('#book-moves-db-lichess-unchecked').click(function(e) {
    e.preventDefault();
    $(this).hide()
    $('#book-moves-db-lichess-checked').show()
    $('#book-moves-db-master-checked').hide()
    $('#book-moves-db-master-unchecked').show()
    updateBookMoves()
});

['#book-moves-rating-1600', '#book-moves-rating-1800', '#book-moves-rating-2000', '#book-moves-rating-2200', '#book-moves-rating-2500',
    '#book-moves-speed-bullet', '#book-moves-speed-blitz', '#book-moves-speed-rapid', '#book-moves-speed-classical'
].forEach(element => {
    $(`${element}-unchecked`).click(function(e) {
        e.preventDefault();
        $(this).hide()
        $(`${element}-checked`).show()
        updateBookMoves()
    });

    $(`${element}-checked`).click(function(e) {
        e.preventDefault();
        $(this).hide()
        $(`${element}-unchecked`).show()
        updateBookMoves()
    });
})

// functions
function resetBoard() {
    game.reset()
    board.position(game.fen())
    if (uiMode == 'opening') {
        openingUpdated()
        updateReport()
    }
    if (uiMode == 'replay') {
        highlightMoveOnReplayList()
    }
}


Array.prototype.remove = function() {
    var what, a = arguments,
        L = a.length,
        ax;
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
    } else {
        values = control.val().trim().split(',')
        if (values.indexOf(name) == -1) {
            values.push(name)
            control.val(values.join(','))
            getNextMoves()
        } else {
            values.remove(name)
            control.val(values.join(','))
            getNextMoves()
        }
    }
}


function getNextMoves() {
    $('#next-moves').html('');
    $.post(`${apiHost}/nextmoves`, {
        pgn: game.pgn(),
        white: $('#white').val(),
        black: $('#black').val(),
        timecontrol: $('#timecontrol').val(),
        simplifyTimecontrol: simplifyTimecontrol,
        from: $('#from').val(),
        to: $('#to').val(),
        minelo: $('#minelo').val(),
        maxelo: $('#maxelo').val(),
        site: $('#site').val()
    }, function(response) {
        var jsonResponse = JSON.parse(response)
        if (jsonResponse.error != undefined && jsonResponse.error != '') {
            showError(jsonResponse.error)
        } else {
            handleNextMovesResponse(jsonResponse.data);
        }
    }).fail(function() {
        showError('Error connecting to ' + apiHost)
    });
}

function updateReport() {
    clearTimeControlValues()
    $.get(`${apiHost}/report`, {
        white: $('#white').val(),
        black: $('#black').val(),
        from: $('#from').val(),
        to: $('#to').val()
    }, function(response) {
        var jsonResponse = JSON.parse(response);
        if (jsonResponse.error != undefined && jsonResponse.error != '') {
            showError(jsonResponse.error)
        } else {
            handleReportResponse(jsonResponse.data)
        }
    }).fail(function() {
        showError('Error connecting to ' + apiHost)
    });
}


function handleReportResponse(data) {
    if (Array.isArray(data.Sites) != false) {
        $('#siteNames').html(Mustache.render(nameListTpl, data.Sites))
        $('#siteNames a').bind('click', function(e) {
            e.preventDefault();
            handleNameClicked(e, $('#site'), $(this).html())
        });
    }
    if (Array.isArray(data.Users) != false) {
        data.Users.forEach((element) => {
            if (element.sitename == 'lichess.org') {
                element.imgpath = '/img/logos/lichessorg-48.png'
            }
            if (element.sitename == 'chess.com') {
                element.imgpath = '/img/logos/chesscom-48.png'
            }
        })
        $('#userNames').html(Mustache.render(usernameListTpl, data.Users))
        $('#userNames a').bind('click', function(e) {
            e.preventDefault();
            var username = $(this).html()
            if ($(this).data('sitename') == 'chess.com') {
                username = 'c:' + username
            }
            if ($(this).data('sitename') == 'lichess.org') {
                username = 'l:' + username
            }
            if (playerInputMode == 'black') {
                handleNameClicked(e, $('#black'), username)
            } else {
                handleNameClicked(e, $('#white'), username)
            }
            updateReport()
        });
    }
    if (Array.isArray(data.TimeControls) != false) {
        var timeControlList = data.TimeControls
        if (simplifyTimecontrol) {
            timeControlList = reduceTimeControlList(timeControlList)
        }
        timeControlList.sort(compareTimecontrolsByName)
        timeControlList = groupTimecontrols(timeControlList)
            // groups
        for (var key in timeControlList.grouped) {
            $('#' + key + '-timeControlNames').html(Mustache.render(timecontrolListTpl, timeControlList.grouped[key]))
            $('#' + key + '-timeControlNames a').bind('click', function(e) {
                e.preventDefault();
                handleNameClicked(e, $('#timecontrol'), $(this).html())
            });
        }
        $('.timeControlLabel').show()
    }
}

function clearTimeControlValues() {
    $('#ultra-bullet-timeControlNames').html('')
    $('#bullet-timeControlNames').html('')
    $('#blitz-timeControlNames').html('')
    $('#rapid-timeControlNames').html('')
    $('#classic-timeControlNames').html('')
}

function reduceTimeControlList(timecontrolList) {
    var reducedList = []
    timecontrolList.forEach((item) => {
        var baseTimeStr = item.name.split('+')[0]
        if (item.name.indexOf('/') != -1) {
            baseTimeStr = '-'
        }
        var result = reducedList.find(({ name }) => name === baseTimeStr);
        if (result == undefined) {
            reducedList.push({ name: baseTimeStr, count: 0 }) // we don't use count yet (we could reduce the list after sorting to be able to count)
        }
    })
    return reducedList
}

function groupTimecontrols(timecontrolList) {
    timecontrolList.grouped = []
    timecontrolList.forEach((item) => {
        var baseTimeStr = ''
        if (item.name.indexOf('/') != -1) {
            baseTimeStr = '-'
        } else {
            baseTimeStr = item.name.split('+')[0]
        }
        if (!isNormalInteger(baseTimeStr)) {
            baseTime = Number.MAX_SAFE_INTEGER;
        }
        var baseTime = parseInt(baseTimeStr)
        var groupName = ''
        if (baseTime < 60) {
            groupName = 'ultra-bullet'
        } else if (baseTime < 180) {
            groupName = 'bullet'
        } else if (baseTime < 600) {
            groupName = 'blitz'
        } else if (baseTime < 3600) {
            groupName = 'rapid'
        } else {
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
    var a = itemA.name;
    var b = itemB.name;

    var intA = Number.MAX_SAFE_INTEGER;
    var intB = Number.MAX_SAFE_INTEGER;
    var int2A = 0;
    var int2B = 0;

    // Types of time control and there order
    // 0:600 1:600+n 2:1/n 3:-
    var typeA = 3
    var typeB = 3

    if (isNormalInteger(a)) {
        typeA = 0
        intA = a
    } else if (-1 != a.indexOf('+')) {
        typeA = 1
        var splitA = a.split('+')
        if (isNormalInteger(splitA[0]) && isNormalInteger(splitA[1])) {
            intA = parseInt(splitA[0])
            int2A = parseInt(splitA[1])
        }
    } else if (-1 != a.indexOf('/')) {
        typeA = 2
        var splitA = a.split('/')
        if (isNormalInteger(splitA[0]) && isNormalInteger(splitA[1])) {
            intA = parseInt(splitA[0])
            int2A = parseInt(splitA[1])
        }
    }

    if (isNormalInteger(b)) {
        typeB = 0
        intB = b
    } else if (-1 != b.indexOf('+')) {
        typeB = 1
        var splitB = b.split('+')
        if (isNormalInteger(splitB[0]) && isNormalInteger(splitB[1])) {
            intB = parseInt(splitB[0])
            int2B = parseInt(splitB[1])
        }
    } else if (-1 != b.indexOf('/')) {
        typeB = 2
        var splitB = b.split('/')
        if (isNormalInteger(splitB[0]) && isNormalInteger(splitB[1])) {
            intB = parseInt(splitB[0])
            int2B = parseInt(splitB[1])
        }
    }

    //  log('a ' + a + typeA + ' ' + intA + ' ' + int2A + ' ... b ' + b + typeB + ' ' + intB + ' ' + int2B)
    if (typeA != typeB && !(typeA == 0 && typeB == 1) && !(typeB == 0 && typeA == 1)) {
        return typeA - typeB
    }
    if (intA == intB) {
        if (typeA != typeB) {
            return typeA - typeB // case of n compared to n+0
        } else {
            return int2A - int2B
        }
    } else {
        return intA - intB
    }
}

function handleNextMovesResponse(dataObject) {
    mostPopularMove = ''
    if (Array.isArray(dataObject) == false) {
        console.log('not an array')
        $('#total-games').html(0)
        return
    }

    var moves = []
    var grandTotal = 0
    var grandWhite = 0
    var grandBlack = 0

    dataObject.forEach(element => {

        grandWhite += element.white
        grandBlack += element.black
        var whitePercent = Math.round(100 * element.white / element.total)
        var blackPercent = Math.round(100 * element.black / element.total)
        var drawPercent = 100 - whitePercent - blackPercent
        var whitePercentText = ''
        if (whitePercent > 12) {
            whitePercentText = '' + whitePercent + '%'
        }
        var blackPercentText = ''
        if (blackPercent > 12) {
            blackPercentText = '' + blackPercent + '%'
        }
        var drawPercentText = ''
        if (drawPercent > 12) {
            drawPercentText = '' + drawPercent + '%'
        }

        var openingLink = false
        var replayLink = false
        if (element.total == 1) {
            replayLink = true
            element.game.userlink = 'https://www.chess.com/member/'
            if (element.game.site == 'lichess.org') {
                element.game.userlink = 'https://lichess.org/@/'
            }
            // white,draw,black
            var white = false
            var black = false
            var draw = false
            if (element.game.result == '1-0') {
                white = true
            } else if (element.game.result == '0-1') {
                black = true
            } else {
                element.game.result = '1/2'
                draw = true
            }
            // date
            element.game.date = new Date(Date.parse(element.game.datetime)).toLocaleDateString()
            moves.push({
                openingLink: openingLink,
                replayLink: replayLink,
                white: white,
                black: black,
                draw: draw,
                game: element.game,
                move: element.move,
            })
        } else {
            openingLink = true
            if (mostPopularMove == '') {
                mostPopularMove = element.move
            }
            moves.push({
                openingLink: openingLink,
                replayLink: replayLink,
                move: element.move,
                total: element.total,
                whitePercent: whitePercent,
                blackPercent: blackPercent,
                drawPercent: drawPercent,
                whitePercentText: whitePercentText,
                blackPercentText: blackPercentText,
                drawPercentText: drawPercentText,
            })
        }
        grandTotal += element.total
    });

    var grandWhitePercent = Math.round(100 * grandWhite / grandTotal)
    var grandBlackPercent = Math.round(100 * grandBlack / grandTotal)
    var grandDrawPercent = 100 - grandWhitePercent - grandBlackPercent

    $('#next-moves').html(Mustache.render(nextMovesTpl, moves))
    $('.next-move').bind('click', function(e) {
        e.preventDefault();
        move($(this).html())
    });
    $('.replay-game').bind('click', function(e) {
        e.preventDefault();
        replayGame($(this).attr('data-gameid'))
    });
    $('#total-games').html(grandTotal + ' (' + grandWhitePercent + '/' + grandDrawPercent + '/' + grandBlackPercent + ')&percnt;')
}

function setReplayMode() {
    uiMode = 'replay'
    $('#back-to-opening-link').show()
    $('#edit-pgn-link').hide()
    $('#edit-pgn').hide()
    $('#filter').hide()
    $('#next-moves').hide()
    $('#values').hide()
    $('#swap').hide()
    $('#reset-all').hide()
    $('#replay').show()
    $('#game-details').show()
    $('#fen-container').hide()
    $('#total-games').hide()
    $('#show-book-moves').show()
    $('#filter').hide()
    $('#show-filter').hide()
    $('#book-moves-panel').hide()

}

function setOpeningMode() {
    uiMode = 'opening'
    $('#back-to-opening-link').hide()
    $('#edit-pgn-link').show()
    $('#filter').show()
    $('#next-moves').show()
    $('#values').show()
    $('#swap').show()
    $('#reset-all').show()
    $('#replay').hide()
    $('#game-details').hide()
    $('#fen-container').hide()
    if ($('#show-fen-checked').is(':visible')) {
        $('#fen-container').show()
    }
    $('#total-games').show()
    $('#show-book-moves').show()
    $('#filter').show()
    $('#show-filter').hide()
    $('#book-moves-panel').hide()

}

function replayGame(gameId) {
    setReplayMode()
        // load data
    $.get(`${apiHost}/game`, { gameId: gameId }, function(response) {
        var jsonResponse = JSON.parse(response)
        if (jsonResponse.error != undefined && jsonResponse.error != '') {
            showError(jsonResponse.error)
        } else {
            handleGameResponse(jsonResponse.data)
        }
    }).fail(function() {
        showError('Error connecting to ' + apiHost)
    });
}

function handleGameResponse(data) {
    var splitPgn = data.pgn.split(' ')
    gameReplaying = []
    splitPgn.forEach((value, index) => {
        var round = Math.floor(index / 3)
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
        } else {
            gameReplaying[round].black = value
            gameReplaying[round].isComplete = true
        }
    })
    data.dateStr = new Date(data.datetime).toGMTString()
    $('#game-details').html(Mustache.render(gameDetailsTpl, data))
    $('#replay').html(Mustache.render(replayBreadcrumbsTpl, gameReplaying))
    $('#replay a').bind('click', function(e) {
        e.preventDefault();
        var round = $(this).attr('data-index')
        var color = $(this).attr('data-color')
        game.reset()
        for (var i = 0; i < round; i++) {
            game.move(gameReplaying[i].white)
            game.move(gameReplaying[i].black)
        }
        game.move(gameReplaying[round].white)
        if (color == 'black') {
            game.move(gameReplaying[round].black)
        }
        board.position(game.fen(), true)
        highlightMoveOnReplayList()
    });
    // replay first move after opening
    replayNext()
}


function showError(error) {
    $('#values').hide()
    $('#error').show()
    $('#error').html('<p>' + error + '</p>')
}

function replayNext() {
    var round = Math.floor(game.history().length / 2)
    if (game.history().length % 2 == 0) {
        move(gameReplaying[round].white)
    } else {
        move(gameReplaying[round].black)
    }
    highlightMoveOnReplayList()
}


function highlightMoveOnReplayList() {
    $('#replay a').parent().removeClass('highlight')
    if (game.history().length > 0) {
        var round = Math.floor((game.history().length - 1) / 2)
        var color = 'white'
        if ((game.history().length - 1) % 2 == 1) {
            color = 'black'
        }
        $('#replay a[data-index="' + round + '"][data-color="' + color + '"]').parent().addClass('highlight')
    }
    updateBookMoves()
    highlightLastMoveOnBoard()
}

function updateOpeningBreadcrumbs() {
    var splitBreadcrumbs = []
    game.history().forEach((value, index) => {
        var round = Math.floor(index / 2)
        if (index % 2 == 0) {
            splitBreadcrumbs.push({
                index: round,
                round: (round + 1) + '.',
                white: value,
                black: '',
                isComplete: false
            })
        } else {
            splitBreadcrumbs[round].black = value
            splitBreadcrumbs[round].isComplete = true
        }
    })

    $('#edit-pgn').val(game.pgn())
    $('#opening').html(Mustache.render(openingBreadcrumbsTpl, splitBreadcrumbs))

    $('#opening a').bind('click', function(e) {
        e.preventDefault();
        var indexInHistory = 2 * $(this).attr('data-index')
        if ($(this).attr('data-color') == 'black') {
            indexInHistory += 1
        }
        var saveHistory = game.history()
        game.reset()
        for (var i = 0; i < indexInHistory + 1; i++) {
            game.move(saveHistory[i])
        }
        setOpeningMode()
        board.position(game.fen())
        openingUpdated()
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

    if (uiMode == 'opening') {
        openingUpdated()
    }
}

// update the board position after the piece snap
// for castling, en passant, pawn promotion
function onSnapEnd() {
    board.position(game.fen())
}

function openingUpdated() {
    highlightLastMoveOnBoard()
    updateOpeningBreadcrumbs()
    getNextMoves()
    $('#fen').html(game.fen())
    updateBookMoves()
}

function updateBookMoves() {
    /* we need to call book moves for opening name even if book moves panel is hidden
    if ($('#book-moves-panel').is(':visible') == false) {
        bookType = 'master'
    }
    */
    // update opening name
    var bookType = 'lichess' // 'lichess', 'master'
    if ($('#book-moves-db-master-checked').is(':visible')) {
        bookType = 'master'
    }

    const allRatings = [1600, 1800, 2000, 2200, 2500]
    const allSpeeds = ['bullet', 'blitz', 'rapid', 'classical']
    var ratings = []
    var speeds = []
    allRatings.forEach(element => {
        if ($(`#book-moves-rating-${element}-checked`).is(':visible')) {
            ratings.push(element)
        }
    })
    allSpeeds.forEach(element => {
        if ($(`#book-moves-speed-${element}-checked`).is(':visible')) {
            speeds.push(element)
        }
    })

    $.get(`https://explorer.lichess.ovh/${bookType}`, {
        fen: game.fen(),
        variant: 'standard',
        'ratings[]': ratings,
        'speeds[]': speeds
    }, function(jsonResponse) {
        if (jsonResponse.opening === null) {
            if (game.pgn() == '') {
                $('#opening-name').html('')
            }
        } else {
            $('#opening-name').html(jsonResponse.opening.name)
        }
        if (jsonResponse.moves === null || Array.isArray(jsonResponse.moves) == false) {
            $('#book-moves').html('')
        } else {
            var moves = []
            jsonResponse.moves.forEach(element => {
                var total = element.white + element.black + element.draws
                var whitePercent = Math.round(100 * element.white / total)
                var blackPercent = Math.round(100 * element.black / total)
                var drawPercent = 100 - whitePercent - blackPercent

                var whitePercentText = ''
                if (whitePercent > 12) {
                    whitePercentText = '' + whitePercent + '%'
                }
                var blackPercentText = ''
                if (blackPercent > 12) {
                    blackPercentText = '' + blackPercent + '%'
                }
                var drawPercentText = ''
                if (drawPercent > 12) {
                    drawPercentText = '' + drawPercent + '%'
                }

                moves.push({
                    openingLink: true,
                    move: element.san,
                    total: total,
                    whitePercent: whitePercent,
                    blackPercent: blackPercent,
                    drawPercent: drawPercent,
                    whitePercentText: whitePercentText,
                    blackPercentText: blackPercentText,
                    drawPercentText: drawPercentText,
                })

                $('#book-moves').html(Mustache.render(nextMovesTpl, moves.slice(0, 12))) // limit to 12 entries
                $('.next-move').bind('click', function(e) {
                    e.preventDefault();
                    if (uiMode == 'opening') {
                        move($(this).html())
                    }
                });
            })
        }
    }).fail(function() {
        console.log('Error connecting to https://explorer.lichess.ovh')
    });

}


// highlight squares 'from' and 'to' on the board
// -> { color: 'w', from: 'e2', to: 'e4', flags: 'b', piece: 'p', san: 'e4' }
function highlightLastMoveOnBoard() {
    var gameHistory = game.history({ verbose: true })
    var lastMove = null
    if (gameHistory.length > 0) {
        lastMove = gameHistory[gameHistory.length - 1]
    }
    $('#myBoard').find('.square-55d63').removeClass('highlight-square')
    if (lastMove) {
        $('#myBoard').find('.square-' + lastMove.from).addClass('highlight-square')
        $('#myBoard').find('.square-' + lastMove.to).addClass('highlight-square')
    }

}

// init
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