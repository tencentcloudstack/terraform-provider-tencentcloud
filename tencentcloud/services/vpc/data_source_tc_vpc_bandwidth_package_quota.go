package vpc

import (
	"context"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	vpc "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/vpc/v20170312"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudVpcBandwidthPackageQuota() *schema.Resource {
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
	defer tccommon.LogElapsed("data_source.tencentcloud_vpc_bandwidth_package_quota.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := VpcService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	var quotaSet []*vpc.Quota

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeVpcBandwidthPackageQuota(ctx)
		if e != nil {
			return tccommon.RetryError(e)
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
		if e := tccommon.WriteToFile(output.(string), tmpList); e != nil {
			return e
		}
	}
	return nil
}
