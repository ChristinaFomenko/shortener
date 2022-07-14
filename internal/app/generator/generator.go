package generator

import (
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
