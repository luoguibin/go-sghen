package helper

import (
	"fmt"
	"math/rand"
	"time"
)

// StrToTimeStamp 把一个时间字符串转为unix时间戳
func StrToTimeStamp(timeStr string) time.Time {
	//	time = "2015-09-14 16:33:00"
	loc, _ := time.LoadLocation("Asia/Shanghai")
	t, _ := time.ParseInLocation("2006-01-02 15:04:05", timeStr, loc)
	return t
}

//GetNowDateTime get now datatime(Y-m-d H:i:s)
func GetNowDateTime() string {
	return time.Now().Format("2006-01-02 15:04:05")
}

// GetMillisecond ...
func GetMillisecond() int64 {
	return time.Now().UnixNano() / 1e6
}

// GetMicrosecond ...
func GetMicrosecond() int64 {
	return time.Now().UnixNano() / 1e3
}

// GetNanosecond ...
func GetNanosecond() int64 {
	return time.Now().UnixNano()
}

// RandomNum4 ...
func RandomNum4() string {
	return fmt.Sprintf("%04v", rand.New(rand.NewSource(time.Now().UnixNano())).Int31n(10000))
}
