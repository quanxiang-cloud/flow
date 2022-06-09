package client

import (
	"context"
	"fmt"
	"github.com/quanxiang-cloud/flow/pkg/config"
	"github.com/quanxiang-cloud/flow/pkg/misc/id2"
	"github.com/quanxiang-cloud/flow/pkg/misc/logger"
	"github.com/quanxiang-cloud/flow/pkg/redis"
	"github.com/quanxiang-cloud/flow/pkg/utils"
	"net/http"
	"reflect"
	"strings"
)

// Form form client
type Form interface {
	BatchGetFormData(c context.Context, req []FormDataConditionModel) (map[string]map[string]interface{}, error)
	GetFormData(c context.Context, req FormDataConditionModel) (map[string]interface{}, error)

	UpdateData(c context.Context, appID string, tableID string, dataID string, req UpdateEntity, eventTrigger bool) error
	CreateData(c context.Context, appID string, tableID string, req CreateEntity, eventTrigger bool) error
	DeleteData(c context.Context, appID string, tableID string, dataIDs []string) error
	FindOneData(c context.Context, appID string, tableID string, dataID string, ref interface{}) (interface{}, error)
	SearchData(c context.Context, appID string, tableID string, req SearchReq) ([]map[string]interface{}, error)
	GetIDs(c context.Context, appID string, tableID string, query map[string]interface{}) ([]string, error)

	BatchGetFormSchema(c context.Context, req []FormSchemaConditionModel) (map[string]interface{}, error)
	GetFormSchema(c context.Context, appID string, tableID string) (interface{}, error)

	GetValueByFieldFormat(fields []string, conditionValue interface{}) interface{}
	GetValue(formData map[string]interface{}, fieldKey string, conditionValue interface{}) (interface{}, interface{})
}

type form struct {
	conf   *config.Configs
	client http.Client
}

// NewForm new
func NewForm(conf *config.Configs) Form {
	return &form{
		conf:   conf,
		client: NewClient(conf.InternalNet),
	}
}

func (f *form) getHeader(c context.Context, eventTrigger bool, method string) map[string]string {
	if !eventTrigger {
		seq := id2.GenID()
		redis.SetValueToRedis(c, seq, "ignore")
		return map[string]string{
			"version": "0.1",
			"method":  method,
			"seq":     seq,
			"magic":   "flow_event_trigger",
		}
	}
	return nil
}

func (f *form) UpdateData(c context.Context, appID string, tableID string, dataID string, req UpdateEntity, eventTrigger bool) error {
	if req.Query == nil {
		req.Query = map[string]interface{}{
			"term": map[string]interface{}{
				"_id": dataID,
			},
		}
	}

	entityMap := utils.ChangeObjectToMap(req.Entity)
	if len(entityMap) == 0 && len(req.Ref) == 0 {
		return nil
	}

	var resp interface{}
	url := fmt.Sprintf("%s%s%s%s%s%s", f.conf.APIHost.FormHost, "api/v1/form/", appID, "/internal/form/", tableID, "/update")
	err := POST2(c, &f.client, url, req, &resp, f.getHeader(c, eventTrigger, "put"))
	if err != nil {
		logger.Logger.Error("Failed to update form data ", err)
	}
	return err
}

func (f *form) CreateData(c context.Context, appID string, tableID string, req CreateEntity, eventTrigger bool) error {
	var resp interface{}
	url := fmt.Sprintf("%s%s%s%s%s%s", f.conf.APIHost.FormHost, "api/v1/form/", appID, "/internal/form/", tableID, "/create")
	err := POST2(c, &f.client, url, req, &resp, f.getHeader(c, eventTrigger, "post"))
	if err != nil {
		logger.Logger.Error("Failed to create form data ", err)
	}
	return err
}

func (f *form) DeleteData(c context.Context, appID string, tableID string, dataIDs []string) error {
	queryMap := map[string]interface{}{
		"query": map[string]interface{}{
			"terms": map[string]interface{}{
				"_id": dataIDs,
			},
		},
	}

	var resp interface{}
	url := fmt.Sprintf("%s%s%s%s%s%s", f.conf.APIHost.FormHost, "api/v1/form/", appID, "/internal/form/", tableID, "/delete")
	err := POST(c, &f.client, url, queryMap, &resp)
	if err != nil {
		logger.Logger.Error("Failed to delete form data ", err)
	}
	return err
}

// FindOneData find one data
func (f *form) FindOneData(c context.Context, appID string, tableID string, dataID string, ref interface{}) (interface{}, error) {
	queryMap := map[string]interface{}{
		"query": map[string]interface{}{
			"term": map[string]interface{}{
				"_id": dataID,
			},
		},
		"ref": ref,
	}

	var resp FindOneDataResp
	url := fmt.Sprintf("%s%s%s%s%s%s", f.conf.APIHost.FormHost, "api/v1/form/", appID, "/internal/form/", tableID, "/get")
	err := POST(c, &f.client, url, queryMap, &resp)
	if err != nil {
		logger.Logger.Error("Failed to get form data ", err)
		return nil, err
	}
	return resp.Entity, nil
}

func (f *form) SearchData(c context.Context, appID string, tableID string, req SearchReq) ([]map[string]interface{}, error) {
	var resp SearchResp
	url := fmt.Sprintf("%s%s%s%s%s%s", f.conf.APIHost.FormHost, "api/v1/form/", appID, "/internal/form/", tableID, "/search")
	err := POST(c, &f.client, url, req, &resp)
	if err != nil {
		logger.Logger.Error("Failed to get form data ", err)
		return nil, err
	}
	return resp.Entities, nil
}

