package helper

import (
	"encoding/base64"
	"strconv"
	"strings"
)

func Bool(i bool) *bool { return &i }

func String(i string) *string { return &i }

func Int(i int) *int { return &i }

func Uint(i uint) *uint { return &i }

func Int64(i int64) *int64 { return &i }

func Float64(i float64) *float64 { return &i }

func Uint64(i uint64) *uint64 { return &i }

func IntInt64(i int) *int64 {
	i64 := int64(i)
	return &i64
}

func IntUint64(i int) *uint64 {
	u := uint64(i)
	return &u
}

func Int64Uint64(i int64) *uint64 {
	u := uint64(i)
	return &u
}

func Strings(strs []string) []*string {
	if len(strs) == 0 {
		return nil
	}

	sp := make([]*string, 0, len(strs))
	for _, s := range strs {
		sp = append(sp, String(s))
	}

	return sp
}

func PString(pointer *string) string {
	if pointer == nil {
		return ""
	}
	return *pointer
}

func PUint64(pointer *uint64) uint64 {
	return *pointer
}

func PInt64(pointer *int64) int64 {
	return *pointer
}

// Takes the result of flatmap.Expand for an array of strings
// and returns a []string
func InterfacesStrings(configured []interface{}) []string {
	vs := make([]string, 0, len(configured))
	for _, v := range configured {
		vs = append(vs, v.(string))
	}
	return vs
}

func InterfacesStringsPoint(configured []interface{}) []*string {
	vs := make([]*string, 0, len(configured))
	for _, v := range configured {
		vs = append(vs, String(v.(string)))
	}
	return vs
}

func StringsStringsPoint(configured []string) []*string {
	vs := make([]*string, 0, len(configured))
	for _, v := range configured {
		vv := v
		vvv := &vv
		vs = append(vs, vvv)
	}
	return vs
}

func InterfacesIntegers(configured []interface{}) []int {
	vs := make([]int, 0, len(configured))
	for _, v := range configured {
		vs = append(vs, v.(int))
	}
	return vs
}

func InterfacesIntInt64Point(configured []interface{}) []*int64 {
	vs := make([]*int64, 0, len(configured))
	for _, v := range configured {
		vs = append(vs, IntInt64(v.(int)))
	}
	return vs
}

func InterfacesUint64Point(configured []interface{}) []*uint64 {
	vs := make([]*uint64, 0, len(configured))
	for _, v := range configured {
		vs = append(vs, Uint64(v.(uint64)))
	}
	return vs
}

// StringsInterfaces Flatten to an array of raw strings and returns a []interface{}
func StringsInterfaces(list []*string) []interface{} {
	vs := make([]interface{}, 0, len(list))
	for _, v := range list {
		vs = append(vs, *v)
	}
	return vs
}

func Uint64sInterfaces(list []*uint64) []interface{} {
	vi := make([]interface{}, 0, len(list))
	for _, v := range list {
		vi = append(vi, int(*v))
	}
	return vi
}

func Int64sInterfaces(list []*int64) []interface{} {
	vi := make([]interface{}, 0, len(list))
	for _, v := range list {
		vi = append(vi, int(*v))
	}
	return vi
}

func BoolToInt64Pointer(s bool) (i *uint64) {
	result := uint64(0)
	if s {
		result = uint64(1)
	}
	i = &result
	return
}

func BoolToInt64Ptr(s bool) (i *int64) {
	result := int64(0)
	if s {
		result = int64(1)
	}
	i = &result
	return
}

func Int64ToStr(s int64) (i string) {
	i = strconv.FormatInt(s, 10)
	return
}

func Int64ToStrPoint(s int64) *string {
	i := Int64ToStr(s)
	return &i
}

func StrToInt64(s string) (i int64) {
	i, _ = strconv.ParseInt(s, 10, 64)
	return
}

func StrToInt64Point(s string) *int64 {
	i := StrToInt64(s)
	return &i
}

func UInt64ToStr(s uint64) (i string) {
	i = strconv.FormatUint(s, 10)
	return
}

func StrToUInt64(s string) (i uint64) {
	intNum, _ := strconv.Atoi(s)
	i = uint64(intNum)
	return
}

func StrToUint64Point(s string) *uint64 {
	i := StrToUInt64(s)
	return &i
}

func StrToBool(s string) (i bool) {
	i = false
	if s == "true" {
		i = true
	}
	return
}

func StrListToStr(strList []*string) string {
	res := ""
	for i, v := range strList {
		res += *v
		if i < len(strList)-1 {
			res += ";"
		}
	}
	return base64.StdEncoding.EncodeToString([]byte(res))
}

func StrToStrList(str string) (res []string, err error) {

	decodeString, err := base64.StdEncoding.DecodeString(str)
	if err != nil {
		return res, err
	}

	res = strings.Split(string(decodeString), ";")
	return res, nil
}
