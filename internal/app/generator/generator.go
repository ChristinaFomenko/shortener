package generator

import (
	"math/rand"
	"time"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

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
