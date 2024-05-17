package players

import (
	"fmt"
	"math/rand"
	"sync"

	"github.com/ytaragin/checkers/pkg/board"
	"github.com/ytaragin/checkers/pkg/game"
)

type MCPlayer struct {
	Color   board.PieceColor
	Verbose bool
}

func (mc MCPlayer) GetMove(g *game.Game) board.Move {

	iterations := 10000
	bestMove := mc.GetBestMove(g, iterations)

	return bestMove
}

func (mc MCPlayer) GetBestMove(g *game.Game, iterations int) board.Move {

	possibleMoves := g.GetLegalMoves()
	scoreMap := make(map[board.Move]float64)

	for _, move := range possibleMoves {
		// scoreMap[move] = mc.RunMonteCarloMT(g, move, iterations, 3)
		scoreMap[move] = mc.RunMonteCarlo(g, move, iterations)
	}

	var bestMove board.Move
	highestScore := -5.0
	for move, score := range scoreMap {
		if mc.Verbose {
			fmt.Printf("%s %.2f\n", move, score)
		}
		if score > highestScore {
			highestScore = score
			bestMove = move
		}
	}
	return bestMove
}

// EvaluateScore evaluates the score or outcome of the given game state
func (mc MCPlayer) EvaluateScore(g *game.Game) float64 {
	score := 0.0
	switch g.GetState() {
	case game.RedWin:
		if mc.Color == board.Red {
			score = 1
		} else {
			score = -1
		}
	case game.BlueWin:
		if mc.Color == board.Blue {
			score = 1
		} else {
			score = -1
		}
	}
	return score
}

func GetRandomMove(g *game.Game) board.Move {
	moves := g.GetLegalMoves()
	randomIndex := rand.Intn(len(moves))

	randomMove := moves[randomIndex]

	return randomMove
}

// RunMonteCarlo runs the Monte Carlo algorithm for a given move and returns the average score
func (mc MCPlayer) RunMonteCarlo(g *game.Game, move board.Move, iterations int) float64 {
	totalScore := 0.0
	for i := 0; i < iterations; i++ {
		gtemp := *g
		gtemp.RunMove(move)
		for gtemp.GetState() == game.Ongoing {
			gtemp.RunMove(GetRandomMove(&gtemp))
		}
		totalScore += mc.EvaluateScore(&gtemp)
	}
	return totalScore / float64(iterations)
}

func (mc MCPlayer) RunMonteCarloMT(g *game.Game, move board.Move, iterations int, numWorkers int) float64 {
	workerCount := numWorkers
	if workerCount < 1 {
		workerCount = 1
	}

	numIterationsPerWorker := iterations / workerCount
	remainingIterations := iterations % workerCount

	scoreChan := make(chan float64, workerCount)
	var wg sync.WaitGroup
	wg.Add(workerCount)

	for i := 0; i < workerCount; i++ {
		go func(workerId int) {
			defer wg.Done()
			iterationsForWorker := numIterationsPerWorker
			if workerId < remainingIterations {
				iterationsForWorker++
			}

			workerScore := 0.0
			for j := 0; j < iterationsForWorker; j++ {
				gtemp := *g
				gtemp.RunMove(move)

				for gtemp.GetState() == game.Ongoing {
					gtemp.RunMove(GetRandomMove(&gtemp))
				}

				workerScore += mc.EvaluateScore(&gtemp)
			}

			scoreChan <- workerScore
		}(i)
	}

	wg.Wait()
	close(scoreChan)

	totalScore := 0.0
	for score := range scoreChan {
		totalScore += score
	}

	return totalScore / float64(iterations)
}
