package client

import (
	"context"
	"github.com/quanxiang-cloud/flow/internal"
	"github.com/quanxiang-cloud/flow/pkg/config"
	"github.com/quanxiang-cloud/flow/pkg/utils"
	"net/http"
	"strings"
)

// Identity interface
type Identity interface {
	FindUsersByGroup(ctx context.Context, groupID string) ([]*UserInfoResp, error)
	FindUsersByGroups(ctx context.Context, groupIDs []string) ([]*UserInfoResp, error)
	FindUserIDsByGroups(ctx context.Context, groupIDs []string) ([]string, error)
	FindGroupsByUserID(ctx context.Context, userID string) (*Group, error)
	FindUsersByIDs(ctx context.Context, userIDs []string) (map[string]*UserInfoResp, error)

	FindDepByUserID(ctx context.Context, userID string) (*Dep, error)
	FindUserByID(ctx context.Context, userID string) (*UserInfoResp, error)
	GetSuperior(ctx context.Context, userID string) (string, error)
	GetLeadOfDepartment(ctx context.Context, userID string) (string, error)
	FindUsersByDepartID(ctx context.Context, departID string) ([]string, error)
	ValidateUserIDs(ctx context.Context, userIDs []string) ([]string, error)

	AddUserInfo(ctx context.Context, data []map[string]interface{}) error
}

type identity struct {
	getGroupsByUser  string
	getUsersByGroups string
	client           http.Client
}

// NewIdentity init
func NewIdentity(conf *config.Configs) Identity {
	i := &identity{
		client:           NewClient(conf.InternalNet),
		getGroupsByUser:  conf.APIHost.OrgHost + "api/v1/org/o/user/ids",
		getUsersByGroups: conf.APIHost.OrgHost + "api/v1/org/o/user/dep/id",
	}
	return i
}

// FindUsersByGroup find users by group
func (i *identity) FindUsersByGroup(ctx context.Context, groupID string) ([]*UserInfoResp, error) {
	var resp UsersResp
	req := map[string]interface{}{"depID": groupID, "isIncludeChild": 0}
	err := POST(ctx, &i.client, i.getUsersByGroups, req, &resp)
	if err != nil {
		return nil, err
	}
	return resp.Users, nil
}

// FindUsersByGroups find users by groups
func (i *identity) FindUsersByGroups(ctx context.Context, groupIDs []string) ([]*UserInfoResp, error) {
	var resp UsersResp

	err := POST(ctx, &i.client, i.getUsersByGroups, groupIDs, &resp)
	if err != nil {
		return nil, err
	}

	return resp.Users, nil
}

// FindUserIDsByGroups find userIDs by groups
func (i *identity) FindUserIDsByGroups(ctx context.Context, groupIDs []string) ([]string, error) {
	users, err := i.FindUsersByGroups(ctx, groupIDs)
	if err != nil {
		return nil, err
	}

	userIDs := make([]string, 0)
	if len(users) > 0 {
		for _, value := range users {
			userIDs = append(userIDs, value.ID)
		}
	}

	return userIDs, nil
}

func (i *identity) FindGroupsByUserID(ctx context.Context, userID string) (*Group, error) {
	if len(userID) == 0 {
		return nil, nil
	}

	var resp UsersResp
	req := map[string][]string{"ids": {userID}}
	err := POST(ctx, &i.client, i.getGroupsByUser, req, &resp)
	if err != nil {
		return nil, err
	}
	if len(resp.Users) > 0 && len(resp.Users[0].Deps) > 0 && len(resp.Users[0].Deps[0]) > 0 {
		dep := resp.Users[0].Deps[0][0]
		group := &Group{
			ID:   strings.Join([]string{internal.Dep, "_", dep.ID}, ""),
			Name: dep.DepartmentName,
			Type: internal.Dep,
		}
		return group, nil
	}
	return nil, nil
}

func (i *identity) FindUsersByIDs(ctx context.Context, userIDs []string) (map[string]*UserInfoResp, error) {
	// remove repeat
	userIDs = utils.RemoveReplicaSliceString(userIDs)
	if len(userIDs) == 0 {
		return nil, nil
	}

	var resp UsersResp
	req := map[string][]string{"ids": userIDs}
	err := POST(ctx, &i.client, i.getGroupsByUser, req, &resp)
	if err != nil {
		return nil, err
	}

	userMap := make(map[string]*UserInfoResp, 0)
	if len(resp.Users) > 0 {
		for _, value := range resp.Users {
			if value.UseStatus == 1 {
				userMap[value.ID] = value
			}
		}
	}
	return userMap, nil
}

func (i *identity) FindDepByUserID(ctx context.Context, userID string) (*Dep, error) {
	if len(userID) == 0 {
		return nil, nil
	}

	var resp UsersResp
	req := map[string][]string{"ids": {userID}}
	err := POST(ctx, &i.client, i.getGroupsByUser, req, &resp)
	if err != nil {
		return nil, err
	}
	if len(resp.Users) > 0 && len(resp.Users[0].Deps) > 0 && len(resp.Users[0].Deps[0]) > 0 {
		return resp.Users[0].Deps[0][0], nil
	}
	return nil, nil
}

