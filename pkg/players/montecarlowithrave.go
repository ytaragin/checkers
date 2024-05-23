package players

import (
	"math"

	"github.com/ytaragin/checkers/pkg/board"
	"github.com/ytaragin/checkers/pkg/game"
)

type MCPlayerRave struct {
	Color   board.PieceColor
	Verbose bool
}

func (mcr MCPlayerRave) GetMove(g *game.Game) board.Move {

	moves := g.GetLegalMoves()
	if len(moves) == 1 {
		return moves[0]
	}

	rootNode := mcr.CreateRootNode(&mcr, g)
	iterations := 3000
	explorationWeight := 1.0
	bestMove := rootNode.findBestMove(iterations, explorationWeight)
	// fmt.Printf("Best move: %v\n", bestMove)

	return bestMove
}

func (mcr MCPlayerRave) EvaluateScore(g *game.Game) float64 {
	score := 0.0
	switch g.GetState() {
	case game.RedWin:
		if mcr.Color == board.Red {
			score = 1
		} else {
			score = -1
		}
	case game.BlueWin:
		if mcr.Color == board.Blue {
			score = 1
		} else {
			score = -1
		}
	}
	return score
}

// Node represents a node in the MCTS search tree
type Node struct {
	player     *MCPlayerRave
	move       board.Move
	state      game.Game
	parent     *Node
	children   []*Node
	visits     int
	value      float64
	rave       float64
	raveVisits int
}
type MCMove struct {
}

// CreateRootNode creates the root node for the MCTS search tree
func (mcr MCPlayerRave) CreateRootNode(player *MCPlayerRave, g *game.Game) *Node {
	return &Node{
		state:  *g,
		player: player,
	}
}

// UCT computes the Upper Confidence Bounds (UCB) value for a node
func (node *Node) UCT(totalVisits int, explorationWeight float64) float64 {
	if node.visits == 0 {
		return math.Inf(1)
	}
	return node.value/float64(node.visits) + explorationWeight*math.Sqrt(2*math.Log(float64(totalVisits))/float64(node.visits))
}

// RAVE computes the Rapid Action Value Estimation (RAVE) value for a node
func (node *Node) RAVE(explorationWeight float64) float64 {
	if node.raveVisits == 0 {
		return node.value / float64(node.visits)
	}
	beta := explorationWeight / float64(1+node.raveVisits)
	return (1-beta)*node.value/float64(node.visits) + beta*node.rave/float64(node.raveVisits)
}

// Select selects the best child node based on UCT and RAVE values
func (node *Node) Select(totalVisits int, explorationWeight float64) *Node {
	var bestChild *Node
	bestValue := -math.Inf(1)

	for _, child := range node.children {
		value := child.UCT(totalVisits, explorationWeight)
		if value > bestValue {
			bestValue = value
			bestChild = child
		}
	}

	if bestChild == nil {
		panic("No child nodes found")
	}

	return bestChild
}

// Expand expands the current node by creating child nodes for each possible move
func (node *Node) expand() {
	possibleMoves := node.state.GetLegalMoves()
	node.children = make([]*Node, len(possibleMoves))

	for i, move := range possibleMoves {
		childState := node.state
		childState.RunMove(move)
		node.children[i] = &Node{
			state:  childState,
			player: node.player,
			move:   move,
			parent: node,
		}
	}
}

// BackPropagateValue backpropagates the value of a terminal state up the tree
func (node *Node) backPropagateValue(value float64) {
	for n := node; n != nil; n = n.parent {
		n.visits++
		n.value += value
		n.rave = (n.rave*float64(n.raveVisits) + value) / float64(n.raveVisits+1)
		n.raveVisits++
	}
}

// Simulate performs a Monte Carlo simulation from the current node
func (node *Node) simulate(explorationWeight float64) float64 {
	currentState := node.state
	totalVisits := node.visits

	for currentState.GetState() == game.Ongoing {
		if len(node.children) == 0 {
			node.expand()
		}
		node = node.Select(totalVisits, explorationWeight)
		currentState = node.state
		totalVisits += node.visits
	}

	value := node.player.EvaluateScore(&node.state)
	node.backPropagateValue(value)
	return value
}

// FindBestMove performs the MCTS search and returns the best move
func (rootNode *Node) findBestMove(iterations int, explorationWeight float64) board.Move {
	for i := 0; i < iterations; i++ {
		rootNode.simulate(explorationWeight)
	}

	var bestChild *Node
	bestVisits := -5

	for _, child := range rootNode.children {
		if child.visits > bestVisits {
			bestVisits = child.visits
			bestChild = child
		}
	}

	if bestChild == nil {
		panic("No best move found")
	}

	return bestChild.move
}
