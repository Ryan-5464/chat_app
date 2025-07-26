package handler

import (
	dto "server/data/DTOs"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParsePayload(t *testing.T) {
	chatSwitchRequestMock := []byte(`{"Type":"SwitchChat","Data":{"ChatId":"2"}}`)
	result, err := parsePayload(chatSwitchRequestMock)
	if err != nil {
		t.Fatalf("result: %v, error: %v", result, err)
	}

	expected := dto.Payload{
		Type: "SwitchChat",
		Data: []byte(`{"ChatId":"2"}`),
	}

	assert.Equal(t, expected, result)
}

func TestParseSwitchChatData(t *testing.T) {
	chatSwitchDataMock := []byte(`{"ChatId":"2"}`)
	result, err := parseSwitchChatData(chatSwitchDataMock)
	if err != nil {
		t.Fatalf("result: %v, error: %v", result, err)
	}

	expected := dto.SwitchChat{
		ChatId: "2",
	}

	assert.Equal(t, expected, result)
}
