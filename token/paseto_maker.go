package token

import (
	"fmt"
	"time"

	"github.com/aead/chacha20poly1305"
	"github.com/o1egl/paseto"
)

type PasetoMaker struct {
	maker        *paseto.V2
	symmetricKey []byte
}

func NewPasetoMaker(symmetricKey string) (Maker, error) {
	if len(symmetricKey) != chacha20poly1305.KeySize {
		err := fmt.Errorf("invalid length : must be %d characters", chacha20poly1305.KeySize)
		return nil, err
	}

	paseto := &PasetoMaker{
		maker:        paseto.NewV2(),
		symmetricKey: []byte(symmetricKey),
	}

	return paseto, nil
}

func (pasetoMaker *PasetoMaker) CreateToken(username string, duration time.Duration) (string, error) {
	payload, err := NewAuthPay(username, duration)
	if err != nil {
		return "", err
	}

	return pasetoMaker.maker.Encrypt([]byte(pasetoMaker.symmetricKey), payload, nil)
}

func (pasetoMaker *PasetoMaker) VerifyToken(accessToken string) (*AuthPay, error) {
	payload := &AuthPay{}
	err := pasetoMaker.maker.Decrypt(accessToken, pasetoMaker.symmetricKey, payload, nil)
	if err != nil {
		return nil, ErrInvalidToken
	}

	err = payload.Valid()
	if err != nil {
		return nil, err
	}

	return payload, nil
}
