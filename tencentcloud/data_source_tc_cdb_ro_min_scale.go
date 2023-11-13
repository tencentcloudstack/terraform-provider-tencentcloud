/*
Use this data source to query detailed information of cdb ro_min_scale

Example Usage

```hcl
data "tencentcloud_cdb_ro_min_scale" "ro_min_scale" {
  ro_instance_id = ""
  master_instance_id = ""
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

func dataSourceTencentCloudCdbRoMinScale() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudCdbRoMinScaleRead,
		Schema: map[string]*schema.Schema{
			"ro_instance_id": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "The read-only instance ID, in the format: cdbro-c1nl9rpv, is the same as the instance ID displayed on the cloud database console page. This parameter and the MasterInstanceId parameter cannot be empty at the same time.",
			},

			"master_instance_id": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "The primary instance ID, in the format: cdb-c1nl9rpv, is the same as the instance ID displayed on the cloud database console page. This parameter and the RoInstanceId parameter cannot be empty at the same time. Note that when the input parameter contains RoInstanceId, the return value is the minimum specification when the read-only instance is upgraded; when the input parameter only contains MasterInstanceId, the return value is the minimum specification when the read-only instance is purchased.",
			},

			"memory": {
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "Memory specification size, unit: MB.",
			},

			"volume": {
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "Disk specification size, unit: GB.",
			},

			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results.",
			},
		},
	}
}

func dataSourceTencentCloudCdbRoMinScaleRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_cdb_ro_min_scale.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("ro_instance_id"); ok {
		paramMap["RoInstanceId"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("master_instance_id"); ok {
		paramMap["MasterInstanceId"] = helper.String(v.(string))
	}

	service := CdbService{client: meta.(*TencentCloudClient).apiV3Conn}

	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeCdbRoMinScaleByFilter(ctx, paramMap)
		if e != nil {
			return retryError(e)
		}
		memory = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(memory))
	if memory != nil {
		_ = d.Set("memory", memory)
	}

	if volume != nil {
		_ = d.Set("volume", volume)
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
