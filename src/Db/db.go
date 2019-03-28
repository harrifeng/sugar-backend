package db

import (
	"fmt"
	"github.com/gomodule/redigo/redis"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"time"
)

var mysqlDb *gorm.DB
var redisPool *redis.Pool
var mysqlConfig *MysqlConfiguration
var redisConfig *RedisConfiguration

func mysqlLinkString() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		mysqlConfig.User, mysqlConfig.Password, mysqlConfig.Host, mysqlConfig.Port, mysqlConfig.DatabaseName)
}

func InitConfiguration(MysqlConfig *MysqlConfiguration, RedisConfig *RedisConfiguration) {
	mysqlConfig = MysqlConfig
	redisConfig = RedisConfig
}

func InitRedis() *redis.Pool {
	redisHost := fmt.Sprintf("%s:%d", redisConfig.Host, redisConfig.Port)
	redisPool = &redis.Pool{
		MaxIdle:     3,
		MaxActive:   500,
		Wait:        true,
		IdleTimeout: 240 * time.Second,
		Dial:        func() (redis.Conn, error) { return redis.Dial("tcp", redisHost) },
	}

	return redisPool
}

func InitMysql() (*gorm.DB, error) {
	var err error
	mysqlDb, err = gorm.Open("mysql", mysqlLinkString())
	return mysqlDb, err
}
