package helper

import (
	"bytes"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"reflect"
	"sort"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/hashcode"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/pkg/errors"
)

// Generates a hash for the set hash function used by the IDs
func DataResourceIdsHash(ids []string) string {
	var buf bytes.Buffer

	for _, id := range ids {
		buf.WriteString(fmt.Sprintf("%s-", id))
	}

	return fmt.Sprintf("%d", hashcode.String(buf.String()))
}

// Generates a hash for the set hash function used by the ID
func DataResourceIdHash(id string) string {
	return fmt.Sprintf("%d", hashcode.String(id))
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
