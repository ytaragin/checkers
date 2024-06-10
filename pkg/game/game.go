package game

import (
	"fmt"

	"github.com/ytaragin/checkers/pkg/board"
)

type GameState int

const (
	Ongoing GameState = iota
	BlueWin
	RedWin
	Draw
)

type Game struct {
	gameboard                 board.Board
	nextTurn                  board.PieceColor
	nextLegalMoves            []board.Move
	countSinceLastInteresting int
	moveCount                 int
	lastMove                  board.Move
}

func NewGame() *Game {
	b := board.NewBoard()
	game := &Game{
		gameboard:                 *b,
		nextTurn:                  board.Red,
		nextLegalMoves:            b.GetAllLegalMovesForColor(board.Red),
		countSinceLastInteresting: 0,
		moveCount:                 0,
	}

	return game
}

func (g *Game) Copy() *Game {
	return &Game{
		gameboard:                 g.gameboard,
		nextTurn:                  g.nextTurn,
		nextLegalMoves:            g.nextLegalMoves,
		countSinceLastInteresting: g.countSinceLastInteresting,
		moveCount:                 g.moveCount,
		lastMove:                  g.lastMove,
	}
}

func (g *Game) GetState() GameState {
	if g.countSinceLastInteresting > 80 {
		return Draw
	}
	if g.isCurrentLosing() {
		if g.nextTurn == board.Red {
			return BlueWin
		} else {
			return RedWin
		}
	}
	return Ongoing
}

func (g *Game) MoveCount() int {
	return g.moveCount
}

func (g *Game) NextTurn() board.PieceColor {
	return g.nextTurn
}

func (g *Game) GetWinner() board.PieceColor {
	return g.nextTurn.NextColor()
}

func (g *Game) RunMove(m board.Move) {
	if !m.IsValid(&g.gameboard, g.nextTurn) {
		return
	}
	pos := m.DoMove(&g.gameboard)
	if pos.Row == 0 || pos.Row == board.BoardRows-1 {
		g.gameboard.KingMe(pos)
	}

	if m.IsInteresting(&g.gameboard, true) {
		g.countSinceLastInteresting = 0
	} else {
		g.countSinceLastInteresting++
	}

	g.nextTurn = g.nextTurn.NextColor()
	g.nextLegalMoves = g.gameboard.GetAllLegalMovesForColor(g.nextTurn)
	g.moveCount++
	g.lastMove = m

}

func (g *Game) GetLegalMoves() []board.Move {
	return g.nextLegalMoves
}

func (g *Game) isCurrentLosing() bool {
	return len(g.nextLegalMoves) == 0
}

func (g *Game) Dump() {
	g.gameboard.Dump(g.lastMove.GetStart(), g.lastMove.GetEnd())
	fmt.Printf("Next: %s Count: %d Int: %d\n", g.nextTurn.Name(), g.moveCount, g.countSinceLastInteresting)

}
