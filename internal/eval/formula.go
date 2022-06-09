package eval

import (
	"errors"
	"github.com/shopspring/decimal"
	"math"
	"sort"
	"strings"
)

// Formula Formula define
type Formula func(args []Expr, env Env) Substance

var (
	// ErrNoFormula ErrNoFormula
	ErrNoFormula = errors.New("no formula like this")

	errCodeMissParam  = "-90054020001"
	errCodeWithoutFun = "-90054020002"

	// formulas Formula map
	formulas = map[string]Formula{
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

	op = map[string]string{
		"\"": "'",
		"{":  "(",
		"}":  ")",
		"∈":  "in",
		"∉":  "nin",
	}
)

// 公式实现
func pow(args []Expr, env Env) Substance {
	return &Float64{math.Pow(args[0].Eval(env).Float64(), args[1].Eval(env).Float64())}
}
func sin(args []Expr, env Env) Substance {
	return &Float64{math.Sin(args[0].Eval(env).Float64())}
}
func sqrt(args []Expr, env Env) Substance {
	return &Float64{math.Sqrt(args[0].Eval(env).Float64())}
}
func sum(args []Expr, env Env) Substance {
	var ans float64
	for _, arg := range args {
		ans += arg.Eval(env).Float64()
	}
	return &Float64{ans}
}
func lower(args []Expr, env Env) Substance {
	return &String{strings.ToLower(args[0].Eval(env).String())}
}

func average(args []Expr, env Env) Substance {
	var ans float64
	for _, arg := range args {
		ans += arg.Eval(env).Float64()
	}
	if len(args) > 0 {
		total := float64(len(args))
		ans = ans / total
	}
	return &Float64{ans}
}

func max(args []Expr, env Env) Substance {
	var ans []float64
	for _, arg := range args {
		ans = append(ans, arg.Eval(env).Float64())
	}
	sort.Sort(sort.Float64Slice(ans))
	if len(ans) > 0 {
		return &Float64{ans[len(ans)-1]}
	}
	return &Float64{}
}

func min(args []Expr, env Env) Substance {
	var ans []float64
	for _, arg := range args {
		ans = append(ans, arg.Eval(env).Float64())
	}
	sort.Sort(sort.Float64Slice(ans))
	if len(ans) > 0 {
		return &Float64{ans[0]}
	}
	return &Float64{}
}

func count(args []Expr, env Env) Substance {
	return &Float64{float64(len(args))}
}

func abs(args []Expr, env Env) Substance {
	return &Float64{math.Abs(args[0].Eval(env).Float64())}
}

// 四舍五入
func round(args []Expr, env Env) Substance {
	rnd := int32(args[1].Eval(env).Float64())
	v4, _ := decimal.NewFromFloat(args[0].Eval(env).Float64()).Round(rnd).Float64()
	return &Float64{v4}
}

func ceil(args []Expr, env Env) Substance {
	var ans float64
	ans = math.Ceil(args[0].Eval(env).Float64())
	return &Float64{ans}
}

func floor(args []Expr, env Env) Substance {
	var ans float64
	ans = math.Floor(args[0].Eval(env).Float64())
	return &Float64{ans}
}

func mod(args []Expr, env Env) Substance {
	var ans float64
	ans = math.Mod(args[0].Eval(env).Float64(), args[1].Eval(env).Float64())
	return &Float64{ans}
}
