package cipher

// Encryptor is an interface for encrypting and decrypting strings.
type Encryptor interface {
	// Type returns the type of the encryptor.
	Type() string
	// Encrypt encrypts the input.
	Encrypt(input string) (string, error)
	// Decrypt decrypts the input.
	Decrypt(input string) (string, error)
}
