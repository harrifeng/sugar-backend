package db

import "time"

func CheckInUser(userId int)error{
	userCheckIn:=UserCheckIn{
		CheckTime:time.Now(),
		UserID:uint(userId),
	}
	mysqlDb.Create(&userCheckIn)
	mysqlDb.Save(&userCheckIn)
	return nil
}

func CheckUserCheckIn(userId int)(bool,error){
	checkIn := !mysqlDb.First(&UserCheckIn{UserID:uint(userId),CheckTime:time.Now()}).RecordNotFound()
	return checkIn,nil
}

func GetUserFamilyMemberList(userId int)([]FamilyMember,error){
	var familyMembers []FamilyMember
	user,err:=GetUserFromUserId(userId)
	if err!=nil{
		return familyMembers,err
	}
	err = mysqlDb.Model(&user).Association("FamilyMembers").Find(&familyMembers).Error
	return familyMembers,err
}

func AddFamilyMember(userId int,callName string,phoneNumber string)error{
	user,err:=GetUserFromUserId(userId)
	if err!=nil{
		return err
	}
	familyMember:=FamilyMember{
		CallName:callName,
		PhoneNumber:phoneNumber,
	}
	err = mysqlDb.Model(&user).Association("FamilyMembers").
		Append(&familyMember).Error
	return err
}