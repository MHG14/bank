package token

import (
	"fmt"
	"time"

	"github.com/o1egl/paseto"
	"golang.org/x/crypto/chacha20poly1305"
)

type PASETOMaker struct {
	paseto       *paseto.V2
	symmetricKey []byte
}

func NewPASETOMaker(symmetricKey string) (Maker, error) {
	if len(symmetricKey) != chacha20poly1305.KeySize {
		return nil, fmt.Errorf("invalid key size, must be exactly %d characters", chacha20poly1305.KeySize)
	}

	maker := &PASETOMaker{
		paseto:       paseto.NewV2(),
		symmetricKey: []byte(symmetricKey),
	}

	return maker, nil
}

func (pasetoMaker *PASETOMaker) CreateToken(username string, duration time.Duration) (string, error) {
	payload, err := NewPayload(username, duration)
	if err != nil {
		return "", err
	}

	return pasetoMaker.paseto.Encrypt(pasetoMaker.symmetricKey, payload, nil)
}

func (pasetoMaker *PASETOMaker) VerifyToken(token string) (*Payload, error) {
	payload := &Payload{}

	if err := pasetoMaker.paseto.Decrypt(token, pasetoMaker.symmetricKey, payload, nil); err != nil {
		return nil, ErrInvalidToken
	}

	err := payload.Valid()

	if err != nil {
		return nil, err
	}
	return payload, nil
}