func (f *form) GetIDs(c context.Context, appID string, tableID string, query map[string]interface{}) ([]string, error) {
	req := SearchReq{
		Page:  1,
		Size:  99999,
		Query: query,
	}
	list, err := f.SearchData(c, appID, tableID, req)
	if err != nil {
		return nil, err
	}

	IDs := make([]string, 0)
	for _, value := range list {
		IDs = append(IDs, value["_id"].(string))
	}
	return IDs, nil
}

// BatchGetFormData batch get form data
func (f *form) BatchGetFormData(c context.Context, req []FormDataConditionModel) (map[string]map[string]interface{}, error) {
	if len(req) == 0 {
		return nil, nil
	}

	result := make(map[string]map[string]interface{}, 0)
	for _, value := range req {
		data, _ := f.GetFormData(c, value)

		if data != nil {
			result[value.DataID] = data
		}
	}

	return result, nil
}

// GetFormData get form data
func (f *form) GetFormData(c context.Context, req FormDataConditionModel) (map[string]interface{}, error) {
	data, err := f.FindOneData(c, req.AppID, req.TableID, req.DataID, req.Ref)
	if err != nil {
		return nil, err
	}
	return utils.ChangeObjectToMap(data), nil
}

// BatchGetFormSchema batch get form schema
func (f *form) BatchGetFormSchema(c context.Context, req []FormSchemaConditionModel) (map[string]interface{}, error) {
	if len(req) == 0 {
		return nil, nil
	}

	result := make(map[string]interface{}, 0)
	for _, value := range req {
		data, _ := f.GetFormSchema(c, value.AppID, value.TableID)

		if data != nil {
			result[value.TableID] = data
		}
	}

	return result, nil
}

// GetFormSchema get form schema
func (f *form) GetFormSchema(c context.Context, appID string, tableID string) (interface{}, error) {
	var resp GetFormSchemaResp
	url := fmt.Sprintf("%s%s%s%s%s", f.conf.APIHost.FormHost, "api/v1/form/", appID, "/internal/schema/", tableID)
	err := POST(c, &f.client, url, nil, &resp)
	if err != nil {
		logger.Logger.Error("Failed to get form schema ", err)
		return nil, err
	}
	return resp.Schema, nil
}

func (f *form) GetValueByFieldFormat(fields []string, conditionValue interface{}) interface{} {
	if len(fields) > 1 {
		if fields[1] == "[]" {
			cValues := make([]interface{}, 0)
			if conditionValue != nil {
				cvs := reflect.ValueOf(conditionValue)
				for i := 0; i < cvs.Len(); i++ {
					v := cvs.Index(i).Interface()
					vMap := utils.ChangeObjectToMap(v)
					cValues = append(cValues, vMap[fields[2]])
				}
			}
			return cValues
		}

		if conditionValue != nil {
			cValueMap := utils.ChangeObjectToMap(conditionValue)
			return cValueMap[fields[1]]
		}
	}
	return conditionValue
}

// GetValue get value from form
func (f *form) GetValue(formData map[string]interface{}, fieldKey string, conditionValue interface{}) (interface{}, interface{}) {
	if len(fieldKey) == 0 {
		return nil, nil
	}

	// fieldxxx, fieldxxx.value, fieldxxx.[].value
	fields := strings.Split(fieldKey, ".")

	var v1 interface{}
	if formData != nil {
		v1 = f.GetValueByFieldFormat(fields, formData[fields[0]])
	}

	v2 := f.GetValueByFieldFormat(fields, conditionValue)
	return v1, v2
}

// GetFormSchemaResp GetFormSchema resp
type GetFormSchemaResp struct {
	ID     string      `json:"id"`
	Schema interface{} `json:"schema"`
}

// FormSchemaConditionModel find form schema condition model
type FormSchemaConditionModel struct {
	AppID   string `json:"appID"`
	TableID string `json:"tableID"`
}

// FindOneDataResp FindOneData resp
type FindOneDataResp struct {
	Entity     interface{} `json:"entity"`
	ErrorCount int64       `json:"errorCount"`
}

// FormDataConditionModel find form data condition model
type FormDataConditionModel struct {
	AppID   string      `json:"appID"`
	TableID string      `json:"tableID"`
	DataID  string      `json:"dataID"`
	Ref     interface{} `json:"ref"`
}

// CreateEntity create
type CreateEntity struct {
	Entity interface{}        `json:"entity"`
	Ref    map[string]RefData `json:"ref"`
}

// UpdateEntity update
type UpdateEntity struct {
	Entity interface{}        `json:"entity"`
	Query  interface{}        `json:"query"`
	Ref    map[string]RefData `json:"ref"`
}

// RefData sub form operation model
type RefData struct {
	AppID   string         `json:"appID"`
	TableID string         `json:"tableID"`
	Type    string         `json:"type"`
	New     interface{}    `json:"new"` // []string    or   []CreateEntity
	Deleted []string       `json:"deleted"`
	Updated []UpdateEntity `json:"updated"`
}

// SearchReq search req
type SearchReq struct {
	Page  int         `json:"page"`
	Size  int         `json:"size"`
	Query interface{} `json:"query"`
}

// SearchResp search resp
type SearchResp struct {
	Entities     []map[string]interface{} `json:"entities"`
	Total        int                      `json:"total"`
	Aggregations int                      `json:"aggregations"`
}
