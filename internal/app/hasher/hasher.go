package hasher

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
)

type hasher struct {
	secret []byte
}

func NewHasher(secret []byte) *hasher {
	return &hasher{
		secret: secret,
	}
}

func (e *hasher) Sign(value string) (string, error) {
	encoded := hex.EncodeToString([]byte(value))

	h := hmac.New(sha256.New, e.secret)
	_, err := h.Write([]byte(value))
	if err != nil {
		return "", fmt.Errorf("sign error: %w", err)
	}

	sign := h.Sum(nil)

	return encoded + hex.EncodeToString(sign), nil
}

func (e *hasher) Validate(value string, dataLength int64) (string, error) {
	decoded, err := hex.DecodeString(value)
	if err != nil {
		return "", fmt.Errorf("decode value error: %w", err)
	}

	valueData := decoded[:dataLength]
	valueSign := decoded[dataLength:]

	h := hmac.New(sha256.New, e.secret)
	h.Write(valueData)
	sign := h.Sum(nil)

	if !hmac.Equal(sign, valueSign) {
		fmt.Println(sign, decoded)
		return "", errors.New("sign validation error")
	}

	return string(valueData), nil
}
