package players

import (
	"github.com/ytaragin/checkers/pkg/board"
	"github.com/ytaragin/checkers/pkg/game"
)

type Player interface {
	GetMove(g *game.Game) board.Move
}
