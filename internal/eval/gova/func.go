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

import (
	"errors"
	"fmt"
	"github.com/shopspring/decimal"
	"math"
	"regexp"
	"sort"
	"strings"
)

const (
	paramsCntErr  = "parameter counts error for '%v'"
	paramsTypeErr = "parameter type error for '%v'"
)

func sum(arguments ...interface{}) (interface{}, error) {
	var result float64
	for _, arg := range arguments {
		if !isFloat64(arg) {
			errorMsg := fmt.Sprintf(paramsTypeErr, "sum")
			return nil, errors.New(errorMsg)
		}
		result += arg.(float64)
	}
	return result, nil
}

func pow(arguments ...interface{}) (interface{}, error) {
	if len(arguments) < 2 {
		errorMsg := fmt.Sprintf(paramsCntErr, "pow")
		return nil, errors.New(errorMsg)
	}
	if !isFloat64(arguments[0]) || !isFloat64(arguments[1]) {
		errorMsg := fmt.Sprintf(paramsTypeErr, "pow")
		return nil, errors.New(errorMsg)
	}
	return math.Pow(arguments[0].(float64), arguments[1].(float64)), nil
}

func sin(arguments ...interface{}) (interface{}, error) {
	if !isFloat64(arguments[0]) {
		errorMsg := fmt.Sprintf(paramsTypeErr, "sin")
		return nil, errors.New(errorMsg)
	}
	return math.Sin(arguments[0].(float64)), nil
}

func sqrt(arguments ...interface{}) (interface{}, error) {
	if !isFloat64(arguments[0]) {
		errorMsg := fmt.Sprintf(paramsTypeErr, "sqrt")
		return nil, errors.New(errorMsg)
	}
	return math.Sqrt(arguments[0].(float64)), nil
}

func lower(arguments ...interface{}) (interface{}, error) {
	if !isString(arguments[0]) {
		errorMsg := fmt.Sprintf(paramsTypeErr, "lower")
		return nil, errors.New(errorMsg)
	}
	return strings.ToLower(arguments[0].(string)), nil
}

func average(arguments ...interface{}) (interface{}, error) {
	var ans float64
	for _, arg := range arguments {
		if !isFloat64(arg) {
			errorMsg := fmt.Sprintf(paramsTypeErr, "average")
			return nil, errors.New(errorMsg)
		}
		ans += arg.(float64)
	}
	if len(arguments) > 0 {
		total := float64(len(arguments))
		ans = ans / total
	}
	return ans, nil
}

func max(arguments ...interface{}) (interface{}, error) {
	var ans []float64
	for _, arg := range arguments {
		if !isFloat64(arg) {
			errorMsg := fmt.Sprintf(paramsTypeErr, "max")
			return nil, errors.New(errorMsg)
		}
		ans = append(ans, arg.(float64))
	}
	sort.Sort(sort.Float64Slice(ans))
	if len(ans) > 0 {
		return ans[len(ans)-1], nil
	}
	return 0, nil
}

func min(arguments ...interface{}) (interface{}, error) {
	var ans []float64
	for _, arg := range arguments {
		if !isFloat64(arg) {
			errorMsg := fmt.Sprintf(paramsTypeErr, "min")
			return nil, errors.New(errorMsg)
		}
		ans = append(ans, arg.(float64))
	}
	sort.Sort(sort.Float64Slice(ans))
	if len(ans) > 0 {
		return ans[0], nil
	}
	return 0, nil
}

func count(arguments ...interface{}) (interface{}, error) {
	return len(arguments), nil
}

func abs(arguments ...interface{}) (interface{}, error) {
	if !isFloat64(arguments[0]) {
		errorMsg := fmt.Sprintf(paramsTypeErr, "abs")
		return nil, errors.New(errorMsg)
	}
	return math.Abs(arguments[0].(float64)), nil
}

// 四舍五入
func round(arguments ...interface{}) (interface{}, error) {
	var rnd int32
	if len(arguments) < 2 {
		arguments = append(arguments, 2)
	}
	if !isFloat64(arguments[0]) || !isFloat64(arguments[1]) {
		errorMsg := fmt.Sprintf(paramsTypeErr, "round")
		return nil, errors.New(errorMsg)
	}
	rnd = int32(arguments[1].(float64))
	v4, _ := decimal.NewFromFloat(arguments[0].(float64)).Round(rnd).Float64()
	return v4, nil
}

func ceil(arguments ...interface{}) (interface{}, error) {
	var ans float64
	if !isFloat64(arguments[0]) {
		errorMsg := fmt.Sprintf(paramsTypeErr, "ceil")
		return nil, errors.New(errorMsg)
	}
	ans = math.Ceil(arguments[0].(float64))
	return ans, nil
}

func floor(arguments ...interface{}) (interface{}, error) {
	var ans float64
	if !isFloat64(arguments[0]) {
		errorMsg := fmt.Sprintf(paramsTypeErr, "floor")
		return nil, errors.New(errorMsg)
	}
	ans = math.Floor(arguments[0].(float64))
	return ans, nil
}

func mod(arguments ...interface{}) (interface{}, error) {
	var ans float64
	if len(arguments) < 2 {
		errorMsg := fmt.Sprintf(paramsCntErr, "mod")
		return nil, errors.New(errorMsg)
	}
	if !isFloat64(arguments[0]) || !isFloat64(arguments[1]) {
		errorMsg := fmt.Sprintf(paramsTypeErr, "mod")
		return nil, errors.New(errorMsg)
	}
	ans = math.Mod(arguments[0].(float64), arguments[1].(float64))
	return ans, nil
}

func isString(value interface{}) bool {

	switch value.(type) {
	case string:
		return true
	}
	return false
}

func isRegexOrString(value interface{}) bool {

	switch value.(type) {
	case string:
		return true
	case *regexp.Regexp:
		return true
	}
	return false
}

func isBool(value interface{}) bool {
	switch value.(type) {
	case bool:
		return true
	}
	return false
}

func isFloat64(value interface{}) bool {
	switch value.(type) {
	case float64:
		return true
	}
	return false
}
