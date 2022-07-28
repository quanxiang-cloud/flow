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
	"context"
	"fmt"
	"math"
	"testing"
)

func TestEval(t *testing.T) {
	tests := []struct {
		expr string
		env  Env
		want interface{}
	}{
		{"lower(A)", Env{"A": &String{"ASS"}}, "ass"},
		{"sum(sum(a,b),c)", Env{"a": &Float64{1}, "b": &Float64{2}, "c": &Float64{3}}, 6},
		{"sum(sum(1,2),c)", Env{"c": &Float64{3}}, 6},
		{"average(a,b,c)", Env{"a": &Float64{1}, "b": &Float64{2}, "c": &Float64{3}}, 2},
		{"average(sum(a,b),c)", Env{"a": &Float64{1}, "b": &Float64{2}, "c": &Float64{3}}, 3},
		{"max(a,b,c,d)", Env{"a": &Float64{6}, "b": &Float64{2}, "c": &Float64{3}, "d": &Float64{5}}, 6},
		{"min(a,b,c,d)", Env{"a": &Float64{6}, "b": &Float64{2}, "c": &Float64{3}, "d": &Float64{5}}, 2},
		{"count(a,b,c,d)", Env{"a": &Float64{6}, "b": &Float64{2}, "c": &Float64{3}, "d": &Float64{5}}, 4},
		{"abs(a)", Env{"a": &Float64{-6}}, 6},
		{"round(a/b,c)", Env{"a": &Float64{19}, "b": &Float64{3}, "c": &Float64{3}}, 6.333},
		{"ceil(a/b)", Env{"a": &Float64{19}, "b": &Float64{3}}, 7},
		{"floor(a/b)", Env{"a": &Float64{19}, "b": &Float64{3}}, 6},
		{"mod(a,b)", Env{"a": &Float64{19}, "b": &Float64{3}}, 1},
		{"sqrt(A / pi)", Env{"A": &Float64{87616}, "pi": &Float64{math.Pi}}, 167},
		{"pow(x, 3) + pow(y, 3)", Env{"x": &Float64{12}, "y": &Float64{1}}, 1729},
		{"pow(x, 3) + pow(y, 3)", Env{"x": &Float64{9}, "y": &Float64{10}}, 1729},
		{"5 / 9 * (F - 32)", Env{"F": &Float64{-40}}, -40},
		{"5 / 9 * (F - 32)", Env{"F": &Float64{32}}, "0"},
		{"5 / 9 * (F - 32)", Env{"F": &Float64{212}}, "100"},
		{"-1 + -x", Env{"x": &Float64{1}}, "-2"},
		{"-1 - x", Env{"x": &Float64{1}}, "-2"},
	}
	var prevExpr string
	for _, test := range tests {
		if test.expr != prevExpr {
			prevExpr = test.expr
		}
		expr, err := Parse(test.expr)
		if err != nil {
			t.Error(err)
			continue
		}
		var got, want string
		switch test.want.(type) {
		case string:
			got = expr.Eval(test.env).String()
			want = test.want.(string)
		case int, int8, int32, int64, uint, uint8, uint32, uint64, float32, float64:
			got = fmt.Sprintf("%.6g", expr.Eval(test.env).Float64())
			want = fmt.Sprintf("%v", test.want)
		}

		if got != want {
			fmt.Printf("\n%s  ", test.expr)
			t.Errorf("%s.Eval() in %v = %q, want %q\n",
				test.expr, test.env, got, want)
		}
	}
}

func TestHandler(t *testing.T) {
	p := map[string]interface{}{
		"a":            35,
		"b":            5,
		"c":            3,
		"name":         "mark",
		"mark":         "mark",
		"$.one":        "1",
		"fb_one:tb_56": "1",
		"two":          2,
		"one":          2,
		"three":        3,
		"four":         4,
		"five":         5,
	}
	// r,err := Handler("sum(sum(a,b),c)",p)
	// expr := "five / [four] * (a - 32)"
	// expr := "1 % 1 * 2 + 3"
	// expr := "1 % 1 * 2 / 3"
	// expr := "1 - 1 * 2 / 3"
	// expr := "1 > 2 / 3"
	// expr := "1 == 1"
	// expr := "abs(-20 + 30) "
	// expr := "sum(sum(a,b),c)"
	// expr := "2 == 1 && 3 > 1"
	// expr := "2 == 1 || 3 > 1"
	// expr := "(fb_one:tb_56 == two || one < 22) && (3 == 4 || 3 <= 4) || (five == 6 && 5 < 6)"
	// expr := "(1 == 2 && 1 < 2) || (1 == 2 && 1 < 2) || (1 == 2 || 1 < 2)"
	// expr := "name == \"mark\""
	// expr := "round([a] / [c],4)"
	// expr := "max(b,20,c,10)"
	expr := "min(b,20,c,10)"
	// expr := "1 ∉ {one, 2 , 3, 4}"
	// expr := "1 ∈ {[$.one], 2 , 3, 4}"
	req := &FormulaReq{
		Expression: expr,
		Parameter:  p,
	}
	r, err := Handler(context.TODO(), req)
	fmt.Println(r, err)
}
