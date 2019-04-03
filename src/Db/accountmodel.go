package db

import (
	"errors"
	"fmt"
	"math/rand"
	"strconv"
)

func CreateNewUser(PhoneNumber string, UserName string, Password string) error {
	imageNumber := rand.Intn(35) + 10
	HeadPictureUrl := fmt.Sprintf("/static/userImg/usertile%d.jpg", imageNumber)
	user := User{
		PhoneNumber:     PhoneNumber,
		Password:        Password,
		UserName:        UserName,
		HeadPortraitUrl: HeadPictureUrl,
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
	err = mysqlDb.Model(&user).Association("UserPrivacySetting").Find(&privacySetting).Error
	return privacySetting, err
}

func AlterUserPrivacySettingFromUserId(UserId string, ShowPhoneNumber bool, ShowGender bool,
	ShowAge bool, ShowHeight bool, ShowWeight bool, ShowArea bool, ShowJob bool) error {
	user, err := GetUserFromUserId(UserId)
	if err != nil {
		return err
	}
	tx := mysqlDb.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	var privacySetting UserPrivacySetting
	err = tx.Model(&user).Association("UserPrivacySetting").Find(&privacySetting).Error
	if err != nil {
		return err
	}
	privacySetting.ShowPhoneNumber = ShowPhoneNumber
	privacySetting.ShowGender = ShowGender
	privacySetting.ShowAge = ShowAge
	privacySetting.ShowHeight = ShowHeight
	privacySetting.ShowWeight = ShowWeight
	privacySetting.ShowArea = ShowArea
	privacySetting.ShowJob = ShowJob
	err = tx.Model(&user).Association("UserPrivacySetting").Replace(&privacySetting).Error
	if err != nil {
		return err
	}
	return tx.Commit().Error
}

func AlterUserPasswordFromPhoneNumber(PhoneNumber string, NewPassword string) error {
	var user User
	err := mysqlDb.Where(&User{PhoneNumber: PhoneNumber}).First(&user).Error
	if err != nil {
		return err
	}
	user.Password = NewPassword
	err = mysqlDb.Save(&user).Error
	return err
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
	err = mysqlDb.Save(&user).Error
	return err
}

func AddUserFollowing(UserId string, TargetUserId string) error {
	var user1, user2 User
	tx := mysqlDb.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	tx.Preload("FollowingUsers").First(&user1, UserId)
	user1To, err := GetUserFromUserId(TargetUserId)
	if err != nil {
		return err
	}
	tx.Model(&user1).Association("FollowingUsers").Append(&user1To)

	tx.Preload("FollowerUsers").First(&user2, TargetUserId)
	user2To, err := GetUserFromUserId(UserId)
	if err != nil {
		return err
	}
	tx.Model(&user2).Association("FollowerUsers").Append(&user2To)
	return tx.Commit().Error
}

func RemoveUserFollowing(UserId string, TargetUserId string) error {
	var user1, user2 User
	tx := mysqlDb.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	tx.Preload("FollowingUsers").First(&user1, UserId)
	user1To, err := GetUserFromUserId(TargetUserId)
	if err != nil {
		return err
	}
	tx.Model(&user1).Association("FollowingUsers").Delete(&user1To)

	tx.Preload("FollowerUsers").First(&user2, TargetUserId)
	user2To, err := GetUserFromUserId(UserId)
	if err != nil {
		return err
	}
	tx.Model(&user2).Association("FollowerUsers").Delete(&user2To)
	err = tx.Commit().Error
	return err
}

func GetUserFollowerList(UserId string, BeginId string, NeedNumber string) ([]User, int, error) {
	var user User
	mysqlDb.Preload("FollowerUsers").First(&user, UserId)
	var users []User
	beginId, _ := strconv.Atoi(BeginId)
	needNumber, _ := strconv.Atoi(NeedNumber)
	err := mysqlDb.Model(&user).Offset(beginId).Limit(needNumber).Related(&users, "FollowerUsers").Error
	count := mysqlDb.Model(&user).Association("FollowerUsers").Count()
	return users, count, err
}

func GetUserFollowingList(UserId string, BeginId string, NeedNumber string) ([]User, int, error) {
	var user User
	mysqlDb.Preload("FollowingUsers").First(&user, UserId)
	var users []User
	beginId, _ := strconv.Atoi(BeginId)
	needNumber, _ := strconv.Atoi(NeedNumber)
	err := mysqlDb.Model(&user).Offset(beginId).Limit(needNumber).Related(&users, "FollowingUsers").Error
	count := mysqlDb.Model(&user).Association("FollowingUsers").Count()
	return users, count, err
}

func CheckUserFollowingOtherUser(UserId string,OtherUserId string)(bool,error){
	var recordCount int
	err :=mysqlDb.Table("user_following_ships").
		Where("user_id = ? and following_user_id = ?",UserId,OtherUserId).Count(&recordCount).Error
	if err!=nil{
		return false,err
	}
	return recordCount > 0,nil
}