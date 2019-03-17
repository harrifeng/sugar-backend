package db

import "fmt"

func SetNewVerificationCode(PhoneNumber string, Code string) error {
	conn := redisPool.Get()
	_ = conn.Send("MULTI")
	key := fmt.Sprintf("ptc_%s", PhoneNumber)
	_ = conn.Send("SET", key, Code)
	_ = conn.Send("EXPIRE", key, 180)
	_, err := conn.Do("EXEC")
	return err
}
