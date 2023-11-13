/*
Use this data source to query detailed information of cdb instance_reboot_time

Example Usage

```hcl
data "tencentcloud_cdb_instance_reboot_time" "instance_reboot_time" {
  instance_ids =
  }
```
*/
package tencentcloud

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func dataSourceTencentCloudCdbInstanceRebootTime() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudCdbInstanceRebootTimeRead,
		Schema: map[string]*schema.Schema{
			"instance_ids": {
				Required: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "The instance ID, in the format: cdb-c1nl9rpv, is the same as the instance ID displayed on the cloud database console page.",
			},

			"items": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Returned parameter information.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"instance_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Instance ID, the format is: cdb-c1nl9rpv, which is the same as the instance ID displayed on the cloud database console page.",
						},
						"time_in_seconds": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Expected restart time.",
						},
					},
				},
			},

			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results.",
			},
		},
	}
}

func dataSourceTencentCloudCdbInstanceRebootTimeRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_cdb_instance_reboot_time.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("instance_ids"); ok {
		instanceIdsSet := v.(*schema.Set).List()
		paramMap["InstanceIds"] = helper.InterfacesStringsPoint(instanceIdsSet)
	}

	service := CdbService{client: meta.(*TencentCloudClient).apiV3Conn}

	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeCdbInstanceRebootTimeByFilter(ctx, paramMap)
		if e != nil {
			return retryError(e)
		}
		totalCount = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(totalCount))
	if items != nil {
		for _, instanceRebootTime := range items {
			instanceRebootTimeMap := map[string]interface{}{}

			if instanceRebootTime.InstanceId != nil {
				instanceRebootTimeMap["instance_id"] = instanceRebootTime.InstanceId
			}

			if instanceRebootTime.TimeInSeconds != nil {
				instanceRebootTimeMap["time_in_seconds"] = instanceRebootTime.TimeInSeconds
			}

			ids = append(ids, *instanceRebootTime.InstanceIds)
			tmpList = append(tmpList, instanceRebootTimeMap)
		}

		_ = d.Set("items", tmpList)
	}

	d.SetId(helper.DataResourceIdsHash(ids))
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := writeToFile(output.(string)); e != nil {
			return e
		}
	}
	return nil
}
