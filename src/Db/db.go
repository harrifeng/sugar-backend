package Db

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var db *gorm.DB

func mysqlLinkString() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		MysqlUserName, MysqlPassword, MysqlIp, MysqlPort, MysqlBbName)
}

func Init() (*gorm.DB, error) {
	var err error
	db, err = gorm.Open("mysql", mysqlLinkString())
	return db, err
}
