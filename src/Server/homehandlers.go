package server

import (
	"db"
	"github.com/gin-gonic/gin"
)

func checkinUser(userId int)responseBody{
	err:=db.CheckInUser(userId)
	if err!=nil{
		return responseInternalServerError(err)
	}
	return responseOK()
}

func getUserFamilyMemberList(userId int) responseBody {
	familyMembers,err:=db.GetUserFamilyMemberList(userId)
	if err!=nil{
		return responseInternalServerError(err)
	}
	respFamilyMembers:=make([]gin.H,len(familyMembers))
	for i,member:=range familyMembers{
		respFamilyMembers[i] = gin.H{
			"nickname": member.CallName,
			"tel":member.PhoneNumber,
			"familyId":member.ID,
		}
	}
	return responseOKWithData(respFamilyMembers)
}

func linkNewFamilyMember(userId int ,callName string ,phoneNumber string,code string)responseBody{
	nowCode ,err:= db.GetNowVerificationCode(phoneNumber)
	if err!=nil{
		return responseInternalServerError(err)
	}
	if nowCode !=code {
		return responseNormalError("验证码错误")
	}
	err = db.AddFamilyMember(userId,callName,phoneNumber)
	if err!=nil{
		return responseInternalServerError(err)
	}
	return responseOK()
}