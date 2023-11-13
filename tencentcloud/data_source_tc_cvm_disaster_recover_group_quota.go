/*
Use this data source to query detailed information of cvm disaster_recover_group_quota

Example Usage

```hcl
data "tencentcloud_cvm_disaster_recover_group_quota" "disaster_recover_group_quota" {
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

func dataSourceTencentCloudCvmDisasterRecoverGroupQuota() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudCvmDisasterRecoverGroupQuotaRead,
		Schema: map[string]*schema.Schema{
			"group_quota": {
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "The maximum number of placement groups that can be created.",
			},

			"current_num": {
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "The number of placement groups that have been created by the current user.",
			},

			"cvm_in_host_group_quota": {
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "Quota on instances in a physical-machine-type disaster recovery group.",
			},

			"cvm_in_sw_group_quota": {
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "Quota on instances in a switch-type disaster recovery group.",
			},

			"cvm_in_rack_group_quota": {
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "Quota on instances in a rack-type disaster recovery group.",
			},

			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results.",
			},
		},
	}
}

func dataSourceTencentCloudCvmDisasterRecoverGroupQuotaRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_cvm_disaster_recover_group_quota.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	paramMap := make(map[string]interface{})
	service := CvmService{client: meta.(*TencentCloudClient).apiV3Conn}

	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeCvmDisasterRecoverGroupQuotaByFilter(ctx, paramMap)
		if e != nil {
			return retryError(e)
		}
		groupQuota = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(groupQuota))
	if groupQuota != nil {
		_ = d.Set("group_quota", groupQuota)
	}

	if currentNum != nil {
		_ = d.Set("current_num", currentNum)
	}

	if cvmInHostGroupQuota != nil {
		_ = d.Set("cvm_in_host_group_quota", cvmInHostGroupQuota)
	}

	if cvmInSwGroupQuota != nil {
		_ = d.Set("cvm_in_sw_group_quota", cvmInSwGroupQuota)
	}

	if cvmInRackGroupQuota != nil {
		_ = d.Set("cvm_in_rack_group_quota", cvmInRackGroupQuota)
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
