/*
Use this data source to query detailed information of monitor statistic_data

Example Usage

```hcl
data "tencentcloud_monitor_statistic_data" "statistic_data" {
  module       = "monitor"
  namespace    = "QCE/TKE2"
  metric_names = ["cpu_usage"]
  conditions {
    key      = "tke_cluster_instance_id"
    operator = "="
    value    = ["cls-mw2w40s7"]
  }
}
```
*/
package tencentcloud

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	monitor "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/monitor/v20180724"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func dataSourceTencentCloudMonitorStatisticData() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudMonitorStatisticDataRead,
		Schema: map[string]*schema.Schema{
			"module": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Module, whose value is fixed at monitor.",
			},

			"namespace": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Namespace. Valid values: QCE, TKE2.",
			},

			"metric_names": {
				Required: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "Metric name list.",
			},

			"conditions": {
				Optional:    true,
				Type:        schema.TypeList,
				Description: "Dimension condition. The = and in operators are supported.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Dimension.",
						},
						"operator": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Operator. Valid values: eq (equal to), ne (not equal to), in.",
						},
						"value": {
							Type: schema.TypeSet,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Required:    true,
							Description: "Dimension value. If Operator is eq or ne, only the first element will be used.",
						},
					},
				},
			},

			"period": {
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "Statistical period in seconds. Default value: 300. Optional values: 60, 300, 3,600, and 86,400.Due to the storage period limit, the statistical period is subject to the time range of statistics:60s: The time range is less than 12 hours, and the timespan between StartTime and the current time cannot exceed 15 days.300s: The time range is less than three days, and the timespan between StartTime and the current time cannot exceed 31 days.3,600s: The time range is less than 30 days, and the timespan between StartTime and the current time cannot exceed 93 days.86,400s: The time range is less than 186 days, and the timespan between StartTime and the current time cannot exceed 186 days.",
			},

			"start_time": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Start time, which is the current time by default, such as 2020-12-08T19:51:23+08:00.",
			},

			"end_time": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "End time, which is the current time by default, such as 2020-12-08T19:51:23+08:00.",
			},

			"group_bys": {
				Optional: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "GroupBy by the specified dimension.",
			},

			"data": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Monitoring data.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"metric_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Metric name.",
						},
						"points": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Monitoring data point.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"dimensions": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Combination of instance object dimensions.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"name": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Instance dimension name.",
												},
												"value": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Instance dimension value.",
												},
											},
										},
									},
									"values": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Data point list.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"timestamp": {
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "Time point when this monitoring data point is generated.",
												},
												"value": {
													Type:        schema.TypeFloat,
													Computed:    true,
													Description: "Monitoring data point valueNote: this field may return null, indicating that no valid values can be obtained.",
												},
											},
										},
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

func dataSourceTencentCloudMonitorStatisticDataRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_monitor_statistic_data.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("module"); ok {
		paramMap["Module"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("namespace"); ok {
		paramMap["Namespace"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("metric_names"); ok {
		metricNamesSet := v.(*schema.Set).List()
		paramMap["MetricNames"] = helper.InterfacesStringsPoint(metricNamesSet)
	}

	if v, ok := d.GetOk("conditions"); ok {
		conditionsSet := v.([]interface{})
		tmpSet := make([]*monitor.MidQueryCondition, 0, len(conditionsSet))

		for _, item := range conditionsSet {
			midQueryCondition := monitor.MidQueryCondition{}
			midQueryConditionMap := item.(map[string]interface{})

			if v, ok := midQueryConditionMap["key"]; ok {
				midQueryCondition.Key = helper.String(v.(string))
			}
			if v, ok := midQueryConditionMap["operator"]; ok {
				midQueryCondition.Operator = helper.String(v.(string))
			}
			if v, ok := midQueryConditionMap["value"]; ok {
				valueSet := v.(*schema.Set).List()
				midQueryCondition.Value = helper.InterfacesStringsPoint(valueSet)
			}
			tmpSet = append(tmpSet, &midQueryCondition)
		}
		paramMap["conditions"] = tmpSet
	}

	if v, ok := d.GetOk("group_bys"); ok {
		groupBysSet := v.(*schema.Set).List()
		paramMap["GroupBys"] = helper.InterfacesStringsPoint(groupBysSet)
	}

	service := MonitorService{client: meta.(*TencentCloudClient).apiV3Conn}

	var statistic *monitor.DescribeStatisticDataResponseParams
	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeMonitorStatisticDataByFilter(ctx, paramMap)
		if e != nil {
			return retryError(e)
		}
		statistic = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0)
	if statistic.Period != nil {
		_ = d.Set("period", statistic.Period)
	}

	if statistic.StartTime != nil {
		_ = d.Set("start_time", statistic.StartTime)
	}

	if statistic.EndTime != nil {
		_ = d.Set("end_time", statistic.EndTime)
	}

	if statistic.Data != nil {
		tmpList := make([]map[string]interface{}, 0, len(statistic.Data))
		for _, metricData := range statistic.Data {
			metricDataMap := map[string]interface{}{}

			if metricData.MetricName != nil {
				metricDataMap["metric_name"] = metricData.MetricName
			}

			if metricData.Points != nil {
				pointsList := []interface{}{}
				for _, points := range metricData.Points {
					pointsMap := map[string]interface{}{}

					if points.Dimensions != nil {
						dimensionsList := []interface{}{}
						for _, dimensions := range points.Dimensions {
							dimensionsMap := map[string]interface{}{}

							if dimensions.Name != nil {
								dimensionsMap["name"] = dimensions.Name
							}

							if dimensions.Value != nil {
								dimensionsMap["value"] = dimensions.Value
							}

							dimensionsList = append(dimensionsList, dimensionsMap)
						}

						pointsMap["dimensions"] = dimensionsList
					}

					if points.Values != nil {
						valuesList := []interface{}{}
						for _, values := range points.Values {
							valuesMap := map[string]interface{}{}

							if values.Timestamp != nil {
								valuesMap["timestamp"] = values.Timestamp
							}

							if values.Value != nil {
								valuesMap["value"] = values.Value
							}

							valuesList = append(valuesList, valuesMap)
						}

						pointsMap["values"] = valuesList
					}

					pointsList = append(pointsList, pointsMap)
				}

				metricDataMap["points"] = pointsList
			}

			ids = append(ids, *metricData.MetricName)
			tmpList = append(tmpList, metricDataMap)
		}

		_ = d.Set("data", tmpList)
	}

	d.SetId(helper.DataResourceIdsHash(ids))
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := writeToFile(output.(string), d); e != nil {
			return e
		}
	}
	return nil
}
