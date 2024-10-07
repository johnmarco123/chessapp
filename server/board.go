package main

import (
	"fmt";
	"math";
	"strconv";
)

type piece struct {
	Kind string `json:"kind"`;  // the type of piece, e.g., pawn, queen
	Moved bool   `json:"moved"`; // if the piece has Moved at least once
	White bool   `json:"white"`; // true if the piece is White, otherwise black
}

func makepiece(t string, m bool, w bool) piece{
	return piece{Kind: t, Moved: m, White: w};
}


func empty() piece {
	return makepiece("", false, false);
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

// checks if every cell from (x1, y1) to (x2, y2) is empty, except for (x2, y2)
// itself
func piecebetween(x1 int, y1 int, x2 int, y2 int) bool {
	// 1. ensure slope is either abs(1) or 0
	// 2. get the cells between with a loop and check to make sure they are all
	// empty

	fmt.Println("at: ("+strconv.Itoa(x1)+" "+strconv.Itoa(y1)+")");
	var denominator = float64(x2) - float64(x1);
	// no dividing by zero
	if (denominator == 0) {
		denominator = 1;
	}

	var slope float64 = (float64(y2) - float64(y1)) / float64(denominator);
	// if there is a decimal place, then we cannot check if there are pieces
	// between! And therefore its a invalid placement (knights are dealt with
	// somewhere else)
	if (slope != float64(int(slope))) { 
		fmt.Println("piecebetween: BAD SLOPE");
		return true;
	}

	var intslope int = int(slope);
	var yinc int = 0;
	if (intslope > 0) {
		yinc = 1;
	} else if (intslope < 0) {
		yinc = -1;
	}

	var xinc = 0;
	if (x2 > x1) {
		xinc = 1;
	} else if (x2 < x1) {
		xinc = -1;
	}
	
	fmt.Println("==================================================");
	fmt.Println("xinc: " + strconv.Itoa(xinc));
	fmt.Println("yinc: " + strconv.Itoa(yinc));
	for (x1 != x2 || y1 != y2) {
		x1 += xinc;
		y1 += yinc;
		if (board[y1][x1].Kind != "") { // there is a piece between the two places
			fmt.Println(board[y1][x1].Kind);
			return true;
		}
	}

	return false;
}

func legal(p1 piece, p2 piece, x1 int, y1 int, x2 int, y2 int) (string, string) {
	var col = "white";
	if (!p1.White) {
		col = "black";
	}
	if (p1.Kind == "pawn") {
		if (p1.White) {
			if (x1 == x2) { // moving forwards 
				if (y1 - y2 == 1 && board[y2][x2].Kind == "") { // moving one square forward
					return "white pawn moved forward one", "";
				} else if (y1 - y2 == 2 && !p1.Moved && !piecebetween(x1, y1, x2, y2) && board[y2][x2].Kind == "") { // moving two squares forward 
					return "white pawn moved forward twice", "";
				} else {
					return "", "PAWN MOVE ERROR";
				}

			} else if (y1 - y2 == 1 && int(math.Abs(float64(x1 - x2))) == 1 && p2.Kind != "" && !p2.White) { // normal diagonal capture
				return col+" "+p1.Kind+" captured a " + p2.Kind, "";

			} // TODO ELSE IF ENPASSANTABLE 
		} else if (!p1.White){
			if (x1 == x2) { // moving forwards 
				if (y2 - y1 == 1 && board[y2][x2].Kind == "") { // moving one square forward
					return "black pawn moved forward one", "";
					
				} else if (y2 - y1 == 2 && !p1.Moved && !piecebetween(x1, y1, x2, y2) && board[y2][x2].Kind == "") { // moving two squares forward 
					return "black pawn moved forward twice", "";

				} else {
					return "", "PAWN MOVE ERROR";

				}

			} else if (y2 - y1 == 1 && int(math.Abs(float64(x1 - x2))) == 1 && p2.Kind != "" && p2.White) { // normal diagonal capture
				return col+" "+p1.Kind+" captured a " + p2.Kind, "";

			} // TODO ELSE IF ENPASSANTABLE 
		}
	} else if (p1.Kind == "rook") {
		if (x1 != x2 && y1 != y2) { // ensure rook only moves in straight lines
			return "", "ILLEGAL ROOK MOVE" ;
		} else if (!piecebetween(x1, y1, x2, y2)) {
			if (p2.Kind == "") {
				return col+" "+p1.Kind+" moved successfully", "";
			} else if(p2.Kind != "" && p1.White != p2.White){
				return col+" "+p1.Kind+" captured a " + p2.Kind, "";
			}
		}
	} else if (p1.Kind == "bishop") {
		if (x1 == x2 || y1 == y2) { // ensure bishop never moves in straight line
			return "", "ILLEGAL BISHOP MOVE";
		} else if (!piecebetween(x1, y1, x2, y2)) {
			if (p2.Kind == "") {
				return col+" "+p1.Kind+" moved successfully", "";
			} else if(p2.Kind != "" && p1.White != p2.White){
				return col+" "+p1.Kind+" captured a " + p2.Kind, "";
			}
		}
	}
	return "", "GENERIC ILLEGAL MOVE MESSAGE";
}

func boardlogic(down string, up string) ([8][8]piece, string, string){
	if (down == "" || up == "") { // ERROR
		return board, "",  "User did not drag to two valid squares";
	}
	var ux, _ = strconv.Atoi(string(up[0]));
	var uy, _ = strconv.Atoi(string(up[2]));
	var uppiece = board[uy][ux];

	var dx, _ = strconv.Atoi(string(down[0]));
	var dy, _ = strconv.Atoi(string(down[2]));
	var downpiece = board[dy][dx];
	msg, err := legal(downpiece, uppiece, dx, dy, ux, uy);
	if (len(msg) > 0) {
		// the move was legal, so this piece has now Moved at least once
		downpiece.Moved = true;
		board[dy][dx] = empty();
		board[uy][ux] = downpiece;
	}
	return board, msg, err;
}
