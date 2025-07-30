package types

import "server/lib"

type ChatId int64

func (c ChatId) Int64() int64 {
	return int64(c)
}

type MessageId int64

func (m MessageId) Int64() int64 {
	return int64(m)
}

type UserId int64

func (u UserId) Int64() int64 {
	return int64(u)
}

func (u UserId) String() string {
	return lib.ConvertInt64ToString(u.Int64())
}

type Rows []map[string]any
