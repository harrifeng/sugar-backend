package server

import (
	"db"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"utils"
)

func SendVerificationCode(PhoneNumber string) error {
	nowCode, err := db.GetNowVerificationCode(PhoneNumber)
	if err != nil {
		return err
	}
	var code string
	if nowCode != "" {
		code = nowCode
	} else {
		code = utils.RandCode()
	}
	url := fmt.Sprintf("http://127.0.0.1:7799/send_message?phone_number=%s&code=%s",
		PhoneNumber, code)
	resp, err := http.Get(url)

	if err != nil {
		fmt.Printf("%s", err)
		return err
	}
	defer func() {
		err := resp.Body.Close()
		if err != nil {
			fmt.Printf("%s", err)
			return
		}
	}()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("%s", err)
		return err
	}
	if body[0] == 'o' {
		return db.SetNewVerificationCode(PhoneNumber, code)
	}
	return errors.New("external service(send message) error")
}

func CheckVerificationCodeCorrection(PhoneNumber string, Code string) (bool, error) {
	return db.CheckPhoneCodeCorrection(PhoneNumber, Code)
}
