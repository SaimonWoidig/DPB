package numbers

import (
	"fmt"
	"strconv"
)

// CensoredSequence represents a list of numbers that can contain censored numbers as asterisks.
type CensoredSequence []string

// CensoredSequence implements the fmt.Stringer interface.
var _ fmt.Stringer = CensoredSequence{}

// String returns a string representation of the censored sequence.
func (cs CensoredSequence) String() string {
	var censoredNumbers string
	for _, num := range cs {
		censoredNumbers += num + "\n"
	}
	return censoredNumbers
}

// CensorNumber returns a list of numbers from 1 to upperBound that do not contain numToCensor. The numbers that do not contain numToCensor are censored - marked with an asterisk.
// The input must be a positive integer.
func CensorNumber(upperBound, numToCensor int) (CensoredSequence, error) {
	if upperBound < 1 {
		return nil, fmt.Errorf("invalid input: %d", upperBound)
	}
	var censoredNumbers CensoredSequence

	for i := 1; i <= upperBound; i++ {
		if !doesNumberContainNumber(i, numToCensor) {
			censoredNumbers = append(censoredNumbers, strconv.Itoa(i))
		} else {
			censoredNumbers = append(censoredNumbers, "*")
		}
	}

	return censoredNumbers, nil
}

// doesNumberContainNumber returns true if the number n contains the number numToFilter.
func doesNumberContainNumber(n, numToFilter int) bool {
	numDigitsN := len(strconv.Itoa(n))
	numDigitsFilter := len(strconv.Itoa(numToFilter))
	for i := 0; i < numDigitsN; i++ {
		if i+numDigitsFilter > numDigitsN {
			return false
		}
		substr := strconv.Itoa(n)[i : i+numDigitsFilter]
		if substr == strconv.Itoa(numToFilter) {
			return true
		}
	}
	return false
}
