const NUMBERSON = false;

function test() {
    ajx('con', 'services.go?cmd=lol')
}

function gid(d) {
    return document.getElementById(d);
}

function ajx(url, callback) {
    var xhr = new XMLHttpRequest()
    xhr.open("GET", "http://localhost:3333/"+url)
    xhr.onreadystatechange = function() {
        if (xhr.readyState == 4) {
            var res = JSON.parse(xhr.response);
            // if there is a callback we use it with the response we get 
            if (res.err) {
                alert(res.err);
                return;
            } 
            if (callback) { 
                callback(res.board);
            }
            console.log(res);
        }
    }
    xhr.send()
}

function getcell(x, y) {
    return document.elementFromPoint(x, y)
}

var down = null;

function initboard() {
    var boardcontainer = gid('boardcontainer');
    var board = document.createElement('table');

    board.addEventListener("mouseup", function(e) {
        var upcell = getcell(e.clientX, e.clientY)
        if (!downcell && !upcell) return;
        ajx("services?up="+upcell.id+"&down="+downcell.id, function(newboard) {
            for (var i = 0; i < 8; i++) {
                for (var j = 0; j < 8; j++) {
                    var piece = newboard[i][j];
                    var id = j + '-' + i;
                    if (piece.kind != "") {
                        var piececolor = "black"
                        if (piece.white) piececolor = "white"
                        gid(id).innerHTML = piececolor + " " + piece.kind;
                    } else {
                        gid(id).innerHTML = "";
                    }
                }
            }
        });
        downcell = null;
    })

    board.addEventListener("mousedown", function(e) {
        var cell = getcell(e.clientX, e.clientY)
        downcell = cell;
    })


    ajx("getboard", function(boardcells) {
        document.board = boardcells
        for (var i = 0; i < 8; i++) {
            var row = document.createElement('tr');
            for (var j = 0; j < 8; j++) {
                var cell = document.createElement('th');
                var id = j + '-' + i;
                if (NUMBERSON) { // for debugging purposes
                    cell.innerHTML = id; 
                } else {
                    var piece = boardcells[i][j];
                    if (piece.kind != "") {
                        var piececolor = "black"
                        if (piece.white) piececolor = "white"
                        cell.innerHTML = piececolor + " " + piece.kind;
                    }
                }

                cell.style.width = '100px';
                cell.style.height = '100px';

                var background = "white";
                var color = "black";

                if ((i + j) % 2 == 1) {
                    background = "black"
                    color = "white";
                }
                cell.style.background = background;
                cell.style.color = color;
                cell.classList.add('cell');
                cell.id = id;
                row.appendChild(cell);
            }
            board.appendChild(row);
        }
        boardcontainer.appendChild(board);
    })
}
