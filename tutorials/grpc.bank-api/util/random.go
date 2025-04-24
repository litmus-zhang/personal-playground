package util

import (
	"math/rand"
	"time"
)

func init() {
	rand.NewSource(time.Now().UnixNano())
}

// RandomString returns a random string of n characters
func RandomString(n int) string {
	var letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}
