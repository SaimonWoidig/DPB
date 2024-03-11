package cipher

// PolybiusEncryptor implements Encryptor.
// https://en.wikipedia.org/wiki/Polybius
type PolybiusEncryptor struct {
	PolybiusSquare [][]string
}

// Interface guard for Encryptor.
var _ Encryptor = (*PolybiusEncryptor)(nil)

// NewPolybiusEncryptor returns a new PolybiusEncryptor.
func NewPolybiusEncryptor(alphabet string) *PolybiusEncryptor {
	alphabetLetters := []rune(alphabet)
	square := make([][]rune, len(alphabet)/2)
	for idx := range square {
		for idx2 := range square[idx] {
			square[idx][idx2] = alphabetLetters[idx*2+idx2]
		}
	}
	return &PolybiusEncryptor{}
}

func (e *PolybiusEncryptor) Type() string {
	return "polybius"
}

func (e *PolybiusEncryptor) Encrypt(input string) string {
	return input
}

func (e *PolybiusEncryptor) Decrypt(input string) string {
	return input
}

const defaultAlphabet = "abcdefghijklmnopqrstuvwxyz"

// Default package global PolybiusEncryptor.
var DefaultPolybiusEncryptor = NewPolybiusEncryptor(defaultAlphabet)
