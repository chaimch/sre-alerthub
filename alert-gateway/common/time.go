package common

import (
	"time"
)

func GetNowTimeStrForHHMMSS() string{
	return time.Now().Format("15:04:05")
}

func GetTodayDateForYYmmDD() time.Time {
	todayZero, _ := time.ParseInLocation("2006-01-02", "2006-01-02 15:04:05", time.Local)
	return todayZero
}
