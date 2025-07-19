package dbServiceModels

import (
	typ "server/types"
)

type User struct {
	Id      typ.UserId
	Name    string
	Email   typ.Email
	PwdHash typ.PwdHash
}
