/*
Use this data source to query detailed information of antiddos bgp_biz_trend

# Example Usage

```hcl

data "tencentcloud_antiddos_bgp_biz_trend" "bgp_biz_trend" {
  business = "bgp-multip"
  start_time = "2023-11-22 09:25:00"
  end_time = "2023-11-22 10:25:00"
  metric_name = "intraffic"
  instance_id = "bgp-00000ry7"
  flag = 0
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

func dataSourceTencentCloudAntiddosBgpBizTrend() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudAntiddosBgpBizTrendRead,
		Schema: map[string]*schema.Schema{
			"business": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Dayu sub product code (bgpip represents advanced defense IP; net represents professional version of advanced defense IP).",
			},

			"start_time": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Statistic start time.",
			},

			"end_time": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Statistic end time.",
			},

			"metric_name": {
				Required:     true,
				Type:         schema.TypeString,
				ValidateFunc: validateAllowedStringValue([]string{"intraffic", "outtraffic", "inpkg", "outpkg"}),
				Description:  "Statistic metric name，for example: intraffic, outtraffic, inpkg, outpkg.",
			},

			"instance_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Antiddos InstanceId.",
			},

			"flag": {
				Required:     true,
				Type:         schema.TypeInt,
				ValidateFunc: validateAllowedIntValue([]int{0, 1}),
				Description:  "0 represents fixed time, 1 represents custom time.",
			},

			"data_list": {
				Computed: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
				Description: "Values at various time points on the graph.",
			},

			"total": {
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "Number of values in the curve graph.",
			},

			"max_data": {
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "Returns the maximum value of an array.",
			},

			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results.",
			},
		},
	}
}

func dataSourceTencentCloudAntiddosBgpBizTrendRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_antiddos_bgp_biz_trend.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("business"); ok {
		paramMap["Business"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("start_time"); ok {
		paramMap["StartTime"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("end_time"); ok {
		paramMap["EndTime"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("metric_name"); ok {
		paramMap["MetricName"] = helper.String(v.(string))
	}

	var instanceId string
	if v, ok := d.GetOk("instance_id"); ok {
		instanceId = v.(string)
		paramMap["InstanceId"] = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("flag"); ok {
		paramMap["Flag"] = helper.IntUint64(v.(int))
	}

	service := AntiddosService{client: meta.(*TencentCloudClient).apiV3Conn}

	var bgpBizTrend *antiddos.DescribeBgpBizTrendResponseParams
	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeAntiddosBgpBizTrendByFilter(ctx, paramMap)
		if e != nil {
			return retryError(e)
		}
		bgpBizTrend = result
		return nil
	})
	if err != nil {
		return err
	}

	if bgpBizTrend.DataList != nil {
		_ = d.Set("data_list", bgpBizTrend.DataList)
	}

	if bgpBizTrend.Total != nil {
		_ = d.Set("total", bgpBizTrend.Total)
	}

	if bgpBizTrend.MetricName != nil {
		_ = d.Set("metric_name", bgpBizTrend.MetricName)
	}

	if bgpBizTrend.MaxData != nil {
		_ = d.Set("max_data", bgpBizTrend.MaxData)
	}

	tmpList := make([]map[string]interface{}, 0)
	mapping := map[string]interface{}{
		"data_list":   bgpBizTrend.DataList,
		"total":       bgpBizTrend.Total,
		"metric_name": bgpBizTrend.MetricName,
		"max_data":    bgpBizTrend.MaxData,
	}
	tmpList = append(tmpList, mapping)

	d.SetId(instanceId)
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := writeToFile(output.(string), tmpList); e != nil {
			return e
		}
	}
	return nil
}
