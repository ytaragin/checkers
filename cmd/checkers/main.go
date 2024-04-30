package main

import (
    "fmt"

    "github.com/ytaragin/checkers/pkg/board"
    "github.com/ytaragin/checkers/pkg/game"
)



func main() {
    fmt.Println("Hello, World!")

    // Create a new board
    game := game.NewGame()

    // Print the initial board state
    game.Dump()


    // m := board.PlainMove{   }
    m := board.CreateMove(2, 1, 3, 2)

    game.RunMove(m)


    game.Dump()

    m = board.CreateMove(5, 4, 4, 3)

    game.RunMove(m)


    game.Dump()


    j := board.CreateJump(3,2,4,3,5,4)
    game.RunMove(j)
    game.Dump()

}

