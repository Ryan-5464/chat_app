package types

func NewSecretKey(s string) SecretKey {
	return SecretKey([32]byte([]byte(s)))
}

type SecretKey [32]byte

func (s SecretKey) Bytes() []byte {
	return s[:]
}
