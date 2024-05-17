package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"runtime/pprof"

	"github.com/ytaragin/checkers/pkg/board"
	"github.com/ytaragin/checkers/pkg/game"
	"github.com/ytaragin/checkers/pkg/players"
)

var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")
var memprofile = flag.String("memprofile", "", "write memory profile to this file")

func main() {
	flag.Parse()
	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal(err)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}
	runGame()
	// hardcode()

	if *memprofile != "" {
		f, err := os.Create(*memprofile)
		if err != nil {
			log.Fatal(err)
		}
		pprof.WriteHeapProfile(f)
		f.Close()
		return
	}

}

func runGame() {

	// redPlayer := players.RandomPlayer{Color: board.Red}
	redPlayer := players.MCPlayer{Color: board.Red}
	// bluePlayer := players.RandomPlayer{Color: board.Blue}
	bluePlayer := players.MCPlayer{Color: board.Blue}

	// g := game.NewGame()
	// g.Dump()
	// runner := players.RunGame(g, redPlayer, bluePlayer)
	// runner.RunTillEnd()
	// g.Dump()

	players.RunMultiple(redPlayer, bluePlayer, 1)

}

func hardcode() {
	fmt.Println("Hello, World!")

	// Create a new board
	game := game.NewGame()
	game.Dump()

	moves := game.GetLegalMoves()
	fmt.Printf("%+v", moves)

	m := board.CreateMove(2, 1, 3, 2)
	game.RunMove(m)
	game.Dump()

	m = board.CreateMove(5, 4, 4, 3)
	game.RunMove(m)
	game.Dump()

	j := board.CreateJump(3, 2, 4, 3, 5, 4)
	game.RunMove(j)
	game.Dump()

	m = board.CreateMove(5, 6, 4, 7)
	game.RunMove(m)
	game.Dump()

	m = board.CreateMove(2, 5, 3, 4)
	game.RunMove(m)
	game.Dump()

	m = board.CreateMove(6, 5, 5, 6)
	game.RunMove(m)
	game.Dump()

	m = board.CreateMove(5, 4, 6, 5)
	game.RunMove(m)
	game.Dump()

	m = board.CreateMove(5, 2, 4, 1)
	game.RunMove(m)
	game.Dump()

	m = board.CreateMove(3, 4, 4, 3)
	game.RunMove(m)
	game.Dump()

	m = board.CreateMove(4, 1, 3, 0)
	game.RunMove(m)
	game.Dump()

	m = board.CreateMove(1, 4, 2, 5)
	game.RunMove(m)
	game.Dump()

	m = board.CreateMove(6, 3, 5, 2)
	game.RunMove(m)
	game.Dump()

	m = board.CreateMove(1, 0, 2, 1)
	game.RunMove(m)
	game.Dump()

	moves = game.GetLegalMoves()
	fmt.Printf("%s\n\n", moves)

	m1 := moves[2]

	fmt.Printf("%s\n", m1)

	game.RunMove(m1)
	game.Dump()
}
