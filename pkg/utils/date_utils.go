package utils

import (
	"github.com/quanxiang-cloud/flow/pkg/misc/time2"
	"time"
)

// AddDaysHoursMinutes add
func AddDaysHoursMinutes(dateStr string, days int, hours int, minutes int) string {
	date, err := time2.ISO8601ToUnix(dateStr)
	if err != nil {
		return ""
	}
	date += int64(days*24*60*60) + int64(hours*60*60) + int64(minutes*60)
	return time2.UnixToISO8601(date)
}

// ChangeBjTimeToISO8601 date format 2006-01-02 15:04:05
func ChangeBjTimeToISO8601(value string) string {
	t, _ := time.Parse("2006-01-02 15:04:05", value)
	return time2.UnixToISO8601(t.Unix() - 8*60*60)
}

// ChangeBjTimeToUtcUnix date format 2006-01-02 15:04:05
// func ChangeBjTimeToUtcUnix(value string) int64 {
// 	t, _ := time.Parse("2006-01-02 15:04:05", value)
// 	return t.Unix() - 8*60*60
// }

// ChangeISO8601ToBjTime change to bj time
func ChangeISO8601ToBjTime(value string) string {
	unix, err := time2.ISO8601ToUnix(value)
	if err != nil {
		return ""
	}

	unix += 8 * 60 * 60
	return time.Unix(unix, 0).Format("2006-01-02 15:04:05")
}
