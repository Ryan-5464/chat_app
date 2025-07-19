package types

type UserId int64

func (u UserId) Int64() int64 {
	return int64(u)
}
