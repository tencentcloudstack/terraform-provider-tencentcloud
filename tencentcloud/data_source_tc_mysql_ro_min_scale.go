/*
Use this data source to query detailed information of mysql ro_min_scale

Example Usage

```hcl
data "tencentcloud_mysql_ro_min_scale" "ro_min_scale" {
  # ro_instance_id = ""
  master_instance_id = "cdb-fitq5t9h"
}
```
*/
package tencentcloud

import (
	"context"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	cdb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cdb/v20170320"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func dataSourceTencentCloudMysqlRoMinScale() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudMysqlRoMinScaleRead,
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

func dataSourceTencentCloudMysqlRoMinScaleRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_mysql_ro_min_scale.read")()
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

	service := MysqlService{client: meta.(*TencentCloudClient).apiV3Conn}
	var minScale *cdb.DescribeRoMinScaleResponseParams
	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeMysqlRoMinScaleByFilter(ctx, paramMap)
		if e != nil {
			return retryError(e)
		}
		minScale = result
		return nil
	})
	if err != nil {
		return err
	}

	if minScale.Memory != nil {
		_ = d.Set("memory", minScale.Memory)
	}

	if minScale.Volume != nil {
		_ = d.Set("volume", minScale.Volume)
	}

	d.SetId(helper.DataResourceIdsHash([]string{strconv.FormatInt(*minScale.Memory, 10), strconv.FormatInt(*minScale.Volume, 10)}))
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := writeToFile(output.(string), map[string]interface{}{
			"memory": minScale.Memory,
			"volume": minScale.Volume,
		}); e != nil {
			return e
		}
	}
	return nil
}
