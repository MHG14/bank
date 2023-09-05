package util

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
)

const alphabet = "abcdefghijklmnopqrstuvwxyz"

func init() {
	rand.New(rand.NewSource(time.Now().UnixNano()))
}

// generates a random int64 between min and max
func RandomInt(min, max int64) int64 {
	return min + rand.Int63n(max-min+1)
}

// generates a random string of length n
func RandomString(n int) string {
	var sb strings.Builder
	k := len(alphabet)

	for i := 0; i < n; i++ {
		c := alphabet[rand.Intn(k)]
		sb.WriteByte(c)
	}

	return sb.String()
}

// generates a random owner name
func RandomOwner() string {
	return RandomString(6)
}

// generates a random amount of money
func RandomMoney() int64 {
	return RandomInt(0, 1000)
}

// generates a random currency code
func RandomCurrency() string {
	currencies := []string{
		USD, CAD, EUR,
	}

	n := len(currencies)
	return currencies[rand.Intn(n)]
}

// generates a random email address
func RandomEmail() string {
	return fmt.Sprintf("%s@email.com", RandomString(6))
}
