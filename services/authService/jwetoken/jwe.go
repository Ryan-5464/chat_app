package jwetoken

import (
	"encoding/json"
	"errors"
	"fmt"
	skey "server/services/authService/secretKeys"
	typ "server/types"
	"time"

	"gopkg.in/square/go-jose.v2"
)

type Claims struct {
	UserId      typ.UserId
	TokenExpiry time.Time
	IssuedAt    time.Time
}

type JWE struct {
	claims Claims
	token  []byte
}

func (j *JWE) String() string {
	return string(j.token)
}

func (j *JWE) Bytes() []byte {
	return j.token
}

func (j *JWE) Claims() *Claims {
	return &j.claims
}

func (j *JWE) UserId() typ.UserId {
	return j.claims.UserId
}

func (j *JWE) TokenExpiry() time.Time {
	return j.claims.TokenExpiry
}

func (j *JWE) IssuedAt() time.Time {
	return j.claims.IssuedAt
}

func NewJWE(userId typ.UserId, key skey.SecretKey) (JWE, error) {
	tokenExpiry := time.Now().Add(time.Hour)
	claims := Claims{
		UserId:      userId,
		TokenExpiry: tokenExpiry,
		IssuedAt:    time.Now(),
	}
	jwe, err := generateJWE(claims, key)
	if err != nil {
		return JWE{}, fmt.Errorf("failed to generate JWE : %w", err)
	}
	return jwe, nil
}

func ParseAndVerifyJWE(token string, key skey.SecretKey) (JWE, error) {
	verifiedToken, err := decryptAndVerifyToken(token, key)
	if err != nil {
		return JWE{}, fmt.Errorf("failed to decrypt token: %w", err)
	}

	var claims Claims
	if err := json.Unmarshal(verifiedToken, &claims); err != nil {
		return JWE{}, fmt.Errorf("failed to unmarshal claims; %w", err)
	}

	if err := validateClaims(claims); err != nil {
		return JWE{}, fmt.Errorf("invalid claims: %w", err)
	}

	return updateJWE(claims.UserId, key)
}

func updateJWE(userId typ.UserId, key skey.SecretKey) (JWE, error) {
	return NewJWE(userId, key)
}

func decryptAndVerifyToken(token string, key skey.SecretKey) ([]byte, error) {

	jweObject, err := jose.ParseEncrypted(token)
	if err != nil {
		return nil, fmt.Errorf("failed to parse encrypted token: %w", err)
	}

	signedJWT, err := jweObject.Decrypt(key.Bytes())
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt token: %w", err)
	}

	jwsObject, err := jose.ParseSigned(string(signedJWT))
	if err != nil {
		return nil, fmt.Errorf("failed to parse decrypted token: %w", err)
	}

	verifiedToken, err := jwsObject.Verify(key.Bytes())
	if err != nil {
		return nil, fmt.Errorf("failed to verify decrypted token : %w", err)
	}

	return verifiedToken, nil
}

func validateClaims(claims Claims) error {
	now := time.Now()

	if claims.UserId == 0 {
		return errors.New("invalid claims: subject (UserId) is empty")
	}

	if claims.IssuedAt.After(now) {
		return errors.New("invalid claims: issued-at time is in the future")
	}

	if now.After(claims.TokenExpiry) {
		return errors.New("invalid claims: token has Uxpired")
	}
	return nil
}

func generateJWE(claims Claims, key skey.SecretKey) (JWE, error) {

	claimsJSON, err := json.Marshal(claims)
	if err != nil {
		return JWE{}, fmt.Errorf("failed to marshal claims: %w", err)
	}

	signer, err := jose.NewSigner(
		jose.SigningKey{
			Algorithm: jose.HS256,
			Key:       key.Bytes(),
		},
		(&jose.SignerOptions{}).WithType("JWT"),
	)
	if err != nil {
		return JWE{}, fmt.Errorf("failed to create signer: %w", err)
	}

	jwsObject, err := signer.Sign(claimsJSON)
	if err != nil {
		return JWE{}, fmt.Errorf("failed to sign claims: %w", err)
	}

	signedJWT, err := jwsObject.CompactSerialize()
	if err != nil {
		return JWE{}, fmt.Errorf("failed to serialize signed claims: %w", err)
	}

	encrypter, err := jose.NewEncrypter(
		jose.A256GCM,
		jose.Recipient{
			Algorithm: jose.DIRECT,
			Key:       key.Bytes(),
		},
		(&jose.EncrypterOptions{}).WithContentType("JWT"),
	)
	if err != nil {
		return JWE{}, fmt.Errorf("failed to initialize encrypter: %w", err)
	}

	jweObject, err := encrypter.Encrypt([]byte(signedJWT))
	if err != nil {
		return JWE{}, fmt.Errorf("failed to encrypt JWT: %w", err)
	}

	token, err := jweObject.CompactSerialize()
	if err != nil {
		return JWE{}, fmt.Errorf("failed to serialize JWE: %w", err)
	}

	jwe := JWE{
		claims: claims,
		token:  []byte(token),
	}

	return jwe, nil
}
