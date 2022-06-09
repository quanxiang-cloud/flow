package restful

import (
	"github.com/gin-gonic/gin"
	"github.com/quanxiang-cloud/flow/internal/flow"
	"github.com/quanxiang-cloud/flow/internal/server/options"
	"github.com/quanxiang-cloud/flow/pkg"
	"github.com/quanxiang-cloud/flow/pkg/config"
	"github.com/quanxiang-cloud/flow/pkg/misc/error2"
	"github.com/quanxiang-cloud/flow/pkg/misc/logger"
	"github.com/quanxiang-cloud/flow/pkg/misc/resp"
	"net/http"
)

// Comment info
type Comment struct {
	comment flow.Comment
}

// NewComment new
func NewComment(c *config.Configs, opts ...options.Options) (*Comment, error) {
	cm, err := flow.NewComment(c, opts...)
	if err != nil {
		return nil, err
	}

	return &Comment{
		comment: cm,
	}, nil
}

func (c *Comment) addComment(ctx *gin.Context) {
	req := &flow.CommentModel{}
	if err := ctx.ShouldBind(req); err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	r, err := c.comment.AddComment(pkg.CTXTransfer(ctx), req)
	if err != nil {
		logger.Logger.Error(err)
		resp.Format(nil, err).Context(ctx)
		return
	}
	resp.Format(r, nil).Context(ctx)
}

func (c *Comment) getComments(ctx *gin.Context) {
	flowInstanceID, ok := ctx.Params.Get("flowInstanceID")
	if !ok {
		error2.NewErrorWithString(error2.Internal, "invalid URI")
		return
	}
	comments, err := c.comment.GetComments(pkg.CTXTransfer(ctx), flowInstanceID)
	resp.Format(comments, err).Context(ctx)

}
