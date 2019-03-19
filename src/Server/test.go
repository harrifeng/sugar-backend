package server

import "fmt"

func SendMessageTest(PhoneNumber string) {
	err := SendVerificationCode(PhoneNumber)
	if err != nil {
		fmt.Println(err)
	}
}

func CheckVerificationCodeTest() {
	PhoneNumber := "18061532353"
	SendMessageTest(PhoneNumber)
	correct, err := CheckVerificationCodeCorrection(PhoneNumber, "498081")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(correct)
}
