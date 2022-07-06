package generator

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"math/rand"
	"time"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

var secretkey = []byte("dsfewf64jwlj6so4difslkdj321")

type generator struct {
}

func NewGenerator() *generator {
	return &generator{}
}

func (g *generator) GenerateID() string {
	return generate(10)

}

func generate(n int) string {
	rand.Seed(time.Now().UnixNano())
	bytes := make([]byte, n)
	for i := range bytes {
		bytes[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(bytes)
}

func GenerateNewUserCookie() (string, string) {
	newID := make([]byte, 4)

	rand.Read(newID)
	encodedID := hex.EncodeToString(newID)

	h := hmac.New(sha256.New, secretkey)
	h.Write(newID)
	sign := h.Sum(nil)
	return encodedID, encodedID + hex.EncodeToString(sign)
}

func GetUserIDFromCookie(cookie string) (string, error) {
	data, err := hex.DecodeString(cookie)
	if err != nil {
		return "", err
	}
	id := data[:4]
	h := hmac.New(sha256.New, secretkey)
	h.Write(data[:4])
	sign := h.Sum(nil)

	if hmac.Equal(sign, data[4:]) {
		return hex.EncodeToString(id), nil
	} else {
		err := errors.New("sign check error")
		return "", err
	}
}
