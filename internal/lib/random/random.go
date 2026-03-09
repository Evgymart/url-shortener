package random

import (
	"math/rand"
	"time"
)

func NewRandomString(size int) string {
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	alphabet := []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789")
	runes := make([]rune, size)
	for i := range runes {
		runes[i] = alphabet[rnd.Intn(len(alphabet))]
	}

	return string(runes)
}
