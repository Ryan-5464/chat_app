package secretkeys

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerateFirstKey(t *testing.T) {
	sks := NewSecretKeyService(100000)
	<-sks.Signal
	assert.True(t, !sks.currentKey.IsZero(), "x ck shouldnt be zero")
	assert.True(t, sks.previousKey.IsZero(), "x pk should be zero")
}

func TestKeyRotation(t *testing.T) {
	sks := NewSecretKeyService(1)
	<-sks.Signal
	assert.True(t, !sks.currentKey.IsZero(), "1 ck shouldnt be zero")
	assert.True(t, sks.previousKey.IsZero(), "1 pk should be zero")

	<-sks.Signal
	assert.True(t, !sks.currentKey.IsZero(), "2 ck shouldnt be zero")
	assert.True(t, !sks.previousKey.IsZero(), "2 pk shouldnt be zero")
	ckey := sks.currentKey

	<-sks.Signal
	assert.True(t, !sks.currentKey.IsZero(), "3 ck shouldnt be zero")
	assert.True(t, !sks.previousKey.IsZero(), "3 pk shouldnt be zero")
	pkey := sks.previousKey

	assert.Equal(t, pkey.Bytes(), ckey.Bytes())
}
