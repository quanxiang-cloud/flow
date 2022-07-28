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

// GetMapKeys get map keys
func GetMapKeys(data map[string]interface{}) []string {
	var keys = make([]string, 0)
	for key := range data {
		keys = append(keys, key)
	}
	return keys
}

// MergeMap merge map
func MergeMap(mObj ...map[string]interface{}) map[string]interface{} {
	newObj := map[string]interface{}{}
	for _, m := range mObj {
		if m != nil {
			for k, v := range m {
				newObj[k] = v
			}
		}
	}
	return newObj
}
