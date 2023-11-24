package helper

import (
	"bytes"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"hash/crc32"
	"reflect"
	"sort"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/pkg/errors"
)

// Generates a hash for the set hash function used by the IDs
func DataResourceIdsHash(ids []string) string {
	var buf bytes.Buffer

	for _, id := range ids {
		buf.WriteString(fmt.Sprintf("%s-", id))
	}

	return fmt.Sprintf("%d", HashString(buf.String()))
}

// Generates a hash for the resource
func ResourceIdsHash(ids []string) string {
	return DataResourceIdsHash(ids)
}

// HashString hashes a string to a unique hashcode.
//
// This will be removed in v2 without replacement. So we place here instead of import.
func HashString(s string) int {
	v := int(crc32.ChecksumIEEE([]byte(s)))
	if v >= 0 {
		return v
	}
	if -v >= 0 {
		return -v
	}
	// v == MinInt
	return 0
}

// Strings hashes a list of strings to a unique hashcode.
func HashStrings(strings []string) string {
	var buf bytes.Buffer

	for _, s := range strings {
		buf.WriteString(fmt.Sprintf("%s-", s))
	}

	return fmt.Sprintf("%d", String(buf.String()))
}

// Generates a hash for the set hash function used by the ID
func DataResourceIdHash(id string) string {
	return fmt.Sprintf("%d", HashString(id))
}

func GetTags(d *schema.ResourceData, k string) map[string]string {
	tags := make(map[string]string)
	if raw, ok := d.GetOk(k); ok {
		for k, v := range raw.(map[string]interface{}) {
			tags[k] = v.(string)
		}
	}
	return tags
}

func BuildToken() string {
	buf := make([]byte, 16)
	_, _ = rand.Read(buf)
	return base64.StdEncoding.EncodeToString(buf)
}

func FormatUnixTime(n uint64) string {
	return time.Unix(int64(n), 0).UTC().Format("2006-01-02T03:04:05Z")
}
func ParseTime(s string) (time.Time, error) {
	return time.ParseInLocation("2006-01-02T03:04:05Z", s, time.UTC)
}

// compose all schema.SchemaValidateFunc to a schema.SchemaValidateFunc,
// like resource.ComposeTestCheckFunc, so that we can reuse exist schema.SchemaValidateFunc
// and reduce custom schema.SchemaValidateFunc codes size.
func ComposeValidateFunc(fns ...schema.SchemaValidateFunc) schema.SchemaValidateFunc {
	return func(v interface{}, k string) (wssRet []string, errsRet []error) {
		for _, fn := range fns {
			wss, errs := fn(v, k)
			if len(errs) > 0 {
				errsRet = append(errsRet, errs...)
				return
			}
			wssRet = append(wssRet, wss...)
		}

		return
	}
}

// CheckIfSetTogether will check all args, they should be all nil or not nil.
//
// Such as vpc_id and subnet_id should set together, or don't set them.
func CheckIfSetTogether(d *schema.ResourceData, args ...string) error {
	var notNil bool

	for _, arg := range args {
		if _, ok := d.GetOk(arg); ok {
			notNil = true
			continue
		}
		if notNil {
			return errors.Errorf("%v must be set together", args)
		}
	}

	return nil
}

func StringsContain(ss []string, str string) bool {
	for _, s := range ss {
		if str == s {
			return true
		}
	}

	return false
}

func DiffSupressJSON(k, olds, news string, d *schema.ResourceData) bool {
	var oldJson interface{}
	err := json.Unmarshal([]byte(olds), &oldJson)
	if err != nil {
		return olds == news
	}
	var newJson interface{}
	err = json.Unmarshal([]byte(news), &newJson)
	if err != nil {
		return olds == news
	}
	flag := reflect.DeepEqual(oldJson, newJson)
	return flag
}

/*
    Serialize slice into the usage document
	eg["status_change","abnormal"] will be "`abnormal`,`status_change`"
*/
func SliceFieldSerialize(slice []string) string {
	types := []string{}
	for _, v := range slice {
		types = append(types, "`"+v+"`")
	}
	sort.Strings(types)
	return strings.Trim(strings.Join(types, ","), ",")
}

// InterfacesHeadMap returns string key map if argument is MaxItem: 1 List Type
func InterfacesHeadMap(d *schema.ResourceData, key string) (result map[string]interface{}, ok bool) {
	v, ok := d.GetOk(key)
	if !ok {
		return
	}
	interfaces, ok := v.([]interface{})
	if !ok || len(interfaces) == 0 {
		ok = false
		return
	}
	head := interfaces[0]
	result, ok = head.(map[string]interface{})
	return
}

// ConvertInterfacesHeadToMap returns string key map if argument is MaxItem: 1 List Type
func ConvertInterfacesHeadToMap(v interface{}) (result map[string]interface{}, ok bool) {
	interfaces, ok := v.([]interface{})
	if !ok || len(interfaces) == 0 {
		ok = false
		return
	}
	head := interfaces[0]
	result, ok = head.(map[string]interface{})
	return
}

// CovertInterfaceMapToStrPtr returns [string:string] map from a [string:interface] map
func CovertInterfaceMapToStrPtr(m map[string]interface{}) map[string]*string {
	result := make(map[string]*string)
	for k, v := range m {
		if s, ok := v.(string); ok {
			result[k] = &s
		}
	}
	return result
}

func SetMapInterfaces(d *schema.ResourceData, key string, values ...map[string]interface{}) error {
	val := make([]interface{}, 0, len(values))
	for i := range values {
		item := values[i]
		val = append(val, item)
	}
	return d.Set(key, val)
}

func InterfaceToMap(d map[string]interface{}, key string) (result map[string]interface{}, ok bool) {
	if v, ok := d[key]; ok {
		if len(v.([]interface{})) != 1 || v.([]interface{})[0] == nil {
			return nil, false
		}
		return v.([]interface{})[0].(map[string]interface{}), true
	}
	return nil, false
}

func ImportWithDefaultValue(defaultValues map[string]interface{}) schema.StateFunc {
	return func(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
		for k, v := range defaultValues {
			_ = d.Set(k, v)
		}
		return []*schema.ResourceData{d}, nil
	}
}

func ImmutableArgsChek(d *schema.ResourceData, arguments ...string) error {
	for _, v := range arguments {
		if d.HasChange(v) {
			o, _ := d.GetChange(v)
			_ = d.Set(v, o)
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}
	return nil
}

func IsEmptyStr(s *string) bool {
	if s == nil {
		return true
	}
	return *s == ""
}

func MapToString(param map[string]interface{}) (string, bool) {
	data, err := json.Marshal(param)
	if err != nil {
		return "", false
	}

	return string(data), true
}

func JsonToMap(str string) (map[string]interface{}, error) {
	var temp map[string]interface{}

	err := json.Unmarshal([]byte(str), &temp)
	if err != nil {
		return nil, err
	}

	return temp, nil
}

func CheckElementsExist(slice1 []string, slice2 []string) (bool, []string) {
	exist := true
	diff := make([]string, 0)

	slice1Map := make(map[string]bool)
	slice2Map := make(map[string]bool)
	for _, element := range slice1 {
		slice1Map[element] = true
	}
	for _, element := range slice2 {
		slice2Map[element] = true
	}

	for _, element := range slice1 {
		if _, ok := slice2Map[element]; !ok {
			exist = false
			break
		}
	}
	if exist {
		for _, element := range slice2 {
			if _, ok := slice1Map[element]; !ok {
				diff = append(diff, element)
			}
		}
	}
	return exist, diff
}
