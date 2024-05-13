package board

import "fmt"

// PieceColor represents the color of a piece
type PieceColor int

const (
	Red PieceColor = iota
	Blue
)

func (pc PieceColor) Name() string {
	if pc == Red {
		return "Red"
	}
	return "Blue"
}

func (pc PieceColor) NextColor() PieceColor {
	if pc == Red {
		return Blue
	}
	return Red
}

// Piece *represents a game piece
type Piece struct {
	Color  PieceColor
	IsKing bool
}

func (p *Piece) Dump() {
	var pieceString string
	switch p.Color {
	case Red:
		pieceString = "R"
	case Blue:
		pieceString = "B"
	}

	if p.IsKing {
		pieceString += pieceString
	} else {
		pieceString += " "
	}

	fmt.Print(pieceString)
}

var (
	RedNormalPiece  = &Piece{Color: Red, IsKing: false}
	RedKingPiece    = &Piece{Color: Red, IsKing: true}
	BlueNormalPiece = &Piece{Color: Blue, IsKing: false}
	BlueKingPiece   = &Piece{Color: Blue, IsKing: true}
)

type Position struct {
	Row, Col int
}

func (p Position) String() string {
	return fmt.Sprintf("(%d,%d)", p.Row, p.Col)
}

// Moves contains arrays of possible moves for a spot
type Moves struct {
	Moves     map[PieceColor][]PlainMove
	Jumps     map[PieceColor][]JumpMove
	KingMoves []PlainMove
	KingJumps []JumpMove
}

// SpotState is an enum representing the state of a spot
type SpotState int

const (
	Invalid SpotState = iota
	Valid
	// Add more spot states if needed
)

// Spot represents a single spot on the checkers board
type Spot struct {
	State         SpotState
	Piece         *Piece
	PossibleMoves Moves
}
