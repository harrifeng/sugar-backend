package utils

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

// 接受xxxx-xx-xx格式
func DateTimeParser(t string)time.Time{
	timeSplit:=strings.Split(t,"-")
	year ,_:= strconv.Atoi(timeSplit[0])
	month,_:= strconv.Atoi(timeSplit[1])
	day,_:= strconv.Atoi(timeSplit[2])
	return time.Date(year,time.Month(month),day,0,0,0,0,time.UTC)
}

func GoTimeToDateTime(tt time.Time)string{
	year,month,day:=tt.Date()
	return fmt.Sprintf("%d-%d-%d",year,month,day)
}