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
