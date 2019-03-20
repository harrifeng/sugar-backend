package db

import (
	"errors"
	"fmt"
	"github.com/gomodule/redigo/redis"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"time"
)

var mysqlDb *gorm.DB
var redisPool *redis.Pool

func mysqlLinkString() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		MysqlUserName, MysqlPassword, MysqlHost, MysqlPort, MysqlBbName)
}

func InitRedis() *redis.Pool {
	redisPool = &redis.Pool{
		MaxIdle:     3,
		MaxActive:   500,
		Wait:        true,
		IdleTimeout: 240 * time.Second,
		Dial:        func() (redis.Conn, error) { return redis.Dial("tcp", RedisHost) },
	}
	return redisPool
}

func InitMysql() (*gorm.DB, error) {
	var err error
	mysqlDb, err = gorm.Open("mysql", mysqlLinkString())
	return mysqlDb, err
}

func CreateNewUser(user User) error {
	var userTmp User
	if mysqlDb.Where(&User{PhoneNumber: user.PhoneNumber}).First(&userTmp).RecordNotFound() {
		mysqlDb.Create(&user)
		mysqlDb.Save(&user)
		return nil
	}
	return errors.New("this user has existed")
}

func GetUserFromPhoneNumber(PhoneNumber string) (User, error) {
	var userTmp User
	if mysqlDb.Where(&User{PhoneNumber: PhoneNumber}).First(&userTmp).RecordNotFound() {
		return userTmp, errors.New("this user has not existed")
	}
	return userTmp, nil

}
