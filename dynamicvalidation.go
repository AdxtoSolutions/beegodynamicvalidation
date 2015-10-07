// Copyright 2014 beego Author. All Rights Reserved.
//
// Copyright 2014 Donal Byrne.
//
// Copyright 2015 David V. Wallin
//
//      Most of the code is as was in beego.validation. The changes
//      I made are mostly the renaming of the Valid function to
//      ValidByStrings and chaning it to take explicit strings
//      instead of a tagged struct in the normal beego case.
//
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package beegodynamicvalidation

import (
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"

	. "github.com/dvwallin/beego/validation"
)

type DynamicValidation struct {
	Validation
}

var (
	// key: function name
	// value: the number of parameters
	funcs = make(Funcs)

	// doesn't belong to validation functions
	unFuncs = map[string]bool{
		"Clear":     true,
		"HasErrors": true,
		"ErrorMap":  true,
		"Error":     true,
		"apply":     true,
		"Check":     true,
		"Valid":     true,
		"NoMatch":   true,
	}
)

func init() {
	v := &Validation{}
	t := reflect.TypeOf(v)
	for i := 0; i < t.NumMethod(); i++ {
		m := t.Method(i)
		if !unFuncs[m.Name] {
			funcs[m.Name] = m.Func
		}
	}
}

// Validate a FieldData object
// the obj parameter must be a struct or a struct pointer
func (v *DynamicValidation) ValidByStrings(fieldName string, validationString string, value interface{}) (b bool, err error) {

	var vfs []ValidFunc
	if vfs, _ = GetValidFuncs(fieldName, validationString); err != nil {
		return
	}

	for _, vf := range vfs {
		//fmt.Printf("%+v", vf.Params...)
		if _, err = funcs.Call(vf.Name,
			mergeParam(&v.Validation, value, vf.Params)...); err != nil {
			return
		}
	}

	return !v.HasErrors(), nil
}

// copied staight out of beego/validation
func GetValidFuncs(fieldName string, validationString string) (vfs []ValidFunc, err error) {
	if len(validationString) == 0 {
		return
	}

	tag := ""
	if vfs, tag, err = getRegFuncs(validationString, fieldName); err != nil {
		fmt.Printf("%+v\n", err)
		return
	}
	fs := strings.Split(tag, ";")
	for _, vfunc := range fs {
		var vf ValidFunc
		if len(vfunc) == 0 {
			continue
		}
		vf, err = parseFunc(vfunc, fieldName)
		if err != nil {
			return
		}
		vfs = append(vfs, vf)
	}
	return
}

// Get Match function
// copied staight out of beego/validation
// May be get NoMatch function in the future
func getRegFuncs(tag, key string) (vfs []ValidFunc, str string, err error) {
	tag = strings.TrimSpace(tag)
	index := strings.Index(tag, "Match(/")
	if index == -1 {
		str = tag
		return
	}
	end := strings.LastIndex(tag, "/)")
	if end < index {
		err = fmt.
			Errorf("invalid Match function")
		return
	}
	reg, err := regexp.Compile(tag[index+len("Match(/") : end])
	if err != nil {
		return
	}
	vfs = []ValidFunc{ValidFunc{"Match", []interface{}{reg, key + ".Match"}}}
	str = strings.TrimSpace(tag[:index]) + strings.TrimSpace(tag[end+len("/)"):])
	return
}

func parseFunc(vfunc, key string) (v ValidFunc, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("%v", r)
		}
	}()

	vfunc = strings.TrimSpace(vfunc)
	start := strings.Index(vfunc, "(")
	var num int

	// doesn't need parameter valid function
	if start == -1 {
		if num, err = numIn(vfunc); err != nil {
			return
		}
		if num != 0 {
			err = fmt.Errorf("%s require %d parameters", vfunc, num)
			return
		}
		v = ValidFunc{vfunc, []interface{}{key + "." + vfunc}}
		//fmt.Printf("\n%v", v)
		return
	}

	//fmt.Printf("\nThis one takes a value %+v", vfunc)

	end := strings.Index(vfunc, ")")
	if end == -1 {
		err = fmt.Errorf("invalid valid function")
		return
	}

	name := strings.TrimSpace(vfunc[:start])
	//fmt.Printf("\nName %s", name)
	if num, err = numIn(name); err != nil {
		return
	}

	params := strings.Split(vfunc[start+1:end], ",")
	// the num of param must be equal
	if num != len(params) {
		err = fmt.Errorf("%s require %d parameters", name, num)
		return
	}

	tParams, err := trim(name, key+"."+name, params)
	if err != nil {
		return
	}
	v = ValidFunc{name, tParams}
	//fmt.Printf("\n%v", v)
	return
}

func numIn(name string) (num int, err error) {
	fn, ok := funcs[name]
	if !ok {
		err = fmt.Errorf("doesn't exsits %s valid function", name)
		return
	}
	// sub *Validation obj and key
	num = fn.Type().NumIn() - 3
	return
}

func trim(name, key string, s []string) (ts []interface{}, err error) {
	ts = make([]interface{}, len(s), len(s)+1)
	fn, ok := funcs[name]
	if !ok {
		err = fmt.Errorf("doesn't exsits %s valid function", name)
		return
	}
	for i := 0; i < len(s); i++ {
		var param interface{}
		// skip *Validation and obj params
		if param, err = magic(fn.Type().In(i+2), strings.TrimSpace(s[i])); err != nil {
			return
		}
		ts[i] = param
	}
	ts = append(ts, key)
	return
}

// modify the parameters's type to adapt the function input parameters' type
func magic(t reflect.Type, s string) (i interface{}, err error) {
	switch t.Kind() {
	case reflect.Int:
		i, err = strconv.Atoi(s)
	case reflect.String:
		i = s
	case reflect.Ptr:
		if t.Elem().String() != "regexp.Regexp" {
			err = fmt.Errorf("does not support %s", t.Elem().String())
			return
		}
		i, err = regexp.Compile(s)
	default:
		err = fmt.Errorf("does not support %s", t.Kind().String())
	}
	return
}

func mergeParam(v *Validation, obj interface{}, params []interface{}) []interface{} {
	return append([]interface{}{v, obj}, params...)
}

func (v *DynamicValidation) Clear() {
	v.Errors = []*ValidationError{}
	v.ErrorsMap = nil
}
