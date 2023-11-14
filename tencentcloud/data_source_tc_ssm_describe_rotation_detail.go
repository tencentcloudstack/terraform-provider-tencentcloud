/*
Use this data source to query detailed information of ssm describe_rotation_detail

Example Usage

```hcl
data "tencentcloud_ssm_describe_rotation_detail" "describe_rotation_detail" {
  secret_name = ""
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

func dataSourceTencentCloudSsmDescribeRotationDetail() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudSsmDescribeRotationDetailRead,
		Schema: map[string]*schema.Schema{
			"secret_name": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Secret name.",
			},

			"enable_rotation": {
				Computed:    true,
				Type:        schema.TypeBool,
				Description: "Whether to allow rotation.",
			},

			"frequency": {
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "The rotation frequency, in days, defaults to 1 day.",
			},

			"latest_rotate_time": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Time of last rotation.",
			},

			"next_rotate_begin_time": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "The time to start the next rotation.",
			},

			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results.",
			},
		},
	}
}

func dataSourceTencentCloudSsmDescribeRotationDetailRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_ssm_describe_rotation_detail.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("secret_name"); ok {
		paramMap["SecretName"] = helper.String(v.(string))
	}

	service := SsmService{client: meta.(*TencentCloudClient).apiV3Conn}

	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeSsmDescribeRotationDetailByFilter(ctx, paramMap)
		if e != nil {
			return retryError(e)
		}
		enableRotation = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(enableRotation))
	if enableRotation != nil {
		_ = d.Set("enable_rotation", enableRotation)
	}

	if frequency != nil {
		_ = d.Set("frequency", frequency)
	}

	if latestRotateTime != nil {
		_ = d.Set("latest_rotate_time", latestRotateTime)
	}

	if nextRotateBeginTime != nil {
		_ = d.Set("next_rotate_begin_time", nextRotateBeginTime)
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
