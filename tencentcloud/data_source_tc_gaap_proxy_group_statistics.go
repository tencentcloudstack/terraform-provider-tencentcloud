package tencentcloud

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	gaap "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/gaap/v20180529"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func dataSourceTencentCloudGaapProxyGroupStatistics() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudGaapProxyGroupStatisticsRead,
		Schema: map[string]*schema.Schema{
			"group_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Group Id.",
			},

			"start_time": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Start Time.",
			},

			"end_time": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "End Time.",
			},

			"metric_names": {
				Required: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "Metric Names. support, InBandwidth, OutBandwidth, Concurrent, InPackets, OutPackets.",
			},

			"granularity": {
				Required:    true,
				Type:        schema.TypeInt,
				Description: "Monitoring granularity, currently supporting 60 300 3600 86400, in seconds.When the time range does not exceed 1 day, support a minimum granularity of 60 seconds;When the time range does not exceed 7 days, support a minimum granularity of 3600 seconds;When the time range does not exceed 30 days, the minimum granularity supported is 86400 seconds.",
			},

			"statistics_data": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "proxy Group Statistics.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"metric_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Metric Name.",
						},
						"metric_data": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Metric Data.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"time": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Time.",
									},
									"data": {
										Type:        schema.TypeFloat,
										Computed:    true,
										Description: "DataNote: This field may return null, indicating that a valid value cannot be obtained.",
									},
								},
							},
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

func dataSourceTencentCloudGaapProxyGroupStatisticsRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_gaap_proxy_group_statistics.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	paramMap := make(map[string]interface{})
	var groupId string
	if v, ok := d.GetOk("group_id"); ok {
		groupId = v.(string)
		paramMap["GroupId"] = helper.String(groupId)
	}

	if v, ok := d.GetOk("start_time"); ok {
		paramMap["StartTime"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("end_time"); ok {
		paramMap["EndTime"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("metric_names"); ok {
		metricNamesSet := v.(*schema.Set).List()
		paramMap["MetricNames"] = helper.InterfacesStringsPoint(metricNamesSet)
	}

	if v, _ := d.GetOk("granularity"); v != nil {
		paramMap["Granularity"] = helper.IntUint64(v.(int))
	}

	service := GaapService{client: meta.(*TencentCloudClient).apiV3Conn}

	var statisticsData []*gaap.MetricStatisticsInfo

	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeGaapProxyGroupStatisticsByFilter(ctx, paramMap)
		if e != nil {
			return retryError(e)
		}
		statisticsData = result
		return nil
	})
	if err != nil {
		return err
	}

	tmpList := make([]map[string]interface{}, 0, len(statisticsData))

	if statisticsData != nil {
		for _, metricStatisticsInfo := range statisticsData {
			metricStatisticsInfoMap := map[string]interface{}{}

			if metricStatisticsInfo.MetricName != nil {
				metricStatisticsInfoMap["metric_name"] = metricStatisticsInfo.MetricName
			}

			if metricStatisticsInfo.MetricData != nil {
				metricDataList := []interface{}{}
				for _, metricData := range metricStatisticsInfo.MetricData {
					metricDataMap := map[string]interface{}{}

					if metricData.Time != nil {
						metricDataMap["time"] = metricData.Time
					}

					if metricData.Data != nil {
						metricDataMap["data"] = metricData.Data
					}

					metricDataList = append(metricDataList, metricDataMap)
				}

				metricStatisticsInfoMap["metric_data"] = metricDataList
			}

			tmpList = append(tmpList, metricStatisticsInfoMap)
		}

		_ = d.Set("statistics_data", tmpList)
	}

	d.SetId(groupId)
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := writeToFile(output.(string), tmpList); e != nil {
			return e
		}
	}
	return nil
}
