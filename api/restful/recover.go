package restful

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/quanxiang-cloud/flow/pkg/misc/error2"
	"github.com/quanxiang-cloud/flow/pkg/misc/resp"
	"log"
	"runtime/debug"
)

// Recover gin error recover
func Recover(c *gin.Context) {
	defer func() {
		if r := recover(); r != nil {
			if gin.Mode() == DebugMode {
				log.Printf("panic: %v\n", r)
				debug.PrintStack()
			}

			(&resp.R{
				Code: error2.Unknown,
				Msg:  fmt.Sprintf("error msg: %v", r),
			}).Context(c)

			// abort next execution
			c.Abort()
		}
	}()
	// execute next after load defer recover
	c.Next()
}
