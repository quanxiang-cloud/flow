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
	return &result{
		Res: res,
	}, nil
}
