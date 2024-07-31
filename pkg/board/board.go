package board

import (
	"fmt"
)

const BoardRows int = 8
const BoardCols int = 8
const RedDirection int = 1
const BlueDirection int = -1

var AllPositionList []*Position
var AllMovesList [64]Moves

// Board represents the checkers board
type Board struct {
	// Grid [BoardRows][BoardCols]Spot
	RedMask  uint64
	BlueMask uint64
	Kings    uint64
	Invalid  uint64
	// AllMoves     [64]Moves
	// AllPositions []*Position
}

// NewBoard creates and initializes a new checkers board
func NewBoard() *Board {
	board := &Board{}
	board.initializeBoard()
	return board
}

func NewEmptyBoard() *Board {
	return &Board{}
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

	for _, pos := range AllPositionList {
		if b.isSpotEmpty(pos) {
			continue
		}
		if !isBitSetByPos(*b.getColorMask(color), pos) {
			continue
		}
		m, jumps := b.GetMovesForPosition(pos, color, onlyJumps)
		if len(m) > 0 {
			if jumps == onlyJumps {
				moves = append(moves, m...)
			} else if jumps {
				moves = m
				onlyJumps = true
			}
		}
	}
	return moves
}

func (b *Board) getColorMask(color PieceColor) *uint64 {
	if color == Blue {
		return &b.BlueMask
	}
	return &b.RedMask
}

func (b *Board) GetMovesForPosition(pos *Position, color PieceColor, onlyJumps bool) (MoveList, bool) {

	p := b.GetPiece(pos)
	//
	// if p == nil || p.Color != color {
	// 	return MoveList{}, false
	// }

	mvs := &AllMovesList[getMaskSpot(pos.Row, pos.Col)]
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
	// var validJumps = make(MoveList, 0, 2)
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
	// var validPlainMoves = make(MoveList, 0, 2)
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
	return b.getPieceByMask(pos.mask)
}

func (b *Board) KingMe(pos *Position) bool {
	setBitByPos(&b.Kings, pos)
	return true
}

func (b *Board) SetPiece(pos *Position, piece *Piece) bool {
	mask := b.getColorMask(piece.Color)

	setBitByPos(mask, pos)
	if piece.IsKing {
		setBitByPos(&b.Kings, pos)
	}

	return true
}

func (b *Board) getPieceByMask(mask uint64) *Piece {
	var p *Piece
	if isBitSetByMask(b.RedMask, mask) {
		if isBitSetByMask(b.Kings, mask) {
			p = RedKingPiece
		} else {
			p = RedNormalPiece
		}
	} else if isBitSetByMask(b.BlueMask, mask) {
		if isBitSetByMask(b.Kings, mask) {
			p = BlueKingPiece
		} else {
			p = BlueNormalPiece
		}
	}
	return p
}

func (b *Board) removePiece(pos *Position) *Piece {
	var p *Piece
	if isBitSetByPos(b.RedMask, pos) {
		if isBitSetByPos(b.Kings, pos) {
			p = RedKingPiece
		} else {
			p = RedNormalPiece
		}
		clearBitByPos(&b.RedMask, pos)
	} else if isBitSetByPos(b.BlueMask, pos) {
		if isBitSetByPos(b.Kings, pos) {
			p = BlueKingPiece
		} else {
			p = BlueNormalPiece
		}
		clearBitByPos(&b.BlueMask, pos)
	}

	clearBitByPos(&b.Kings, pos)
	return p
}

func (b *Board) movePiece(start *Position, end *Position) *Piece {
	// fmt.Printf("Moving %+v to %+v\n", start, end)
	p := b.removePiece(start)
	if p != nil {
		// fmt.Printf("Setting %+v to %+v\n", p, end)
		b.SetPiece(end, p)
	}

	return p
}

func (b *Board) isSpotEmpty(pos *Position) bool {
	// posSetred := isBitSetByPos(b.Red, pos)
	// posSetblue := isBitSetByPos(b.Blue, pos)
	// return posSetred && posSetblue
	return !isBitSetByPos(b.RedMask, pos) && !isBitSetByPos(b.BlueMask, pos)
}

func (b *Board) addPieces(row_start, row_end int, piece *Piece) {

	for row := row_start; row <= row_end; row++ {
		j_start := 0
		if row%2 == 0 {
			j_start = 1
		}
		for col := j_start; col < BoardCols; col += 2 {
			setBit(b.getColorMask(piece.Color), row, col)
		}
	}

}

// initializeBoard initializes the checkers board with the proper checkers setup
func (b *Board) initializeBoard() {

	AllPositionList = make([]*Position, 0, 32)

	// Initialize spots
	for row := 0; row < BoardRows; row++ {
		for col := 0; col < BoardCols; col++ {
			if (row+col)%2 != 0 {
				AllMovesList[getMaskSpot(row, col)] = b.createMoves(row, col)
				AllPositionList = append(AllPositionList, NewPosition(row, col))
			} else {
				setBitByMask(&b.Invalid, getMaskForPosition(row, col))
			}
		}
	}

	// mvs := &b.Grid[5][2]
	// fmt.Printf("Moves are: %+v\n", mvs)

	b.addPieces(0, 2, RedNormalPiece)
	b.addPieces(BoardRows-3, BoardRows-1, BlueNormalPiece)
}

func (b *Board) Dump(start *Position, last *Position) {

	fmt.Printf("  ")
	for i := 0; i < BoardRows; i++ {
		fmt.Printf("%d  ", i)
	}
	fmt.Println()
	for i := 0; i < BoardRows; i++ {
		fmt.Printf("%d ", i)
		for j := 0; j < BoardCols; j++ {
			mask := getMaskForPosition(i, j)
			p := b.getPieceByMask(mask)
			if isBitSetByMask(b.Invalid, mask) {
				fmt.Print("   ")
			} else if p == nil {
				fmt.Print("- ")
				printMark(i, j, start)
			} else {
				p.Dump() // Call Piece.Dump() method
				printMark(i, j, last)
			}
		}
		fmt.Println()

	}
}

func printMark(row, col int, pos *Position) {
	if pos != nil && pos.Row == row && pos.Col == col {
		fmt.Print("<")
	} else {
		fmt.Print(" ")
	}
}
