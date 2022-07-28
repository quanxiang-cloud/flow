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

package node

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/quanxiang-cloud/flow/internal/convert"
	"github.com/quanxiang-cloud/flow/internal/models"
	"github.com/quanxiang-cloud/flow/pkg"
	"github.com/quanxiang-cloud/flow/pkg/client"
	"github.com/quanxiang-cloud/flow/pkg/config"
	"github.com/quanxiang-cloud/flow/pkg/misc/error2"
	"github.com/quanxiang-cloud/flow/pkg/misc/id2"
	"github.com/quanxiang-cloud/flow/pkg/misc/time2"
	"github.com/quanxiang-cloud/flow/pkg/utils"
	"github.com/quanxiang-cloud/flow/rpc/pb"
)

var (
	// MaxaTimeLime 延时一段时间最大支持的天数
	MaxaTimeLime int64 = 86400 * 365
	// codeFmt 创建dispatcher任务的ID格式: {address}:{callback_type}_{callback_tasks.id}
	codeFmt = "flow:%s_%s"
)

// DelayPolicyData DelayPolicyData
type DelayPolicyData struct {
	DelayPolicy DelayPolicy `json:"delayPolicy"`
}

// DelayPolicy DelayPolicy
type DelayPolicy struct {
	Type string     `json:"type"`
	Data PolicyData `json:"data"`
}

// PolicyData PolicyData
type PolicyData struct {
	TimeFmt interface{}  `json:"timeFmt"`
	Column  PolicyColumn `json:"column"`
}

// PolicyColumn PolicyColumn
type PolicyColumn struct {
	TableID  string `json:"tableID"`
	ColumnID string `json:"columnID"`
}

// Delay struct
type Delay struct {
	*Node
}

// NewDelay 延时节点
func NewDelay(conf *config.Configs, node *Node) *Delay {
	return &Delay{
		Node: node,
	}
}

// InitBegin event
func (d *Delay) InitBegin(ctx context.Context, eventData *EventData) (*pb.NodeEventRespData, error) {

	return nil, nil
}

// InitEnd 初始化节点
func (d *Delay) InitEnd(ctx context.Context, eventData *EventData) (*pb.NodeEventRespData, error) {
	nowTimeTs := time2.NowUnix()
	var execTimeString string
	var skipDelay bool = false
	resp := pb.NodeEventRespData{ExecuteType: convert.PauseExecution}
	instance, err := d.InstanceRepo.GetEntityByProcessInstanceID(d.Db, eventData.ProcessInstanceID)
	if err != nil {
		return nil, err
	}
	dp, err := d.ConvertPolicy2Struct(&eventData.Shape.Data.BusinessData)
	if err != nil {
		return nil, err
	}
	if dp == nil {
		return nil, errors.New("businessData is nil")
	}
	// 时间格式检查
	err = d.TimeCheck(dp)
	if err != nil {
		return nil, nil
	}
	switch dp.DelayPolicy.Type {
	// execTime = nowTime + delayTime
	case convert.DelayTypeOfaTime:
		execTimeTs := d.ExecTimeOfaTime(dp)
		execTimeString = time2.UnixToISO8601(*execTimeTs)
	// execTime = ColumnTime + delayTime
	case convert.DelayTypeOfTableColumn:
		dataReq := client.FormDataConditionModel{
			AppID:   instance.AppID,
			TableID: instance.FormID,
			DataID:  instance.FormInstanceID,
		}
		formResp, err := d.FormAPI.GetFormData(ctx, dataReq)
		if err != nil {
			return nil, err
		}
		if formResp == nil {
			return nil, nil
		}
		fieldTime := formResp[dp.DelayPolicy.Data.Column.ColumnID]
		if fieldTime == nil {
			return nil, fmt.Errorf("column id %s is nil", dp.DelayPolicy.Data.Column.ColumnID)
		}
		s := fieldTime.(string)
		execTimeTs := d.ExecTimeOfTableColumn(dp, &s)
		// 小于当前时间直接跳过该节点
		if *execTimeTs <= nowTimeTs {
			skipDelay = true
		} else {
			execTimeString = time2.UnixToISO8601(*execTimeTs)
		}
	// execTime = user define time
	case convert.DelayTypeOfSpecTime:
		// 小于当前时间需跳过
		execTimeTs := d.ExecTimeOfSpecTime(dp)
		if *execTimeTs <= nowTimeTs {
			skipDelay = true
		} else {
			execTimeString = time2.UnixToISO8601(*execTimeTs)
		}
	}
	if skipDelay == false {
		otherInfo, err := d.completeNodeReq(ctx, instance, eventData)
		if err != nil {
			return nil, err
		}
		// 回调时需要用到的数据
		dbData := models.DispatcherCallback{
			Type:              dp.DelayPolicy.Type,
			OtherInfo:         *otherInfo,
			ProcessInstanceID: eventData.ProcessInstanceID,
		}
		ID := id2.GenID()
		dbData.ID = *d.genNewCode(&ID)
		dbData.CreateTime = time2.Now()
		tx := d.Db.Begin()
		if err := d.createDelayJob(ctx, &execTimeString, &dbData.ID); err != nil {
			return nil, err
		}
		if err := d.DispatcherCallbackRepo.Create(d.Db, &dbData); err != nil {
			return nil, err
		}
		tx.Commit()
		// close the transaction
		defer func() {
			if err != nil {
				tx.Rollback()
			}
		}()
		return &resp, nil
	}
	return nil, nil
}

