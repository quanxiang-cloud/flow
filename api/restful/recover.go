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
