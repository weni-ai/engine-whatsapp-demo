package utils

import (
	"fmt"
	"math/rand"
	"time"
)

const chars = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ-_"
const sufixLength = 10

const crumb = "weni-demo"

func genTokenSufix() string {
	rand.Seed(time.Now().UTC().UnixNano())
	b := make([]byte, sufixLength)
	for i := range b {
		b[i] = chars[rand.Intn(len(chars))]
	}
	return string(b)
}

func GenToken() string {
	sufix := genTokenSufix()
	return fmt.Sprintf("%s-%s", crumb, sufix)
}
