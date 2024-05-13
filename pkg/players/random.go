package players

import (
	"math/rand"

	"github.com/ytaragin/checkers/pkg/board"
	"github.com/ytaragin/checkers/pkg/game"
)

type RandomPlayer struct {
	Color board.PieceColor
}

func (r RandomPlayer) GetMove(g *game.Game) board.Move {
	moves := g.GetLegalMoves()

	randomIndex := rand.Intn(len(moves))

	randomMove := moves[randomIndex]

	return randomMove
}
