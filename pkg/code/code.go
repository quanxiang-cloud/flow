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

package code

import "github.com/quanxiang-cloud/flow/pkg/misc/error2"

func init() {
	error2.CodeTable = CodeTable
}

// 应用标识：标注错误来自哪个应用，用三位数表示
// 功能区域：错误来自区域的哪个功能模块，用三位数表示
// 状态码/错误类型：错误是哪种类型，比如001 参数输入错误，002: 业务错误，003: 权限错误，用三位数字表示。
// 具体错误：错误类型下的，具体错误信息，用4位数字表示。
const (
	// InvalidURI 无效的URI
	InvalidURI = 70014000000
	// InvalidParams 无效的参数
	InvalidParams = 70014000001
	// InvalidTimestamp 无效的时间格式
	InvalidTimestamp = 70014000002
	// NameExist 名字已经存在
	NameExist = 70014000003
	// InvalidDel 无效的删除
	InvalidDel = 70014000004
	// InvalidProcessID process id no exist
	InvalidProcessID = 70014000005
	// InvalidTaskID task id no exist
	InvalidTaskID = 70014000006
	// NoMessagesAllowed No messages allowed
	NoMessagesAllowed = 70014000007
	// InvalidInstanceID process id no exist
	InvalidInstanceID = 70014000008

	// flow review error
	// CannotCcToSelf can not cc to self
	CannotCcToSelf = 70020020001
	// CannotInviteReadToSelf can not invite read to self
	CannotInviteReadToSelf = 70020020002
	// CannotRepeatInviteRead can not repeat invite read
	CannotRepeatInviteRead = 70020020003

	// CannotRepeatUrge can not repeat urge
	CannotRepeatUrge = 70030020001

	// Task handle error begin------------------
	TaskCannotFind = 70040020001
	// Task handle error end------------------
)

// CodeTable 码表
var CodeTable = map[int64]string{
	InvalidURI:       "无效的URI.",
	InvalidParams:    "无效的参数.",
	InvalidTimestamp: "无效的时间格式.",
	NameExist:        "名称已被使用！请检查后重试！",
	InvalidDel:       "删除无效！对象不存在或请检查参数！",

	CannotCcToSelf:         "不能抄送给自己",
	CannotInviteReadToSelf: "不能邀请自己阅示",
	CannotRepeatInviteRead: "不能重复邀请",
	CannotRepeatUrge:       "已经催办过了",

	TaskCannotFind: "任务已被处理或者不存在，请刷新重试",
}
