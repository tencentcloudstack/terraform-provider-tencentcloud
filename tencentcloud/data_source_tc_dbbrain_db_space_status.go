/*
Use this data source to query detailed information of dbbrain db_space_status

Example Usage

```hcl
data "tencentcloud_dbbrain_db_space_status" "db_space_status" {
  instance_id = "%s"
  range_days = 7
  product = "mysql"
}
```
*/
package tencentcloud

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	dbbrain "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dbbrain/v20210527"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func dataSourceTencentCloudDbbrainDbSpaceStatus() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudDbbrainDbSpaceStatusRead,
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "instance id.",
			},

			"range_days": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "The number of days in the time period, the deadline is the current day, and the default is 7 days.",
			},

			"product": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Service product type, supported values include: mysql - cloud database MySQL, cynosdb - cloud database CynosDB for MySQL, the default is mysql.",
			},

			"growth": {
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "Disk growth (MB).",
			},

			"remain": {
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "Disk remaining (MB).",
			},

			"total": {
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "Total disk size (MB).",
			},

			"available_days": {
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "Estimated number of days available.",
			},

			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results.",
			},
		},
	}
}

func dataSourceTencentCloudDbbrainDbSpaceStatusRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_dbbrain_db_space_status.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var instanceId string

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("instance_id"); ok {
		paramMap["InstanceId"] = helper.String(v.(string))
		instanceId = v.(string)
	}

	if v, _ := d.GetOk("range_days"); v != nil {
		paramMap["RangeDays"] = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("product"); ok {
		paramMap["Product"] = helper.String(v.(string))
	}

	service := DbbrainService{client: meta.(*TencentCloudClient).apiV3Conn}

	var rows *dbbrain.DescribeDBSpaceStatusResponseParams

	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeDbbrainDbSpaceStatusByFilter(ctx, paramMap)
		if e != nil {
			return retryError(e)
		}
		rows = result
		return nil
	})
	if err != nil {
		return err
	}

	tmpList := []map[string]interface{}{}

	if rows != nil {

		if rows.Growth != nil {
			_ = d.Set("growth", rows.Growth)
		}

		if rows.Remain != nil {
			_ = d.Set("remain", rows.Remain)
		}

		if rows.Total != nil {
			_ = d.Set("total", rows.Total)
		}

		if rows.AvailableDays != nil {
			_ = d.Set("available_days", rows.AvailableDays)
		}
		tmpList = append(tmpList, map[string]interface{}{
			"growth":         rows.Growth,
			"remain":         rows.Remain,
			"total":          rows.Total,
			"available_days": rows.AvailableDays,
		})

	}

	d.SetId(helper.DataResourceIdHash(instanceId))
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := writeToFile(output.(string), tmpList); e != nil {
			return e
		}
	}
	return nil
}