// getFormInstance 从form获取表单实例
func (d *Delay) getFormInstance(ctx *context.Context, instanceInfo *models.Instance) (*map[string]interface{}, error) {
	if instanceInfo == nil {
		return nil, nil
	}
	dataReq := client.FormDataConditionModel{
		AppID:   instanceInfo.AppID,
		TableID: instanceInfo.FormID,
		DataID:  instanceInfo.FormInstanceID,
	}

	dataResp, err := d.FormAPI.GetFormData(*ctx, dataReq)
	if err != nil {
		return nil, err
	}
	if dataResp == nil {
		return nil, error2.NewErrorWithString(error2.Internal, "not found from data")
	}
	return &dataResp, nil
}

// ConvertPolicy2Struct 将业务数据转换为结构体
func (d *Delay) ConvertPolicy2Struct(bd interface{}) (*DelayPolicyData, error) {
	var resp DelayPolicyData
	b, err := json.Marshal(bd)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(b, &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

// TimeCheck 请求的时间有效性检查
func (d *Delay) TimeCheck(dp *DelayPolicyData) error {
	Type := dp.DelayPolicy.Type
	timeFmt := dp.DelayPolicy.Data.TimeFmt
	switch Type {
	// "2022-04-22T07:49:12+0000"
	case convert.DelayTypeOfSpecTime:
		t := timeFmt.(string)
		result := utils.ISO8601FmtCheck(&t)
		if result == false {
			return fmt.Errorf("invalid time format. except time format: `2022-04-22T07:49:12+0000` got: %s", timeFmt)
		}
		return nil
	// 86400 延时多少秒
	case convert.DelayTypeOfaTime:
		t := timeFmt.(int64)
		if t <= 0 || t > MaxaTimeLime {
			return fmt.Errorf("invalid time format. $time <= 0 || $time > %d", MaxaTimeLime)
		}
		return nil
	// 	86400 延时/提前多少秒
	case convert.DelayTypeOfTableColumn:
		t := utils.Abs2(timeFmt.(int64))
		if t > MaxaTimeLime {
			return fmt.Errorf("invalid time format. $time >= %d", MaxaTimeLime)
		}
		return nil
	default:
		return fmt.Errorf("unknown type of: %s", Type)
	}
}

// ExecTimeOfaTime 延时一段时间的执行时间戳
func (d *Delay) ExecTimeOfaTime(dp *DelayPolicyData) *int64 {
	t := dp.DelayPolicy.Data.TimeFmt.(int64)
	now := time2.NowUnix()
	// 延时到某个时间
	execTime := now + t
	return &execTime
}

// ExecTimeOfTableColumn 按日期字段延时执行时间
func (d *Delay) ExecTimeOfTableColumn(dp *DelayPolicyData, ColumnTime *string) *int64 {
	t := dp.DelayPolicy.Data.TimeFmt.(int64)
	cTime, _ := time2.ISO8601ToUnix(*ColumnTime)
	execTime := cTime + t
	return &execTime
}

// ExecTimeOfSpecTime 延时到指定时间执行时间
func (d *Delay) ExecTimeOfSpecTime(dp *DelayPolicyData) *int64 {
	t := dp.DelayPolicy.Data.TimeFmt.(string)
	execTime, _ := time2.ISO8601ToUnix(t)

	return &execTime
}

// createDelayJob 创建延时任务
// @params
// execTime:  	iso8601 format
// id: 			callback_tasks.id
func (d *Delay) createDelayJob(ctx context.Context, execTime *string, id *string) error {
	code := d.genNewCode(id)
	req := client.TaskPostReq{
		Code:       *code,
		Type:       1,
		TimeBar:    *execTime,
		State:      1,
		Retry:      3,
		RetryDelay: 60,
	}
	return d.Dispatcher.TakePost(ctx, req)
}

// completeNodeParams 生成params
func (d *Delay) completeNodeParams(ctx context.Context, instance *models.Instance) (*map[string]interface{}, error) {
	// 生成params
	params, err := d.Flow.GetInstanceVariableValues(ctx, instance)
	if err != nil {
		return nil, err
	}
	p := d.Flow.FormatFormValue(instance, params)
	return &p, nil
}

// completeNodeReq 生成调用process completeNode接口请求数据
// return: json string
func (d *Delay) completeNodeReq(ctx context.Context, instance *models.Instance, eventData *EventData) (*string, error) {
	params, err := d.completeNodeParams(ctx, instance)
	if err != nil {
		return nil, err
	}
	req := client.CompleteNodeReq{
		ProcessID:   eventData.ProcessID,
		InstanceID:  eventData.ProcessInstanceID,
		NodeDefKey:  eventData.NodeDefKey,
		ExecutionID: eventData.ExecutionID,
		NextNodes:   "", // not need
		Params:      *params,
		UserID:      pkg.STDUserID(ctx),
	}
	if b, err := json.Marshal(&req); err != nil {
		return nil, err
	} else {
		s := string(b)
		return &s, nil
	}
}

// genID 生成ID
func (d *Delay) genNewCode(ID *string) *string {
	code := fmt.Sprintf(codeFmt, convert.CallbackOfDelay, *ID)

	return &code
}
