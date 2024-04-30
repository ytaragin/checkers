package game

import (
	"fmt"

	"github.com/ytaragin/checkers/pkg/board"
)
type Game struct {
	gameboard board.Board
    nextTurn board.PieceColor
}


func NewGame() *Game{
    game := &Game{
        gameboard:  *board.NewBoard(),
        nextTurn: board.Red,
    }


    return game
}


func (g *Game) RunMove(m board.Move) {
    if m.IsValid(&g.gameboard, g.nextTurn) {
        m.DoMove(&g.gameboard)
        g.nextTurn = g.nextTurn.NextColor()
    }
    
}

func (g *Game) Dump() {
    g.gameboard.Dump();
    fmt.Printf("Next Turn: %s\n", g.nextTurn.Name());

}
