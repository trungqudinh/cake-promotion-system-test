package time

import (
	"fmt"
	"time"
)

var (
	LocLocal *time.Location
)

type TimeFormatString string

const APIClientDateTimeFormat TimeFormatString = "02/01/2006 15:04:05"
const ApiClientTimeFormat TimeFormatString = "15:04:05"

func (t TimeFormatString) IsValid() bool {
	return t == APIClientDateTimeFormat
}

func (t TimeFormatString) String() string {
	return string(t)
}

const (
	FormatDateTime = "2006-01-02 15:04:05"
)

func init() {
	l, err := time.LoadLocation("Asia/Bangkok")
	if err != nil {
		panic(err)
	}
	LocLocal = l
}

func ParseLocal(layout, value string) (time.Time, error) {
	t, err := time.ParseInLocation(layout, value, LocLocal)
	if err != nil {
		return time.Time{}, err
	}
	return t, nil
}

func GetTimeZone() string {
	_, offsetSeconds := time.Now().Zone()
	offsetHours := offsetSeconds / 3600
	offsetMinutes := (offsetSeconds % 3600) / 60
	zoneString := fmt.Sprintf("%+03d%02d %+03d", offsetHours, offsetMinutes, offsetHours)
	return zoneString
}

func FromTimeToString(value time.Time, format TimeFormatString) string {
	if !format.IsValid() {
		return ""
	}
	return value.Format(string(format))
}

func FromTimeToClientString(value time.Time) string {
	return FromTimeToString(value, APIClientDateTimeFormat)
}
