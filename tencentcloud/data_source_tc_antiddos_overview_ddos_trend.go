package tencentcloud

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	antiddos "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/antiddos/v20200309"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func dataSourceTencentCloudAntiddosOverviewDdosTrend() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudAntiddosOverviewDdosTrendRead,
		Schema: map[string]*schema.Schema{
			"period": {
				Required:    true,
				Type:        schema.TypeInt,
				Description: "Statistical granularity, values [300 (5 minutes), 3600 (hours), 86400 (days)].",
			},

			"start_time": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "StartTime.",
			},

			"end_time": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "EndTime.",
			},

			"metric_name": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Indicator, value [bps (attack traffic bandwidth, pps (attack packet rate)].",
			},

			"business": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Dayu sub product code (bgpip represents advanced defense IP; net represents professional version of advanced defense IP).",
			},

			"ip_list": {
				Optional: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "instance IpList.",
			},

			"data": {
				Computed: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
				Description: "Array, attack traffic bandwidth in Mbps, packet rate in pps.",
			},

			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results.",
			},
		},
	}
}

func dataSourceTencentCloudAntiddosOverviewDdosTrendRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_antiddos_overview_ddos_trend.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, _ := d.GetOk("period"); v != nil {
		paramMap["Period"] = helper.IntInt64(v.(int))
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

	if v, ok := d.GetOk("business"); ok {
		paramMap["Business"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("ip_list"); ok {
		ipListSet := v.(*schema.Set).List()
		paramMap["IpList"] = helper.InterfacesStringsPoint(ipListSet)
	}

	service := AntiddosService{client: meta.(*TencentCloudClient).apiV3Conn}

	var describeOverviewDDoSTrendResponseParams *antiddos.DescribeOverviewDDoSTrendResponseParams
	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeAntiddosOverviewDdosTrendByFilter(ctx, paramMap)
		if e != nil {
			return retryError(e)
		}
		describeOverviewDDoSTrendResponseParams = result
		return nil
	})
	if err != nil {
		return err
	}

	resultMap := make(map[string]interface{})
	if describeOverviewDDoSTrendResponseParams.Data != nil {
		resultMap["data"] = describeOverviewDDoSTrendResponseParams.Data
		_ = d.Set("data", describeOverviewDDoSTrendResponseParams.Data)
	}

	d.SetId(helper.BuildToken())
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := writeToFile(output.(string), resultMap); e != nil {
			return e
		}
	}
	return nil
}
