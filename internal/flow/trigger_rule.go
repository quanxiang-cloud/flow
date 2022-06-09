package flow

import (
	"context"
	"encoding/json"
	"github.com/quanxiang-cloud/flow/internal/models"
	"github.com/quanxiang-cloud/flow/internal/models/mysql"
	"github.com/quanxiang-cloud/flow/internal/server/options"
	"github.com/quanxiang-cloud/flow/pkg/client"
	"github.com/quanxiang-cloud/flow/pkg/config"
	"github.com/quanxiang-cloud/flow/pkg/misc/logger"
	"github.com/quanxiang-cloud/flow/pkg/utils"
	"gorm.io/gorm"
	"strings"
)

// TriggerRule service
type TriggerRule interface {
	CheckDataModify(ctx context.Context, req *FormMsg) error
}

type triggerRule struct {
	db              *gorm.DB
	triggerRuleRepo models.TriggerRuleRepo
	instance        Instance
	formAPI         client.Form
}

// NewTriggerRule init
func NewTriggerRule(conf *config.Configs, opts ...options.Options) (TriggerRule, error) {
	instance, _ := NewInstance(conf, opts...)
	tr := &triggerRule{
		triggerRuleRepo: mysql.NewTriggerRuleRepo(),
		instance:        instance,
		formAPI:         client.NewForm(conf),
	}

	for _, opt := range opts {
		opt(tr)
	}
	return tr, nil
}

func (tr *triggerRule) checkTrigger(ctx context.Context, rule *models.TriggerRule, req *FormMsg, triggerType string) {
	if req.Entity == nil {
		return
	}

	triggerModel := &TriggerModel{}
	err := json.Unmarshal([]byte(rule.Rule), triggerModel)
	if err != nil {
		logger.Logger.Errorw(err.Error())
		return
	}

	entityMap := utils.ChangeObjectToMap(req.Entity)
	if tr.checkRule(entityMap, triggerModel, req.Method) {
		userID := ""
		if triggerType == "post" {
			userID = entityMap["creator_id"].(string)
		} else {
			userID = entityMap["modifier_id"].(string)
		}

		startFlowModel := &StartFlowModel{
			FlowID:   rule.FlowID,
			FormData: entityMap,
			UserID:   userID,
		}
		_, err := tr.instance.StartFlow(ctx, startFlowModel)
		if err != nil {
			logger.Logger.Errorw(err.Error())
		}
	}
}

func (tr *triggerRule) checkRule(formData map[string]interface{}, model *TriggerModel, t string) bool {
	if "post" == t {
		if !utils.Contain(model.TriggerWay, "whenAdd") {
			return false
		}
	} else if "put" == t {
		if !utils.Contain(model.TriggerWay, "whenAlter") {
			return false
		}

		if len(model.WhenAlterFields) > 0 {
			keyList := utils.GetMapKeys(formData)
			if !utils.IntersectString(keyList, model.WhenAlterFields) {
				return false
			}
		}
	}

	if len(model.TriggerCondition.Op) == 0 {
		model.TriggerCondition.Op = "or"
	}

	return tr.executeCondition(model.TriggerCondition, formData)
}

