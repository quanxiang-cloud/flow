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
	"fmt"
	"testing"
)

// func Test_GetFieldValueArr(t *testing.T) {
// 	arr:=make([]interface{},0)
// 	arr=append(arr, map[string]interface{}{
// 		"openUserId":111,
// 	})
// 	arr=append(arr, map[string]interface{}{
// 		"openUserId":222,
// 	})
// 	data:=map[string]interface{}{
// 		"$webhook0ea9vOcw8z2LWnDA_rZm2.empList":arr,
// 	}
// 	result:=GetFieldValue(data,"$webhook0ea9vOcw8z2LWnDA_rZm2.empList.[0].openUserId")
// 	fmt.Println(result)
//
// }

func Test_GetFieldValue(t *testing.T) {
	data := map[string]interface{}{
		"$webhook0ea9vOcw8z2LWnDA_rZm2.empList": map[string]interface{}{
			"openUserId": 111,
		},
	}
	result := GetFieldValue(data, "$webhook0ea9vOcw8z2LWnDA_rZm2.empList.openUserId")
	fmt.Println(result)

}
