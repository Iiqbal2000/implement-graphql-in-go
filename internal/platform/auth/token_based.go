package platform

import (
	"errors"
	"fmt"
	"log"
	"time"

	// "github.com/aead/chacha20poly1305"
	"github.com/google/uuid"
	"github.com/o1egl/paseto"
)

var (
	ErrInvalidToken = errors.New("token is invalid")
	ErrExpiredToken = errors.New("token has expired")
	ErrInternal = errors.New("internal server error")
)

var secret_key = "ahluPW6HKmzAsL3Rirs7LSQ7orBKHP0f"

// var secret_key = os.Getenv("secret_key")
var pasetoV2 = paseto.NewV2()

type Payload struct {
	ID        uuid.UUID `json:"id"`
	UserID    string    `json:"user_id"`
	IssuedAt  time.Time `json:"issued_at"`
	ExpiredAt time.Time `json:"expired_at"`
}

func (p Payload) checkExpired() error {
	if time.Now().After(p.ExpiredAt) {
		return ErrExpiredToken
	}
	return nil
}

// func init() {
// 	// checking length of secret_key.
// 	checkSecretKey(chacha20poly1305.KeySize)
// }

func newPayload(userID string, duration time.Duration) (*Payload, error) {
	tokenID, err := uuid.NewRandom()
	if err != nil {
		return nil, fmt.Errorf("failure when generating uuid: %s", err.Error())
	}

	return &Payload{
		ID:        tokenID,
		UserID:    userID,
		IssuedAt:  time.Now(),
		ExpiredAt: time.Now().Add(duration),
	}, nil
}

func checkSecretKey(lenSecureKey int) {
	if len(secret_key) != lenSecureKey {
		log.Fatalf("invalid key size: must be exactly %d characters, got %d", lenSecureKey, len(secret_key))
	}
}

func createToken(userId string, duration time.Duration) (string, error) {
	paylod, err := newPayload(userId, duration)
	if err != nil {
		log.Println(err.Error())
		return "", ErrInternal
	}

	token, err := pasetoV2.Encrypt([]byte(secret_key), paylod, nil)
	if err != nil {
		return "", ErrInternal
	}
	return token, nil
}

func verifyToken(token string) (*Payload, error) {
	payload := &Payload{}
	err := pasetoV2.Decrypt(token, []byte(secret_key), payload, nil)
	if err != nil {
		return nil, ErrInvalidToken
	}

	err = payload.checkExpired()
	if err != nil {
		return nil, err
	}

	return payload, nil
}