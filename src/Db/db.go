package db

import (
	"errors"
	"fmt"
	"github.com/gomodule/redigo/redis"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"strconv"
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

func CreateNewUser(PhoneNumber string, Password string) error {
	user := User{
		PhoneNumber: PhoneNumber,
		Password:    Password,
		UserPrivacySetting: UserPrivacySetting{
			ShowPhoneNumber: true,
			ShowGender:      true,
			ShowAge:         true,
			ShowHeight:      true,
			ShowWeight:      true,
			ShowArea:        true,
			ShowJob:         true,
		},
	}
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

func GetUserFromUserId(UserId string) (User, error) {
	var userTmp User
	userId, _ := strconv.Atoi(UserId)
	if mysqlDb.First(&userTmp, userId).RecordNotFound() {
		return userTmp, errors.New("this user has not existed")
	}
	return userTmp, nil
}

func GetPrivacySettingFromUserId(UserId string) (UserPrivacySetting, error) {
	user, err := GetUserFromUserId(UserId)
	if err != nil {
		return UserPrivacySetting{}, err
	}
	var privacySetting UserPrivacySetting
	mysqlDb.Model(&user).Association("UserPrivacySetting").Find(&privacySetting)
	return privacySetting, nil
}

func AlterUserPrivacySettingFromUserId(UserId string, ShowPhoneNumber bool, ShowGender bool,
	ShowAge bool, ShowHeight bool, ShowWeight bool, ShowArea bool, ShowJob bool) error {
	user, err := GetUserFromUserId(UserId)
	if err != nil {
		return err
	}
	var privacySetting UserPrivacySetting
	mysqlDb.Model(&user).Association("UserPrivacySetting").Find(&privacySetting)
	privacySetting.ShowPhoneNumber = ShowPhoneNumber
	privacySetting.ShowGender = ShowGender
	privacySetting.ShowAge = ShowAge
	privacySetting.ShowHeight = ShowHeight
	privacySetting.ShowWeight = ShowWeight
	privacySetting.ShowArea = ShowArea
	privacySetting.ShowJob = ShowJob
	mysqlDb.Model(&user).Association("UserPrivacySetting").Replace(&privacySetting)
	return nil
}

func AlterUserPasswordFromPhoneNumber(PhoneNumber string, NewPassword string) error {
	var user User
	err := mysqlDb.Where(&User{PhoneNumber: PhoneNumber}).First(&user).Error
	if err != nil {
		return err
	}
	user.Password = NewPassword
	mysqlDb.Save(&user)
	return nil
}

func AlterUserInformationFromUserId(UserId string, UserName string, Gender string, Height float64,
	Weight float64, Area string, Job string, Age int) error {
	userId, _ := strconv.Atoi(UserId)
	var user User
	err := mysqlDb.First(&user, userId).Error
	if err != nil {
		return err
	}
	user.UserName = UserName
	user.Gender = Gender
	user.Height = Height
	user.Weight = Weight
	user.Area = Area
	user.Job = Job
	user.Age = Age
	mysqlDb.Save(&user)
	return nil
}

func AddUserFollowing(UserId string, TargetUserId string) error {
	var user User
	mysqlDb.Preload("FollowUsers").First(&user, UserId)
	userTo, err := GetUserFromUserId(TargetUserId)
	if err != nil {
		return err
	}
	mysqlDb.Model(&user).Association("FollowUsers").Append(&userTo)
	return nil
}

func RemoveUserFollowing(UserId string, TargetUserId string) error {
	var user User
	mysqlDb.Preload("FollowUsers").First(&user, UserId)
	userTo, err := GetUserFromUserId(TargetUserId)
	if err != nil {
		return err
	}
	mysqlDb.Model(&user).Association("FollowUsers").Delete(&userTo)
	return nil
}

func GetUserFollowerList(UserId string, BeginId string, NeedNumber string) error {
	//var user User
	//mysqlDb.Preload("FollowUsers").First(&user,UserId)
	//var users []User
	//beginId,_ :=strconv.Atoi(BeginId)
	//needNumber,_:=strconv.Atoi(NeedNumber)
	//mysqlDb.Model(&user).Related(&users,"FollowUsers").Offset(beginId).Limit(needNumber).Find(&users)
	//fmt.Println(users)
	return nil
}

func AddNewArticle(Title string, Content string) error {
	article := Article{
		Title:   Title,
		Content: Content,
	}
	mysqlDb.Create(&article)
	mysqlDb.Save(&article)
	return nil
}
