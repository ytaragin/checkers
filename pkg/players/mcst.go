package players

import (
	"fmt"
	"math"
	"math/rand"
	"time"

	"github.com/ytaragin/checkers/pkg/board"
	"github.com/ytaragin/checkers/pkg/game"
)

type ChildStatSelecter interface {
	SelectStat(child *MCSTNode) float64
	StatName() string
}

type MostVisitsSelector struct{}

func (mv MostVisitsSelector) SelectStat(child *MCSTNode) float64 {
	return float64(child.VisitCount)
}
func (mv MostVisitsSelector) StatName() string {
	return "Visits"
}

type WinRateSelector struct{}

func (mv WinRateSelector) SelectStat(child *MCSTNode) float64 {
	return child.WinCount / float64(child.VisitCount)
}
func (mv WinRateSelector) StatName() string {
	return "WinRate"
}

var (
	WinRate    ChildStatSelecter = WinRateSelector{}
	MostVisits ChildStatSelecter = MostVisitsSelector{}
)

type MCSTPlayer struct {
	Color              board.PieceColor
	SelectionAlgorithm ChildStatSelecter
	Iterations         int
	Duration           time.Duration
	Verbose            bool
}

func (mc MCSTPlayer) GetMove(g *game.Game) board.Move {

	moves := g.GetLegalMoves()
	if len(moves) == 1 {
		return moves[0]
	}

	if mc.Duration == 0 && mc.Iterations == 0 {
		mc.Iterations = 50000
		fmt.Printf("No Iterations or Duration set. Will run %d iterations", mc.Iterations)
	}

	// iterations := 50000
	// bestMove := mc.GetBestMove(g, iterations)
	bestMove := mc.GetBestMove(g)

	return bestMove
}

// func (mc MCSTPlayer) GetBestMove(g *game.Game, iterations int, d time.Duration) board.Move {
func (mc MCSTPlayer) GetBestMove(g *game.Game) board.Move {
	rootNode := &MCSTNode{
		State: g,
		// Player: node.player,
		Move:     nil,
		Parent:   nil,
		Children: nil,
	}
	count := 0

	if mc.Iterations > 0 {
		for i := 0; i < mc.Iterations; i++ {
			rootNode.RunLoop()
		}
		count = mc.Iterations
	} else {
		endTime := time.Now().Add(mc.Duration)
		for time.Now().Before(endTime) {
			rootNode.RunLoop()
			count++
		}
	}

	bestChild := rootNode.Children[0]
	for _, child := range rootNode.Children {

		// if mc.Verbose {
		// 	fmt.Printf("%s %.2f %d\n", child.Move, child.WinCount, child.VisitCount)
		// }
		if mc.SelectionAlgorithm.SelectStat(child) > mc.SelectionAlgorithm.SelectStat(bestChild) {
			bestChild = child
		}
	}
	if mc.Verbose {
		fmt.Printf("Iterations: %d, Visits: %d WinCount: %d %s: %.4f\n",
			count,
			bestChild.VisitCount,
			bestChild.WinCount,
			mc.SelectionAlgorithm.StatName(),
			mc.SelectionAlgorithm.SelectStat(bestChild))
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
			n.WinCount += 0.0
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
