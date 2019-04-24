package utils

import "strings"

func StringCut(str string, length int) string {
	rs := []rune(str)
	if len(rs) <= length {
		return str
	}
	return string(rs[:length])
}

func StringContains(o string, strList []string) (string, bool) {
	for _, str := range strList {
		if strings.Contains(o, str) {
			return str, true
		}
	}
	return "", false
}

func StringHasPrefixs(o string, strList []string) (string, bool) {
	for _, str := range strList {
		if strings.HasPrefix(o, str) {
			return str, true
		}
	}
	return "", false
}
