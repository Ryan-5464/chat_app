package secretkeys

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGenerateSecretKey(t *testing.T) {
	sks := &SecretKeyService{}
	assert.True(t, sks.currentKey.IsZero())
	assert.True(t, sks.previousKey.IsZero())
	sks.generateNewKey()
	assert.True(t, !sks.currentKey.IsZero())
	assert.True(t, sks.previousKey.IsZero())
	sks.generateNewKey()
	assert.True(t, !sks.currentKey.IsZero())
	assert.True(t, !sks.previousKey.IsZero())
}

func TestGenerateFirstKey(t *testing.T) {
	interval := 100000
	mInterval := 500
	sks := NewSecretKeyService(interval)
	time.Sleep(time.Duration(mInterval) * time.Millisecond)
	assert.True(t, !sks.currentKey.IsZero(), "x ck shouldnt be zero")
	assert.True(t, sks.previousKey.IsZero(), "x pk should be zero")
}

func TestKeyRotation(t *testing.T) {
	// Keys generated every second.
	// Sleep 0.5 seconds to give time to generate first key
	// check update every second to be in sync with generation
	interval := 1
	mInterval := 500
	sks := NewSecretKeyService(interval)
	time.Sleep(time.Duration(mInterval) * time.Millisecond)
	assert.True(t, !sks.currentKey.IsZero(), "1 ck shouldnt be zero")
	assert.True(t, sks.previousKey.IsZero(), "1 pk should be zero")

	time.Sleep(time.Duration(interval) * time.Second)
	assert.True(t, !sks.currentKey.IsZero(), "2 ck shouldnt be zero")
	assert.True(t, !sks.previousKey.IsZero(), "2 pk shouldnt be zero")
	ckey := sks.currentKey

	time.Sleep(time.Duration(interval) * time.Second)
	assert.True(t, !sks.currentKey.IsZero(), "3 ck shouldnt be zero")
	assert.True(t, !sks.previousKey.IsZero(), "3 pk shouldnt be zero")
	pkey := sks.previousKey

	assert.Equal(t, pkey.Bytes(), ckey.Bytes())
}
