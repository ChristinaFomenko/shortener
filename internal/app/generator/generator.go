package generator

import (
	"math/rand"
	"time"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

type generator struct{}

func NewGenerator() *generator {
	return &generator{}
}

func (g *generator) Letters(n int64) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Int63()%int64(len(letterBytes))]
	}

	return string(b)
}

func (g *generator) Random(n int64) string {
	b := make([]byte, n)
	rand.Seed(time.Now().Unix())
	_, _ = rand.Read(b)

	return string(b)
}
