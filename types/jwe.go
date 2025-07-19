package types

import (
	"encoding/json"
	"fmt"
	"time"

	"gopkg.in/square/go-jose.v2"
)

type Claims struct {
	UserId      UserId
	TokenExpiry time.Time
}

func NewJWEToken(claims Claims, key SecretKey) (JWE, error) {
	jwe, err := generateJWE(claims, key)
	if err != nil {
		return nil, fmt.Errorf("failed to generate JWE : %w", err)
	}
	return jwe, nil
}

func generateJWE(claims Claims, key SecretKey) (JWE, error) {

	claimsJSON, err := json.Marshal(claims)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal claims: %w", err)
	}

	signer, err := jose.NewSigner(
		jose.SigningKey{
			Algorithm: jose.HS256,
			Key:       key.Bytes(),
		},
		(&jose.SignerOptions{}).WithType("JWT"),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create signer: %w", err)
	}

	jwsObject, err := signer.Sign(claimsJSON)
	if err != nil {
		return nil, fmt.Errorf("failed to sign claims: %w", err)
	}

	signedJWT, err := jwsObject.CompactSerialize()
	if err != nil {
		return nil, fmt.Errorf("failed to serialize signed claims: %w", err)
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
		return nil, fmt.Errorf("failed to initialize encrypter: %w", err)
	}

	jweObject, err := encrypter.Encrypt([]byte(signedJWT))
	if err != nil {
		return nil, fmt.Errorf("failed to encrypt JWT: %w", err)
	}

	token, err := jweObject.CompactSerialize()
	if err != nil {
		return nil, fmt.Errorf("failed to serialize JWE: %w", err)
	}

	return JWE(token), nil
}

type JWE []byte
