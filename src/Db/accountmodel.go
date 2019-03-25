package db

import (
	"errors"
	"strconv"
)

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
	var user1, user2 User
	tx := mysqlDb.Begin()
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
	tx.Commit()
	return nil
}

func RemoveUserFollowing(UserId string, TargetUserId string) error {
	var user User
	mysqlDb.Preload("FollowingUsers").First(&user, UserId)
	userTo, err := GetUserFromUserId(TargetUserId)
	if err != nil {
		return err
	}
	mysqlDb.Model(&user).Association("FollowingUsers").Delete(&userTo)

	mysqlDb.Preload("FollowerUsers").First(&user, TargetUserId)
	userTo, err = GetUserFromUserId(UserId)
	if err != nil {
		return err
	}
	mysqlDb.Model(&user).Association("FollowerUsers").Delete(&userTo)
	return nil
}

// WARNING
func GetUserFollowerList(UserId string, BeginId string, NeedNumber string) ([]User, error) {
	var user User
	mysqlDb.Preload("FollowerUsers").First(&user, UserId)
	var users []User
	beginId, _ := strconv.Atoi(BeginId)
	needNumber, _ := strconv.Atoi(NeedNumber)
	mysqlDb.Model(&user).Related(&users, "FollowerUsers").Offset(beginId).
		Limit(needNumber).Find(&users)
	return users, nil
}

// WARNING
func GetUserFollowingList(UserId string, BeginId string, NeedNumber string) ([]User, error) {
	var user User
	mysqlDb.Preload("FollowingUsers").First(&user, UserId)
	var users []User
	beginId, _ := strconv.Atoi(BeginId)
	needNumber, _ := strconv.Atoi(NeedNumber)
	mysqlDb.Model(&user).Related(&users, "FollowingUsers").Offset(beginId).
		Limit(needNumber).Find(&users)
	return users, nil
}
