package helpers

import (
	"math/rand"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

type generator struct {
}

func NewGenerator() *generator {
	return &generator{}
}

func (g *generator) ID() string {
	return generate(10)

}

func generate(n int) string {
	bytes := make([]byte, n)
	for i := range bytes {
		bytes[i] = letterBytes[rand.Int63()%int64(len(letterBytes))]
	}
	return string(bytes)
}
