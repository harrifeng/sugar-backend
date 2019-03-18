package db

import (
	"fmt"
	"github.com/gomodule/redigo/redis"
)

func checkKeyExistFromConn(conn redis.Conn, key string) (bool,error) {
	r,err:=conn.Do("EXISTS",key)
	if err !=nil{
		return false,err
	}
	return r.(int64)==1,err
}

func checkKeyExist(key string) (bool,error){
	conn := redisPool.Get()
	defer func() {
		err:=conn.Close()
		if err!=nil{
			fmt.Printf("%s",err)
			return
		}
	}()

	return checkKeyExistFromConn(conn,key)
}

func CheckPhoneCodeExist(PhoneNumber string) (bool,error){
	return checkKeyExist(fmt.Sprintf("ptc_%s", PhoneNumber))
}

func SetNewVerificationCode(PhoneNumber string, Code string) error {
	conn := redisPool.Get()
	defer func() {
		err:=conn.Close()
		if err!=nil{
			fmt.Printf("%s",err)
			return
		}
	}()
	key := fmt.Sprintf("ptc_%s", PhoneNumber)

	_ = conn.Send("MULTI")
	_ = conn.Send("SET", key, Code)
	_ = conn.Send("EXPIRE", key, 180)
	_, err := conn.Do("EXEC")
	return err
}
