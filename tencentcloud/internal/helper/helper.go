package helper

import (
	"bytes"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/hashcode"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/pkg/errors"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	emr "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/emr/v20190103"
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
func ParseMultiDisks(_multiDisks []map[string]interface{}) []*emr.MultiDisk {
	multiDisks := make([]*emr.MultiDisk, len(_multiDisks))
	for _, item := range _multiDisks {
		var diskType string
		var volume int64
		var count int64
		for subK, subV := range item {
			if subK == "disk_type" {
				diskType = subV.(string)
			} else if subK == "volume" {
				volume = subV.(int64)
			} else if subK == "count" {
				count = subV.(int64)
			}
		}
		multiDisks = append(multiDisks,
			&emr.MultiDisk{
				DiskType: common.StringPtr(diskType),
				Volume:   common.Int64Ptr(volume),
				Count:    common.Int64Ptr(count),
			})
	}

	return multiDisks
}

func ParseTags(_tags []map[string]string) []*emr.Tag {
	tags := make([]*emr.Tag, len(_tags))
	for _, item := range _tags {
		for k, v := range item {
			tags = append(tags, &emr.Tag{TagKey: &k, TagValue: &v})
		}
	}
	return tags
}

func ParseResource(_resource map[string]interface{}) *emr.Resource {
	resultResource := &emr.Resource{}
	for k, v := range _resource {
		if k == "spec" {
			resultResource.Spec = common.StringPtr(v.(string))
		} else if k == "storage_type" {
			resultResource.StorageType = common.Int64Ptr(v.(int64))
		} else if k == "disk_type" {
			resultResource.DiskType = common.StringPtr(v.(string))
		} else if k == "mem_size" {
			resultResource.MemSize = common.Int64Ptr(v.(int64))
		} else if k == "cpu" {
			resultResource.Cpu = common.Int64Ptr(v.(int64))
		} else if k == "disk_size" {
			resultResource.DiskSize = common.Int64Ptr(v.(int64))
		} else if k == "root_size" {
			resultResource.RootSize = common.Int64Ptr(v.(int64))
		} else if k == "multi_disks" {
			multiDisks := v.([]map[string]interface{})
			resultResource.MultiDisks = ParseMultiDisks(multiDisks)
		} else if k == "tags" {
			tags := v.([]map[string]string)
			resultResource.Tags = ParseTags(tags)
		} else if k == "instance_type" {
			resultResource.InstanceType = common.StringPtr(v.(string))
		} else if k == "local_disk_num" {
			resultResource.LocalDiskNum = common.Uint64Ptr(v.(uint64))
		} else if k == "DiskNum" {
			resultResource.DiskNum = common.Uint64Ptr(v.(uint64))
		}
	}
	return resultResource
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
