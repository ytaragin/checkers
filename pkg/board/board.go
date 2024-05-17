package board

import (
	"fmt"
)

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

func (b *Board) createMoves(row, col int) Moves {
	// Create a new Moves instance
	moves := Moves{
		Moves:     make(map[PieceColor][]PlainMove),
		Jumps:     make(map[PieceColor][]JumpMove),
		KingMoves: []PlainMove{},
		KingJumps: []JumpMove{},
	}

	redJumps := createJumps(row, col, RedDirection, []JumpMove{})
	blueJumps := createJumps(row, col, BlueDirection, []JumpMove{})
	moves.Jumps[Red] = redJumps
	moves.Jumps[Blue] = blueJumps

	redMoves := createMovesInRow(row, col, RedDirection, []PlainMove{})
	blueMoves := createMovesInRow(row, col, BlueDirection, []PlainMove{})
	moves.Moves[Red] = redMoves
	moves.Moves[Blue] = blueMoves

	moves.KingMoves = createKingMoves(row, col, moves.KingMoves)
	moves.KingJumps = createKingJumps(row, col, moves.KingJumps)

	return moves
}

func (b *Board) GetAllLegalMovesForColor(color PieceColor) MoveList {
	moves := make(MoveList, 0, 3)
	onlyJumps := false

	for row := 0; row < BoardRows; row++ {
		for col := 0; col < BoardCols; col++ {
			m, jumps := b.GetMovesForPosition(&Position{row, col}, color, onlyJumps)
			if len(m) > 0 {
				if jumps == onlyJumps {
					moves = append(moves, m...)
				} else if jumps {
					moves = m
					onlyJumps = true
				}
			}
		}
	}
	return moves
}

func (b *Board) GetMovesForPosition(pos *Position, color PieceColor, onlyJumps bool) (MoveList, bool) {
	p := b.Grid[pos.Row][pos.Col].Piece
	if p == nil {
		return MoveList{}, false
	}
	if p.Color != color {
		return MoveList{}, false
	}

	mvs := &b.Grid[pos.Row][pos.Col].PossibleMoves
	// fmt.Printf("Moves are: %+v\n", mvs)

	return b.getValidMoves(p, mvs, onlyJumps)
}

func (b *Board) getValidMoves(piece *Piece, moves *Moves, onlyJumps bool) (MoveList, bool) {
	var validMoves MoveList
	// var jumps []Move

	// isKing := piece.IsKing
	// color := piece.Color
	//
	// if isKing {
	// 	jumps = b.filterJumps(moves.KingJumps, piece.Color)
	// } else {
	// 	jumps = b.filterJumps(moves.Jumps[color], piece.Color)
	// }
	jumps := b.getJumpMoves(piece, moves)

	if len(jumps) > 0 || onlyJumps {
		return jumps, true
	}

	if piece.IsKing {
		validMoves = b.filterPlainMoves(moves.KingMoves, piece.Color)
	} else {
		validMoves = b.filterPlainMoves(moves.Moves[piece.Color], piece.Color)
	}

	// fmt.Printf("Moves: %+v\n", validMoves)
	return validMoves, false
}

func (b *Board) getJumpMoves(piece *Piece, moves *Moves) MoveList {
	var jumps MoveList

	if piece.IsKing {
		jumps = b.filterJumps(moves.KingJumps, piece.Color)
	} else {
		jumps = b.filterJumps(moves.Jumps[piece.Color], piece.Color)
	}

	ret_jumps := MoveList{}

	for _, jump := range jumps {
		b2 := *b
		jump.DoMove(&b2)
		jumps2, _ := b2.GetMovesForPosition(jump.GetEnd(), piece.Color, true)
		if len(jumps2) > 0 {
			for _, j2 := range jumps2 {
				mjm := NewMultiMove()
				mjm.AddMove(jump)
				mjm.AddMove(j2)
				ret_jumps = append(ret_jumps, mjm)
			}
		} else {
			ret_jumps = append(ret_jumps, jump)
		}
	}

	return ret_jumps
}

func (b *Board) filterJumps(jumps []JumpMove, color PieceColor) MoveList {
	var validJumps MoveList
	// for _, jump := range jumps {
	for i := 0; i < len(jumps); i++ {
		if jumps[i].IsValid(b, color) {
			validJumps = append(validJumps, &jumps[i])
		}
	}
	return validJumps
}

func (b *Board) filterPlainMoves(plainMoves []PlainMove, color PieceColor) MoveList {
	// fmt.Printf("Filtering %+v\n", plainMoves)
	var validPlainMoves MoveList
	// var validPlainMoves []Move
	for i := 0; i < len(plainMoves); i++ {
		// for _, move := range plainMoves {
		if plainMoves[i].IsValid(b, color) {
			validPlainMoves = append(validPlainMoves, &plainMoves[i])
		}
	}
	return validPlainMoves
}

func (b *Board) GetPiece(pos *Position) *Piece {
	return b.Grid[pos.Row][pos.Col].Piece
}

func (b *Board) KingMe(pos *Position) bool {
	if b.Grid[pos.Row][pos.Col].Piece.IsKing {
		return false
	}
	switch b.Grid[pos.Row][pos.Col].Piece.Color {
	case Red:
		b.Grid[pos.Row][pos.Col].Piece = RedKingPiece
	case Blue:
		b.Grid[pos.Row][pos.Col].Piece = BlueKingPiece
	}

	return true
}

func (b *Board) setPiece(pos *Position, piece *Piece) bool {
	f := b.Grid[pos.Row][pos.Col].Piece
	b.Grid[pos.Row][pos.Col].Piece = piece
	return f == nil
}

func (b *Board) removePiece(pos *Position) *Piece {
	f := b.Grid[pos.Row][pos.Col].Piece
	b.Grid[pos.Row][pos.Col].Piece = nil
	return f
}

func (b *Board) movePiece(start *Position, end *Position) *Piece {
	// fmt.Printf("Moving %+v to %+v\n", start, end)
	p := b.removePiece(start)
	if p != nil {
		// fmt.Printf("Setting %+v to %+v\n", p, end)
		b.setPiece(end, p)
	}

	return p
}

func (b *Board) isSpotEmpty(spot *Position) bool {
	return b.Grid[spot.Row][spot.Col].Piece == nil
}

func (b *Board) addPieces(row_start, row_end int, piece *Piece) {

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
	for row := 0; row < BoardRows; row++ {
		for col := 0; col < BoardCols; col++ {
			if (row+col)%2 == 0 {
				b.Grid[row][col] = Spot{
					State:         Invalid,
					Piece:         nil,
					PossibleMoves: Moves{},
				}
			} else {
				b.Grid[row][col] = Spot{
					State:         Valid,
					Piece:         nil,
					PossibleMoves: b.createMoves(row, col),
				}
			}
		}
	}

	// mvs := &b.Grid[5][2]
	// fmt.Printf("Moves are: %+v\n", mvs)

	b.addPieces(0, 2, RedNormalPiece)
	b.addPieces(BoardRows-3, BoardRows-1, BlueNormalPiece)
}

func (b *Board) Dump() {

	fmt.Printf("  ")
	for i := 0; i < BoardRows; i++ {
		fmt.Printf("%d  ", i)
	}
	fmt.Println()
	for i := 0; i < BoardRows; i++ {
		fmt.Printf("%d ", i)
		for j := 0; j < BoardCols; j++ {
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
