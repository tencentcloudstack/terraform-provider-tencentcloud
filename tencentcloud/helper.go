package tencentcloud

import (
	"bytes"
	"crypto/md5"
	"encoding/base64"
	"fmt"
	"math/rand"
	"time"

	"github.com/hashicorp/terraform/helper/hashcode"
	"github.com/hashicorp/terraform/helper/schema"
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
