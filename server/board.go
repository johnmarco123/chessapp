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
	if (x1 == x2 && y1 == y2) { // same piece
		return false;
	}

	var xinc = 0;
	var yinc = 0;
	if (x1 == x2) { // forward/backward move
		if (y2 > y1) {
			yinc = 1;
		} else {
			yinc = -1;
		}

		for (y1 != y2) {
			y1 += yinc;
			if (y1 == y2) {
				break;
			} else if (y1 > 7 || x1 < 0) {
				fmt.Println("FATAL ERROR, OUT OF BOUNDS y1");
				return true;
			}
			if (board[y1][x1].Kind != "") {
				fmt.Println("Piece between found! Piece was: " + board[y1][x1].Kind);
				return true;
			}
		}

	} else if (y1 == y2) { // left/right move
		if (x2 > x1) {
			xinc = 1;
		} else {
			xinc = -1;
		}

		for (x1 != x2) {
			x1 += xinc;
			if (x1 == x2) {
				break;
			} else if (x1 > 7 || x1 < 0) {
				fmt.Println("FATAL ERROR, OUT OF BOUNDS x1");
				return true;
			}
			if (board[y1][x1].Kind != "") {
				fmt.Println("Piece between found! Piece was: " + board[y1][x1].Kind);
				return true; 
			}
		}

	} else { // ensure that it is a diagonal, and then check if there is pieces between
		// simply check if the slope is not a decimal, if that is the case,
		// ignore slope 
		var slope float64 = (float64(y2) - float64(y1)) / (float64(x2) - float64(x1))
		if (slope != float64(int(slope))) { // this is not a diagonal line
			return true;
		}

		if (x2 > x1) {
			xinc = 1;
		} else {
			xinc = -1;
		}

		if (y2 > y1) {
			yinc = 1;
		} else {
			yinc = -1;
		}
		for (x2 != x1 && y2 != y1) {
			x1 += xinc;
			y1 += yinc;
			if (x1 == x2 && y2 == y1) {
				break;
			} else if (x1 < 0 || x1 > 7 || y1 < 0 || y1 > 7) {
				fmt.Println("FATAL ERROR, OUT OF BOUNDS x1 OR y1");
			}
			if (board[y1][x1].Kind != "") {
				fmt.Println("Piece between found! Piece was: " + board[y1][x1].Kind);
				return true;
			}
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
	} else if (p1.Kind == "queen") {
		if (!piecebetween(x1, y1, x2, y2)) {
			if (p2.Kind == "") {
				return col+" "+p1.Kind+" moved successfully", "";
			} else if(p2.Kind != "" && p1.White != p2.White){
				return col+" "+p1.Kind+" captured a " + p2.Kind, "";
			}
		}
	} else if (p1.Kind == "king") {
		// cannot move into check
		// can only move one square at a time unless castling
		// can only castle if havent moved and the castle does not put the king in check
		if (math.Abs(float64(x2-x1)) > 1 || math.Abs(float64(y2-y1)) > 1) { // king can only move one square at a time
			return "", "king cannot move more then one square at a time";
		} else if (p2.Kind == "") {
			return col+" "+p1.Kind+" moved successfully", "";
		} else if(p2.Kind != "" && p1.White != p2.White){
			return col+" "+p1.Kind+" captured a " + p2.Kind, "";
		}
	} else if (p1.Kind == "knight") { // this is gonna be fun...
		// knights have 8 valid moves. 2 moves out and then to left one each way
		// for all 4 dirs
		var xknightmoves[8] int; // gather all 8 
		var yknightmoves[8] int; // gather all 8 
		var idx int = 0;
		if (y1 > 1) {
			if (x1 != 0) {
				xknightmoves[idx] = x1-1
				yknightmoves[idx] = y1-2
				idx++;
			}
			if (x1 != 7) {
				xknightmoves[idx] = x1+1
				yknightmoves[idx] = y1-2
				idx++;
			}
		}
		if (y1 < 6) {
			if (x1 > 0) {
				xknightmoves[idx] = x1-1
				yknightmoves[idx] = y1+2
				idx++;
			}
			if (x1 < 7) {
				xknightmoves[idx] = x1+1
				yknightmoves[idx] = y1+2
				idx++;
			}
		}
		if (x1 > 1) {
			if (y1 > 0) {
				xknightmoves[idx] = x1-2
				yknightmoves[idx] = y1-1
				idx++;
			}
			if (x1 != 7) {
				xknightmoves[idx] = x1-2
				yknightmoves[idx] = y1+1
				idx++;
			}
		}
		if (x1 < 6) {
			if (y1 > 0) {
				xknightmoves[idx] = x1+2
				yknightmoves[idx] = y1-1
				idx++;
			}
			if (x1 != 7) {
				xknightmoves[idx] = x1+2
				yknightmoves[idx] = y1+1
				idx++;
			}
		}
		for (idx < 7) {
			// set all invalid knight moves to a impossible value
			xknightmoves[idx] = 1337;
			yknightmoves[idx] = 1337;
			idx++;
		}
		var knightmovevalid bool = false;
		for i := 0; i < len(xknightmoves); i++ {
			if (x2 == xknightmoves[i] && y2 == yknightmoves[i]) {
				knightmovevalid = true;
				break;
				// move is valid!
			}
		}
		if (!knightmovevalid) {
			return "", "The square provided is not a valid knight move"
		}
		// now we got all valid knight moves, so lets check if the x2, y2 move
		// is in here
		if (p2.Kind == "") {
			return col+" "+p1.Kind+" moved successfully", "";
		} else if(p2.Kind != "" && p1.White != p2.White){
			return col+" "+p1.Kind+" captured a " + p2.Kind, "";
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
