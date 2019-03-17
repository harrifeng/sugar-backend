package server

func SendMessageTest(PhoneNumber string) {
	err := AddVerificationCode("18061532353")
	if err != nil {
		return
	}
}
