package server

import (
	"db"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid"
	"io/ioutil"
	"net/http"
	"strconv"
	"utils"
)

func sendVerificationCode(phoneNumber string) responseBody {
	nowCode, err := db.GetNowVerificationCode(phoneNumber)
	if err != nil {
		return responseInternalServerError(err)
	}
	var code string
	if nowCode != "" {
		code = nowCode
	} else {
		code = utils.RandCode()
	}
	url := fmt.Sprintf("http://127.0.0.1:19987/send_message?phone_number=%s&code=%s",
		phoneNumber, code)
	resp, err := http.Get(url)

	if err != nil {
		return responseInternalServerError(err)
	}
	defer func() {
		err := resp.Body.Close()
		if err != nil {
			fmt.Printf("%s", err)
			return
		}
	}()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return responseInternalServerError(err)
	}
	if body[0] == 'o' {
		err := db.SetNewVerificationCode(phoneNumber, code)
		if err != nil {
			return responseInternalServerError(err)
		} else {
			return responseOK()
		}
	}
	return responseInternalServerError(errors.New("external service(send message) error"))

}

func registerNewUser(phoneNumber string, password string, userName string, code string) responseBody {
	if phoneNumber == "" {
		return responseNormalError("手机号码不能为空")
	}
	if password == "" {
		return responseNormalError("密码不能为空")
	}
	if code == "" {
		return responseNormalError("验证码不能为空")
	}
	codeCheck, err := db.CheckPhoneCodeCorrection(phoneNumber, code)
	if err != nil {
		return responseInternalServerError(err)
	}
	if !codeCheck {
		return responseNormalError("验证码错误或已失效")
	}
	err = db.CreateNewUser(phoneNumber, userName, password)
	if err != nil {
		return responseNormalError("用户已经存在")
	}
	return responseOK()
}

func loginUser(phoneNumber string, password string) responseBody {
	if phoneNumber == "" {
		return responseNormalError("手机号码不能为空")
	}
	if password == "" {
		return responseNormalError("密码不能为空")
	}
	user, err := db.GetUserFromPhoneNumber(phoneNumber)
	if err != nil {
		return responseNormalError("用户不存在")
	}
	if user.Password != password {
		return responseNormalError("密码错误")
	}
	sessionId, err := db.GetNowSessionId(phoneNumber)
	var userId string
	if err != nil {
		return responseInternalServerError(err)
	}
	if sessionId == "" {
		userId = strconv.Itoa(int(user.ID))
		sessionId = uuid.NewV5(uuid.NamespaceDNS, userId).String()
		err := db.SetNewSessionId(sessionId, userId)
		if err != nil {
			return responseInternalServerError(err)
		}
	}
	checkIn, err := db.CheckUserCheckIn(int(user.ID))
	if err != nil {
		return responseInternalServerError(err)
	}
	return responseOKWithData(gin.H{
		"userId":     user.ID,
		"session_id": sessionId,
		"username":   user.UserName,
		"iconUrl":    user.HeadPortraitUrl,
		"exp":        user.Exp,
		"level":      user.Level,
		"isCheck":    checkIn,
	})
}

func alterUserInformation(userId int, userName string, gender string, height float64,
	weight float64, area string, job string, age int) responseBody {
	err := db.AlterUserInformationFromUserId(userId, userName, gender, height, weight, area, job, age)
	if err != nil {
		return responseInternalServerError(err)
	}
	return responseOK()
}

func getUserInfoFromUserId(userId int) responseBody {
	user, err := db.GetUserFromUserId(userId)
	if err != nil {
		return responseInternalServerError(err)
	}
	return responseOKWithData(gin.H{
		"username": user.UserName,
		"iconUrl":  user.HeadPortraitUrl,
		"age":      user.Age,
		"gender":   user.Gender,
		"job":      user.Job,
		"area":     user.Area,
		"height":   user.Height,
		"weight":   user.Weight,
		"exp":      user.Exp,
		"level":    user.Level,
	})
}

