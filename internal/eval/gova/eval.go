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

package gova

import "github.com/xpsl/govaluate"

// 自定义函数
var functions = map[string]govaluate.ExpressionFunction{
	"pow":     pow,
	"sin":     sin,
	"sqrt":    sqrt,
	"sum":     sum,
	"lower":   lower,
	"average": average,
	"max":     max,
	"min":     min,
	"count":   count,
	"round":   round,
	"abs":     abs,
	"ceil":    ceil,
	"floor":   floor,
	"mod":     mod,
}

type result struct {
	Res interface{} `json:"result"`
}

// EvalFunc EvalFunc
func EvalFunc(exprString string, parameters map[string]interface{}) (interface{}, error) {
	// in exprString,parameter use [xxx] style
	expr, err := govaluate.NewEvaluableExpressionWithFunctions(exprString, functions)
	if err != nil {
		return nil, err
	}
	res, err := expr.Evaluate(parameters)
	if err != nil {
		return nil, err
	}
	return res, nil
}
