/*
Use this data source to query detailed information of gaap proxy statistics

Example Usage

```hcl
data "tencentcloud_gaap_proxy_statistics" "proxy_statistics" {
  proxy_id = "link-8lpyo88p"
  start_time = "2023-10-09 00:00:00"
  end_time = "2023-10-09 23:59:59"
  metric_names = ["InBandwidth", "OutBandwidth", "InFlow", "OutFlow", "InPackets", "OutPackets", "Concurrent", "HttpQPS", "HttpsQPS", "Latency", "PacketLoss"]
  granularity = 300
}
```
*/
package tencentcloud

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	gaap "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/gaap/v20180529"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func dataSourceTencentCloudGaapProxyStatistics() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudGaapProxyStatisticsRead,
		Schema: map[string]*schema.Schema{
			"proxy_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Proxy Id.",
			},

			"start_time": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Start Time(2019-03-25 12:00:00).",
			},

			"end_time": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "End Time(2019-03-25 12:00:00).",
			},

			"metric_names": {
				Required: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "Metric Names. Valid values: InBandwidth,OutBandwidth, Concurrent, InPackets, OutPackets, PacketLoss, Latency, HttpQPS, HttpsQPS.",
			},

			"granularity": {
				Required:    true,
				Type:        schema.TypeInt,
				Description: "Monitoring granularity, currently supporting 60 300 3600 86400, in seconds.When the time range does not exceed 3 days, support a minimum granularity of 60 seconds;When the time range does not exceed 7 days, support a minimum granularity of 300 seconds;When the time range does not exceed 30 days, the minimum granularity supported is 3600 seconds.",
			},

			"isp": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Operator (valid when the proxy is a three network proxy), supports CMCC, CUCC, CTCC, and merges data from the three operators if null values are passed or not passed.",
			},

			"statistics_data": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "proxy Statistics.",
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

func dataSourceTencentCloudGaapProxyStatisticsRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_gaap_proxy_statistics.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	paramMap := make(map[string]interface{})

	var proxyId string
	if v, ok := d.GetOk("proxy_id"); ok {
		proxyId = v.(string)
		paramMap["ProxyId"] = helper.String(proxyId)
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

	if v, ok := d.GetOk("isp"); ok {
		paramMap["Isp"] = helper.String(v.(string))
	}

	service := GaapService{client: meta.(*TencentCloudClient).apiV3Conn}

	var statisticsData []*gaap.MetricStatisticsInfo

	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeGaapProxyStatisticsByFilter(ctx, paramMap)
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

	d.SetId(proxyId)
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := writeToFile(output.(string), tmpList); e != nil {
			return e
		}
	}
	return nil
}
