package main

import (
	"dpb03/pkg/text"
	"fmt"
)

func main() {
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
}
