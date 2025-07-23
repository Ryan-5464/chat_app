package DTO

import "encoding/json"

type SwitchChat struct {
	ChatId string `json:"ChatId"`
}

type Payload struct {
	Type string          `json:"Type"`
	Data json.RawMessage `json:"Data"`
}
