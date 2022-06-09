package error2

// 10[模块代码] 4XX[状态码] 1000[错误码]
const (
	// Unknown 未知错误
	Unknown = -1
	// Internal 内部错误
	Internal = -2

	// Success 成功
	Success = 0

	// ErrParams 参数错误
	ErrParams = 1
)

var baseCode = map[int64]string{
	Unknown:  "unknown err.",
	Internal: "internal err.",

	Success: "sucess",

	ErrParams: "param error",
}

// CodeTable 码表
var CodeTable map[int64]string

// Translation translation code to message
func Translation(code int64) string {
	if CodeTable != nil {
		if text, ok := CodeTable[code]; ok {
			return text
		}
	}
	if text, ok := baseCode[code]; ok {
		return text
	}
	return "unknown code."
}
