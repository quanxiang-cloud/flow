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

package eval

import (
	"bytes"
	"fmt"
)

// Format formats an expression as a string.
// It does not attempt to remove unnecessary parens.
func Format(e Expr) string {
	var buf bytes.Buffer
	write(&buf, e)
	return buf.String()
}

func write(buf *bytes.Buffer, e Expr) {
	switch e := e.(type) {
	case literal:
		fmt.Fprintf(buf, "%g", e)

	case Var:
		fmt.Fprintf(buf, "%s", e)

	case unary:
		fmt.Fprintf(buf, "(%c", e.op)
		write(buf, e.x)
		buf.WriteByte(')')

	case binary:
		buf.WriteByte('(')
		write(buf, e.x)
		fmt.Fprintf(buf, " %c ", e.op)
		write(buf, e.y)
		buf.WriteByte(')')

	case call:
		fmt.Fprintf(buf, "%s(", e.fn)
		for i, arg := range e.args {
			if i > 0 {
				buf.WriteString(", ")
			}
			write(buf, arg)
		}
		buf.WriteByte(')')

	default:
		panic(fmt.Sprintf("unknown Expr: %T", e))
	}
}
