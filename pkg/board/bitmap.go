package board

func boardToBitmap(board Board, rows, cols int) (uint64, uint64) {
	var red uint64
	var blue uint64
	for r := 0; r < rows; r++ {
		for c := 0; c < cols; c++ {
			pos := NewPosition(r, c)
			if board.GetPiece(pos).Color == Red {
				red |= 1 << (r*cols + c)
			} else if board.GetPiece(pos).Color == Blue {
				blue |= 1 << (r*cols + c)
			}
		}
	}
	return red, blue
}

func getMaskForPosition(row, col int) uint64 {
	return (1 << (row*BoardCols + col))

}

func getMaskSpot(row, col int) int {
	return (row*BoardCols + col)
}

func isBitSet(bitmap uint64, row, col int) bool {
	return (bitmap & (1 << (row*BoardCols + col))) != 0
}

func setBit(bitmap *uint64, row, col int) {
	*bitmap = *bitmap | (1 << (row*BoardCols + col))
}

func clearBit(bitmap *uint64, row, col int) {
	*bitmap = *bitmap & ^(1 << (row*BoardCols + col))
}

func isBitSetByMask(bitmap uint64, mask uint64) bool {
	return (bitmap & mask) != 0
}

func setBitByMask(bitmap *uint64, mask uint64) {
	*bitmap = *bitmap | mask
}

func clearBitByMask(bitmap *uint64, mask uint64) {
	*bitmap = *bitmap & ^mask
}
func isBitSetByPos(bitmap uint64, pos *Position) bool {
	return (bitmap & pos.mask) != 0
}

func setBitByPos(bitmap *uint64, pos *Position) {
	*bitmap = *bitmap | pos.mask

}

func clearBitByPos(bitmap *uint64, pos *Position) {
	*bitmap = *bitmap & ^pos.mask
}
func bitmapToBoard(bitmap uint64, rows, cols int) [][]bool {
	board := make([][]bool, rows)
	for r := 0; r < rows; r++ {
		board[r] = make([]bool, cols)
		for c := 0; c < cols; c++ {
			board[r][c] = isBitSet(bitmap, r, c)
		}
	}
	return board
}
