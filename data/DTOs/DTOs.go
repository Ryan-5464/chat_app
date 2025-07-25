package DTO

import "encoding/json"

type SwitchChat struct {
	ChatId string `json:"ChatId"`
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
