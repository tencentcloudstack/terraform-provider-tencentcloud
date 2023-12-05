package tencentcloud

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	organization "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/organization/v20210331"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func dataSourceTencentCloudOrganizationOrgFinancialByMonth() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudOrganizationOrgFinancialByMonthRead,
		Schema: map[string]*schema.Schema{
			"end_month": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Query for the end month. Format:yyyy-mm, for example:2021-01.",
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

			"items": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Organization financial info by month.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Record ID.",
						},
						"month": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Month.",
						},
						"total_cost": {
							Type:        schema.TypeFloat,
							Computed:    true,
							Description: "Total cost of the month.",
						},
						"growth_rate": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Growth rate compared to last month.",
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

func dataSourceTencentCloudOrganizationOrgFinancialByMonthRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_organization_org_financial_by_month.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("end_month"); ok {
		paramMap["EndMonth"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("member_uins"); ok {
		memberUinsSet := v.(*schema.Set).List()
		paramMap["MemberUins"] = helper.InterfacesIntInt64Point(memberUinsSet)
	}

	if v, ok := d.GetOk("product_codes"); ok {
		productCodesSet := v.(*schema.Set).List()
		paramMap["ProductCodes"] = helper.InterfacesStringsPoint(productCodesSet)
	}

	service := OrganizationService{client: meta.(*TencentCloudClient).apiV3Conn}

	var items []*organization.OrgFinancialByMonth

	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeOrganizationOrgFinancialByMonthByFilter(ctx, paramMap)
		if e != nil {
			return retryError(e)
		}
		items = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(items))
	tmpList := make([]map[string]interface{}, 0, len(items))

	if items != nil {
		for _, orgFinancialByMonth := range items {
			orgFinancialByMonthMap := map[string]interface{}{}

			if orgFinancialByMonth.Id != nil {
				orgFinancialByMonthMap["id"] = orgFinancialByMonth.Id
			}

			if orgFinancialByMonth.Month != nil {
				orgFinancialByMonthMap["month"] = orgFinancialByMonth.Month
			}

			if orgFinancialByMonth.TotalCost != nil {
				orgFinancialByMonthMap["total_cost"] = orgFinancialByMonth.TotalCost
			}

			if orgFinancialByMonth.GrowthRate != nil {
				orgFinancialByMonthMap["growth_rate"] = orgFinancialByMonth.GrowthRate
			}

			ids = append(ids, *orgFinancialByMonth.Month)
			tmpList = append(tmpList, orgFinancialByMonthMap)
		}

		_ = d.Set("items", tmpList)
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
