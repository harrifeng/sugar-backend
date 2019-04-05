package db

import (
	"errors"
	"fmt"
	"github.com/jinzhu/gorm"
	"math/rand"
)

func CreateNewUser(phoneNumber string, userName string, password string) error {
	imageNumber := rand.Intn(35) + 10
	headPictureUrl := fmt.Sprintf("/static/userImg/usertile%d.jpg", imageNumber)
	user := User{
		PhoneNumber:     phoneNumber,
		Password:        password,
		UserName:        userName,
		HeadPortraitUrl: headPictureUrl,
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

func GetUserFromPhoneNumber(phoneNumber string) (User, error) {
	var userTmp User
	if mysqlDb.Where(&User{PhoneNumber: phoneNumber}).First(&userTmp).RecordNotFound() {
		return userTmp, errors.New("this user has not existed")
	}
	return userTmp, nil
}

func GetUserFromUserId(userId int) (User, error) {
	var userTmp User
	if mysqlDb.First(&userTmp, userId).RecordNotFound() {
		return userTmp, errors.New("this user has not existed")
	}
	return userTmp, nil
}

func GetPrivacySettingFromUserId(userId int) (UserPrivacySetting, error) {
	user, err := GetUserFromUserId(userId)
	if err != nil {
		return UserPrivacySetting{}, err
	}
	var privacySetting UserPrivacySetting
	err = mysqlDb.Model(&user).Association("UserPrivacySetting").Find(&privacySetting).Error
	return privacySetting, err
}

func AlterUserPrivacySettingFromUserId(userId int, showPhoneNumber bool, showGender bool,
	showAge bool, showHeight bool, showWeight bool, showArea bool, showJob bool) error {
	user, err := GetUserFromUserId(userId)
	if err != nil {
		return err
	}
	var privacySetting UserPrivacySetting
	err = mysqlDb.Model(&user).Association("UserPrivacySetting").Find(&privacySetting).Error
	if err != nil {
		return err
	}
	privacySetting.ShowPhoneNumber = showPhoneNumber
	privacySetting.ShowGender = showGender
	privacySetting.ShowAge = showAge
	privacySetting.ShowHeight = showHeight
	privacySetting.ShowWeight = showWeight
	privacySetting.ShowArea = showArea
	privacySetting.ShowJob = showJob
	err = mysqlDb.Model(&user).Association("UserPrivacySetting").Replace(&privacySetting).Error
	return err
}

func AlterUserPasswordFromPhoneNumber(phoneNumber string, newPassword string) error {
	var user User
	err := mysqlDb.Where(&User{PhoneNumber: phoneNumber}).First(&user).Error
	if err != nil {
		return err
	}
	user.Password = newPassword
	err = mysqlDb.Save(&user).Error
	return err
}

func AlterUserInformationFromUserId(userId int, userName string, gender string, height float64,
	weight float64, area string, job string, age int) error {
	var user User
	err := mysqlDb.First(&user, userId).Error
	if err != nil {
		return err
	}
	user.UserName = userName
	user.Gender = gender
	user.Height = height
	user.Weight = weight
	user.Area = area
	user.Job = job
	user.Age = age
	err = mysqlDb.Save(&user).Error
	return err
}

func AddUserFollowing(userId int, targetUserId int) error {
	return Transaction(func(db *gorm.DB) error {
		var user1, user2 User
		db.Preload("FollowingUsers").First(&user1, userId)
		user1To, err := GetUserFromUserId(targetUserId)
		if err != nil {
			return err
		}
		db.Model(&user1).Association("FollowingUsers").Append(&user1To)

		db.Preload("FollowerUsers").First(&user2, targetUserId)
		user2To, err := GetUserFromUserId(userId)
		if err != nil {
			return err
		}
		db.Model(&user2).Association("FollowerUsers").Append(&user2To)
		return nil
	})
}

func RemoveUserFollowing(userId int, targetUserId int) error {
	return Transaction(func(db *gorm.DB) error {
		var user1, user2 User
		db.Preload("FollowingUsers").First(&user1, userId)
		user1To, err := GetUserFromUserId(targetUserId)
		if err != nil {
			return err
		}
		db.Model(&user1).Association("FollowingUsers").Delete(&user1To)

		db.Preload("FollowerUsers").First(&user2, targetUserId)
		user2To, err := GetUserFromUserId(userId)
		if err != nil {
			return err
		}
		db.Model(&user2).Association("FollowerUsers").Delete(&user2To)
		return nil
	})
}

func GetUserFollowerList(userId int, beginId int, needNumber int) ([]User, int, error) {
	var user User
	mysqlDb.Preload("FollowerUsers").First(&user, userId)
	var users []User
	err := mysqlDb.Model(&user).Offset(beginId).Limit(needNumber).Related(&users, "FollowerUsers").Error
	count := mysqlDb.Model(&user).Association("FollowerUsers").Count()
	return users, count, err
}

func GetUserFollowingList(userId int, beginId int, needNumber int) ([]User, int, error) {
	var user User
	mysqlDb.Preload("FollowingUsers").First(&user, userId)
	var users []User
	err := mysqlDb.Model(&user).Offset(beginId).Limit(needNumber).Related(&users, "FollowingUsers").Error
	count := mysqlDb.Model(&user).Association("FollowingUsers").Count()
	return users, count, err
}

func CheckUserFollowingOtherUser(userId int, otherUserId int) (bool, error) {
	var recordCount int
	err := mysqlDb.Table("user_following_ships").
		Where("user_id = ? and following_user_id = ?", userId, otherUserId).Count(&recordCount).Error
	return recordCount > 0, err
}
