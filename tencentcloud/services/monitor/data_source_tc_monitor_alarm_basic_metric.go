package monitor

import (
	"context"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	monitor "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/monitor/v20180724"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudMonitorAlarmBasicMetric() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudMonitorAlarmBasicMetricRead,
		Schema: map[string]*schema.Schema{
			"namespace": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "The business namespace is different for each cloud product. To obtain the business namespace, please go to the product monitoring indicator documents, such as the namespace of the cloud server, which can be found in [Cloud Server Monitoring Indicators](https://cloud.tencent.com/document/product/248/6843 ).",
			},

			"metric_name": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Indicator names are different for each cloud product. To obtain indicator names, please go to the monitoring indicator documents of each product, such as the indicator names of cloud servers, which can be found in [Cloud Server Monitoring Indicators]( https://cloud.tencent.com/document/product/248/6843).",
			},

			"dimensions": {
				Optional: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "Optional parameters, filtered by dimension.",
			},

			"metric_set": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "List of indicator descriptions obtained from query.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"namespace": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Namespaces, each cloud product will have a namespace.",
						},
						"metric_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Indicator Name.",
						},
						"unit": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Units used for indicators.",
						},
						"unit_cname": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Units used for indicators.",
						},
						"period": {
							Type: schema.TypeSet,
							Elem: &schema.Schema{
								Type: schema.TypeInt,
							},
							Computed:    true,
							Description: "The statistical period supported by the indicator, in seconds, such as 60, 300.",
						},
						"periods": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Indicator method within the statistical cycle.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"period": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Cycle.",
									},
									"stat_type": {
										Type: schema.TypeSet,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
										Computed:    true,
										Description: "Statistical methods.",
									},
								},
							},
						},
						"meaning": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Explanation of the meaning of statistical indicators.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"en": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Explanation of indicators in English.",
									},
									"zh": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Chinese interpretation of indicators.",
									},
								},
							},
						},
						"dimensions": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Dimension description information.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"dimensions": {
										Type: schema.TypeSet,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
										Computed:    true,
										Description: "Dimension name array.",
									},
								},
							},
						},
						"metric_c_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Indicator Chinese Name.",
						},
						"metric_e_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Indicator English name.",
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

func dataSourceTencentCloudMonitorAlarmBasicMetricRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_monitor_alarm_metric.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("namespace"); ok {
		paramMap["Namespace"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("metric_name"); ok {
		paramMap["MetricName"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("dimensions"); ok {
		dimensionsSet := v.(*schema.Set).List()
		paramMap["Dimensions"] = helper.InterfacesStringsPoint(dimensionsSet)
	}

	service := MonitorService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	var metricSet []*monitor.MetricSet

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeMonitorAlarmBasicMetricByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}
		metricSet = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(metricSet))
	tmpList := make([]map[string]interface{}, 0, len(metricSet))

	if metricSet != nil {
		for _, metricSet := range metricSet {
			metricSetMap := map[string]interface{}{}

			if metricSet.Namespace != nil {
				metricSetMap["namespace"] = metricSet.Namespace
			}

			if metricSet.MetricName != nil {
				metricSetMap["metric_name"] = metricSet.MetricName
			}

			if metricSet.Unit != nil {
				metricSetMap["unit"] = metricSet.Unit
			}

			if metricSet.UnitCname != nil {
				metricSetMap["unit_cname"] = metricSet.UnitCname
			}

			if metricSet.Period != nil {
				metricSetMap["period"] = metricSet.Period
			}

			if metricSet.Periods != nil {
				periodsList := []interface{}{}
				for _, periods := range metricSet.Periods {
					periodsMap := map[string]interface{}{}

					if periods.Period != nil {
						periodsMap["period"] = periods.Period
					}

					if periods.StatType != nil {
						periodsMap["stat_type"] = periods.StatType
					}

					periodsList = append(periodsList, periodsMap)
				}

				metricSetMap["periods"] = periodsList
			}

			if metricSet.Meaning != nil {
				meaningMap := map[string]interface{}{}

				if metricSet.Meaning.En != nil {
					meaningMap["en"] = metricSet.Meaning.En
				}

				if metricSet.Meaning.Zh != nil {
					meaningMap["zh"] = metricSet.Meaning.Zh
				}

				metricSetMap["meaning"] = []interface{}{meaningMap}
			}

			if metricSet.Dimensions != nil {
				dimensionsList := []interface{}{}
				for _, dimensions := range metricSet.Dimensions {
					dimensionsMap := map[string]interface{}{}

					if dimensions.Dimensions != nil {
						dimensionsMap["dimensions"] = dimensions.Dimensions
					}

					dimensionsList = append(dimensionsList, dimensionsMap)
				}

				metricSetMap["dimensions"] = dimensionsList
			}

			if metricSet.MetricCName != nil {
				metricSetMap["metric_c_name"] = metricSet.MetricCName
			}

			if metricSet.MetricEName != nil {
				metricSetMap["metric_e_name"] = metricSet.MetricEName
			}

			ids = append(ids, *metricSet.MetricName)
			tmpList = append(tmpList, metricSetMap)
		}

		_ = d.Set("metric_set", tmpList)
	}

	d.SetId(helper.DataResourceIdsHash(ids))
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), tmpList); e != nil {
			return e
		}
	}
	return nil
}
