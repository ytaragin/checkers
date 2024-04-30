package board

import "fmt"

const BoardRows int = 8
const BoardCols int = 8
const RedDirection int = 1
const BlueDirection int = -1


// Board represents the checkers board
type Board struct {
	Grid [BoardRows][BoardCols]Spot
}

// NewBoard creates and initializes a new checkers board
func NewBoard() *Board {
	board := &Board{}
	board.initializeBoard()
	return board
}

func (b *Board) createMoves(row, col int) {
// Create a new Moves instance
	moves := Moves{
		Moves:  make(map[PieceColor][]PlainMove),
		Jumps:  make(map[PieceColor][]JumpMove),
		KingMoves:   []PlainMove{}, // Initialize as nil if king moves are handled differently
		KingJumps:   []JumpMove{},  // Initialize as nil if king jumps are handled differently
	}

	// Define some sample moves for Red and Blue pieces
	redJumps:= createJumps(row, col, RedDirection, []JumpMove{ })
	blueJumps := createJumps(row, col, BlueDirection, []JumpMove{ })
    
	// Add moves to the Moves struct
	moves.Jumps[Red] = redJumps
	moves.Jumps[Blue] =blueJumps 

	redMoves := createMovesInRow(row, col, RedDirection, []PlainMove{ })
	blueMoves := createMovesInRow(row, col, BlueDirection, []PlainMove{ })
    
	// Add moves to the Moves struct
	moves.Moves[Red] = redMoves
	moves.Moves[Blue] = blueMoves

	moves.KingMoves = createKingMoves(row, col, moves.KingMoves)
    moves.KingJumps = createKingJumps(row, col, moves.KingJumps)


    // fmt.Printf("Row: %d Col: %d\n", row, col)
    // fmt.Printf("%+v\n", moves)

}

func (b *Board) getMovesForPosition(pos *Position, color PieceColor) *Move[]{
    p := b.Grid[pos.Row][pos.Col].Piece
    if (p == nil) {
        return &Move[]
    }

    mvs := b.Grid[pos.Row][pos.Col].PossibleMoves


   var jumps Move[]
   if p.IsKing {
       jumps = mvs.KingJumps
   } else {
       jumps = mvs.Jumps[p.Color]
   }



}

func GetValidMoves(piece Piece, moves Moves) []Move {
    var validMoves []Move
    var jumps, plainMoves []Move

    isKing := piece.IsKing
    color := piece.Color

    if isKing {
        jumps = filterJumps(moves.KingJumps, piece.Color)
    } else {
        jumps = filterJumps(moves.Jumps[color], piece.Color)
    }

    if len(jumps) > 0 {
        validMoves = jumps
    } else {
        if isKing {
            plainMoves = moves.KingMoves
        } else {
            plainMoves = moves.Moves[color]
        }
        validMoves = filterPlainMoves(plainMoves, piece.Color)
    }

    return validMoves
}

func filterJumps(jumps []JumpMove, color PieceColor) []Move {
    var validJumps []Move
    for _, jump := range jumps {
        if jump.IsValid(color) {
            validJumps = append(validJumps, jump)
        }
    }
    return validJumps
}

func filterPlainMoves(plainMoves []PlainMove, color PieceColor) []Move {
    var validPlainMoves []Move
    for _, move := range plainMoves {
        if move.IsValid(color) {
            validPlainMoves = append(validPlainMoves, move)
        }
    }
    return validPlainMoves
}

func (b *Board) getPiece(pos *Position) *Piece {
    return b.Grid[pos.Row][pos.Col].Piece
}
 
func (b *Board) setPiece(pos *Position, piece *Piece) bool  {
    f := b.Grid[pos.Row][pos.Col].Piece
    b.Grid[pos.Row][pos.Col].Piece = piece
    return f == nil
}

func (b *Board) removePiece(pos *Position) *Piece {
    f := b.Grid[pos.Row][pos.Col].Piece
    b.Grid[pos.Row][pos.Col].Piece =nil 
    return f
}
 
func (b *Board) movePiece(start *Position, end *Position) *Piece {
    // fmt.Printf("Moving %+v to %+v\n", start, end)
    p := b.removePiece(start)
    if (p != nil ) {
        // fmt.Printf("Setting %+v to %+v\n", p, end)
        b.setPiece(end, p)
    }

    return p
}

func (b *Board) isSpotEmpty(spot *Position) bool {
    return b.Grid[spot.Row][spot.Col].Piece==nil
}
 
func (b *Board) addPieces(row_start, row_end  int, piece *Piece) {

	for row := row_start; row <= row_end; row++ {
        j_start := 0
        if row%2 == 0 {
            j_start = 1
        }
		for col := j_start; col < BoardCols; col += 2 {
				b.Grid[row][col].Piece = piece
		}
	}


}

// initializeBoard initializes the checkers board with the proper checkers setup
func (b *Board) initializeBoard() {
	// Initialize spots
	for row := 0; row < 8; row++ {
		for col := 0; col < 8; col++ {
			if (row+col)%2 == 0 {
				b.Grid[row][col] = Spot{
					State: Invalid,
					Piece: nil,
					PossibleMoves: Moves{},
				}
			} else {
				b.Grid[row][col] = Spot{
					State: Valid,
					Piece: nil,
                    PossibleMoves: Moves{
						Moves: make(map[PieceColor][]PlainMove),
					},
				}
                b.createMoves(row, col)
			}
		}
	}


    b.addPieces(0, 2, RedNormalPiece);
    b.addPieces(BoardRows-3, BoardRows-1, BlueNormalPiece);
}

func (b *Board) Dump() {
	for i := 0; i < 8; i++ {
		for j := 0; j < 8; j++ {
			if b.Grid[i][j].State == Invalid {
				fmt.Print("   ")
			} else if b.Grid[i][j].Piece == nil {
				fmt.Print("-  ")
			} else {
				b.Grid[i][j].Piece.Dump() // Call Piece.Dump() method
				fmt.Print(" ")
			}
		}
		fmt.Println()


        
	}
}

