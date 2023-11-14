/*
Use this data source to query detailed information of dbbrain slow_log_time_series_stats

Example Usage

```hcl
data "tencentcloud_dbbrain_slow_log_time_series_stats" "slow_log_time_series_stats" {
  instance_id = ""
  start_time = ""
  end_time = ""
  product = ""
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

func dataSourceTencentCloudDbbrainSlowLogTimeSeriesStats() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudDbbrainSlowLogTimeSeriesStatsRead,
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Instance ID.",
			},

			"start_time": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Start time, such as 2019-09-10 12:13:14.",
			},

			"end_time": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "End time, such as 2019-09-10 12:13:14, the interval between the end time and the start time can be up to 7 days.",
			},

			"product": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Service product type, supported values include： mysql - cloud database MySQL, cynosdb - cloud database CynosDB for MySQL, the default is mysql.",
			},

			"period": {
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "The unit time interval between bars, in seconds.",
			},

			"time_series": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Statistics on the number of slow logs in a unit time interval.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"count": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Total.",
						},
						"timestamp": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Statistics start time.",
						},
					},
				},
			},

			"series_data": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Instance cpu utilization monitoring data within a unit time interval.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"series": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Monitor metrics.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"metric": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Indicator name.",
									},
									"unit": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Indicator unit.",
									},
									"values": {
										Computed:    true,
										Description: "Index value. Note： This field may return null, indicating that no valid value can be obtained.",
									},
								},
							},
						},
						"timestamp": {
							Type: schema.TypeSet,
							Elem: &schema.Schema{
								Type: schema.TypeInt,
							},
							Computed:    true,
							Description: "The timestamp corresponding to the monitoring indicator.",
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

func dataSourceTencentCloudDbbrainSlowLogTimeSeriesStatsRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_dbbrain_slow_log_time_series_stats.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("instance_id"); ok {
		paramMap["InstanceId"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("start_time"); ok {
		paramMap["StartTime"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("end_time"); ok {
		paramMap["EndTime"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("product"); ok {
		paramMap["Product"] = helper.String(v.(string))
	}

	service := DbbrainService{client: meta.(*TencentCloudClient).apiV3Conn}

	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeDbbrainSlowLogTimeSeriesStatsByFilter(ctx, paramMap)
		if e != nil {
			return retryError(e)
		}
		period = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(period))
	if period != nil {
		_ = d.Set("period", period)
	}

	if timeSeries != nil {
		for _, timeSlice := range timeSeries {
			timeSliceMap := map[string]interface{}{}

			if timeSlice.Count != nil {
				timeSliceMap["count"] = timeSlice.Count
			}

			if timeSlice.Timestamp != nil {
				timeSliceMap["timestamp"] = timeSlice.Timestamp
			}

			ids = append(ids, *timeSlice.InstanceId)
			tmpList = append(tmpList, timeSliceMap)
		}

		_ = d.Set("time_series", tmpList)
	}

	if seriesData != nil {
		monitorMetricSeriesDataMap := map[string]interface{}{}

		if seriesData.Series != nil {
			seriesList := []interface{}{}
			for _, series := range seriesData.Series {
				seriesMap := map[string]interface{}{}

				if series.Metric != nil {
					seriesMap["metric"] = series.Metric
				}

				if series.Unit != nil {
					seriesMap["unit"] = series.Unit
				}

				if series.Values != nil {
					seriesMap["values"] = series.Values
				}

				seriesList = append(seriesList, seriesMap)
			}

			monitorMetricSeriesDataMap["series"] = []interface{}{seriesList}
		}

		if seriesData.Timestamp != nil {
			monitorMetricSeriesDataMap["timestamp"] = seriesData.Timestamp
		}

		ids = append(ids, *seriesData.InstanceId)
		_ = d.Set("series_data", monitorMetricSeriesDataMap)
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
