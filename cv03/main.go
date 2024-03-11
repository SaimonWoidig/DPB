package main

import (
	"dpb03/pkg/chessboard"
	"dpb03/pkg/numbers"
	"fmt"
)

func main() {
	numToFactor := 1
	factors, err := numbers.Factorize(numToFactor)
	if err != nil {
		panic(err.Error())
	}
	fmt.Printf("The prime factors of %d are: %v.\n", numToFactor, factors)

	boardSize := 8
	queenX := 1
	queenY := 1
	queenBoard, err := chessboard.Queen(boardSize, boardSize, queenX, queenY)
	if err != nil {
		panic(err.Error())
	}
	fmt.Printf("The chessboard with the queen on (%d, %d) and the squares under attack:\n%v\n", queenX, queenY, queenBoard)

	upperBound := 200
	numToCensor := 17
	censoredNumbers, err := numbers.CensorNumber(upperBound, numToCensor)
	if err != nil {
		panic(err.Error())
	}
	fmt.Printf("The numbers from 1 to %d that do not contain %d are:\n%v\n", upperBound, numToCensor, censoredNumbers)
}
