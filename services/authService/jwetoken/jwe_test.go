package jwetoken

import (
	skey "server/services/authService/secretKeys"
	"testing"
	"time"

	typ "server/types"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// refactor with channels and remove sleep
func secretKey() skey.SecretKey {
	s := skey.NewSecretKeyService(100000)
	// keys generated via go routine. Sleep gives time for keys te be generated
	time.Sleep(time.Duration(500) * time.Millisecond)
	return s.GetCurrentKey()
}

func TestJWEGenerationVerificationPipeline(t *testing.T) {
	key := secretKey()
	userId := typ.UserId(1)
	jwe, err := NewJWE(userId, key)
	require.NoError(t, err, "Expected NewJWE to succeed")

	assert.NotEmpty(t, jwe.Bytes(), "Token should not be empty")
	assert.Equal(t, userId, jwe.UserId(), "UserId should match")
	assert.WithinDuration(t, time.Now(), jwe.IssuedAt(), time.Second, "IssuedAt should be recent")
	assert.WithinDuration(t, time.Now().Add(time.Hour), jwe.TokenExpiry(), time.Second, "TokenExpiry should be 1 hour from now")

	// Parse and verify
	parsedJWE, err := ParseAndVerifyJWE(jwe.String(), key)
	require.NoError(t, err, "Expected ParseAndVerifyJWE to succeed")

	assert.Equal(t, userId, parsedJWE.UserId(), "Parsed UserId should match")
	assert.WithinDuration(t, jwe.IssuedAt(), parsedJWE.IssuedAt(), time.Second, "IssuedAt should match")

}

func TestExpiredToken(t *testing.T) {
	key := secretKey()
	userId := typ.UserId(1)

	// Manually construct claims with past expiry
	claims := Claims{
		UserId:      userId,
		TokenExpiry: time.Now().Add(-1 * time.Hour),
		IssuedAt:    time.Now().Add(-2 * time.Hour),
	}
	jwe, err := generateJWE(claims, key)
	require.NoError(t, err)

	_, err = ParseAndVerifyJWE(jwe.String(), key)
	assert.Error(t, err, "Expected error due to expired token")
	assert.Contains(t, err.Error(), "token has Uxpired") // Note: "Uxpired" typo in your code!
}

func TestInvalidKey(t *testing.T) {
	key := secretKey()
	newKey := secretKey()
	userId := typ.UserId(1)

	jwe, err := NewJWE(userId, key)
	require.NoError(t, err)

	_, err = ParseAndVerifyJWE(jwe.String(), newKey)
	assert.Error(t, err, "Expected error due to invalid decryption key")
}

func TestMalformedToken(t *testing.T) {
	key := secretKey()
	userId := typ.UserId(1)

	jwe, err := NewJWE(userId, key)
	token := jwe.Bytes()
	idx := 240
	assert.Less(t, idx, len(token), "Index should be less than token length")
	assert.Equal(t, token[idx], token[idx], "Byte should match")
	newByte := byte(0)
	token[240] = newByte
	assert.Equal(t, token[idx], newByte, "Byte should match")

	tamperedToken := string(token)
	_, err = ParseAndVerifyJWE(tamperedToken, key)
	assert.Error(t, err, "Expected error for malformed token")
}
