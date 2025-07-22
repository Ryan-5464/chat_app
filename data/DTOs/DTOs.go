package DTO

import "encoding/json"

type SwitchChat struct {
	ChatId int `json:"ChatId"`
}

type payload struct {
	Type string          `json:"Type"`
	Data json.RawMessage `json:"Data"`
}
