package helper

import (
	"bytes"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"hash/crc32"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

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

func InterfaceToMap(d map[string]interface{}, key string) (result map[string]interface{}, ok bool) {
	if v, ok := d[key]; ok {
		if len(v.([]interface{})) != 1 {
			return nil, false
		}
		return v.([]interface{})[0].(map[string]interface{}), true
	}
	return nil, false
}

// String hashes a string to a unique hashcode.
//
// Deprecated: This will be removed in v2 without replacement. If you need
// its functionality, you can copy it, import crc32 directly, or reference the
// v1 package.
//
// crc32 returns a uint32, but for our use we need
// and non negative integer. Here we cast to an integer
// and invert it if the result is negative.
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
//
// Deprecated: This will be removed in v2 without replacement. If you need
// its functionality, you can copy it, import crc32 directly, or reference the
// v1 package.
func HashStrings(strings []string) string {
	var buf bytes.Buffer

	for _, s := range strings {
		buf.WriteString(fmt.Sprintf("%s-", s))
	}

	return fmt.Sprintf("%d", String(buf.String()))
}

// CheckSchemaSetResourceAttr can be used for checking attributes which type is schema.TypeSet
// @params
// name:  resource name e.g. `tencentcloud_cos_bucket.foo`
// path:  path to TypeSet arguments e.g. `lifecycle_rules.0.expiration`
// key:   elem key without index, because TypeSet is unordered. e.g. `days`, this can set to empty "" to test as primitive elems
// value: expect value includes in argument set
func CheckSchemaSetResourceAttr(name, path, key, value string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		ms := s.RootModule()
		rs, ok := ms.Resources[name]
		if !ok {
			return fmt.Errorf("Not found: %s in %s", name, ms.Path)
		}

		is := rs.Primary
		if is == nil {
			return fmt.Errorf("No primary instance: %s in %s", name, ms.Path)
		}
		mapSize, ok := is.Attributes[fmt.Sprintf("%s.#", path)]
		length, err := strconv.Atoi(mapSize)
		if !ok || err != nil {
			return fmt.Errorf("cannot read atribute %s.%s.\\# , got %s", name, path, mapSize)
		}
		if length == 0 {
			return fmt.Errorf("%s.%s has no elements", name, path)
		}
		values := make([]string, 0)
		hit := false
		for i := 0; i < length; i++ {
			fullKey := fmt.Sprintf("%s.%d.%s", path, i, key)
			if key == "" {
				fullKey = fmt.Sprintf("%s.%d", path, i)
			}
			val, ok := is.Attributes[fullKey]
			if ok && val == value {
				hit = true
				break
			}
			values = append(values, val)
		}

		if !hit {
			return fmt.Errorf("unexpected assert of %s.%s, expect: %v, got %s", name, key, value, values)
		}

		return nil
	}
}
