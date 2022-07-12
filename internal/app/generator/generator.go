package generator

import (
	crypto "crypto/rand"
	math "math/rand"
	"time"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func init() {
	math.Seed(time.Now().Unix())
}

type generator struct{}

func NewGenerator() *generator {
	return &generator{}
}

func (g *generator) Letters(n int64) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[math.Int63()%int64(len(letterBytes))]
	}

	return string(b)
}

func (g *generator) Random(n int64) string {
	b := make([]byte, n)
	_, _ = crypto.Read(b)

	return string(b)
}
