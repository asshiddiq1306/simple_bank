package util

import (
	"math/rand"
	"strings"
	"time"
)

const alphabet = "abcdefghijklmnopqrstuvwxyz"

func init() {
	rand.Seed(time.Now().UnixNano())
}

func RandomNumber(min, max int64) int64 {
	return min + rand.Int63n(max-min+1)
}

func RandomMoney() int64 {
	return RandomNumber(1, 1000)
}

func RandomString(n int) string {
	var sb strings.Builder
	k := len(alphabet)
	for i := 0; i < n; i++ {
		c := alphabet[rand.Intn(k)]
		sb.WriteByte(c)
	}
	return sb.String()
}

func RandomName() string {
	return RandomString(6)
}

func RandomCurrency() string {
	currency := []string{"USD", "IDR", "EUR"}
	k := len(currency)
	return currency[rand.Intn(k)]
}
