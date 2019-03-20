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
	if PhoneNumber == "" {
		return responseNormalError("密码不能为空")
	}
	if PhoneNumber == "" {
		return responseNormalError("验证码不能为空")
	}
	codeCheck, err := db.CheckPhoneCodeCorrection(PhoneNumber, Code)
	if err != nil {
		return responseInternalServerError(err)
	}
	if !codeCheck {
		return responseNormalError("验证码错误或已失效")
	}
	user := db.User{
		PhoneNumber: PhoneNumber,
		Password:    Password,
	}
	err = db.CreateNewUser(user)
	if err != nil {
		return responseNormalError("用户已经存在")
	}
	return responseOK()
}

func loginUser(PhoneNumber string, Password string) responseBody {
	user, err := db.GetUserFromPhoneNumber(PhoneNumber)
	if err != nil {
		return responseNormalError("用户已存在")
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

func checkPasswordCorrection(PhoneNumber string, Password string) (bool, error) {
	return true, nil
}
