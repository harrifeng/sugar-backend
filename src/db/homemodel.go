package db

import (
	"encoding/json"
	"fmt"
	"time"
)

func CheckInUser(userId int) error {
	userCheckIn := UserCheckIn{
		CheckTime: time.Now(),
		UserID:    uint(userId),
	}
	mysqlDb.Create(&userCheckIn)
	mysqlDb.Save(&userCheckIn)
	return nil
}

func CheckUserCheckIn(userId int) (bool, error) {
	checkIn := !mysqlDb.First(&UserCheckIn{UserID: uint(userId), CheckTime: time.Now()}).RecordNotFound()
	return checkIn, nil
}

func GetUserFamilyMemberList(userId int) ([]FamilyMember, error) {
	var familyMembers []FamilyMember
	user, err := GetUserFromUserId(userId)
	if err != nil {
		return familyMembers, err
	}
	err = mysqlDb.Model(&user).Association("FamilyMembers").Find(&familyMembers).Error
	return familyMembers, err
}

func AddFamilyMember(userId int, callName string, phoneNumber string) error {
	user, err := GetUserFromUserId(userId)
	if err != nil {
		return err
	}
	familyMember := FamilyMember{
		CallName:    callName,
		PhoneNumber: phoneNumber,
	}
	err = mysqlDb.Model(&user).Association("FamilyMembers").
		Append(&familyMember).Error
	return err
}

func GetBloodSugarRecordFromRecordDate(userId int, recordDate time.Time) (BloodRecord, bool, error) {
	var bloodRecord BloodRecord
	year, month, day := recordDate.Date()
	exist := !mysqlDb.Where("record_date=? and user_id=?", fmt.Sprintf("%d-%d-%d", year, month, day), userId).
		First(&bloodRecord).RecordNotFound()
	return bloodRecord, exist, nil
}

func AddBloodSugarRecord(userId int, period string, level string, recordTime string, recordDate time.Time) error {
	bloodRecord, exist, err := GetBloodSugarRecordFromRecordDate(userId, recordDate)
	if err != nil {
		return err
	}
	if !exist {
		levelJson := map[string]string{
			"0": "0", "1": "0", "2": "0", "3": "0", "4": "0", "5": "0", "6": "0",
		}
		timeJson := map[string]string{
			"0": "0", "1": "0", "2": "0", "3": "0", "4": "0", "5": "0", "6": "0",
		}
		levelJson[period] = level
		timeJson[period] = recordTime
		levelBytes, err := json.Marshal(levelJson)
		if err != nil {
			return err
		}
		timeBytes, err := json.Marshal(timeJson)
		if err != nil {
			return err
		}
		bloodRecord = BloodRecord{
			UserID:     uint(userId),
			Level:      string(levelBytes),
			RecordTime: string(timeBytes),
			RecordDate: recordDate,
		}
		mysqlDb.Create(&bloodRecord)
		mysqlDb.Save(&bloodRecord)
		return nil
	}
	levelJson := make(map[string]string)
	timeJson := make(map[string]string)
	err = json.Unmarshal([]byte(bloodRecord.Level), &levelJson)
	if err != nil {
		return err
	}
	err = json.Unmarshal([]byte(bloodRecord.RecordTime), &timeJson)
	if err != nil {
		return err
	}
	levelJson[period] = level
	timeJson[period] = recordTime
	levelBytes, err := json.Marshal(levelJson)
	if err != nil {
		return err
	}
	timeBytes, err := json.Marshal(timeJson)
	if err != nil {
		return err
	}
	bloodRecord.Level = string(levelBytes)
	bloodRecord.RecordTime = string(timeBytes)
	err = mysqlDb.Save(&bloodRecord).Error
	return nil
}

func GetHealthRecordFromRecordDate(userId int, recordDate time.Time) (HealthRecord, bool, error) {
	var healthRecord HealthRecord
	year, month, day := recordDate.Date()
	exist := !mysqlDb.Where("record_date=? and user_id=?", fmt.Sprintf("%d-%d-%d", year, month, day), userId).
		First(&healthRecord).RecordNotFound()
	return healthRecord, exist, nil
}

func AddHealthRecord(userId int, insulin string, sportTime string, weight string, bloodPressure string,
	recordTime string, recordDate time.Time) error {
	healthRecord, exist, err := GetHealthRecordFromRecordDate(userId, recordDate)
	if err != nil {
		return err
	}
	if exist {
		healthRecord.Insulin = insulin
		healthRecord.SportTime = sportTime
		healthRecord.Weight = weight
		healthRecord.BloodPressure = bloodPressure
		healthRecord.RecordTime = recordTime
		return mysqlDb.Save(&healthRecord).Error
	}
	healthRecord = HealthRecord{
		UserID:        uint(userId),
		Insulin:       insulin,
		SportTime:     sportTime,
		Weight:        weight,
		BloodPressure: bloodPressure,
		RecordTime:    recordTime,
		RecordDate:    recordDate,
	}
	mysqlDb.Create(&healthRecord)
	mysqlDb.Save(&healthRecord)
	return nil
}

func GetHealthRecordList(userId int, beginId int, needNumber int) ([]HealthRecord, error) {
	var healthRecords []HealthRecord
	err := mysqlDb.Where("user_id=?", userId).Order("record_date desc").
		Offset(beginId).Limit(needNumber).Find(&healthRecords).Error
	return healthRecords, err
}

func GetBloodSugarRecordList(userId int, beginId int, needNumber int) ([]BloodRecord, error) {
	var bloodRecords []BloodRecord
	err := mysqlDb.Where("user_id=?", userId).Order("record_date desc").
		Offset(beginId).Limit(needNumber).Find(&bloodRecords).Error
	return bloodRecords, err
}
