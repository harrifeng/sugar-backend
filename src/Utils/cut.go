package utils

func StringCut(str string, length int) string {
	if len(str) <= length {
		return str
	}
	return str[:length]
}
