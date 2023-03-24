/*
Use this data source to query detailed information of dbbrain slow_log_time_series_stats

Example Usage

```hcl
data "tencentcloud_dbbrain_slow_log_time_series_stats" "test" {
  instance_id = "%s"
  start_time = "%s"
  end_time = "%s"
  product = "mysql"
}
```
*/
package tencentcloud

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	dbbrain "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dbbrain/v20210527"
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
				Description: "Start time, such as `2019-09-10 12:13:14`.",
			},

			"end_time": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "End time, such as `2019-09-10 12:13:14`, the interval between the end time and the start time can be up to 7 days.",
			},

			"product": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Service product type, supported values include: `mysql` - cloud database MySQL, `cynosdb` - cloud database CynosDB for MySQL, the default is `mysql`.",
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
							Description: "total.",
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
				Computed: true,
				Type:     schema.TypeList,
				// MaxItems:    1,
				Description: "Instan1ce cpu utilization monitoring data within a unit time interval.",
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
										Type: schema.TypeSet,
										Elem: &schema.Schema{
											Type: schema.TypeInt,
										},
										Computed:    true,
										Description: "Index value. Note: This field may return null, indicating that no valid value can be obtained.",
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

	var id string
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("instance_id"); ok {
		paramMap["instance_id"] = helper.String(v.(string))
		id = v.(string)
	}

	if v, ok := d.GetOk("start_time"); ok {
		paramMap["start_time"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("end_time"); ok {
		paramMap["end_time"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("product"); ok {
		paramMap["product"] = helper.String(v.(string))
	}

	var result *dbbrain.DescribeSlowLogTimeSeriesStatsResponseParams
	service := DbbrainService{client: meta.(*TencentCloudClient).apiV3Conn}

	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		var e error
		result, e = service.DescribeDbbrainSlowLogTimeSeriesStatsByFilter(ctx, paramMap)
		if e != nil {
			return retryError(e)
		}
		return nil
	})
	if err != nil {
		return err
	}

	if result != nil {
		period := result.Period
		if period != nil {
			_ = d.Set("period", period)
		}

		timeSeries := result.TimeSeries
		if timeSeries != nil {
			tmpList := make([]map[string]interface{}, 0, len(timeSeries))
			for _, timeSlice := range timeSeries {
				timeSliceMap := map[string]interface{}{}

				if timeSlice.Count != nil {
					timeSliceMap["count"] = timeSlice.Count
				}

				if timeSlice.Timestamp != nil {
					timeSliceMap["timestamp"] = timeSlice.Timestamp
				}

				tmpList = append(tmpList, timeSliceMap)
			}

			_ = d.Set("time_series", tmpList)
		}

		seriesData := result.SeriesData
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

				monitorMetricSeriesDataMap["series"] = seriesList
			}

			if seriesData.Timestamp != nil {
				monitorMetricSeriesDataMap["timestamp"] = seriesData.Timestamp
			}

			_ = d.Set("series_data", []interface{}{monitorMetricSeriesDataMap})
		}
	}

	d.SetId(id)
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := writeToFile(output.(string), result); e != nil {
			return e
		}
	}
	return nil
}
