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
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	cvm "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cvm/v20170312"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/ratelimit"
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

	var response *cvm.DescribeDisasterRecoverGroupQuotaResponse

	request := cvm.NewDescribeDisasterRecoverGroupQuotaRequest()
	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseCvmClient().DescribeDisasterRecoverGroupQuota(request)
		if e != nil {
			return retryError(e)
		}
		response = result
		return nil
	})
	if err != nil {
		return err
	}
	if response == nil || response.Response == nil {
		d.SetId("")
		return fmt.Errorf("Response is nil")
	}
	_ = d.Set("group_quota", response.Response.GroupQuota)
	_ = d.Set("current_num", response.Response.CurrentNum)
	_ = d.Set("cvm_in_host_group_quota", response.Response.CvmInHostGroupQuota)
	_ = d.Set("cvm_in_sw_group_quota", response.Response.CvmInSwGroupQuota)
	_ = d.Set("cvm_in_rack_group_quota", response.Response.CvmInRackGroupQuota)

	d.SetId(helper.BuildToken())
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := writeToFile(output.(string), map[string]interface{}{
			"group_quota":             response.Response.GroupQuota,
			"current_num":             response.Response.CurrentNum,
			"cvm_in_host_group_quota": response.Response.CvmInHostGroupQuota,
			"cvm_in_sw_group_quota":   response.Response.CvmInSwGroupQuota,
			"cvm_in_rack_group_quota": response.Response.CvmInRackGroupQuota,
		}); e != nil {
			return e
		}
	}
	return nil
}
