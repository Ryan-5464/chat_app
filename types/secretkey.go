package types

type SecretKey []byte

func (s SecretKey) Bytes() []byte {
	return []byte(s)
}
