package utils

import (
	"math/rand"
	"time"
)

var lowerChars = "0123456789abcdefghijklmnopqrstuvwxyz"
var tokenLength = 10

func GenTokenSufix() string {
	rand.Seed(time.Now().UTC().UnixNano())
	b := make([]byte, tokenLength)
	for i := range b {
		b[i] = lowerChars[rand.Intn(len(lowerChars))]
	}
	return string(b)
}
