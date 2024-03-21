package cipher

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

// PolybiusSquare is a 2D slice of runes representing a polybius square.
type PolybiusSquare [][]rune

// Interface guard for fmt.Stringer.
var _ fmt.Stringer = PolybiusSquare{}

// LookupLetter returns the index of the letter in the polybius square.
// Returns an error if the letter is not found in the polybius square.
func (s PolybiusSquare) LookupLetter(letter rune) (int, int, error) {
	for i, row := range s {
		for j, l := range row {
			if l == letter {
				return i, j, nil
			}
		}
	}
	return 0, 0, fmt.Errorf("letter not found in polybius square")
}

// GetLetter returns the letter at the given index in the polybius square.
// Returns an error if the index is out of range.
func (s PolybiusSquare) GetLetter(i, j int) (rune, error) {
	if i < 0 || i >= len(s) || j < 0 || j >= len(s[i]) {
		return 0, fmt.Errorf("index out of range")
	}
	return s[i][j], nil
}

// PolybiusSquare implements the fmt.Stringer interface.
func (s PolybiusSquare) String() string {
	var sb strings.Builder
	for _, row := range s {
		for letterIdx, letter := range row {
			sb.WriteString(string(letter))
			if letterIdx != len(row) {
				sb.WriteString(" ")
			}
		}
		sb.WriteString("\n")
	}
	return sb.String()
}

// PolybiusEncryptor implements Encryptor.
// It uses a polybius square to encrypt and decrypt strings.
// https://en.wikipedia.org/wiki/Polybius
type PolybiusEncryptor struct {
	PolybiusSquare       PolybiusSquare
	IgnoreUnknownLetters bool
}

// Interface guard for Encryptor.
var _ Encryptor = (*PolybiusEncryptor)(nil)

// NewPolybiusEncryptor returns a new PolybiusEncryptor with the given alphabet.
// If the alphabet is empty, it returns an error.
func NewPolybiusEncryptor(alphabet string, ignoreUnknownLetters bool) (*PolybiusEncryptor, error) {
	polybiusSquare, err := CreatePolybiusSquare(alphabet)
	if err != nil {
		return nil, err
	}
	return &PolybiusEncryptor{
		PolybiusSquare:       polybiusSquare,
		IgnoreUnknownLetters: ignoreUnknownLetters,
	}, nil
}

// MustNewPolybiusEncryptor is like NewPolybiusEncryptor, but panics if an error occurs.
func MustNewPolybiusEncryptor(alphabet string, ignoreUnknownLetters bool) *PolybiusEncryptor {
	polybiusEncryptor, err := NewPolybiusEncryptor(alphabet, ignoreUnknownLetters)
	if err != nil {
		panic(err)
	}
	return polybiusEncryptor
}

// Type returns the type of the encryptor.
func (e *PolybiusEncryptor) Type() string {
	return "polybius"
}

// Encrypt encrypts the input using the polybius square.
// It ignores (omits) unknown letters if IgnoreUnknownLetters is true.
// Returns an error if the input contains unknown letters and IgnoreUnknownLetters is false.
// The encryption format is "i-j i-j ..." for each letter in the input which is the index of the letter in the polybius square.
func (e *PolybiusEncryptor) Encrypt(input string) (string, error) {
	var encrypted string
	for _, letter := range input {
		i, j, err := e.PolybiusSquare.LookupLetter(letter)
		if err != nil {
			if e.IgnoreUnknownLetters {
				continue
			}
			return "", err
		}
		encrypted += fmt.Sprintf("%d-%d ", i+1, j+1)
	}
	return strings.TrimRight(encrypted, " "), nil
}

// Decrypt decrypts the input using the polybius square.
// It returns an error if the input is not in the correct format.
func (e *PolybiusEncryptor) Decrypt(input string) (string, error) {
	var decrypted string
	for _, encryptedLetter := range strings.Split(input, " ") {
		pair := strings.Split(encryptedLetter, "-")
		if len(pair) != 2 {
			return "", fmt.Errorf("invalid encrypted letter: %s", encryptedLetter)
		}
		i, err := strconv.Atoi(pair[0])
		if err != nil {
			return "", err
		}
		j, err := strconv.Atoi(pair[1])
		if err != nil {
			return "", err
		}
		letter, err := e.PolybiusSquare.GetLetter(i-1, j-1)
		if err != nil {
			return "", err
		}
		decrypted += string(letter)
	}
	return decrypted, nil
}

// DefaultAlphabet is the default alphabet used by the default PolybiusEncryptor.
const DefaultAlphabet = "abcdefghijklmnopqrstuvwxyz"

// DefaultIgnoreUnknownLetters is the default value of ignoreUnknownLetters used by the default PolybiusEncryptor.
const DefaultIgnoreUnknownLetters = false

// Default package global PolybiusEncryptor.
var DefaultPolybiusEncryptor = MustNewPolybiusEncryptor(DefaultAlphabet, DefaultIgnoreUnknownLetters)

// CalculateMinPolybiusSquareSize calculates the minimum size of the polybius square.
// It returns 0 if the alphabet length is 0.
func CalculateMinPolybiusSquareSize(alphabetLength int) int {
	if alphabetLength <= 0 {
		return 0
	}

	// get square size and add 1 if the polybius square couldn't fit all the letters
	size := int(math.Sqrt(float64(alphabetLength)))
	if size*size < alphabetLength {
		size++
	}

	return size
}

// CreatePolybiusSquare creates a new polybius square with the given alphabet.
// If the alphabet is empty, it returns an error.
func CreatePolybiusSquare(alphabet string) (PolybiusSquare, error) {
	// check if alphabet is empty
	if len(alphabet) == 0 {
		return nil, fmt.Errorf("empty alphabet")
	}

	// check if alphabet is unique
	for i := 0; i < len(alphabet); i++ {
		for j := i + 1; j < len(alphabet)-1; j++ {
			if alphabet[i] == alphabet[j] {
				return nil, fmt.Errorf("alphabet must be unique")
			}
		}
	}

	// calculate square size and initialize
	size := CalculateMinPolybiusSquareSize(len(alphabet))
	square := make([][]rune, size)
	for i := range square {
		square[i] = make([]rune, size)
	}

	// fill square with letters
	var idx int
	for i := 0; i < size; i++ {
		for j := 0; j < size; j++ {
			if idx < len(alphabet) {
				square[i][j] = rune(alphabet[idx])
				idx++
			}
		}
	}

	return square, nil
}
