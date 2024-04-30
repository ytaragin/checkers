package board

// import "fmt"

type Move interface {
    DoMove(b *Board) bool
    IsValid(b *Board, color PieceColor) bool
    
}

type PlainMove struct {
    Start,End Position
}

func(pm *PlainMove) IsValid(b *Board, color PieceColor) bool { 
    p :=  b.getPiece(&pm.Start)
    return b.isSpotEmpty(&pm.End) && 
           (p != nil && p.Color == color) 
}

func (p *PlainMove) DoMove(b *Board) bool {
    // Implementation for PlainMove's DoMove
    // fmt.Println("Performing a plain move")

    return b.movePiece(&p.Start, &p.End) != nil
}

// JumpMove represents a possible move for a piece
type JumpMove struct {
    Jump Position
    Move PlainMove
}

func(j* JumpMove) IsValid(b *Board, color PieceColor) bool { 
    p :=  b.getPiece(&j.Move.Start)
    jspot := b.getPiece(&j.Jump)
    return b.isSpotEmpty(&j.Move.End) && 
           (p != nil && p.Color == color) &&
           (jspot != nil && jspot.Color != color)
}

func (j *JumpMove) DoMove(b *Board) bool{
    // Implementation for JumpMove's DoMove
    // fmt.Println("Performing a jump move")

    b.movePiece(&j.Move.Start, &j.Move.End)
    b.removePiece(&j.Jump)

    return true
}

func CreateMove(row, col, newrow, newcol int) *PlainMove {
   if newcol >= 0 && newcol < BoardCols && newrow >= 0 && newrow < BoardRows {
           return   &PlainMove{ 
                  Start: Position{Row: row, Col: col}, 
                  End: Position{Row:newrow, Col: newcol},
              }
          
   }

   return nil;
}

func createAppendMove(row, col, newrow, newcol int, moves []PlainMove) []PlainMove {
    m := CreateMove(row, col, newrow, newcol) 
    if m != nil {
       moves = append(moves, 
              PlainMove{ 
                  Start: Position{Row: row, Col: col}, 
                  End: Position{Row:newrow, Col: newcol},
              },
          )
   }

   return moves;
}

func createMovesInRow(row, col, dir int, moves []PlainMove) []PlainMove   {
    moves = createAppendMove(row, col, row+dir, col-1, moves)
    moves = createAppendMove(row, col, row+dir, col+1, moves)

    return moves;
}

func createKingMoves(row, col int, moves []PlainMove) []PlainMove {
    moves = createMovesInRow(row, col, RedDirection, moves);
    moves = createMovesInRow(row, col, BlueDirection, moves);

    return moves
}

func CreateJump(row, col, skiprow, skipcol, newrow, newcol int) *JumpMove {
    
    m := CreateMove(row, col, newrow, newcol) 
    if (m != nil) {
    
    j := &JumpMove{ 
           Move: *m,
           Jump: Position{Row: skiprow, Col:skipcol}, 
       }

       return j
   }
   return nil

}

func createAppendJump(row, col, skiprow, skipcol, newrow, newcol int, jumps []JumpMove) []JumpMove {
    j := CreateJump(row, col, skiprow, skipcol, newrow, newcol ) 
    if j != nil {
       jumps = append(jumps,  *j)
   }

   return jumps;
}
func createJumps(row, col, dir int, jumps[] JumpMove) []JumpMove {
    jumps = createAppendJump(row, col, row+dir, col+1, row+2*dir, col+2, jumps)
    jumps = createAppendJump(row, col, row+dir, col-1, row+2*dir, col-2, jumps)

    return jumps
}

func createKingJumps(row, col int, jumps []JumpMove) []JumpMove {
    jumps = createJumps(row, col, RedDirection, jumps);
    jumps = createJumps(row, col, BlueDirection, jumps);

    return jumps
}

