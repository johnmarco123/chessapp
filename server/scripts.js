const NUMBERSON = false;

function test() {
    ajx('con', 'services.go?cmd=lol')
}

function gid(d) {
    return document.getElementById(d);
}

function ajx(container1, container2, url, callback) {
    var xhr = new XMLHttpRequest()
    xhr.open("GET", "http://localhost:3333/"+url)
    xhr.onreadystatechange = function() {
        if (xhr.readyState == 4) {
            var res = JSON.parse(xhr.response);
            // if there is a callback we use it with the response we get 
            if (callback) { 
                callback(res);
            }
            console.log(res);
        }
    }
    xhr.send()
}

function hovercell(i, j) {
    return function() {
        console.log(i + '-' + j);
    }
}


function getcell(x, y) {
    return document.elementFromPoint(x, y)
}

var down = null;

function initboard() {
    var boardcontainer = gid('boardcontainer');
    var board = document.createElement('table');

    board.addEventListener("mouseup", function(e) {
        var cell = getcell(e.clientX, e.clientY)
        var up = cell.innerHTML;
        if (down == "" && up == "") return;
        ajx(up, down, "services?up="+up+"&down="+down);
        down = null;
    })

    board.addEventListener("mousedown", function(e) {
        var cell = getcell(e.clientX, e.clientY)
        down = cell.innerHTML;
    })


    ajx("", "", "getboard", function(boardcells) {
        for (var i = 0; i < 8; i++) {
            var row = document.createElement('tr');
            for (var j = 0; j < 8; j++) {
                var cell = document.createElement('th');
                if (NUMBERSON) { // for debugging purposes
                    var val = j + '-' + i;
                    cell.innerHTML = val; 
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
                cell.id = val;
                row.appendChild(cell);
            }
            board.appendChild(row);
        }
        boardcontainer.appendChild(board);
    })
}
