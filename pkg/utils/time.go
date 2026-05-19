package utils

import (
	"time"
)

const LayoutDateTime = "2006-01-02 15:04"
const TimezoneLocation = "Asia/Jakarta"

// ParseTime mengonversi string ke time.Time dengan format WIB
func ParseTime(timeStr string) (time.Time, error) {
	loc, _ := time.LoadLocation(TimezoneLocation)
	return time.ParseInLocation(LayoutDateTime, timeStr, loc)
}

// FormatTime mengonversi time.Time ke string dengan format WIB
func FormatTime(t time.Time) string {
	loc, _ := time.LoadLocation(TimezoneLocation)
	return t.In(loc).Format(LayoutDateTime)
}
