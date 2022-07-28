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
	"fmt"
)

// Contain Is arr has data
func Contain(arr []string, data string) bool {
	for _, value := range arr {
		if value == data {
			return true
		}
	}
	return false
}

// Intersect arr1 and arr2 is intersect
func Intersect(arr1 []interface{}, arr2 []interface{}) bool {
	for _, key1 := range arr1 {
		for _, key2 := range arr2 {
			if key1 == key2 {
				return true
			}
		}
	}

	return false
}

// IntersectString arr1 and arr2 is intersect
func IntersectString(arr1 []string, arr2 []string) bool {
	for _, key1 := range arr1 {
		for _, key2 := range arr2 {
			if key1 == key2 {
				return true
			}
		}
	}

	return false
}

// Contains Is arr1 contain arr2
func Contains(arr1 []interface{}, arr2 []interface{}) bool {
	for _, key2 := range arr2 {
		flag := false
		for _, key1 := range arr1 {
			if key1 == key2 {
				flag = true
				break
			}
		}
		if !flag {
			return false
		}
	}

	return true
}

// StringJoin string join
func StringJoin(arr []string) string {
	result := ""
	for _, value := range arr {
		if len(result) == 0 {
			result = value
		} else {
			result += ";" + value
		}
	}
	return result
}

// RemoveReplicaSliceString slice(string类型)元素去重
func RemoveReplicaSliceString(slc []string) []string {
	result := make([]string, 0)
	tempMap := make(map[string]bool, len(slc))
	for _, e := range slc {
		if tempMap[e] == false && len(e) > 0 {
			tempMap[e] = true
			result = append(result, e)
		}
	}
	return result
}

// ChangeObjectToMap change object to map
func ChangeObjectToMap(obj interface{}) map[string]interface{} {
	data, err := json.Marshal(obj)
	if err != nil {
		return nil
	}

	var mapData map[string]interface{}
	json.Unmarshal(data, &mapData)
	return mapData
}

// ChangeObjectToStringMap change object to string map
func ChangeObjectToStringMap(obj interface{}) map[string]string {
	data, err := json.Marshal(obj)
	if err != nil {
		return nil
	}

	var mapData map[string]string
	json.Unmarshal(data, &mapData)
	return mapData
}

// ChangeObjectToMapList change object to list map
func ChangeObjectToMapList(obj interface{}) []map[string]interface{} {
	data, err := json.Marshal(obj)
	if err != nil {
		return nil
	}

	fmt.Println("----------------------------" + string(data))
	var mapData []map[string]interface{}
	json.Unmarshal(data, &mapData)
	return mapData
}

// ChangeStringArrayToString change string array to string
func ChangeStringArrayToString(strs []string) string {
	result := ""
	for _, str := range strs {
		if result == "" {
			result = str
		} else {
			result += "," + str
		}
	}
	return result
}

// ChangeInterfaceToIDArray change interface to string array
func ChangeInterfaceToIDArray(obj interface{}) []string {
	array := make([]string, 0)
	arr := ChangeObjectToMapList(obj)
	for _, value := range arr {
		array = append(array, value["_id"].(string))
	}

	return array
}

// ChangeInterfaceToValueArray change interface to string array
func ChangeInterfaceToValueArray(obj interface{}) []string {
	array := make([]string, 0)
	arr := ChangeObjectToMapList(obj)
	for _, value := range arr {
		array = append(array, value["value"].(string))
	}

	return array
}

// SliceRemoveElement remove
func SliceRemoveElement(objs []string, obj string) []string {
	for i, e := range objs {
		if e == obj {
			return append(objs[:i], objs[i+1:]...)
		}
	}
	return objs
}

// SliceInsert 在Slice的指定位置插入元素
func SliceInsert(s []interface{}, index int, value interface{}) []interface{} {
	rear := append([]interface{}{}, s[index:]...)
	return append(append(s[:index], value), rear...)
}
