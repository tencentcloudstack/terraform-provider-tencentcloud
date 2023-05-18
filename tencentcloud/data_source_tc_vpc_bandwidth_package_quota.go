/*
Use this data source to query detailed information of vpc bandwidth_package_quota

Example Usage

```hcl
data "tencentcloud_vpc_bandwidth_package_quota" "bandwidth_package_quota" {
  }
```
*/
package tencentcloud

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	vpc "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/vpc/v20170312"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func dataSourceTencentCloudVpcBandwidthPackageQuota() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudVpcBandwidthPackageQuotaRead,
		Schema: map[string]*schema.Schema{
			"quota_set": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Bandwidth Package Quota Details.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"quota_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Quota type.",
						},
						"quota_current": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "current amount.",
						},
						"quota_limit": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "quota amount.",
						},
					},
				},
			},

			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results.",
			},
		},
	}
}

func dataSourceTencentCloudVpcBandwidthPackageQuotaRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_vpc_bandwidth_package_quota.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := VpcService{client: meta.(*TencentCloudClient).apiV3Conn}

	var quotaSet []*vpc.Quota

	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeVpcBandwidthPackageQuota(ctx)
		if e != nil {
			return retryError(e)
		}
		quotaSet = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(quotaSet))
	tmpList := make([]map[string]interface{}, 0, len(quotaSet))

	if quotaSet != nil {
		for _, quota := range quotaSet {
			quotaMap := map[string]interface{}{}

			if quota.QuotaId != nil {
				quotaMap["quota_id"] = quota.QuotaId
			}

			if quota.QuotaCurrent != nil {
				quotaMap["quota_current"] = quota.QuotaCurrent
			}

			if quota.QuotaLimit != nil {
				quotaMap["quota_limit"] = quota.QuotaLimit
			}

			ids = append(ids, *quota.QuotaId)
			tmpList = append(tmpList, quotaMap)
		}

		_ = d.Set("quota_set", tmpList)
	}

	d.SetId(helper.DataResourceIdsHash(ids))
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := writeToFile(output.(string), tmpList); e != nil {
			return e
		}
	}
	return nil
}
