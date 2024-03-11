package cipher

type Encryptor interface {
	Type() string
	Encrypt(input string) string
	Decrypt(input string) string
}