func (tr *triggerRule) executeCondition(model TriggerConditionModel, formData map[string]interface{}) bool {
	if model.Op == "eq" || model.Op == "neq" || model.Op == "lt" || model.Op == "gt" {
		if len(model.Key) == 0 || utils.IsNil(model.Value) {
			return true
		}
	}

	var fieldValue, conditionValue interface{}
	isArray := false
	if len(model.Key) > 0 {
		fieldValue, conditionValue = tr.formAPI.GetValue(formData, model.Key, model.Value)

		fields := strings.Split(model.Key, ".")
		if len(fields) > 1 && fields[1] == "[]" {
			isArray = true
		}
	}

	if model.Op == "eq" {
		return conditionValue == fieldValue
	} else if model.Op == "neq" {
		return conditionValue != fieldValue
	} else if model.Op == "lt" {
		return utils.ExprCompare(conditionValue, fieldValue) > 0
	} else if model.Op == "gt" {
		return utils.ExprCompare(conditionValue, fieldValue) < 0
	} else if model.Op == "or" {
		if len(model.Expr) > 0 {
			for _, conditionModel := range model.Expr {
				if tr.executeCondition(conditionModel, formData) {
					return true
				}
			}
			return false
		}
	} else if model.Op == "and" {
		if len(model.Expr) > 0 {
			for _, conditionModel := range model.Expr {
				if !tr.executeCondition(conditionModel, formData) {
					return false
				}
			}
			return true
		}
	} else if model.Op == "gte" {
		return utils.ExprCompare(conditionValue, fieldValue) <= 0
	} else if model.Op == "lte" {
		return utils.ExprCompare(conditionValue, fieldValue) >= 0
	} else if model.Op == "null" {
		return utils.ExprNull(isArray, fieldValue)
	} else if model.Op == "not-null" {
		return !utils.ExprNull(isArray, fieldValue)
	} else if model.Op == "include" || model.Op == "all" {
		return utils.ExprInclude(isArray, fieldValue, conditionValue)
	} else if model.Op == "not-include" {
		return !utils.ExprInclude(isArray, fieldValue, conditionValue)
	} else if model.Op == "any" {
		return utils.ExprAnyInclude(isArray, fieldValue, conditionValue)
	} else if model.Op == "range" {
		var conditionDataList []string
		err := json.Unmarshal([]byte(utils.Strval(conditionValue)), conditionDataList)
		if err != nil {
			logger.Logger.Errorw(err.Error())
		}
		if len(conditionDataList) == 0 {
			return false
		}

		beginDate := conditionDataList[0]
		endDate := conditionDataList[1]

		return strings.Compare(beginDate, fieldValue.(string)) <= 0 && strings.Compare(endDate, fieldValue.(string)) >= 0
	}
	return true
}

// CheckDataModify check form data modify
func (tr *triggerRule) CheckDataModify(ctx context.Context, req *FormMsg) error {
	if len(req.TableID) == 0 {
		return nil
	}

	rules, err := tr.triggerRuleRepo.FindTriggerRulesByFormID(tr.db, req.TableID)
	if err != nil {
		return err
	}

	if len(rules) == 0 {
		return nil
	}

	switch req.Method {
	case "post":
		{
			for _, value := range rules {
				tr.checkTrigger(ctx, value, req, "post")
			}
		}
	case "put":
		{
			for _, value := range rules {
				tr.checkTrigger(ctx, value, req, "put")
			}
		}
	case "delete":
		{
			entityMap := utils.ChangeObjectToMap(req.Entity)
			formDataIDs := make([]string, 0)
			if entityMap["data"] == nil {
				return nil
			}

			v, ok := entityMap["data"].([]interface{})
			if !ok {
				return nil
			}
			for _, value := range v {
				formDataIDs = append(formDataIDs, value.(string))
			}
			return tr.instance.FormDataDeleted(ctx, formDataIDs, entityMap["delete_id"].(string))
		}
	}

	return nil
}

// SetDB set db
func (tr *triggerRule) SetDB(db *gorm.DB) {
	tr.db = db
}

// FormMsg form msg
type FormMsg struct {
	TableID string      `json:"tableID"`
	Entity  interface{} `json:"entity"`
	Magic   string      `json:"magic"`
	Seq     string      `json:"seq"`
	Version string      `json:"version"`
	Method  string      `json:"method"` // post 对应的就是新增，put 是修改，delete 是删除，
}

// EventModel kafka msg event model
type EventModel struct {
	EventType string    `json:"eventType"`
	EventName string    `json:"eventName"`
	Data      DataModel `json:"data"`
}

// DataModel kafka msg data model
type DataModel struct {
	TableID string                   `json:"tableID"`
	Entity  []map[string]interface{} `json:"entity"`
	Method  string                   `json:"method"` // create,update,delete
	// UserID  string                   `json:"userID"`
	Topic string `json:"topic"`
	Event string `json:"event"` // dataModify
}

// TriggerModel trigger model
type TriggerModel struct {
	TriggerWay       []string              `json:"triggerWay"`
	WhenAlterFields  []string              `json:"whenAlterFields"`
	TriggerCondition TriggerConditionModel `json:"triggerCondition"`
}

// TriggerConditionModel trigger condition model
type TriggerConditionModel struct {
	Op    string                  `json:"op"` // or, and, eq,lt,gt,neq
	Expr  []TriggerConditionModel `json:"expr"`
	Key   string                  `json:"key"`
	Value interface{}             `json:"value"`
}
