package helpers

import (
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

func GetPathCollection(path string) string {
	pathList := strings.Split(path, "/")
	if len(pathList) == 1 {
		return ""
	}
	return pathList[len(pathList)-1]
}

func GetFormattedTime(t time.Time) string {
	return strings.Join(strings.Split(t.String(), " ")[:2], " ")
}
