package time2

import "time"

// time2 用于统一项目时间格式

const (
	// ISO8601 ISO8601
	ISO8601 = "2006-01-02T15:04:05-0700"
)

// NowUnix 获取当前时间戳
func NowUnix() int64 {
	return time.Now().UTC().Unix()
}

// Now 获取当前时间
func Now() string {
	return time.Now().UTC().Format(ISO8601)
}

// ISO8601ToUnix ISO8601 to unix
func ISO8601ToUnix(ts string) (int64, error) {
	t, err := time.Parse(ISO8601, ts)
	if err != nil {
		return 0, err
	}

	return t.Unix(), nil
}

// UnixToISO8601   unix to ISO8601
func UnixToISO8601(ts int64) string {
	return time.Unix(ts, 0).Format(ISO8601)
}

// NowUnixMill 获取当前时间戳(毫秒)
func NowUnixMill() int64 {
	return time.Now().UTC().UnixNano() / 1e6
}
