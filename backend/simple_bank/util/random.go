package util

import (
	"math/rand"
	"strings"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

// RandomInt generates a random number between min and max
const alphabetNum = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func RandomInt(min, max int64) int64 {
	return min + rand.Int63n(max-min+1)
}
func RandomString(n int) string {
	var sb strings.Builder
	k := len(alphabetNum)

	for i := 0; i < n; i++ {
		c := alphabetNum[rand.Intn(k)]
		sb.WriteByte(c)
	}
	return sb.String()
}

func RandomOwner() string {
	return RandomString(8)
}

func RandomMoney(min int64) int64 {
	if min <= 0 {
		min = 1000
		return RandomInt(0, 1000)
	}
	return RandomInt(0, min)
}

func RandomCurrency() string {
	currencies := []string{USD, EUR, CAD}
	n := len(currencies)
	return currencies[rand.Intn(n)]
}
