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

func sendVerificationCode(PhoneNumber string) responseBody {
	nowCode, err := db.GetNowVerificationCode(PhoneNumber)
	if err != nil {
		return responseInternalServerError(err)
	}
	var code string
	if nowCode != "" {
		code = nowCode
	} else {
		code = utils.RandCode()
	}
	url := fmt.Sprintf("http://127.0.0.1:7799/send_message?phone_number=%s&code=%s",
		PhoneNumber, code)
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
		err := db.SetNewVerificationCode(PhoneNumber, code)
		if err != nil {
			return responseInternalServerError(err)
		} else {
			return responseOK()
		}
	}
	return responseInternalServerError(errors.New("external service(send message) error"))

}

func registerNewUser(PhoneNumber string, Password string, Code string) responseBody {
	if PhoneNumber == "" {
		return responseNormalError("手机号码不能为空")
	}
	if Password == "" {
		return responseNormalError("密码不能为空")
	}
	if Code == "" {
		return responseNormalError("验证码不能为空")
	}
	codeCheck, err := db.CheckPhoneCodeCorrection(PhoneNumber, Code)
	if err != nil {
		return responseInternalServerError(err)
	}
	if !codeCheck {
		return responseNormalError("验证码错误或已失效")
	}
	err = db.CreateNewUser(PhoneNumber, Password)
	if err != nil {
		return responseNormalError("用户已经存在")
	}
	return responseOK()
}

func loginUser(PhoneNumber string, Password string) responseBody {
	if PhoneNumber == "" {
		return responseNormalError("手机号码不能为空")
	}
	if Password == "" {
		return responseNormalError("密码不能为空")
	}
	user, err := db.GetUserFromPhoneNumber(PhoneNumber)
	if err != nil {
		return responseNormalError("用户不存在")
	}
	if user.Password != Password {
		return responseNormalError("密码错误")
	}
	sessionId, err := db.GetNowSessionId(PhoneNumber)
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
	return responseOKWithData(gin.H{
		"userId":     userId,
		"session_id": sessionId,
		"username":   user.UserName,
		"iconUrl":    user.HeadPortraitUrl,
		"exp":        user.Exp,
		"level":      user.Level,
	})
}

func alterUserInformation(SessionId string, UserName string, Gender string, Height float64,
	Weight float64, Area string, Job string, Age int) responseBody {
	if SessionId == "" {
		return responseNormalError("请先登录")
	}
	userId, err := db.GetNowSessionId(SessionId)
	if err != nil {
		return responseInternalServerError(err)
	}
	err = db.AlterUserInformationFromUserId(userId, UserName, Gender, Height, Weight, Area, Job, Age)
	if err != nil {
		return responseInternalServerError(err)
	}
	return responseOK()
}

