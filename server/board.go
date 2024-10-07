package main

import (
	"strconv"
	"math"
)

type piece struct {
	Kind string `json:"kind"`  // the type of piece, e.g., pawn, queen
	Moved bool   `json:"moved"` // if the piece has Moved at least once
	White bool   `json:"white"` // true if the piece is White, otherwise black
}

func makepiece(t string, m bool, w bool) piece{
	return piece{Kind: t, Moved: m, White: w}
}


func empty() piece {
	return makepiece("", false, false)
}

var board = [8][8] piece {
	{
		makepiece("rook", false, false),
		makepiece("knight", false, false),
		makepiece("bishop", false, false),
		makepiece("queen", false, false),
		makepiece("king", false, false),
		makepiece("bishop", false, false),
		makepiece("knight", false, false),
		makepiece("rook", false, false),
	},
	{
		makepiece("pawn", false, false),
		makepiece("pawn", false, false),
		makepiece("pawn", false, false),
		makepiece("pawn", false, false),
		makepiece("pawn", false, false),
		makepiece("pawn", false, false),
		makepiece("pawn", false, false),
		makepiece("pawn", false, false),
	},
	{ empty(), empty(), empty(), empty(), empty(), empty(), empty(), empty() },
	{ empty(), empty(), empty(), empty(), empty(), empty(), empty(), empty() },
	{ empty(), empty(), empty(), empty(), empty(), empty(), empty(), empty() },
	{ empty(), empty(), empty(), empty(), empty(), empty(), empty(), empty() },
	{
		makepiece("pawn", false, true),
		makepiece("pawn", false, true),
		makepiece("pawn", false, true),
		makepiece("pawn", false, true),
		makepiece("pawn", false, true),
		makepiece("pawn", false, true),
		makepiece("pawn", false, true),
		makepiece("pawn", false, true),
	},
	{
		makepiece("rook", false, true),
		makepiece("knight", false, true),
		makepiece("bishop", false, true),
		makepiece("queen", false, true),
		makepiece("king", false, true),
		makepiece("bishop", false, true),
		makepiece("knight", false, true),
		makepiece("rook", false, true),
	},
}

// have an array of enpassantable pawns that gets flushed after each move if it
// contains anything

func legal(p1 piece, p2 piece, x1 int, y1 int, x2 int, y2 int) (bool, string) {
	if (p1.Kind == "pawn") {
		if (p1.White) {
			if (x1 == x2 && p2.Kind == "") { // moving forwards with nothing blocking your way
				if (y1 - y2 == 1) { // moving one square forward
					return true, "pawn move forward one"
					
				} else if (y1 - y2 == 2 && !p1.Moved) { // moving two squares forward 
					return true, "pawn move forward twice (first move only)"

				} else {
					return false, "pawns cannot move this way"

				}

			} else if (y1 - y2 == 1 && int(math.Abs(float64(x1 - x2))) == 1) { // normal diagonal capture
				return true, "captured a pawn diagonally"; 

			} // TODO ELSE IF ENPASSANTABLE 
		}
	}
	return false, "what the fuck";
}

func boardlogic(down string, up string) (string, string, string){
	var upafter string = ""
	var downafter string = ""
	var msg string = "illegal"
	if (down == "" || up == "") { // user didnt grab a piece to begin with
		return upafter, downafter, msg;
	}
	var ux, _ = strconv.Atoi(string(up[0]));
	var uy, _ = strconv.Atoi(string(up[2]));
	var uppiece = board[uy][ux]

	var dx, _ = strconv.Atoi(string(down[0]));
	var dy, _ = strconv.Atoi(string(down[2]));
	var downpiece = board[dy][dx]
	if (legal(downpiece, uppiece, dx, dy, ux, uy)) {
		// the move was legal, so this piece has now Moved at least once
		downpiece.Moved = true; 
		board[dy][dx] = empty(); 
		board[uy][ux] = downpiece;
		return "", strconv.Itoa(ux) + "-" + strconv.Itoa(uy), "legal"
	}
	return "", "", "illegal"
}
