/*
Use this data source to query detailed information of antiddos overview_attack_trend

Example Usage

```hcl
data "tencentcloud_antiddos_overview_attack_trend" "overview_attack_trend" {
  type       = "ddos"
  dimension  = "attackcount"
  period     = 86400
  start_time = "2023-11-21 10:28:31"
  end_time   = "2023-11-22 10:28:31"
}
```
*/
package tencentcloud

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	antiddos "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/antiddos/v20200309"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func dataSourceTencentCloudAntiddosOverviewAttackTrend() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudAntiddosOverviewAttackTrendRead,
		Schema: map[string]*schema.Schema{
			"type": {
				Required:     true,
				Type:         schema.TypeString,
				ValidateFunc: validateAllowedStringValue([]string{"cc", "ddos"}),
				Description:  "Attack type: cc, ddos.",
			},

			"dimension": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Latitude, currently only attackcount is supported.",
			},

			"period": {
				Required:    true,
				Type:        schema.TypeInt,
				Description: "Period, currently only 86400 is supported.",
			},

			"start_time": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Protection Overview Attack Trend Start Time.",
			},

			"end_time": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Protection Overview Attack Trend End Time.",
			},

			"data": {
				Computed: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
				Description: "Number of attacks per cycle point.",
			},

			"period_point_count": {
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "Number of period points included.",
			},

			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results.",
			},
		},
	}
}

func dataSourceTencentCloudAntiddosOverviewAttackTrendRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_antiddos_overview_attack_trend.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("type"); ok {
		paramMap["Type"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("dimension"); ok {
		paramMap["Dimension"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("period"); ok {
		paramMap["Period"] = helper.Uint64(uint64(v.(int)))
	}

	if v, ok := d.GetOk("start_time"); ok {
		paramMap["StartTime"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("end_time"); ok {
		paramMap["EndTime"] = helper.String(v.(string))
	}

	service := AntiddosService{client: meta.(*TencentCloudClient).apiV3Conn}

	var overviewAttackTrend *antiddos.DescribeOverviewAttackTrendResponseParams
	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeAntiddosOverviewAttackTrendByFilter(ctx, paramMap)
		if e != nil {
			return retryError(e)
		}
		overviewAttackTrend = result
		return nil
	})
	if err != nil {
		return err
	}

	if overviewAttackTrend.Type != nil {
		_ = d.Set("type", overviewAttackTrend.Type)
	}

	if overviewAttackTrend.Period != nil {
		_ = d.Set("period", overviewAttackTrend.Period)
	}

	if overviewAttackTrend.StartTime != nil {
		_ = d.Set("start_time", overviewAttackTrend.StartTime)
	}

	if overviewAttackTrend.EndTime != nil {
		_ = d.Set("end_time", overviewAttackTrend.EndTime)
	}

	if overviewAttackTrend.Data != nil {
		_ = d.Set("data", overviewAttackTrend.Data)
	}

	if overviewAttackTrend.Count != nil {
		_ = d.Set("period_point_count", overviewAttackTrend.Count)
	}

	tmpList := make([]map[string]interface{}, 0)
	mapping := map[string]interface{}{
		"type":               overviewAttackTrend.Type,
		"period":             overviewAttackTrend.Period,
		"start_time":         overviewAttackTrend.StartTime,
		"end_time":           overviewAttackTrend.EndTime,
		"data":               overviewAttackTrend.Data,
		"period_point_count": overviewAttackTrend.Count,
	}
	tmpList = append(tmpList, mapping)

	d.SetId(helper.BuildToken())
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := writeToFile(output.(string), tmpList); e != nil {
			return e
		}
	}
	return nil
}
