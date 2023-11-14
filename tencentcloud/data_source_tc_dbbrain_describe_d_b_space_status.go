/*
Use this data source to query detailed information of dbbrain describe_d_b_space_status

Example Usage

```hcl
data "tencentcloud_dbbrain_describe_d_b_space_status" "describe_d_b_space_status" {
  instance_id = ""
  range_days =
  product = ""
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

func dataSourceTencentCloudDbbrainDescribeDBSpaceStatus() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudDbbrainDescribeDBSpaceStatusRead,
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Instance id.",
			},

			"range_days": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "The number of days in the time period, the deadline is the current day, and the default is 7 days.",
			},

			"product": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Service product type, supported values include： mysql - cloud database MySQL, cynosdb - cloud database CynosDB for MySQL, the default is mysql.",
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

func dataSourceTencentCloudDbbrainDescribeDBSpaceStatusRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_dbbrain_describe_d_b_space_status.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("instance_id"); ok {
		paramMap["InstanceId"] = helper.String(v.(string))
	}

	if v, _ := d.GetOk("range_days"); v != nil {
		paramMap["RangeDays"] = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("product"); ok {
		paramMap["Product"] = helper.String(v.(string))
	}

	service := DbbrainService{client: meta.(*TencentCloudClient).apiV3Conn}

	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeDbbrainDescribeDBSpaceStatusByFilter(ctx, paramMap)
		if e != nil {
			return retryError(e)
		}
		growth = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(growth))
	if growth != nil {
		_ = d.Set("growth", growth)
	}

	if remain != nil {
		_ = d.Set("remain", remain)
	}

	if total != nil {
		_ = d.Set("total", total)
	}

	if availableDays != nil {
		_ = d.Set("available_days", availableDays)
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
