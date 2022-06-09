package eval

import (
	"context"
	"errors"
	"github.com/quanxiang-cloud/flow/internal/eval/gova"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

// 公式最大长度
const exprMaxLength = 5000

// Result response
type Result interface{}

// FormulaReq FormulaReq
type FormulaReq struct {
	Expression string                 `json:"expression"`
	Parameter  map[string]interface{} `json:"parameter"`
}

// Handler Handler
func Handler(c context.Context, req *FormulaReq) (Result, error) {
	expr, err := arrayReplace(req.Expression, req.Parameter)
	if err != nil {
		return nil, err
	}
	expr = symbolReplace(expr)

	r, err := gova.EvalFunc(expr, req.Parameter)
	if err != nil {
		return nil, err
	}
	return r, nil
}

func arrayReplace(expr string, param map[string]interface{}) (string, error) {
	if len(expr) > exprMaxLength {
		return "", errors.New("expr too long")
	}
	re, _ := regexp.Compile(`\{([^}]+)\}`)
	arrs := re.FindAllStringSubmatch(expr, -1)
	arr := delRepeat(&arrs)
	for _, ar := range arr {
		if val, ok := param[ar]; ok {
			switch t := reflect.TypeOf(val); t.Kind() {
			case reflect.Slice, reflect.Array:
				of := reflect.ValueOf(val)
				keys := make([]string, 0, of.Len())
				for i := 0; i < of.Len(); i++ {
					str := ar + strconv.Itoa(i)
					keys = append(keys, str)
					if of.Index(i).CanInterface() {
						param[str] = of.Index(i).Interface()
					}
				}
				expr = strings.ReplaceAll(expr, ar, strings.Join(keys, ","))
			}
		}
	}
	return expr, nil
}

func delRepeat(arr *[][]string) []string {
	res := make([]string, 0)
	for _, ar := range *arr {
		// 取不包含匹配符号的
		str := ar[1]
		strArr := strings.Split(str, ",")
		for _, sar := range strArr {
			if !in(sar, &res) {
				res = append(res, sar)
			}
		}
	}
	return res
}

func symbolReplace(expr string) string {
	for k, v := range op {
		if strings.Contains(expr, k) {
			expr = strings.ReplaceAll(expr, k, v)
		}
	}
	return expr
}
