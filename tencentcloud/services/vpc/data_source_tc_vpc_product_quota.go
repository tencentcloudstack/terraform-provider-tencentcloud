package vpc

import (
	"context"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	vpc "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/vpc/v20170312"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudVpcProductQuota() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudVpcProductQuotaRead,
		Schema: map[string]*schema.Schema{
			"product": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "The name of the network product to be queried. The products that can be queried are:vpc, ccn, vpn, dc, dfw, clb, eip.",
			},

			"product_quota_set": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "ProductQuota Array.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"quota_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Quota Id.",
						},
						"quota_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Quota name.",
						},
						"quota_current": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Current Quota.",
						},
						"quota_limit": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Quota limit.",
						},
						"quota_region": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Quota region.",
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

func dataSourceTencentCloudVpcProductQuotaRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_vpc_product_quota.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("product"); ok {
		paramMap["Product"] = helper.String(v.(string))
	}

	service := VpcService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	var productQuotaSet []*vpc.ProductQuota

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeVpcProductQuota(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}
		productQuotaSet = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(productQuotaSet))
	tmpList := make([]map[string]interface{}, 0, len(productQuotaSet))

	if productQuotaSet != nil {
		for _, productQuota := range productQuotaSet {
			productQuotaMap := map[string]interface{}{}

			if productQuota.QuotaId != nil {
				productQuotaMap["quota_id"] = productQuota.QuotaId
			}

			if productQuota.QuotaName != nil {
				productQuotaMap["quota_name"] = productQuota.QuotaName
			}

			if productQuota.QuotaCurrent != nil {
				productQuotaMap["quota_current"] = productQuota.QuotaCurrent
			}

			if productQuota.QuotaLimit != nil {
				productQuotaMap["quota_limit"] = productQuota.QuotaLimit
			}

			if productQuota.QuotaRegion != nil {
				productQuotaMap["quota_region"] = productQuota.QuotaRegion
			}

			ids = append(ids, *productQuota.QuotaId)
			tmpList = append(tmpList, productQuotaMap)
		}

		_ = d.Set("product_quota_set", tmpList)
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
