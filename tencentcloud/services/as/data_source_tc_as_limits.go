package as

import (
	"context"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	as "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/as/v20180419"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudAsLimits() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudAsLimitsRead,
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

func dataSourceTencentCloudAsLimitsRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_as_limits.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := AsService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	var limit *as.DescribeAccountLimitsResponseParams

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeAsLimits(ctx)
		if e != nil {
			return tccommon.RetryError(e)
		}
		limit = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0)
	asLimitMap := map[string]interface{}{}
	if limit.MaxNumberOfLaunchConfigurations != nil {
		_ = d.Set("max_number_of_launch_configurations", limit.MaxNumberOfLaunchConfigurations)
		asLimitMap["max_number_of_launch_configurations"] = limit.MaxNumberOfLaunchConfigurations
	}

	if limit.NumberOfLaunchConfigurations != nil {
		_ = d.Set("number_of_launch_configurations", limit.NumberOfLaunchConfigurations)
		asLimitMap["number_of_launch_configurations"] = limit.NumberOfLaunchConfigurations
	}

	if limit.MaxNumberOfAutoScalingGroups != nil {
		_ = d.Set("max_number_of_auto_scaling_groups", limit.MaxNumberOfAutoScalingGroups)
		asLimitMap["max_number_of_auto_scaling_groups"] = limit.MaxNumberOfAutoScalingGroups
	}

	if limit.NumberOfAutoScalingGroups != nil {
		_ = d.Set("number_of_auto_scaling_groups", limit.NumberOfAutoScalingGroups)
		asLimitMap["number_of_auto_scaling_groups"] = limit.NumberOfAutoScalingGroups
	}

	d.SetId(helper.DataResourceIdsHash(ids))
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), asLimitMap); e != nil {
			return e
		}
	}
	return nil
}
