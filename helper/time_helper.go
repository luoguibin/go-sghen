package helper

import "time"

//把一个时间字符串转为unix时间戳
func StrToTimeStamp(timeStr string) time.Time {
	//	time = "2015-09-14 16:33:00"
	loc, _ := time.LoadLocation("Asia/Shanghai")
	t, _ := time.ParseInLocation("2006-01-02 15:04:05", timeStr, loc)
	return t
}