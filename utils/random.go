package utils

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
)

const alphabet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func init() {
	rand.Seed(time.Now().UnixNano())
	
}

func RandomInt(min, max int64) int64 {
	return min + rand.Int63n(max-min+1)
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

func RandomEmail() string {
    return fmt.Sprintf("%s@example.com", RandomString(10))
}


func RandomDescription() string {
	return RandomString(1000)
}

func RandomPrice() string {
    price := float64(RandomInt(1, 1000)) + float64(RandomInt(0, 99))/100
    return fmt.Sprintf("%.2f", price)
}

func RandomPhoneNumber(countryCode string) string {
    return fmt.Sprintf("+%s%d", countryCode, RandomInt(7000000000, 9999999999))
}

