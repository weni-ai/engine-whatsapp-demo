package utils

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
)

const lowerChars = "0123456789abcdefghijklmnopqrstuvwxyz"
const sufixLength = 10

const crumb = "whatsapp-demo"

func genTokenSufix() string {
	rand.Seed(time.Now().UTC().UnixNano())
	b := make([]byte, sufixLength)
	for i := range b {
		b[i] = lowerChars[rand.Intn(len(lowerChars))]
	}
	return string(b)
}

func GenToken(prefixName string) string {
	formatedName := strings.ToLower(strings.ReplaceAll(prefixName, " ", ""))
	sufix := genTokenSufix()
	return fmt.Sprintf("%s-%s-%s", formatedName, crumb, sufix)
}
