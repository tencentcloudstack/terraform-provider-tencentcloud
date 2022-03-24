package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	emr "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/emr/v20190103"
)

const (
	EmrInternetStatusCreated int64 = 2
	EmrInternetStatusDeleted int64 = 201
)

const (
	DisplayStrategyIsclusterList = "clusterList"
)

const (
	EMR_MASTER_WAN_TYPE_NEED_MASTER_WAN     = "NEED_MASTER_WAN"
	EMR_MASTER_WAN_TYPE_NOT_NEED_MASTER_WAN = "NOT_NEED_MASTER_WAN"
)

var EMR_MASTER_WAN_TYPES = []string{EMR_MASTER_WAN_TYPE_NEED_MASTER_WAN, EMR_MASTER_WAN_TYPE_NOT_NEED_MASTER_WAN}

func buildResourceSpecSchema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Optional: true,
		ForceNew: true,
		MaxItems: 1,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"spec":         {Type: schema.TypeString, Optional: true},
				"storage_type": {Type: schema.TypeInt, Optional: true},
				"disk_type":    {Type: schema.TypeString, Optional: true},
				"mem_size":     {Type: schema.TypeInt, Optional: true},
				"cpu":          {Type: schema.TypeInt, Optional: true},
				"disk_size":    {Type: schema.TypeInt, Optional: true},
				"root_size":    {Type: schema.TypeInt, Optional: true},
			},
		},
	}
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
			resultResource.StorageType = common.Int64Ptr((int64)(v.(int)))
		} else if k == "disk_type" {
			resultResource.DiskType = common.StringPtr(v.(string))
		} else if k == "mem_size" {
			resultResource.MemSize = common.Int64Ptr((int64)(v.(int)))
		} else if k == "cpu" {
			resultResource.Cpu = common.Int64Ptr((int64)(v.(int)))
		} else if k == "disk_size" {
			resultResource.DiskSize = common.Int64Ptr((int64)(v.(int)))
		} else if k == "root_size" {
			resultResource.RootSize = common.Int64Ptr((int64)(v.(int)))
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
