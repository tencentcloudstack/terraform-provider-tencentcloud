/*
Use this data source to query detailed information of autoscaling limits

Example Usage

```hcl
data "tencentcloud_autoscaling_limits" "limits" {
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

func dataSourceTencentCloudAutoscalingLimits() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudAutoscalingLimitsRead,
		Schema: map[string]*schema.Schema{
			"max_number_of_launch_configurations": {
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "Maximum number of launch configurations allowed for creation by the user account.",
			},

			"number_of_launch_configurations": {
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "Current number of launch configurations under the user account.",
			},

			"max_number_of_auto_scaling_groups": {
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "Maximum number of auto scaling groups allowed for creation by the user account.",
			},

			"number_of_auto_scaling_groups": {
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "Current number of auto scaling groups under the user account.",
			},

			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results.",
			},
		},
	}
}

func dataSourceTencentCloudAutoscalingLimitsRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_autoscaling_limits.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	paramMap := make(map[string]interface{})
	service := AutoscalingService{client: meta.(*TencentCloudClient).apiV3Conn}

	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeAutoscalingLimitsByFilter(ctx, paramMap)
		if e != nil {
			return retryError(e)
		}
		maxNumberOfLaunchConfigurations = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(maxNumberOfLaunchConfigurations))
	if maxNumberOfLaunchConfigurations != nil {
		_ = d.Set("max_number_of_launch_configurations", maxNumberOfLaunchConfigurations)
	}

	if numberOfLaunchConfigurations != nil {
		_ = d.Set("number_of_launch_configurations", numberOfLaunchConfigurations)
	}

	if maxNumberOfAutoScalingGroups != nil {
		_ = d.Set("max_number_of_auto_scaling_groups", maxNumberOfAutoScalingGroups)
	}

	if numberOfAutoScalingGroups != nil {
		_ = d.Set("number_of_auto_scaling_groups", numberOfAutoScalingGroups)
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
