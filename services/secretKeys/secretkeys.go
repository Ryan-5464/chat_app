package secretkeys

import (
	"crypto/rand"
	"fmt"
	"time"
)

func NewSecretKeyService(rotationInterval int) *SecretKeyService {
	s := &SecretKeyService{}
	s.startkeyRotation(rotationInterval)
	return s
}

type SecretKey [32]byte

func (s *SecretKey) Set(bytes []byte) {
	*s = [32]byte(bytes)
}

func (s SecretKey) Bytes() []byte {
	return s[:]
}

func (s SecretKey) IsZero() bool {
	for _, b := range s {
		if b != 0 {
			return false
		}
	}
	return true
}

type SecretKeyService struct {
	currentKey  SecretKey
	previousKey SecretKey
	signal      chan int
}

func (s *SecretKeyService) GetCurrentKey() SecretKey {
	return s.currentKey
}

func (s *SecretKeyService) GetPreviousKey() SecretKey {
	return s.previousKey
}

func (s *SecretKeyService) generateNewKey() error {
	bytes := make([]byte, 32)
	_, err := rand.Read(bytes)
	if err != nil {
		return fmt.Errorf("failed to generate bytes for secret key: %w", err)
	}
	s.previousKey = s.currentKey
	s.currentKey.Set(bytes)
	s.signal <- 1
	return nil
}

func (s *SecretKeyService) startkeyRotation(timeIntervalSeconds int) {

	go func() {
		timer := time.NewTicker(time.Duration(timeIntervalSeconds) * time.Second)
		defer timer.Stop()
		generateNewKeyPing := timer.C

		s.generateNewKey()

		for {
			select {
			case <-generateNewKeyPing:
				s.generateNewKey()
				fmt.Println("New secret key generated")
			}
		}
	}()
}
