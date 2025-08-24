package session

import (
	skey "server/services/authService/secretKeys"
	typ "server/types"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewSession(t *testing.T) {
	sks := skey.NewSecretKeyService(100)
	key := sks.CurrentKey()
	userId := typ.UserId(1)
	s, err := NewSession(userId, key)
	if err != nil {
		t.Fatalf("failed to create new session: %v", err)
	}

	assert.Equal(t, s.UserId(), typ.UserId(1), "userId doesn't match expect value")
	assert.NotEmpty(t, s.JWEToken())
	assert.NotEmpty(t, s.TokenExpiry())
	assert.NotEmpty(t, s.Name())
	assert.WithinDuration(t, time.Now(), s.TokenExpiry(), time.Hour, "TokenExpiry should be recent")
	assert.True(t, s.Name() == "session_token")
}
