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
	"bytes"
	"encoding/json"
	"net/mail"
)

// Struct2Bytes 结构体转换为字节
func Struct2Bytes(reqData interface{}) (*bytes.Buffer, error) {
	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(reqData)
	if err != nil {
		return nil, err
	}
	return &buf, err
}

// Abs2 绝对值
func Abs2(x int64) int64 {
	if x < 0 {
		return -x
	}
	return x
}

// EmailAddressValid 邮箱格式检查
func EmailAddressValid(email *string) bool {
	_, err := mail.ParseAddress(*email)
	return err == nil
}
