package util

import (
	"math/rand"
	"strings"
	"time"
)

const alphabet = "abcdefghijklmnopqrstuvwxyz"

func init()  {
	rand.Seed(time.Now().UnixNano())
}

// RandomInt generates a random integer btwn  min and max
func RandomInt(min, max int64) float64 {
	return float64(min + rand.Int63n(max-min+1)) // Random integer btwn min and max
}

// RandomString returns a random string of length n
func RandomString(n int) string  {
	var sb strings.Builder
	k := len(alphabet)

	for i := 0; i < n; i++ {
		c := alphabet[rand.Intn(k)]
		sb.WriteByte(c)
	}
	return sb.String()
}

// RandomOwner generates a random owner
func RandomOwner()string  {
	return RandomString(8)
}

// RandomMoney generates random amount of money
func RandomMoney() float64 {
	return RandomInt(0, 2000)
}

func RandomCurrency() string {
	currencies := []string{"USD", "KES", "UGSH", "TSH", "EUR"}
	n := len(currencies)
	return  currencies[rand.Intn(n)]
}