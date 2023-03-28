package helpers

import (
	"crypto/rand"
	"crypto/sha256"
	"os"
	"strconv"
	"strings"
	"time"
)

func GetFormattedTime(t time.Time) string {
	return strings.Join(strings.Split(t.String(), " ")[:2], " ")
}

func HashString(s string) string {
	hsha256 := sha256.Sum256([]byte(s))
	newStr := strconv.QuoteToASCII(string(hsha256[:]))
	return newStr
}

func GetRandomString(length uint8) string {
	letters := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	buf := make([]byte, length)
	rand.Read(buf)
	s := make([]rune, length)
	for i := range s {
		s[i] = letters[uint8(buf[i])%length]
	}
	return string(s)
}

func GetRequiredEnv(k string) string {
	v := os.Getenv(k)
	if v == "" {
		panic("missing required environment variable")
	}
	return v
}
