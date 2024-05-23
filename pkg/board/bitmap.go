package board

func boardToBitmap(board Board, rows, cols int) (uint64, uint64) {
	var red uint64
	var blue uint64
	for r := 0; r < rows; r++ {
		for c := 0; c < cols; c++ {
			if board.GetPiece(&Position{r, c}).Color == Red {
				red |= 1 << (r*cols + c)
			} else if board.GetPiece(&Position{r, c}).Color == Blue {
				blue |= 1 << (r*cols + c)
			}
		}
	}
	return red, blue
}

func getBit(bitmap uint64, row, col, cols int) bool {
	return (bitmap & (1 << (row*cols + col))) != 0
}

func setBit(bitmap uint64, row, col, cols int) uint64 {
	return bitmap | (1 << (row*cols + col))
}

func clearBit(bitmap uint64, row, col, cols int) uint64 {
	return bitmap & ^(1 << (row*cols + col))
}

func bitmapToBoard(bitmap uint64, rows, cols int) [][]bool {
	board := make([][]bool, rows)
	for r := 0; r < rows; r++ {
		board[r] = make([]bool, cols)
		for c := 0; c < cols; c++ {
			board[r][c] = getBit(bitmap, r, c, cols)
		}
	}
	return board
}
