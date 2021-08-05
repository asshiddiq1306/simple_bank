package token

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

var (
	ErrExpiredToken = errors.New("token already expired")
	ErrInvalidToken = errors.New("token invalid")
)

type AuthPay struct {
	ID        uuid.UUID `json:"id"`
	Username  string    `json:"username"`
	IssuedAt  time.Time `json:"issued_at"`
	ExpiredAt time.Time `json:"expired_at"`
}

func NewAuthPay(username string, duration time.Duration) (*AuthPay, error) {
	generateID, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}

	payload := &AuthPay{
		ID:        generateID,
		Username:  username,
		IssuedAt:  time.Now(),
		ExpiredAt: time.Now().Add(duration),
	}

	return payload, nil
}

func (payload *AuthPay) Valid() error {
	if time.Now().After(payload.ExpiredAt) {
		return ErrExpiredToken
	}
	return nil
}
