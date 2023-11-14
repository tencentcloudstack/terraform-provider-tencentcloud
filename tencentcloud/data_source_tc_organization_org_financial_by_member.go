/*
Use this data source to query detailed information of organization org_financial_by_member

Example Usage

```hcl
data "tencentcloud_organization_org_financial_by_member" "org_financial_by_member" {
  month = &lt;nil&gt;
  end_month = &lt;nil&gt;
  member_uins = &lt;nil&gt;
  product_codes = &lt;nil&gt;
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

func dataSourceTencentCloudOrganizationOrgFinancialByMember() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudOrganizationOrgFinancialByMemberRead,
		Schema: map[string]*schema.Schema{
			"month": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Query for the start month. Format:yyyy-mm, for example:2021-01.",
			},

			"end_month": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Query for the end month. Format:yyyy-mm, for example:2021-01.The default value is the `Month`.",
			},

			"member_uins": {
				Optional: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
				Description: "Member uin list. Up to 100.",
			},

			"product_codes": {
				Optional: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "Product code list. Up to 100.",
			},

			"total_cost": {
				Computed:    true,
				Type:        schema.TypeFloat,
				Description: "Total cost of the member(¥).",
			},

			"items": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Member financial detail.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"member_uin": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Member uin.",
						},
						"member_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Member name.",
						},
						"total_cost": {
							Type:        schema.TypeFloat,
							Computed:    true,
							Description: "Total cost of the member(¥).",
						},
						"ratio": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The percentage of the organization total cost that is accounted for by the member.",
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

func dataSourceTencentCloudOrganizationOrgFinancialByMemberRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_organization_org_financial_by_member.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("month"); ok {
		paramMap["Month"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("end_month"); ok {
		paramMap["EndMonth"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("member_uins"); ok {
		memberUinsSet := v.(*schema.Set).List()
		for i := range memberUinsSet {
			memberUins := memberUinsSet[i].(int)
			paramMap["MemberUins"] = append(paramMap["MemberUins"], helper.IntInt64(memberUins))
		}
	}

	if v, ok := d.GetOk("product_codes"); ok {
		productCodesSet := v.(*schema.Set).List()
		paramMap["ProductCodes"] = helper.InterfacesStringsPoint(productCodesSet)
	}

	service := OrganizationService{client: meta.(*TencentCloudClient).apiV3Conn}

	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeOrganizationOrgFinancialByMemberByFilter(ctx, paramMap)
		if e != nil {
			return retryError(e)
		}
		totalCost = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(totalCost))
	if totalCost != nil {
		_ = d.Set("total_cost", totalCost)
	}

	if items != nil {
		for _, orgMemberFinancial := range items {
			orgMemberFinancialMap := map[string]interface{}{}

			if orgMemberFinancial.MemberUin != nil {
				orgMemberFinancialMap["member_uin"] = orgMemberFinancial.MemberUin
			}

			if orgMemberFinancial.MemberName != nil {
				orgMemberFinancialMap["member_name"] = orgMemberFinancial.MemberName
			}

			if orgMemberFinancial.TotalCost != nil {
				orgMemberFinancialMap["total_cost"] = orgMemberFinancial.TotalCost
			}

			if orgMemberFinancial.Ratio != nil {
				orgMemberFinancialMap["ratio"] = orgMemberFinancial.Ratio
			}

			ids = append(ids, *orgMemberFinancial.Month)
			tmpList = append(tmpList, orgMemberFinancialMap)
		}

		_ = d.Set("items", tmpList)
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
