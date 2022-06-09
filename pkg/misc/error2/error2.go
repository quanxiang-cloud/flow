package error2

import "fmt"

// Error 错误处理封装
type Error struct {
	Code    int64
	Message string
}

// NewError 创建一个错误
func NewError(code int64, params ...interface{}) Error {
	if len(params) > 0 {
		return Error{
			Code:    code,
			Message: fmt.Sprintf(Translation(code), params...),
		}
	}
	return Error{
		Code: code,
	}
}

func (e Error) Error() string {
	if e.Message != "" {
		return e.Message
	}
	return Translation(e.Code)
}

// NewErrorWithString 创建一个错误
func NewErrorWithString(code int64, message string) Error {
	return Error{
		Code:    code,
		Message: message,
	}
}
