package testdata

import (
	ent "server/data/entities"
	cred "server/services/authService/credentials"
)

func TestUser() ent.User {
	pwdHash, err := cred.NewPwdHash([]byte("Testpwd123#!*"))
	if err != nil {
		return ent.User{}
	}

	return ent.User{
		Name:    "testName",
		Email:   "test@outlook.com",
		PwdHash: pwdHash,
	}
}
