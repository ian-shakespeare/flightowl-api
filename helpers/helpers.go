package helpers

import (
	"crypto/rand"
	"crypto/sha256"
	"strconv"
	"strings"
	"time"
)

func Includes(str string, subStr string) bool {
	if len(subStr) > len(str) {
		return false
	}
	for i := range str {
		if len(subStr)+i > len(str) {
			break
		}
		if str[i:len(subStr)+i] == subStr {
			return true
		}
	}
	return false
}

func StartsWith(str string, subStr string) bool {
	if len(subStr) > len(str) {
		return false
	}
	return str[:len(subStr)] == subStr
}

func GetPathResource(path string) string {
	pathList := strings.Split(path, "/")
	if len(pathList) == 1 {
		return ""
	}
	return pathList[len(pathList)-1]
}

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
