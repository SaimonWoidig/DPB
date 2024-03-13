package chessboard

import "fmt"

// ChessBoard represents a chessboard.
type ChessBoard [][]string

// ChessBoard implements the fmt.Stringer interface.
var _ fmt.Stringer = ChessBoard{}

// String returns a string representation of the chessboard.
func (cb ChessBoard) String() string {
	var chessBoard string
	for _, row := range cb {
		for _, square := range row {
			chessBoard += square
		}
		chessBoard += "\n"
	}
	return chessBoard
}

// Queen takes a board size and queen coordinates and returns a chessboard with the queen placed on the board and the squares under attack.
// The queen is symbolized by 'D'. The squares under attack are symbolized by '*' and the empty squares are symbolized by '.'.
func Queen(n, m, queenX, queenY int) (ChessBoard, error) {
	// check if queen coordinates are valid - they must be in the range [0, n-1] and [0, m-1]
	if queenX < 0 || queenX >= n || queenY < 0 || queenY >= m {
		return nil, fmt.Errorf("invalid queen coordinates: (%d, %d)", queenX, queenY)
	}

	var chessBoard ChessBoard
	// create chessboard
	for i := 0; i < n; i++ {
		var row []string
		for j := 0; j < m; j++ {
			// check if queen is on this board tile
			if i == queenX && j == queenY {
				row = append(row, "D")
				continue
			}
			// check if this board tile is under attack
			if i == queenX || j == queenY || i+j == queenX+queenY || i-j == queenX-queenY {
				// this board tile is under attack
				row = append(row, "*")
				continue
			}
			// this board tile is empty and not under attack
			row = append(row, ".")
		}
		// add row to chessboard
		chessBoard = append(chessBoard, row)
	}

	return chessBoard, nil
}
