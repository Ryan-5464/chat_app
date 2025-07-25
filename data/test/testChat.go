package testdata

import (
	ent "server/data/entities"
	typ "server/types"
)

func TestChat() ent.Chat {
	return ent.Chat{
		Name:    "testChat1",
		AdminId: typ.UserId(1),
	}
}
