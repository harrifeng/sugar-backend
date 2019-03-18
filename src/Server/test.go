package server

import "fmt"

func SendMessageTest(PhoneNumber string) {
	err := SendVerificationCode(PhoneNumber)
	if err != nil {
		fmt.Println(err)
	}
}
