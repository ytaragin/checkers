package board

import (
	"fmt"
	"strings"
)

// import "fmt"

type MoveList []Move

type Move interface {
	DoMove(b *Board) *Position
	IsValid(b *Board, color PieceColor) bool
	IsInteresting(b *Board, afterRun bool) bool
	GetStart() *Position
	GetEnd() *Position
	GetJumpedPositions() []*Position
}

type PlainMove struct {
	Start, End Position
}

func CreatePlainMove(startrow, startcol, endrow, endcol int) PlainMove {
	return PlainMove{
		Start: *NewPosition(startrow, startcol),
		End:   *NewPosition(endrow, endcol),
	}

}

func (pm PlainMove) String() string {
	return fmt.Sprintf("{%s>%s}", pm.Start, pm.End)
}

func (pm PlainMove) GetStart() *Position {
	return &pm.Start
}

func (pm PlainMove) GetEnd() *Position {
	return &pm.End
}

func (pm PlainMove) GetJumpedPositions() []*Position {
	return []*Position{}
}

func (pm PlainMove) IsValid(b *Board, color PieceColor) bool {
	p := b.GetPiece(&pm.Start)
	// empty := b.isSpotEmpty(&pm.End)

	return b.isSpotEmpty(&pm.End) &&
		// return empty &&
		(p != nil && p.Color == color)
}

func (pm PlainMove) DoMove(b *Board) *Position {
	// Implementation for PlainMove's DoMove
	// fmt.Println("Performing a plain move")

	b.movePiece(&pm.Start, &pm.End)
	return &pm.End
}

func (pm PlainMove) IsInteresting(b *Board, afterRun bool) bool {
	var p *Piece
	if afterRun {
		p = b.GetPiece(&pm.End)
	} else {
		p = b.GetPiece(&pm.Start)
	}

	return p != nil && !p.IsKing
}

// JumpMove represents a possible move for a piece
type JumpMove struct {
	Jump     Position
	Move     PlainMove
	jumpmask uint64
}

func CreateJump(row, col, skiprow, skipcol, newrow, newcol int) *JumpMove {
	m := CreateMove(row, col, newrow, newcol)
	if m != nil {
		j := &JumpMove{
			Move:     *m,
			Jump:     *NewPosition(skiprow, skipcol),
			jumpmask: getMaskForPosition(skiprow, skipcol),
		}
		return j
	}
	return nil
}

func (j JumpMove) String() string {
	return fmt.Sprintf("{%s[%s]}", j.Move, j.Jump)
}

func (j JumpMove) GetStart() *Position {
	return &j.Move.Start
}

func (j JumpMove) GetEnd() *Position {
	return &j.Move.End
}

func (j JumpMove) GetJumpedPositions() []*Position {
	return []*Position{&j.Jump}
}

func (j JumpMove) IsValid(b *Board, color PieceColor) bool {

	p := b.GetPiece(&j.Move.Start)
	jspot := b.GetPiece(&j.Jump)
	return b.isSpotEmpty(&j.Move.End) &&
		(p != nil && p.Color == color) &&
		(jspot != nil && jspot.Color != color)
}

func (j JumpMove) DoMove(b *Board) *Position {
	b.movePiece(&j.Move.Start, &j.Move.End)
	b.removePiece(&j.Jump)

	return &j.Move.End
}

func (j JumpMove) IsInteresting(b *Board, afterRun bool) bool {
	return true
}

type MultiMove struct {
	Moves MoveList
}

func NewMultiMove() *MultiMove {
	return &MultiMove{
		Moves: MoveList{},
	}
}

func (mjm MultiMove) String() string {
	var sb strings.Builder
	sb.WriteString("{")
	for i, m := range mjm.Moves {
		sb.WriteString(fmt.Sprintf("%s", m))
		if i < len(mjm.Moves)-1 {
			sb.WriteString("/")
		}
	}
	sb.WriteString("}")
	return sb.String()
}

func (mjm *MultiMove) AddMove(move Move) {
	mjm.Moves = append(mjm.Moves, move)
}

func (mjm MultiMove) GetStart() *Position {
	return mjm.Moves[0].GetStart()
}

func (mjm MultiMove) GetEnd() *Position {
	return mjm.Moves[len(mjm.Moves)-1].GetEnd()
}

func (mjm MultiMove) GetJumpedPositions() []*Position {
	var allJumpedPositions []*Position

	for _, move := range mjm.Moves {
		jumpedPositions := move.GetJumpedPositions()
		allJumpedPositions = append(allJumpedPositions, jumpedPositions...)
	}

	return allJumpedPositions
}

func (mjm MultiMove) IsValid(b *Board, color PieceColor) bool {
	valid := true
	// for _, j := range mjm.Moves {
	// 	if !j.IsValid(b, color) {
	// 		valid = false
	// 		break
	// 	}
	// }
	return valid
}

func (mjm MultiMove) DoMove(b *Board) *Position {
	end := (*Position)(nil)
	for _, j := range mjm.Moves {
		end = j.DoMove(b)
	}

	return end
}

func (mjm MultiMove) IsInteresting(b *Board, afterRun bool) bool {
	return true
}

func CreateMove(row, col, newrow, newcol int) *PlainMove {
	if newcol >= 0 && newcol < BoardCols && newrow >= 0 && newrow < BoardRows {
		return &PlainMove{
			Start: *NewPosition(row, col),
			End:   *NewPosition(newrow, newcol),
		}
	}

	return nil
}

func createAppendMove(row, col, newrow, newcol int, moves []PlainMove) []PlainMove {
	m := CreateMove(row, col, newrow, newcol)
	if m != nil {
		moves = append(moves, *m)
	}

	return moves
}

func createMovesInRow(row, col, dir int, moves []PlainMove) []PlainMove {
	moves = createAppendMove(row, col, row+dir, col-1, moves)
	moves = createAppendMove(row, col, row+dir, col+1, moves)

	return moves
}

func createKingMoves(row, col int, moves []PlainMove) []PlainMove {
	moves = createMovesInRow(row, col, RedDirection, moves)
	moves = createMovesInRow(row, col, BlueDirection, moves)

	return moves
}
func createAppendJump(row, col, skiprow, skipcol, newrow, newcol int, jumps []JumpMove) []JumpMove {
	j := CreateJump(row, col, skiprow, skipcol, newrow, newcol)
	if j != nil {
		jumps = append(jumps, *j)
	}

	return jumps
}
func createJumps(row, col, dir int, jumps []JumpMove) []JumpMove {
	jumps = createAppendJump(row, col, row+dir, col+1, row+2*dir, col+2, jumps)
	jumps = createAppendJump(row, col, row+dir, col-1, row+2*dir, col-2, jumps)

	return jumps
}

func createKingJumps(row, col int, jumps []JumpMove) []JumpMove {
	jumps = createJumps(row, col, RedDirection, jumps)
	jumps = createJumps(row, col, BlueDirection, jumps)

	return jumps
}
