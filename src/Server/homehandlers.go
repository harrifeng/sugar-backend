package server

import (
	"db"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"strconv"
	"time"
	"utils"
)

func checkinUser(userId int) responseBody {
	err := db.CheckInUser(userId)
	if err != nil {
		return responseInternalServerError(err)
	}
	return responseOK()
}

func getUserFamilyMemberList(userId int) responseBody {
	familyMembers, err := db.GetUserFamilyMemberList(userId)
	if err != nil {
		return responseInternalServerError(err)
	}
	respFamilyMembers := make([]gin.H, len(familyMembers))
	for i, member := range familyMembers {
		respFamilyMembers[i] = gin.H{
			"nickname": member.CallName,
			"tel":      member.PhoneNumber,
			"familyId": member.ID,
		}
	}
	return responseOKWithData(respFamilyMembers)
}

func linkNewFamilyMember(userId int, callName string, phoneNumber string, code string) responseBody {
	nowCode, err := db.GetNowVerificationCode(phoneNumber)
	if err != nil {
		return responseInternalServerError(err)
	}
	if nowCode != code {
		return responseNormalError("验证码错误")
	}
	err = db.AddFamilyMember(userId, callName, phoneNumber)
	if err != nil {
		return responseInternalServerError(err)
	}
	return responseOK()
}

func recordBloodSugar(userId int, bloodSugarValue string, period string, recordTime string, recordDate time.Time) responseBody {
	periodMap := map[string]string{
		"beforeBF": "0",
		"afterBF":  "1",
		"beforeLC": "2",
		"afterLC":  "3",
		"beforeDN": "4",
		"afterDN":  "5",
		"beforeSP": "6",
	}
	err := db.AddBloodSugarRecord(userId, periodMap[period], bloodSugarValue, recordTime, recordDate)
	if err != nil {
		return responseInternalServerError(err)
	}
	return responseOK()
}

func recordHealth(userId int, insulin string, sportTime string, weight string, bloodPressure string,
	recordTime string, recordDate time.Time) responseBody {
	err := db.AddHealthRecord(userId, insulin, sportTime, weight, bloodPressure, recordTime, recordDate)
	if err != nil {
		return responseInternalServerError(err)
	}
	return responseOK()
}

func getBloodSugarRecord(userId int, recordDate time.Time) responseBody {
	bloodRecord, _, err := db.GetBloodSugarRecordFromRecordDate(userId, recordDate)
	if err != nil {
		return responseInternalServerError(err)
	}
	levelMap := make(map[string]interface{})
	err = json.Unmarshal([]byte(bloodRecord.Level), &levelMap)
	if err != nil {
		return responseInternalServerError(err)
	}
	timeMap := make(map[string]interface{})
	err = json.Unmarshal([]byte(bloodRecord.RecordTime), &timeMap)
	if err != nil {
		return responseInternalServerError(err)
	}
	return responseOKWithData(gin.H{
		"bloodId":   bloodRecord.ID,
		"level":     gin.H(levelMap),
		"bloodTime": gin.H(timeMap),
	})
}

func getHealthRecordList(userId int, beginId int, needNumber int) responseBody {
	healthRecords, err := db.GetHealthRecordList(userId, beginId, needNumber)
	if err != nil {
		return responseInternalServerError(err)
	}
	respRecords := make([]gin.H, len(healthRecords))
	for i, record := range healthRecords {
		respRecords[i] = gin.H{
			"healthId":      record.ID,
			"insulin":       record.Insulin,
			"sportTime":     record.SportTime,
			"weight":        record.Weight,
			"bloodPressure": record.BloodPressure,
			"healthTime":    record.RecordTime,
			"healthDate":    utils.GoTimeToDateTime(record.RecordDate),
		}
	}
	return responseOKWithData(respRecords)
}

func getBloodSugarRecordList(userId int, beginId int, needNumber int) responseBody {
	bloodRecords, err := db.GetBloodSugarRecordList(userId, beginId, needNumber)
	if err != nil {
		return responseInternalServerError(err)
	}
	respRecords := make([]gin.H, len(bloodRecords))
	for i, record := range bloodRecords {
		levelMap := make(map[string]string)
		err := json.Unmarshal([]byte(record.Level), &levelMap)
		if err != nil {
			return responseInternalServerError(err)
		}
		maxBlood := 0.0
		minBlood := 9999.0
		sumBlood := 0.0
		count := 0
		for _, value := range levelMap {
			if value == "0" {
				continue
			}
			lev, _ := strconv.ParseFloat(value, 64)
			if maxBlood < lev {
				maxBlood = lev
			}
			if minBlood > lev {
				minBlood = lev
			}
			sumBlood += lev
			count++
		}
		aveBlood := sumBlood / float64(count)
		respRecords[i] = gin.H{
			"bloodId":      record.ID,
			"averageBlood": aveBlood,
			"maxBlood":     maxBlood,
			"minBlood":     minBlood,
			"bloodDate":    utils.GoTimeToDateTime(record.RecordDate),
		}
	}
	return responseOKWithData(respRecords)
}
