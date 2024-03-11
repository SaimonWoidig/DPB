package text

import (
	"fmt"
	"net/url"
	"os"
	"strconv"
	"strings"
	"unicode"

	"github.com/gobs/sortedmap"
)

// TextAnalysisResult is a map of word frequencies.
// Key is the word, value is the frequency.
type TextAnalysisResult map[string]int

// ToSortedByValue returns a sorted map of word frequencies sorted in descending order by value.
func (t TextAnalysisResult) ToSortedByValue() sortedmap.SortedByValue[string, int] {
	return sortedmap.AsSortedByValue(t, false)
}

// Interface guard for fmt.Stringer.
var _ fmt.Stringer = TextAnalysisResult{}

// String implements fmt.Stringer.
func (t TextAnalysisResult) String() string {
	var str string
	for _, keyval := range sortedmap.AsSortedByValue(t, false) {
		str += fmt.Sprintf("%q: %d\n", keyval.Key, keyval.Value)
	}
	return str
}

// TextAnalysis analyses a text file and returns letter and word frequencies.
// Returns error if file cannot be read.
// This function is not 100% accurate, but it's good enough for this use-case.
func TextAnalysis(filepath string) (TextAnalysisResult, TextAnalysisResult, error) {
	fileContentsBytes, err := os.ReadFile(filepath)
	if err != nil {
		return nil, nil, err
	}
	fileContents := string(fileContentsBytes)

	letterFreq := make(TextAnalysisResult)
	wordFreq := make(TextAnalysisResult)

	for _, word := range strings.Fields(fileContents) {
		trimmedWord := strings.TrimSpace(word)
		// ignore urls
		if _, err := url.ParseRequestURI(trimmedWord); err == nil {
			continue
		}
		cleanWord := cleanWord(trimmedWord)
		// ignore empty words
		if cleanWord == "" {
			continue
		}
		// ignore numbers
		if _, err := strconv.Atoi(cleanWord); err == nil {
			continue
		}

		wordFreq[cleanWord] += 1
		for _, char := range word {
			// ignore non-letters
			if !unicode.IsLetter(char) {
				continue
			}
			letterFreq[string(char)] += 1
		}
	}

	return letterFreq, wordFreq, nil
}

// cleanWord cleans a word by removing unwanted characters.
// This function is not 100% accurate, but it's good enough for this use-case.
func cleanWord(word string) string {
	if word == "" {
		return ""
	}
	charReplacer := strings.NewReplacer(
		"\"", "",
		"?", "",
		".", "",
		"!", "",
		",", "",
		")", "",
		"(", "",
		":", "",
		" ", "",
		";", "",
		"$", "",
		"#", "",
		"*", "",
		"&", "",
		"%", "",
		"^", "",
		"@", "",
		"=", "",
		"\\", "",
		"/", "",
		"--", "",
		"–", "",
		"”", "",
	)
	return strings.TrimSuffix(strings.TrimPrefix(strings.ToLower(charReplacer.Replace(word)), "-"), "-")
}

// GetWords returns the top numWords words of length >= minWordLength in wordFreq.
// The numWords and minWordLength parameters should be > 0. Otherwise, they are set to 1 and 1 respectively.
func GetWords(numWords int, minWordLength int, wordFreq TextAnalysisResult) TextAnalysisResult {
	// set defaults if bogus value
	if numWords < 1 {
		numWords = 1
	}
	if minWordLength < 1 {
		minWordLength = 1
	}
	words := make(TextAnalysisResult)
	for _, wordFreqPair := range wordFreq.ToSortedByValue() {
		word := wordFreqPair.Key
		freq := wordFreqPair.Value
		if len(word) < minWordLength {
			continue
		}
		words[word] = freq
		if len(words) == numWords {
			break
		}
	}
	return words
}
