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
	"encoding/json"
	"github.com/quanxiang-cloud/flow/pkg/misc/logger"
	"strconv"
	"time"
)

// FormatValue format value
func FormatValue(value string, valueType string) interface{} {
	defer func() {
		err := recover()
		if err != nil {
			logger.Logger.Error(err)
		}
	}()

	switch valueType {
	case "datetime":
		{
			t, _ := time.Parse("2006-01-02", value)
			return t
		}
	case "boolean":
		{
			t, _ := strconv.ParseBool(value)
			return t
		}
	case "number":
		{
			t, _ := strconv.ParseFloat(value, 64)
			return t
		}
	case "json":
		{
			var t interface{}
			err := json.Unmarshal([]byte(value), &t)
			if err != nil {
				return nil
			}
			return t
		}
	default:
		return value
	}

	return nil
}

// Strval change to string
func Strval(value interface{}) string {
	// interface 转 string
	val, _ := StrvalAndType(value)
	return val
}

// StrvalAndType change to string and return type
func StrvalAndType(value interface{}) (string, string) {
	// interface 转 string
	var key string
	_type := "string"
	if value == nil {
		return key, _type
	}

	switch value.(type) {
	case float64:
		ft := value.(float64)
		key = strconv.FormatFloat(ft, 'f', -1, 64)
		_type = "float64"
	case float32:
		ft := value.(float32)
		key = strconv.FormatFloat(float64(ft), 'f', -1, 64)
		_type = "float32"
	case int:
		it := value.(int)
		key = strconv.Itoa(it)
		_type = "int"
	case uint:
		it := value.(uint)
		key = strconv.Itoa(int(it))
		_type = "uint"
	case int8:
		it := value.(int8)
		key = strconv.Itoa(int(it))
		_type = "int8"
	case uint8:
		it := value.(uint8)
		key = strconv.Itoa(int(it))
		_type = "uint8"
	case int16:
		it := value.(int16)
		key = strconv.Itoa(int(it))
		_type = "int16"
	case uint16:
		it := value.(uint16)
		key = strconv.Itoa(int(it))
		_type = "uint16"
	case int32:
		it := value.(int32)
		key = strconv.Itoa(int(it))
		_type = "int32"
	case uint32:
		it := value.(uint32)
		key = strconv.Itoa(int(it))
		_type = "uint32"
	case int64:
		it := value.(int64)
		key = strconv.FormatInt(it, 10)
		_type = "int64"
	case uint64:
		it := value.(uint64)
		key = strconv.FormatUint(it, 10)
		_type = "uint64"
	case string:
		key = value.(string)
	case []byte:
		key = string(value.([]byte))
		_type = "[]byte"
	case bool:
		key = strconv.FormatBool(value.(bool))
		_type = "bool"
	default:
		newValue, _ := json.Marshal(value)
		key = string(newValue)
		_type = "json"
	}

	return key, _type
}
