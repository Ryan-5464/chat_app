package DTO

import (
	"encoding/json"
	"server/lib"
	typ "server/types"
)

type SwitchChat struct {
	ChatId string `json:"ChatId"`
}

func (s *SwitchChat) GetChatId() (typ.ChatId, error) {
	cid, err := lib.ConvertStringToInt64(s.ChatId)
	if err != nil {
		return typ.ChatId(0), err
	}
	return typ.ChatId(cid), nil
}

type NewChat struct {
	UserId string `json:"UserId"`
	Name   string `json:"Name"`
}

type NewMessage struct {
	UserId  string `json:"UserId"`
	ChatId  string `json:"ChatId"`
	ReplyId string `json:"ReplyId"`
	MsgText string `json:"MsgText"`
}

type Payload struct {
	Type string          `json:"Type"`
	Data json.RawMessage `json:"Data"`
}
