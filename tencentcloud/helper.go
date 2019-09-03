package tencentcloud

import (
	"bytes"
	"crypto/md5"
	"encoding/base64"
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/hashicorp/terraform/helper/hashcode"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/zqfan/tencentcloud-sdk-go/services/cvm/v20170312"
)

// Generates a hash for the set hash function used by the IDs
func dataResourceIdsHash(ids []string) string {
	var buf bytes.Buffer

	for _, id := range ids {
		buf.WriteString(fmt.Sprintf("%s-", id))
	}

	return fmt.Sprintf("%d", hashcode.String(buf.String()))
}

// Generates a hash for the set hash function used by the ID
func dataResourceIdHash(id string) string {
	return fmt.Sprintf("%d", hashcode.String(id))
}

// Transform filter condition to API's param
func buildFiltersParam(params map[string]string, filterList *schema.Set, maxFiltersLimit int, maxFilterValuesLimit int) error {
	if len(filterList.List()) > maxFiltersLimit {
		return fmt.Errorf("Too many filters, should not be more than %v", maxFiltersLimit)
	}
	for i, v := range filterList.List() {
		paramsKeyFilterValues := fmt.Sprintf("Filters.%v", i)

		m := v.(map[string]interface{})
		name := m["name"].(string)
		filterValues := m["values"].([]interface{})
		if len(filterValues) > maxFilterValuesLimit {
			return fmt.Errorf("Too many filter values, should not be more than %v", maxFilterValuesLimit)
		}

		paramsKeyFilterName := fmt.Sprintf("Filters.%v.Name", i)
		params[paramsKeyFilterName] = name
		for j, e := range filterValues {
			filterValue := e.(string)
			if len(filterValue) == 0 {
				return fmt.Errorf("One of the filter value for name: %v is empty", name)
			}

			paramsKeyFilterValues += fmt.Sprintf(".Values.%v", j)
			params[paramsKeyFilterValues] = e.(string)
		}
	}
	return nil
}

// Transform filter condition to TencentCloud Go SDK's param
func buildFiltersParamForSDK(filterList *schema.Set) (r []*cvm.Filter) {
	for _, v := range filterList.List() {
		m := v.(map[string]interface{})
		name := m["name"].(string)
		filterValues := m["values"].([]interface{})

		filter := &cvm.Filter{}
		filter.Name = &name
		for _, fv := range filterValues {
			filterValue := fv.(string)
			filter.Values = append(filter.Values, &filterValue)
		}
		r = append(r, filter)
	}
	return
}

func retryable(code string, msg string) bool {
	msg = strings.ToLower(msg)
	return code == "InternalError" && strings.Contains(msg, "retry")
}

// Takes the result of flatmap.Expand for an array of strings
// and returns a []string
func expandStringList(configured []interface{}) []string {
	vs := make([]string, 0, len(configured))
	for _, v := range configured {
		vs = append(vs, v.(string))
	}
	return vs
}

// Flatten to an array of raw strings and returns a []interface{}
func flattenStringList(list []*string) []interface{} {
	vs := make([]interface{}, 0, len(list))
	for _, v := range list {
		vs = append(vs, *v)
	}
	return vs
}

func flattenIntList(list []*uint64) []interface{} {
	vi := make([]interface{}, 0, len(list))
	for _, v := range list {
		vi = append(vi, int(*v))
	}
	return vi
}

func pointerToString(pointer *string) string {
	if pointer == nil {
		return ""
	}
	return *pointer
}

func stringToPointer(s string) *string {
	return &s
}

func intToPointer(i int) *uint64 {
	u := uint64(i)
	return &u
}

func uint64Pt(i uint64) *uint64 {
	return &i
}

func int64Pt(i int64) *int64 {
	return &i
}

func getTags(d *schema.ResourceData, k string) map[string]string {
	var tags map[string]string
	if raw, ok := d.GetOk(k); ok {
		rawTags := raw.(map[string]interface{})
		tags = make(map[string]string, len(rawTags))
		for k, v := range rawTags {
			tags[k] = v.(string)
		}
	}
	return tags
}

func buildToken() string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	buf := make([]byte, 16)
	r.Read(buf)
	sum := md5.Sum(buf)
	return base64.StdEncoding.EncodeToString(sum[:])
}

func boolToPointer(b bool) *bool {
	return &b
}

func int64ToPointer(n int) *int64 {
	i64 := int64(n)
	return &i64
}

func formatUnixTime(n uint64) string {
	return time.Unix(int64(n), 0).UTC().Format("2006-01-02T03:04:05Z")
}

func parseTime(s string) (time.Time, error) {
	t, err := time.ParseInLocation("2006-01-02T03:04:05Z", s, time.UTC)
	if err != nil {
		return time.Time{}, err
	}
	return t, nil
}