func getOtherUserInformationFromOtherUserId(userId int, otherUserId int) responseBody {
	user, err := db.GetUserFromUserId(otherUserId)
	if err != nil {
		return responseInternalServerError(err)
	}
	following, err := db.CheckUserFollowingOtherUser(userId, otherUserId)
	return responseOKWithData(gin.H{
		"username": user.UserName,
		"iconUrl":  user.HeadPortraitUrl,
		"age":      user.Age,
		"gender":   user.Gender,
		"job":      user.Job,
		"area":     user.Area,
		"height":   user.Height,
		"weight":   user.Weight,
		"exp":      user.Exp,
		"level":    user.Level,
		"isFollow": following,
	})
}

func alterPassword(phoneNumber string, code string, newPassword string) responseBody {
	if phoneNumber == "" {
		return responseNormalError("手机号码不能为空")
	}
	if newPassword == "" {
		return responseNormalError("新密码不能为空")
	}
	if code == "" {
		return responseNormalError("验证码不能为空")
	}
	codeCheck, err := db.CheckPhoneCodeCorrection(phoneNumber, code)
	if err != nil {
		return responseInternalServerError(err)
	}
	if !codeCheck {
		return responseNormalError("验证码错误或已失效")
	}
	err = db.AlterUserPasswordFromPhoneNumber(phoneNumber, newPassword)
	if err != nil {
		return responseNormalError("用户不存在")
	}
	return responseOK()
}

func getUserPrivacySetting(userId int) responseBody {
	privacySetting, err := db.GetPrivacySettingFromUserId(userId)
	if err != nil {
		return responseInternalServerError(err)
	}
	return responseOKWithData(gin.H{
		"showPhone":  privacySetting.ShowPhoneNumber,
		"showGender": privacySetting.ShowGender,
		"showAge":    privacySetting.ShowAge,
		"showHeight": privacySetting.ShowHeight,
		"showWeight": privacySetting.ShowWeight,
		"showArea":   privacySetting.ShowArea,
		"showJob":    privacySetting.ShowJob,
	})
}

func alterUserPrivacy(userId int, showPhoneNumber bool, showGender bool, showAge bool,
	showHeight bool, showWeight bool, showArea bool, showJob bool) responseBody {
	err := db.AlterUserPrivacySettingFromUserId(userId, showPhoneNumber, showGender, showAge,
		showHeight, showWeight, showArea, showJob)
	if err != nil {
		return responseInternalServerError(err)
	}
	return responseOK()
}

func followUser(userId int, targetUserId int) responseBody {
	err := db.AddUserFollowing(userId, targetUserId)
	if err != nil {
		responseInternalServerError(err)
	}
	return responseOK()
}

func ignoreUser(userId int, targetUserId int) responseBody {
	err := db.RemoveUserFollowing(userId, targetUserId)
	if err != nil {
		responseInternalServerError(err)
	}
	return responseOK()
}

func getUserFollowingList(userId int, beginId int, needNumber int) responseBody {
	users, count, err := db.GetUserFollowingList(userId, beginId, needNumber)
	if err != nil {
		return responseInternalServerError(err)
	}
	respUsers := make([]gin.H, len(users))
	for i := 0; i < len(users); i++ {
		respUsers[i] = gin.H{
			"followId": users[i].ID,
			"username": users[i].UserName,
			"iconUrl":  users[i].HeadPortraitUrl,
		}
	}
	return responseOKWithData(gin.H{
		"data":  respUsers,
		"total": count,
	})
}

func getUserFollowerList(userId int, beginId int, needNumber int) responseBody {
	users, count, err := db.GetUserFollowerList(userId, beginId, needNumber)
	if err != nil {
		return responseInternalServerError(err)
	}
	respUsers := make([]gin.H, len(users))
	for i := 0; i < len(users); i++ {
		respUsers[i] = gin.H{
			"followMeId": users[i].ID,
			"username":   users[i].UserName,
			"iconUrl":    users[i].HeadPortraitUrl,
		}
	}
	return responseOKWithData(gin.H{
		"data":  respUsers,
		"total": count,
	})
}

func logoutUser(SessionId string) responseBody {
	if SessionId == "" {
		return responseNormalError("请先登录")
	}
	err := db.RemoveSessionId(SessionId)
	if err != nil {
		return responseInternalServerError(err)
	}
	return responseOK()
}

