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
	if queenX < 0 || queenX >= n || queenY < 0 || queenY >= m {
		return nil, fmt.Errorf("invalid queen coordinates: (%d, %d)", queenX, queenY)
	}

	var chessBoard ChessBoard

	for i := 0; i < n; i++ {
		var row []string
		for j := 0; j < m; j++ {
			if i == queenX && j == queenY {
				row = append(row, "D")
				continue
			}
			if i == queenX || j == queenY || i+j == queenX+queenY || i-j == queenX-queenY {
				row = append(row, "*")
				continue
			}
			row = append(row, ".")
		}
		chessBoard = append(chessBoard, row)
	}

	return chessBoard, nil
}
