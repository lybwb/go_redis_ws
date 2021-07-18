package commons

import (
	"time"
)

func GetNowUTCTime() time.Time {
	return time.Now().UTC()
}

func GetUTCMsTimeStamp(t time.Time) string {
	TimeStampMsFormat := "2006-01-02T15:04:05.000000Z"
	return t.UTC().Format(TimeStampMsFormat)
}

func GetNowUTCMsTimeStamp() string {
	return GetUTCMsTimeStamp(time.Now())
}
