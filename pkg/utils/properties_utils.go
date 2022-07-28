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

package utils

import (
	"errors"
	"fmt"
	"reflect"
)

// CopyProperties copy properties to dest
func CopyProperties(dst, src interface{}) (err error) {
	// Prevention of accidents panic
	defer func() {
		if e := recover(); e != nil {
			err = fmt.Errorf("%v", e)
		}
	}()

	dstType, dstValue := reflect.TypeOf(dst), reflect.ValueOf(dst)
	srcType, srcValue := reflect.TypeOf(src), reflect.ValueOf(src)

	// dst must be struct pointer
	if dstType.Kind() != reflect.Ptr || dstType.Elem().Kind() != reflect.Struct {
		return errors.New("dst type should be a struct pointer")
	}

	// src must be struct or struct pointer
	if srcType.Kind() == reflect.Ptr {
		srcType, srcValue = srcType.Elem(), srcValue.Elem()
	}
	if srcType.Kind() != reflect.Struct {
		return errors.New("src type should be a struct or a struct pointer")
	}

	dstType, dstValue = dstType.Elem(), dstValue.Elem()

	propertyNums := dstType.NumField()

	for i := 0; i < propertyNums; i++ {
		property := dstType.Field(i)
		propertyValue := srcValue.FieldByName(property.Name)

		if !propertyValue.IsValid() || property.Type != propertyValue.Type() {
			continue
		}

		if dstValue.Field(i).CanSet() {
			dstValue.Field(i).Set(propertyValue)
		}
	}

	return nil
}

// IsNil is nil
func IsNil(i interface{}) bool {
	defer func() {
		recover()
	}()
	if i == nil {
		return true
	}
	vi := reflect.ValueOf(i)

	switch reflect.TypeOf(i).Kind() {
	case reflect.String:
		if vi.String() == "" {
			return true
		}
	}

	// case reflect.Slice, reflect.Array:
	// 	entity, err := json.Marshal(v.Interface())
	// 	if err != nil {
	// 		return err
	// 	}
	// 	variables.ComplexValue = entity
	// 	variables.VarType = "[]string"
	// case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
	// 	variables.Value = fmt.Sprintf("%d", v.Int())
	// 	variables.VarType = "int"
	// case reflect.Float32, reflect.Float64:
	// 	variables.Value = fmt.Sprintf("%f", v.Float())
	// 	variables.VarType = "float"
	// }
	return vi.IsNil()
}

// IsNotNil is not nil
func IsNotNil(i interface{}) bool {
	return !IsNil(i)
}
