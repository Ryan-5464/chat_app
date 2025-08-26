package view

import (
	ent "server/data/entities"
	typ "server/types"
	"server/util"
)

func findLastestMessageId(messages []ent.Message) typ.MessageId {
	util.Log.FunctionInfo()

	var latestMsgId typ.MessageId
	for _, msg := range messages {
		if msg.Id > latestMsgId {
			latestMsgId = msg.Id
		}
	}

	return latestMsgId
}
