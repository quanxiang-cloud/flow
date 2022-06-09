package resp

import (
	"context"
	"encoding/json"
	"io"
	"net/http"

	"github.com/quanxiang-cloud/flow/pkg/misc/error2"
	"github.com/quanxiang-cloud/flow/pkg/misc/logger"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type resp interface{}

// R 返回值
type R struct {
	err  error
	Code int64  `json:"code"`
	Msg  string `json:"msg,omitempty"`
	Data resp   `json:"data"`
}

// Context context
func (r *R) Context(c *gin.Context, code ...int) {
	status := http.StatusOK
	if r.Code == error2.Unknown {
		status = http.StatusInternalServerError
		if r.err != nil {
			logger.Logger.Errorw(r.err.Error(), logger.GINRequestID(c))
		} else {
			logger.Logger.Errorw("unknown err", logger.GINRequestID(c))
		}
	} else if len(code) != 0 {
		status = code[0]
	}

	c.JSON(status, r)
}

// Format 统一返回值格式
func Format(resp resp, err error) *R {
	if err == nil {
		return &R{
			Code: error2.Success,
			Data: resp,
		}
	}

	var fail = func(err error) *R {
		return &R{
			Code: error2.Unknown,
			err:  err,
		}
	}

	var code int64
	switch e := err.(type) {
	case error2.Error:
		code = e.Code
	case *error2.Error:
		if e == nil {
			return fail(nil)
		}
		code = e.Code
	case validator.ValidationErrors:
		code = error2.ErrParams
	default:
		return fail(e)
	}

	return &R{
		Code: code,
		Msg:  err.Error(),
	}
}

// DeserializationResp DeserializationResp
func DeserializationResp(ctx context.Context, response *http.Response, entity interface{}) error {
	if response.StatusCode != http.StatusOK {
		return error2.NewError(error2.Internal)
	}

	r := new(R)
	r.Data = entity

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return err
	}

	err = json.Unmarshal(body, r)
	if err != nil {
		return err
	}

	if r.Code != error2.Success {
		return error2.NewErrorWithString(r.Code, r.Msg)
	}

	return nil
}
