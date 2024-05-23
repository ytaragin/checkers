package players

import (
	"fmt"
	"time"

	"github.com/ytaragin/checkers/pkg/board"
	"github.com/ytaragin/checkers/pkg/game"
)

type GameRunner struct {
	players map[board.PieceColor]Player
	game    *game.Game
}

func RunMultiple(redPlayer Player, bluePlayer Player, amount int) {
	redWins := 0
	blueWins := 0
	draws := 0

	for i := 0; i < amount; i++ {

		g := game.NewGame()
		// g.Dump()

		runner := RunGame(g, redPlayer, bluePlayer)
		runner.RunTillEnd()

		switch g.GetState() {
		case game.Draw:
			draws++
		case game.RedWin:
			redWins++
		case game.BlueWin:
			blueWins++
		}

		if i%100 == 0 {
			println("+")
		} else if i%10 == 0 {
			print("#")
		}

	}

	fmt.Printf("Red Wins: %d Blue Wins: %d Draws: %d\n", redWins, blueWins, draws)

}

func RunGame(game *game.Game, redPlayer Player, bluePlayer Player) *GameRunner {
	players := make(map[board.PieceColor]Player)

	players[board.Red] = redPlayer
	players[board.Blue] = bluePlayer

	runner := &GameRunner{players, game}
	return runner
}

func (gr *GameRunner) RunTillEnd() {
	start := time.Now()
	for gr.game.GetState() == game.Ongoing {
		m := gr.players[gr.game.NextTurn()].GetMove(gr.game)
		gr.game.RunMove(m)

		i := gr.game.MoveCount()
		if i%10 == 0 {
			fmt.Printf(".")
		}

		// fmt.Println(m)
		// gr.game.Dump()

	}

	elapsed := time.Since(start)
	timePerMove := (float64(elapsed) / 1e6) / float64(gr.game.MoveCount())
	fmt.Printf("Time Per Move: %2f\n", timePerMove)

	// switch gr.game.GetState() {
	// case game.RedWin:
	// 	fmt.Printf("Winner is Red\n")
	// case game.BlueWin:
	// 	fmt.Printf("Winner is Blue\n")
	// case game.Draw:
	// 	fmt.Println("It was a draw")
	//
	// }
}
