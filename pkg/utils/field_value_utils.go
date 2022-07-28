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
	"reflect"
	"strconv"
	"strings"
)

// 变量格式
// $a.b.c
// $a.[0].b.c
// $webhook0ea9vOcw8z2LWnDA_rZm2.empList.[0].openUserId
// $formData3Wt4sVFrwJqbsQzqyd8X3.field_UXvvfRBU
// fieldxxx, fieldxxx.value, fieldxxx.[].value

// GetFieldValue get field value from data
func GetFieldValue(data map[string]interface{}, fieldKey string) interface{} {
	if len(fieldKey) == 0 || data == nil {
		return nil
	}

	value, ok := data[fieldKey]
	if ok {
		return value
	}

	fields := strings.Split(fieldKey, ".")
	if len(fields) > 1 {
		for index := range fields {
			tempFieldKey := strings.Join(fields[:index+1], ".")

			tempValue, ok := data[tempFieldKey]
			if ok {
				if index == len(fields)-1 {
					return tempValue
				}
				return RecursionGetValue(tempValue, fields[index+1:])
			}
		}
	}

	return nil
}

// 变量格式
// $a.b.c
// $a.[0].b.c
// $webhook0ea9vOcw8z2LWnDA_rZm2.empList.[0].openUserId
// $formData3Wt4sVFrwJqbsQzqyd8X3.field_UXvvfRBU
// fieldxxx, fieldxxx.value, fieldxxx.[].value

// RecursionGetValue 递归获取值
func RecursionGetValue(data interface{}, fields []string) interface{} {
	if data == nil || len(fields) == 0 {
		return data
	}

	for index, field := range fields {
		if strings.HasPrefix(field, "[") && strings.HasSuffix(field, "]") && len(field) > 2 {
			arrIndexStr := field[1 : len(field)-1]
			arrIndex, err := strconv.Atoi(arrIndexStr)
			if err != nil {
				return nil
			}

			cvs := reflect.ValueOf(data)
			v := cvs.Index(arrIndex).Interface()
			return RecursionGetValue(v, fields[index+1:])
		} else if strings.HasPrefix(field, "[") && strings.HasSuffix(field, "]") && len(field) == 2 {
			cValues := make([]interface{}, 0)
			cvs := reflect.ValueOf(data)
			for i := 0; i < cvs.Len(); i++ {
				v := cvs.Index(i).Interface()
				cValues = append(cValues, RecursionGetValue(v, fields[index+1:]))
			}
			return cValues
		} else {
			cValueMap := ChangeObjectToMap(data)
			value, ok := cValueMap[field]
			if ok {
				return RecursionGetValue(value, fields[index+1:])
			}
		}
	}

	return nil
}
