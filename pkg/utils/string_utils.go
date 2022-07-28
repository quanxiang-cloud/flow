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
	"github.com/quanxiang-cloud/flow/pkg/misc/logger"
	"strconv"
	"strings"
)

// StringJoins string join
func StringJoins(strs ...string) string {
	var build strings.Builder
	for _, s := range strs {
		if len(s) > 0 {
			build.WriteString(s)
		}
	}
	return build.String()
}

// GetAsString 类型转换
func GetAsString(v interface{}) string {
	if s, ok := v.(string); ok {
		return s
	}
	logger.Logger.Errorf("invalid string value %#v", v)
	return ""
}

// String2Int string转int
func String2Int(s *string) *int {
	i, err := strconv.Atoi(*s)
	if err != nil {
		logger.Logger.Errorf("convert string: %s to int error. %v", *s, err)
		return nil
	}
	return &i
}
