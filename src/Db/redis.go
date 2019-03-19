package db

import (
	"fmt"
)

func phoneNumberToKey(PhoneNumber string) string {
	return fmt.Sprintf("ptc_%s", PhoneNumber)
}

func replyToString(r interface{}) string {
	return string(r.([]uint8))
}

func getValueFromKey(key string) (string, error) {
	conn := redisPool.Get()
	defer func() {
		err := conn.Close()
		if err != nil {
			fmt.Printf("%s", err)
			return
		}
	}()

	r, err := conn.Do("GET", key)
	if err != nil {
		return "", err
	}
	if r != nil {
		return replyToString(r), nil
	}
	return "", nil
}

func setKeyToValue(key string, value string) error {
	conn := redisPool.Get()
	defer func() {
		err := conn.Close()
		if err != nil {
			fmt.Printf("%s", err)
			return
		}
	}()
	err := conn.Send("SET", key, value)
	return err
}

func setKeyToValueLimitTime(key string, value string, limitTime int64) error {
	conn := redisPool.Get()
	defer func() {
		err := conn.Close()
		if err != nil {
			fmt.Printf("%s", err)
			return
		}
	}()
	var err error
	err = conn.Send("MULTI")
	if err != nil {
		return err
	}
	err = conn.Send("SET", key, value)
	if err != nil {
		return err
	}
	err = conn.Send("EXPIRE", key, 180)
	if err != nil {
		return err
	}
	_, err = conn.Do("EXEC")
	return err
}

func CheckPhoneCodeCorrection(PhoneNumber string, Code string) (bool, error) {
	value, err := GetNowVerificationCode(PhoneNumber)
	if err != nil {
		return false, err
	}
	if value == Code {
		return true, nil
	}
	return false, nil
}

func GetNowVerificationCode(PhoneNumber string) (string, error) {
	return getValueFromKey(phoneNumberToKey(PhoneNumber))
}

func SetNewVerificationCode(PhoneNumber string, Code string) error {
	return setKeyToValueLimitTime(phoneNumberToKey(PhoneNumber), Code, 180)
}
