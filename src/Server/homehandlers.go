package server

import (
	"db"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
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
	bloodRecord, exist, err := db.GetBloodSugarRecordFromRecordDate(userId, recordDate)
	if !exist {
		return responseOKWithData([]gin.H{})
	}
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

func getVoiceDictationResult(audioBase64 string) (string, error) {
	u := fmt.Sprintf("http://106.15.187.190:19987/voice_dictation?pwd=04bc1911b62299651aa9ce63c8d74770")
	resp, err := http.PostForm(u, url.Values{"audio": {audioBase64}})
	if err != nil {
		return "", err
	}
	defer func() {
		err := resp.Body.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	respMap := make(map[string]interface{})
	err = json.Unmarshal(body, &respMap)
	if err != nil {
		return "", err
	}
	fmt.Println(respMap)
	if int(respMap["code"].(float64)) == 0 {
		return respMap["data"].(string), nil
	} else {
		return "", errors.New(fmt.Sprintf("external service(voice_dictation:%s) error", respMap["msg"].(string)))
	}
}

func parseBloodSugarRecordVoiceInput(audioBase64 string) responseBody {
	dictationResult, err := getVoiceDictationResult(audioBase64)
	if err != nil {
		return responseInternalServerError(err)
	}
	timeString := []string{"早餐前", "早餐后", "午餐前", "午餐后", "晚餐前", "晚餐后", "睡前"}
	conString, conTime := utils.StringContains(dictationResult, timeString)
	if !conTime {
		return responseNormalError("请说明记录的时间段")
	}
	reg := regexp.MustCompile(`\d+[:.]*\d*`)
	conValue := reg.FindAllString(dictationResult, -1)
	var valueResult float64
	var valueMatch bool
	for _, value := range conValue {
		fv, err := strconv.ParseFloat(value, 64)
		if err == nil && fv >= 1.0 && fv <= 33.3 {
			valueResult = fv
			valueMatch = true
			break
		}
	}
	if !valueMatch {
		return responseNormalError("未解析到合法的血糖值（范围在1.0到33.3之间）")
	}
	periodMap := map[string]string{
		"早餐前": "beforeBF",
		"早餐后": "afterBF",
		"午餐前": "beforeLC",
		"午餐后": "afterLC",
		"晚餐前": "beforeDN",
		"晚餐后": "afterDN",
		"睡前":  "beforeSP",
	}
	return responseOKWithData(gin.H{
		"periodValue": periodMap[conString],
		"periodLabel": conString,
		"value":       valueResult,
	})
}

func parseHealthRecordVoiceInput(audioBase64 string) responseBody {
	dictationResult, err := getVoiceDictationResult(audioBase64)
	if err != nil {
		return responseInternalServerError(err)
	}
	prefixStrings := []string{"胰岛素用量", "运动时长", "体重", "舒张压", "收缩压"}
	sentences := strings.Split(dictationResult, "，")
	resultMap := make(gin.H)
	for _, st := range sentences {
		pre, con := utils.StringHasPrefixs(st, prefixStrings)
		if !con {
			continue
		}
		switch pre {
		case "胰岛素用量":
			reg := regexp.MustCompile(`\d+[:.]*\d*`)
			result := reg.FindString(st)
			insulin, err := strconv.ParseFloat(result, 64)
			if err != nil {
				break
			}
			resultMap["insulin"] = insulin
		case "运动时长":
			regData := regexp.MustCompile(`\d+[:.]*\d*`)
			regHour := regexp.MustCompile(`\d+[:.]*\d*小时`)
			regMinute := regexp.MustCompile(`\d+[:.]*\d*分钟?`)
			hour := regHour.FindString(st)
			minute := regMinute.FindString(st)
			timeMap := gin.H{}
			if minute != "" {
				minute = regData.FindString(minute)
				minuteF, err := strconv.ParseFloat(minute, 64)
				if err != nil {
					break
				}
				timeMap["minute"] = int(minuteF)

			} else {
				timeMap["minute"] = 0
			}
			if hour != "" {
				hour = regData.FindString(hour)
				strings.Replace(hour, ":", ".", -1)
				hourF, err := strconv.ParseFloat(hour, 64)
				if err != nil {
					break
				}
				timeMap["hour"] = int(hourF)
				if timeMap["minute"] ==0 {
					timeMap["minute"] = int((hourF - float64(int(hourF))) * 60)
				}
			} else {
				timeMap["hour"] = 0
			}
			fmt.Println(timeMap)
			resultMap["sportsTime"] = timeMap
		case "体重":
			reg := regexp.MustCompile(`\d+[:.]*\d*`)
			result := reg.FindString(st)
			weight, err := strconv.ParseFloat(result, 64)
			if err != nil {
				break
			}
			resultMap["weight"] = weight
		case "舒张压":
			reg := regexp.MustCompile(`\d+[:.]*\d*`)
			result := reg.FindString(st)
			pressure, err := strconv.Atoi(result)
			if err != nil {
				break
			}
			resultMap["pressure1"] = pressure
		case "收缩压":
			reg := regexp.MustCompile(`\d+[:.]*\d*`)
			result := reg.FindString(st)
			pressure, err := strconv.Atoi(result)
			if err != nil {
				break
			}
			resultMap["pressure2"] = pressure
		}
	}
	return responseOKWithData(resultMap)
}
