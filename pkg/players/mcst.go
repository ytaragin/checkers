package players

import (
	"fmt"
	"math"
	"math/rand"

	"github.com/ytaragin/checkers/pkg/board"
	"github.com/ytaragin/checkers/pkg/game"
)

type MCSTPlayer struct {
	Color   board.PieceColor
	Verbose bool
}

func (mc MCSTPlayer) GetMove(g *game.Game) board.Move {

	moves := g.GetLegalMoves()
	if len(moves) == 1 {
		return moves[0]
	}
	iterations := 50000
	bestMove := mc.GetBestMove(g, iterations)

	return bestMove
}

func (mc MCSTPlayer) GetBestMove(g *game.Game, iterations int) board.Move {
	rootNode := &MCSTNode{
		State: g,
		// Player: node.player,
		Move:     nil,
		Parent:   nil,
		Children: nil,
	}

	for i := 0; i < iterations; i++ {
		rootNode.RunLoop()
	}

	bestChild := rootNode.Children[0]
	for _, child := range rootNode.Children {

		if mc.Verbose {
			fmt.Printf("%s %.2f %d\n", child.Move, child.WinCount, child.VisitCount)
		}
		if child.VisitCount > bestChild.VisitCount {
			bestChild = child
		}
	}

	return bestChild.Move
}

type MCSTNode struct {
	State      *game.Game
	Move       board.Move
	VisitCount int
	WinCount   float64
	Children   []*MCSTNode
	Parent     *MCSTNode
}

func (node *MCSTNode) RunLoop() {
	if node.State.GetState() != game.Ongoing {
		node.BackPropagate(node.State.GetState())
		return
	}
	child := node.Select()

	endingStatus := child.Simulate()
	child.BackPropagate(endingStatus)
}

func (node *MCSTNode) Simulate() game.GameState {
	tempGame := node.State.Copy()
	for tempGame.GetState() == game.Ongoing {
		moves := tempGame.GetLegalMoves()
		randomIndex := rand.Intn(len(moves))
		randomMove := moves[randomIndex]
		tempGame.RunMove(randomMove)
	}

	return tempGame.GetState()
}

func (node *MCSTNode) Expand() bool {
	if node.Children != nil {
		return false
	}
	possibleMoves := node.State.GetLegalMoves()
	node.Children = make([]*MCSTNode, len(possibleMoves))

	for i, move := range possibleMoves {
		childState := node.State.Copy()
		childState.RunMove(move)
		node.Children[i] = &MCSTNode{
			State: childState,
			// Player: node.player,
			Move:     move,
			Parent:   node,
			Children: nil,
		}
	}
	return true
}

func (node *MCSTNode) Select() *MCSTNode {

	if node.State.GetState() != game.Ongoing || node.VisitCount == 0 {
		return node
	}
	expanded := node.Expand()

	var bestChild *MCSTNode
	maxUCB := math.Inf(-1)
	for _, child := range node.Children {
		ucbVal := ucb(child, node.VisitCount)
		if ucbVal >= maxUCB {
			maxUCB = ucbVal
			bestChild = child
		}
	}

	if bestChild == nil {
		panic("No child nodes found")
	}

	if expanded {
		return bestChild
	} else {

		return bestChild.Select()
	}
}

func (node *MCSTNode) BackPropagate(endState game.GameState) {
	n := node

	for n != nil {
		n.VisitCount++
		switch endState {
		case game.RedWin:
			if n.State.NextTurn() == board.Blue {
				n.WinCount++
			}
		case game.BlueWin:
			if n.State.NextTurn() == board.Red {
				n.WinCount++
			}
		case game.Draw:
			n.WinCount += 0.5
		}
		n = n.Parent
	}

}

func ucb(node *MCSTNode, parentVisits int) float64 {
	if node.VisitCount == 0 {
		return math.Inf(1)
	}

	nodeVisits := float64(node.VisitCount)
	part1 := node.WinCount / nodeVisits

	part2 := math.Sqrt(2.0 * math.Log(float64(parentVisits)) / nodeVisits)

	return part1 + part2
}
