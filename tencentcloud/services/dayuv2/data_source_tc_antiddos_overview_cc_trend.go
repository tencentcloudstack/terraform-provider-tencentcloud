package dayuv2

import (
	"context"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	svcantiddos "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/antiddos"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	antiddos "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/antiddos/v20200309"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudAntiddosOverviewCcTrend() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudAntiddosOverviewCcTrendRead,
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
				Description: "Indicator, values [inqps (peak total requests, dropqps (peak attack requests)), incount (number of requests), dropcount (number of attacks)].",
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
				Description: "resource id list.",
			},

			"data": {
				Computed: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
				Description: "Data.",
			},

			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results.",
			},
		},
	}
}

func dataSourceTencentCloudAntiddosOverviewCcTrendRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_antiddos_overview_cc_trend.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

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

	service := svcantiddos.NewAntiddosService(meta.(tccommon.ProviderMeta).GetAPIV3Conn())
	var overviewCCTrendResponseParams *antiddos.DescribeOverviewCCTrendResponseParams

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeAntiddosOverviewCcTrendByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}
		overviewCCTrendResponseParams = result
		return nil
	})
	if err != nil {
		return err
	}

	d.SetId(helper.BuildToken())
	resultMap := make(map[string]interface{})

	if overviewCCTrendResponseParams != nil {
		_ = d.Set("data", overviewCCTrendResponseParams.Data)
		resultMap["data"] = overviewCCTrendResponseParams.Data
	}
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), resultMap); e != nil {
			return e
		}
	}
	return nil
}
