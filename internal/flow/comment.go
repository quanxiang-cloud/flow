package flow

import (
	"context"
	"git.internal.yunify.com/qxp/misc/error2"
	"git.internal.yunify.com/qxp/misc/id2"
	"git.internal.yunify.com/qxp/misc/time2"
	"github.com/quanxiang-cloud/flow/internal/models"
	"github.com/quanxiang-cloud/flow/internal/models/mysql"
	"github.com/quanxiang-cloud/flow/internal/server/options"
	"github.com/quanxiang-cloud/flow/pkg"
	"github.com/quanxiang-cloud/flow/pkg/client"
	"github.com/quanxiang-cloud/flow/pkg/code"
	"github.com/quanxiang-cloud/flow/pkg/config"
	"github.com/quanxiang-cloud/flow/pkg/utils"
	"gorm.io/gorm"
)

// Comment inter
type Comment interface {
	AddComment(ctx context.Context, req *CommentModel) (bool, error)
	GetComments(ctx context.Context, flowInstanceID string) ([]*CommentModel, error)
}

type comment struct {
	db                    *gorm.DB
	commentRepo           models.CommentRepo
	flowRepo              models.FlowRepo
	instanceRepo          models.InstanceRepo
	processAPI            client.Process
	commentAttachmentRepo models.CommentAttachmentRepo
	identityAPI           client.Identity
}

// NewComment init
func NewComment(conf *config.Configs, opts ...options.Options) (Comment, error) {
	var c = &comment{
		commentRepo:           mysql.NewCommentRepo(),
		flowRepo:              mysql.NewFlowRepo(),
		instanceRepo:          mysql.NewInstanceRepo(),
		processAPI:            client.NewProcess(conf),
		commentAttachmentRepo: mysql.NewCommentAttachmentRepo(),
		identityAPI:           client.NewIdentity(conf),
	}
	for _, opt := range opts {
		opt(c)
	}
	return c, nil
}

// SetDB set db
func (c *comment) SetDB(db *gorm.DB) {
	c.db = db
}

// CommentModel comment info
type CommentModel struct {
	*models.Comment
	Attachments []*models.CommentAttachment `json:"attachments"`
	*models.BaseModel
}

func (c *comment) AddComment(ctx context.Context, req *CommentModel) (bool, error) {
	err := c.verificationAuthority(ctx, req.FlowInstanceID)
	if err != nil {
		return false, err
	}

	id := id2.GenID()
	err = c.commentRepo.Create(c.db, &models.Comment{
		FlowInstanceID: req.FlowInstanceID,
		CommentUserID:  pkg.STDUserID(ctx),
		Content:        utils.UnicodeEmojiCode(req.Content),
		BaseModel: models.BaseModel{
			ID:         id,
			CreatorID:  pkg.STDUserID(ctx),
			CreateTime: time2.Now(),
			ModifierID: pkg.STDUserID(ctx),
			ModifyTime: time2.Now(),
		},
	})
	if err != nil {
		return false, err
	}

	if req.Attachments != nil {
		for _, e := range req.Attachments {
			e.FlowCommentID = id
			e.CreatorID = pkg.STDUserID(ctx)
			e.CreateTime = time2.Now()
			e.ID = id2.GenID()
			err := c.commentAttachmentRepo.Create(c.db, e)
			if err != nil {
				return false, err
			}
		}
	}
	return true, nil
}

func (c *comment) GetComments(ctx context.Context, flowInstanceID string) ([]*CommentModel, error) {
	err := c.verificationAuthority(ctx, flowInstanceID)
	if err != nil {
		return nil, err
	}

	comments, err := c.commentRepo.FindComments(c.db, map[string]interface{}{
		"flow_instance_id": flowInstanceID,
	}, "create_time DESC")
	if err != nil {
		return nil, err
	}
	if comments == nil {
		return nil, nil
	}

	userIDs := make([]string, 0)
	for _, e := range comments {
		userIDs = append(userIDs, e.CommentUserID)
	}
	users, err := c.identityAPI.FindUsersByIDs(ctx, userIDs)
	if err != nil {
		return nil, err
	}

	resp := make([]*CommentModel, 0)
	for _, e := range comments {
		attachments, err := c.commentAttachmentRepo.FindAttachments(c.db, map[string]interface{}{
			"flow_comment_id": e.ID,
		}, "create_time desc")
		if err != nil {
			return nil, err
		}
		user := users[e.CommentUserID]
		if user == nil {
			user = &client.UserInfoResp{}
		}
		e.Content = utils.UnicodeEmojiDecode(e.Content)
		comment := &CommentModel{
			Comment:     e,
			Attachments: attachments,
			BaseModel: &models.BaseModel{
				CreatorName:   user.UserName,
				CreatorAvatar: user.Avatar,
				CreateTime:    e.CreateTime,
			},
		}
		resp = append(resp, comment)
	}
	return resp, nil
}

func (c *comment) verificationAuthority(ctx context.Context, flowInstanceID string) error {
	if flowInstanceID == "" {
		return error2.NewErrorWithString(code.InvalidProcessID, "flowInstanceId is nil")
	}
	instances, err := c.instanceRepo.FindByProcessInstanceIDs(c.db, []string{flowInstanceID})
	if err != nil {
		return err
	}
	if instances == nil || len(instances) == 0 {
		return error2.NewErrorWithString(code.InvalidProcessID, " process model is nil")
	}

	flow, err := c.flowRepo.FindByID(c.db, instances[0].FlowID)
	if err != nil {
		return err
	}
	if flow == nil {
		return error2.NewErrorWithString(code.InvalidProcessID, " process model is nil")
	}
	if flow.CanMsg != 1 {
		return error2.NewErrorWithString(code.NoMessagesAllowed, "No messages allowed ")
	}

	ps, err := c.processAPI.GetInstances(ctx, client.GetTasksReq{
		InstanceID: []string{flowInstanceID},
		Assignee:   pkg.STDUserID(ctx),
	})
	if err != nil {
		return err
	}
	if ps == nil {
		return error2.NewErrorWithString(code.NoMessagesAllowed, "No messages allowed ")
	}
	return nil
}
