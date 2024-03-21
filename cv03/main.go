package main

import (
	"dpb03/pkg/chessboard"
	"dpb03/pkg/cipher"
	"dpb03/pkg/numbers"
	"dpb03/pkg/text"
	"fmt"
)

func main() {
	numToFactor := 220
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
	dataPath := "data/book.txt"
	letterFreq, wordFreq, err := text.TextAnalysis(dataPath)
	if err != nil {
		panic(err.Error())
	}
	fmt.Printf("letter freq:\n%v\n", letterFreq)
	fmt.Printf("word freq:\n%v\n", wordFreq)

	numWords := 10
	minWordLength := 8
	filteredWordFreq := text.GetWords(numWords, minWordLength, wordFreq)
	fmt.Printf("top %d words of min length %d:\n%v\n", numWords, minWordLength, filteredWordFreq)

	// declare a custom encryptor with a custom alphabet
	alphabet := "abcdefghijklmnopqrstuvwxyz"
	ignoreUnknownLetters := true
	enc, err := cipher.NewPolybiusEncryptor(alphabet, ignoreUnknownLetters)
	if err != nil {
		panic(err.Error())
	}

	// or use the default encryptor with a default alphabet
	// enc := cipher.DefaultPolybiusEncryptor

	fmt.Printf("the polybius square is:\n%v\n", enc.PolybiusSquare)

	toEncrypt := "hello my friend"
	encrypted, err := enc.Encrypt(toEncrypt)
	if err != nil {
		panic(err.Error())
	}

	fmt.Printf("string %q is encrypted to %q\n", toEncrypt, encrypted)
	decrypted, err := enc.Decrypt(encrypted)
	if err != nil {
		panic(err.Error())
	}
	fmt.Printf("string %q is decrypted to %q\n", encrypted, decrypted)
}
