package utils

import (
	"math/rand"
	"strconv"
)

func RandPhoneNumber() string {
	phoneNumber := ""
	for i := 0; i < 11; i++ {
		phoneNumber += strconv.Itoa(rand.Intn(10))
	}
	return phoneNumber
}
