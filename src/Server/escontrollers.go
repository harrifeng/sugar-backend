package server

import (
	"db"
	"utils"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

func SendVerificationCode(PhoneNumber string) error {
	exists,err:=db.CheckPhoneCodeExist(PhoneNumber)
	if err!=nil{
		return err
	}
	if exists {
		return errors.New("phone number key has exists")
	}
	code := utils.RandCode()
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
