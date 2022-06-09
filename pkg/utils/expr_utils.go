package utils

import (
	"github.com/quanxiang-cloud/flow/pkg/misc/logger"
	"reflect"
	"strings"
)

// ExprCompare expr compare
func ExprCompare(value1 interface{}, value2 interface{}) int {
	type1 := reflect.TypeOf(value1).Kind()
	type2 := reflect.TypeOf(value2).Kind()

	if type1 != type2 {
		logger.Logger.Error("type mismatch")
		return 0
	}

	switch type1 {
	case reflect.String:
		return strings.Compare(value1.(string), value2.(string))
	case reflect.Float64:
		if value1.(float64) == value2.(float64) {
			return 0
		} else if value1.(float64) > value2.(float64) {
			return 1
		} else {
			return -1
		}
	case reflect.Float32:
		if value1.(float32) == value2.(float32) {
			return 0
		} else if value1.(float32) > value2.(float32) {
			return 1
		} else {
			return -1
		}
	case reflect.Int:
		if value1.(int) == value2.(int) {
			return 0
		} else if value1.(int) > value2.(int) {
			return 1
		} else {
			return -1
		}
	case reflect.Uint:
		if value1.(uint) == value2.(uint) {
			return 0
		} else if value1.(uint) > value2.(uint) {
			return 1
		} else {
			return -1
		}
	case reflect.Int8:
		if value1.(int8) == value2.(int8) {
			return 0
		} else if value1.(int8) > value2.(int8) {
			return 1
		} else {
			return -1
		}
	case reflect.Uint8:
		if value1.(uint8) == value2.(uint8) {
			return 0
		} else if value1.(uint8) > value2.(uint8) {
			return 1
		} else {
			return -1
		}
	case reflect.Int16:
		if value1.(int16) == value2.(int16) {
			return 0
		} else if value1.(int16) > value2.(int16) {
			return 1
		} else {
			return -1
		}
	case reflect.Uint16:
		if value1.(uint16) == value2.(uint16) {
			return 0
		} else if value1.(uint16) > value2.(uint16) {
			return 1
		} else {
			return -1
		}
	case reflect.Int32:
		if value1.(int32) == value2.(int32) {
			return 0
		} else if value1.(int32) > value2.(int32) {
			return 1
		} else {
			return -1
		}
	case reflect.Uint32:
		if value1.(uint32) == value2.(uint32) {
			return 0
		} else if value1.(uint32) > value2.(uint32) {
			return 1
		} else {
			return -1
		}
	case reflect.Int64:
		if value1.(int64) == value2.(int64) {
			return 0
		} else if value1.(int64) > value2.(int64) {
			return 1
		} else {
			return -1
		}
	case reflect.Uint64:
		if value1.(uint64) == value2.(uint64) {
			return 0
		} else if value1.(uint64) > value2.(uint64) {
			return 1
		} else {
			return -1
		}
	}

	return 0
}

// ExprInclude expr include
func ExprInclude(isArray bool, value1 interface{}, value2 interface{}) bool {
	if isArray {
		cValues1 := make([]interface{}, 0)
		if value1 != nil {
			cvs1 := reflect.ValueOf(value1)
			for i := 0; i < cvs1.Len(); i++ {
				v1 := cvs1.Index(i).Interface()
				cValues1 = append(cValues1, v1)
			}
		}
		cValues2 := make([]interface{}, 0)
		if value2 != nil {
			cvs2 := reflect.ValueOf(value2)
			for i := 0; i < cvs2.Len(); i++ {
				v2 := cvs2.Index(i).Interface()
				cValues2 = append(cValues2, v2)
			}
		}
		return Contains(cValues1, cValues2)
	}

	return strings.Contains(Strval(value1), Strval(value2))
}

// ExprNull expr null
func ExprNull(isArray bool, value1 interface{}) bool {
	if isArray {
		cValues1 := make([]interface{}, 0)
		if value1 != nil {
			cvs1 := reflect.ValueOf(value1)
			for i := 0; i < cvs1.Len(); i++ {
				v1 := cvs1.Index(i).Interface()
				cValues1 = append(cValues1, v1)
			}
		}
		return len(cValues1) == 0
	}
	return len(Strval(value1)) == 0
}

// ExprAnyInclude expr any include
func ExprAnyInclude(isArray bool, value1 interface{}, value2 interface{}) bool {
	if isArray {
		cValues1 := make([]interface{}, 0)
		if value1 != nil {
			cvs1 := reflect.ValueOf(value1)
			for i := 0; i < cvs1.Len(); i++ {
				v1 := cvs1.Index(i).Interface()
				cValues1 = append(cValues1, v1)
			}
		}
		cValues2 := make([]interface{}, 0)
		if value2 != nil {
			cvs2 := reflect.ValueOf(value2)
			for i := 0; i < cvs2.Len(); i++ {
				v2 := cvs2.Index(i).Interface()
				cValues2 = append(cValues2, v2)
			}
		}
		return Intersect(cValues1, cValues2)
	}
	return false
}
