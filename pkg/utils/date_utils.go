/*
Copyright 2022 QuanxiangCloud Authors
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at
     http://www.apache.org/licenses/LICENSE-2.0
Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package utils

import (
	"github.com/quanxiang-cloud/flow/pkg/misc/time2"
	"regexp"
	"time"
)

var (
	ISO8601TimeCompile = regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}\+\d{4}$`)
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

// ISO8601FmtCheck 时间格式检查
// fmt: 2022-04-22T05:42:17+0000
func ISO8601FmtCheck(fmt *string) bool {
	if len(*fmt) != 24 {
		return false
	}
	result := ISO8601TimeCompile.FindAllStringSubmatch(*fmt, -1)
	// not matched
	if len(result) == 0 {
		return false
	}
	return true
}