func getUserInformationFromUserId(SessionId string, TargetUserId string) responseBody {
	if SessionId == "" {
		return responseNormalError("请先登录")
	}
	UserId, err := db.GetNowSessionId(SessionId)
	if err != nil {
		return responseInternalServerError(err)
	}
	if UserId == "" {
		return responseNormalError("请先登录")
	}
	user, err := db.GetUserFromUserId(TargetUserId)
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

func getUserInformationFromSessionId(SessionId string) responseBody {
	if SessionId == "" {
		return responseNormalError("请先登录")
	}
	userId, err := db.GetNowSessionId(SessionId)
	if err != nil {
		return responseInternalServerError(err)
	}
	return getUserInformationFromUserId(SessionId, userId)
}

func alterPassword(PhoneNumber string, Code string, NewPassword string) responseBody {
	if PhoneNumber == "" {
		return responseNormalError("手机号码不能为空")
	}
	if NewPassword == "" {
		return responseNormalError("新密码不能为空")
	}
	if Code == "" {
		return responseNormalError("验证码不能为空")
	}
	codeCheck, err := db.CheckPhoneCodeCorrection(PhoneNumber, Code)
	if err != nil {
		return responseInternalServerError(err)
	}
	if !codeCheck {
		return responseNormalError("验证码错误或已失效")
	}
	err = db.AlterUserPasswordFromPhoneNumber(PhoneNumber, NewPassword)
	if err != nil {
		return responseNormalError("用户不存在")
	}
	return responseOK()
}

func getUserPrivacySetting(SessionId string) responseBody {
	if SessionId == "" {
		return responseNormalError("请先登录")
	}
	userId, err := db.GetNowSessionId(SessionId)
	if err != nil {
		return responseInternalServerError(err)
	}
	privacySetting, err := db.GetPrivacySettingFromUserId(userId)
	return responseOKWithData(gin.H{
		"ShowPhoneNumber": privacySetting.ShowPhoneNumber,
		"ShowGender":      privacySetting.ShowGender,
		"ShowAge":         privacySetting.ShowAge,
		"ShowHeight":      privacySetting.ShowHeight,
		"ShowWeight":      privacySetting.ShowWeight,
		"ShowArea":        privacySetting.ShowArea,
		"ShowJob":         privacySetting.ShowJob,
	})
}

func alterUserPrivacy(SessionId string, ShowPhoneNumber bool, ShowGender bool, ShowAge bool,
	ShowHeight bool, ShowWeight bool, ShowArea bool, ShowJob bool) responseBody {
	if SessionId == "" {
		return responseNormalError("请先登录")
	}
	userId, err := db.GetNowSessionId(SessionId)
	if err != nil {
		return responseInternalServerError(err)
	}
	err = db.AlterUserPrivacySettingFromUserId(userId, ShowPhoneNumber, ShowGender, ShowAge,
		ShowHeight, ShowWeight, ShowArea, ShowJob)
	if err != nil {
		return responseInternalServerError(err)
	}
	return responseOK()
}

func followUser(SessionId string, TargetUserId string) responseBody {
	if SessionId == "" {
		return responseNormalError("请先登录")
	}
	userId, err := db.GetNowSessionId(SessionId)
	if err != nil {
		return responseInternalServerError(err)
	}
	err = db.AddUserFollowing(userId, TargetUserId)
	if err != nil {
		responseInternalServerError(err)
	}
	return responseOK()
}

func ignoreUser(SessionId string, TargetUserId string) responseBody {
	if SessionId == "" {
		return responseNormalError("请先登录")
	}
	userId, err := db.GetNowSessionId(SessionId)
	if err != nil {
		return responseInternalServerError(err)
	}
	err = db.RemoveUserFollowing(userId, TargetUserId)
	if err != nil {
		responseInternalServerError(err)
	}
	return responseOK()
}

func getUserFollowingList(SessionId string, BeginId string, NeedNumber string) responseBody {
	if SessionId == "" {
		return responseNormalError("请先登录")
	}
	userId, err := db.GetNowSessionId(SessionId)
	if err != nil {
		return responseInternalServerError(err)
	}
	users, err := db.GetUserFollowingList(userId, BeginId, NeedNumber)
	respUsers := make([]responseSimpleUser, len(users))
	for i := 0; i < len(users); i++ {
		fmt.Println(users[i])
		respUsers[i] = responseSimpleUser{
			followId: users[i].ID,
			username: users[i].UserName,
			iconUrl:  users[i].HeadPortraitUrl,
		}
	}
	return responseOKWithData(gin.H{
		"data": respUsers,
	})
}

func getUserFollowerList(SessionId string, BeginId string, NeedNumber string) responseBody {
	//userId,err:=db.GetNowSessionId(SessionId)
	//if err!=nil{
	//	return responseInternalServerError(err)
	//}

	return responseBody{}
}

func logoutUser(SessionId string) responseBody {
	err := db.RemoveSessionId(SessionId)
	if err != nil {
		return responseInternalServerError(err)
	}
	return responseOK()
}