func (i *identity) FindUserByID(ctx context.Context, userID string) (*UserInfoResp, error) {
	userMap, err := i.FindUsersByIDs(ctx, []string{userID})
	if err != nil {
		return nil, err
	}

	return userMap[userID], nil
}

func (i *identity) GetSuperior(ctx context.Context, userID string) (string, error) {
	return i.GetLeadOfDepartment(ctx, userID)
}

func (i *identity) GetLeadOfDepartment(ctx context.Context, userID string) (string, error) {
	user, err := i.FindUserByID(ctx, userID)
	if err != nil {
		return "", err
	}

	if user == nil {
		return "", nil
	}

	if len(user.Leaders) > 0 && len(user.Leaders[0]) > 0 && user.Leaders[0][0].UseStatus == 1 {
		return user.Leaders[0][0].ID, nil
	}

	return "", nil
}

func (i *identity) FindUsersByDepartID(ctx context.Context, departID string) ([]string, error) {
	var resp UsersResp
	params := map[string]interface{}{"depID": departID, "includeChildDEPChild": 0}
	err := POST(ctx, &i.client, i.getUsersByGroups, params, &resp)
	if err != nil {
		return nil, err
	}

	users := make([]string, 0)
	for _, value := range resp.Users {
		if value.UseStatus == 1 {
			users = append(users, value.ID)
		}
	}
	return users, nil
}

func (i *identity) ValidateUserIDs(ctx context.Context, userIDs []string) ([]string, error) {
	userMap, err := i.FindUsersByIDs(ctx, userIDs)
	if err != nil {
		return nil, err
	}

	users := make([]string, 0)
	for k := range userMap {
		users = append(users, k)
	}
	return users, nil
}

func (i *identity) AddUserInfo(ctx context.Context, list []map[string]interface{}) error {
	if len(list) > 0 {
		userIDs := make([]string, 0)

		for _, baseEntity := range list {
			if _, ok := baseEntity["creatorId"]; ok && len(baseEntity["creatorId"].(string)) > 0 {
				userIDs = append(userIDs, baseEntity["creatorId"].(string))
			}
			if _, ok := baseEntity["modifierId"]; ok && len(baseEntity["modifierId"].(string)) > 0 {
				userIDs = append(userIDs, baseEntity["modifierId"].(string))
			}

			assignee, ok := baseEntity["assignee"]
			if ok && assignee != nil && len(assignee.(string)) > 0 {
				userIDs = append(userIDs, assignee.(string))
			}
		}
		if len(userIDs) > 0 {
			userMap, err := i.FindUsersByIDs(ctx, userIDs)
			if err != nil {
				return err
			}
			for _, baseEntity := range list {
				if _, ok := baseEntity["creatorId"]; ok && len(baseEntity["creatorId"].(string)) > 0 {
					userEntity := userMap[baseEntity["creatorId"].(string)]
					if userEntity != nil {
						baseEntity["creatorName"] = userEntity.UserName
						baseEntity["creatorAvatar"] = userEntity.Avatar
					}

				}
				if _, ok := baseEntity["modifierId"]; ok && len(baseEntity["modifierId"].(string)) > 0 {
					userEntity := userMap[baseEntity["modifierId"].(string)]
					if userEntity != nil {
						baseEntity["modifierName"] = userEntity.UserName
					}
				}

				assignee, ok := baseEntity["assignee"]
				if ok && assignee != nil && len(assignee.(string)) > 0 {
					userEntity := userMap[assignee.(string)]
					if userEntity != nil {
						baseEntity["assigneeName"] = userEntity.UserName
					}
				}
			}
		}
	}
	return nil
}

// User info
// type User struct {
// 	ID        string `json:"id"`
// 	UserName  string `json:"userName"`
// 	UseStatus int    `json:"useStatus"`
// 	TenantID  string `json:"tenantId"`
// }

// Group info
type Group struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Type     string `json:"type"`
	TenantID string `json:"tenantId"`
}

// UsersResp resp
type UsersResp struct {
	Users []*UserInfoResp `json:"users"`
}

// UserInfoResp info
type UserInfoResp struct {
	ID        string            `json:"id"`
	UserName  string            `json:"name"`
	Avatar    string            `json:"avatar"`
	Email     string            `json:"email"`
	UseStatus int8              `json:"useStatus"` // 1:正常，-2禁用,-1真删除
	Deps      [][]*Dep          `json:"deps"`
	Leaders   [][]*UserInfoResp `json:"leaders"`
}

// Dep Dep
type Dep struct {
	ID             string `json:"id"`
	DepartmentName string `json:"name"`
}
