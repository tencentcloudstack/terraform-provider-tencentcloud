/*
Use this data source to query detailed information of postgres recovery_time

Example Usage

```hcl
data "tencentcloud_postgres_recovery_time" "recovery_time" {
  d_b_instance_id = ""
      tags = {
    "createdBy" = "terraform"
  }
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

func dataSourceTencentCloudPostgresRecoveryTime() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudPostgresRecoveryTimeRead,
		Schema: map[string]*schema.Schema{
			"d_b_instance_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Instance ID.",
			},

			"recovery_begin_time": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "The earliest restoration time (UTC+8).",
			},

			"recovery_end_time": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "The latest restoration time (UTC+8).",
			},

			"tags": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "Tag description list.",
			},
			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results.",
			},
		},
	}
}

func dataSourceTencentCloudPostgresRecoveryTimeRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_postgres_recovery_time.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("d_b_instance_id"); ok {
		paramMap["DBInstanceId"] = helper.String(v.(string))
	}

	service := PostgresService{client: meta.(*TencentCloudClient).apiV3Conn}

	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribePostgresRecoveryTimeByFilter(ctx, paramMap)
		if e != nil {
			return retryError(e)
		}
		recoveryBeginTime = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(recoveryBeginTime))
	if recoveryBeginTime != nil {
		_ = d.Set("recovery_begin_time", recoveryBeginTime)
	}

	if recoveryEndTime != nil {
		_ = d.Set("recovery_end_time", recoveryEndTime)
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
