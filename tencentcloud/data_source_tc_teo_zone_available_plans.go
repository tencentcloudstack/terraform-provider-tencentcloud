/*
Use this data source to query zone available plans.

Example Usage

```hcl
data "tencentcloud_teo_zone_available_plans" "available_plans" {}
```
*/
package tencentcloud

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	teo "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220106"
)

func dataSourceTencentCloudTeoZoneAvailablePlans() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudTeoZoneAvailablePlansRead,

		Schema: map[string]*schema.Schema{
			"plan_info_list": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Available plans for a zone.",
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
							Description: "Currency type. Valid values: `CNY`, `USD`.",
						},
						"flux": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The number of fluxes included in the zone plan. Unit: Byte.",
						},
						"frequency": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Billing cycle. Valid values: `y`: Billed by the year; `m`: Billed by the month; `h`: Billed by the hour; `M`: Billed by the minute; `s`: Billed by the second.",
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
					},
				},
			},
			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used for save results.",
			},
		},
	}
}

func dataSourceTencentCloudTeoZoneAvailablePlansRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_teo_zone_available_plans.read")()

	var (
		logId          = getLogId(contextNil)
		ctx            = context.WithValue(context.TODO(), logIdKey, logId)
		service        = TeoService{client: meta.(*TencentCloudClient).apiV3Conn}
		availablePlans *teo.DescribeAvailablePlansResponseParams
		err            error
	)

	var outErr, inErr error
	availablePlans, outErr = service.DescribeAvailablePlans(ctx)
	if outErr != nil {
		outErr = resource.Retry(readRetryTimeout, func() *resource.RetryError {
			availablePlans, inErr = service.DescribeAvailablePlans(ctx)
			if inErr != nil {
				return retryError(inErr)
			}
			return nil
		})
	}

	planInfos := availablePlans.PlanInfoList
	planInfoList := make([]map[string]interface{}, 0, len(planInfos))
	for _, v := range planInfos {
		planInfo := map[string]interface{}{
			"plan_type":   v.PlanType,
			"currency":    v.Currency,
			"flux":        v.Flux,
			"frequency":   v.Frequency,
			"price":       v.Price,
			"request":     v.Request,
			"site_number": v.SiteNumber,
		}
		planInfoList = append(planInfoList, planInfo)
	}
	if err = d.Set("plan_info_list", planInfoList); err != nil {
		log.Printf("[CRITAL]%s provider set list fail, reason:%s", logId, err.Error())
		return err
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
