package tencentcloud

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	teo "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901"
)

func dataSourceTencentCloudTeoZoneAvailablePlans() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudTeoZoneAvailablePlansRead,

		Schema: map[string]*schema.Schema{
			"plan_info_list": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Zone plans which current account can use.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"plan_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Plan type.",
						},
						"currency": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Settlement Currency Type. Valid values: `CNY`, `USD`.",
						},
						"flux": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The number of fluxes included in the zone plan. Unit: Byte.",
						},
						"frequency": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Billing cycle. Valid values:- `y`: Billed by the year.- `m`: Billed by the month.- `h`: Billed by the hour.- `M`: Billed by the minute.- `s`: Billed by the second.",
						},
						"price": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Price of the plan. Unit: cent.",
						},
						"request": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The number of requests included in the zone plan.",
						},
						"site_number": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The number of zones this zone plan can bind.",
						},
						"area": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Acceleration area of the plan. Valid value: `mainland`, `overseas`.",
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

func dataSourceTencentCloudTeoZoneAvailablePlansRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_teo_zone_available_plans.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	teoService := TeoService{client: meta.(*TencentCloudClient).apiV3Conn}

	var planInfos []*teo.PlanInfo
	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		results, e := teoService.DescribeTeoZoneAvailablePlansByFilter(ctx)
		if e != nil {
			return retryError(e)
		}
		planInfos = results
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s read Teo planInfo failed, reason:%+v", logId, err)
		return err
	}

	planInfoList := []interface{}{}
	if planInfos != nil {
		for _, planInfo := range planInfos {
			planInfoMap := map[string]interface{}{}
			if planInfo.PlanType != nil {
				planInfoMap["plan_type"] = planInfo.PlanType
			}
			if planInfo.Currency != nil {
				planInfoMap["currency"] = planInfo.Currency
			}
			if planInfo.Flux != nil {
				planInfoMap["flux"] = planInfo.Flux
			}
			if planInfo.Frequency != nil {
				planInfoMap["frequency"] = planInfo.Frequency
			}
			if planInfo.Price != nil {
				planInfoMap["price"] = planInfo.Price
			}
			if planInfo.Request != nil {
				planInfoMap["request"] = planInfo.Request
			}
			if planInfo.SiteNumber != nil {
				planInfoMap["site_number"] = planInfo.SiteNumber
			}
			if planInfo.Area != nil {
				planInfoMap["area"] = planInfo.Area
			}

			planInfoList = append(planInfoList, planInfoMap)
		}
		_ = d.Set("plan_info_list", planInfoList)
	}

	d.SetId("zone_available_plans")

	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := writeToFile(output.(string), planInfoList); e != nil {
			return e
		}
	}

	return nil
}
