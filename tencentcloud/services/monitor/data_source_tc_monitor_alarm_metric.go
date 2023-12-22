package monitor

import (
	"context"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	monitor "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/monitor/v20180724"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudMonitorAlarmMetric() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudMonitorAlarmMetricRead,
		Schema: map[string]*schema.Schema{
			"module": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Fixed value, as `monitor`.",
			},

			"monitor_type": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Monitoring Type Filter MT_QCE=Cloud Product Monitoring.",
			},

			"namespace": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Alarm policy type, obtained from DescribeAllNamespaces, such as cvm_device.",
			},

			"metrics": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Alarm indicator list.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"namespace": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Alarm strategy type.",
						},
						"metric_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Indicator Name.",
						},
						"description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Indicator display name.",
						},
						"min": {
							Type:        schema.TypeFloat,
							Computed:    true,
							Description: "Minimum value.",
						},
						"max": {
							Type:        schema.TypeFloat,
							Computed:    true,
							Description: "Maximum value.",
						},
						"dimensions": {
							Type: schema.TypeSet,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Computed:    true,
							Description: "Dimension List.",
						},
						"unit": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Unit.",
						},
						"metric_config": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Indicator configuration.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"operator": {
										Type: schema.TypeSet,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
										Computed:    true,
										Description: "Allowed Operators.",
									},
									"period": {
										Type: schema.TypeSet,
										Elem: &schema.Schema{
											Type: schema.TypeInt,
										},
										Computed:    true,
										Description: "The data period allowed for configuration, in seconds.",
									},
									"continue_period": {
										Type: schema.TypeSet,
										Elem: &schema.Schema{
											Type: schema.TypeInt,
										},
										Computed:    true,
										Description: "Number of allowed duration cycles for configuration.",
									},
								},
							},
						},
						"is_advanced": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Is it a high-level indicator. 1 Yes 0 No.",
						},
						"is_open": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Is the advanced indicator activated. 1 Yes 0 No.",
						},
						"product_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Integration Center Product ID.",
						},
						"operators": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Matching operator.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Operator identification.",
									},
									"name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Operator Display Name.",
									},
								},
							},
						},
						"periods": {
							Type: schema.TypeSet,
							Elem: &schema.Schema{
								Type: schema.TypeInt,
							},
							Computed:    true,
							Description: "Indicator trigger.",
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

func dataSourceTencentCloudMonitorAlarmMetricRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_monitor_alarm_metric.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("module"); ok {
		paramMap["Module"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("monitor_type"); ok {
		paramMap["MonitorType"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("namespace"); ok {
		paramMap["Namespace"] = helper.String(v.(string))
	}

	service := MonitorService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	var metrics []*monitor.Metric

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeMonitorAlarmMetricByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}
		metrics = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(metrics))
	tmpList := make([]map[string]interface{}, 0, len(metrics))

	if metrics != nil {
		for _, metric := range metrics {
			metricMap := map[string]interface{}{}

			if metric.Namespace != nil {
				metricMap["namespace"] = metric.Namespace
			}

			if metric.MetricName != nil {
				metricMap["metric_name"] = metric.MetricName
			}

			if metric.Description != nil {
				metricMap["description"] = metric.Description
			}

			if metric.Min != nil {
				metricMap["min"] = metric.Min
			}

			if metric.Max != nil {
				metricMap["max"] = metric.Max
			}

			if metric.Dimensions != nil {
				metricMap["dimensions"] = metric.Dimensions
			}

			if metric.Unit != nil {
				metricMap["unit"] = metric.Unit
			}

			if metric.MetricConfig != nil {
				metricConfigMap := map[string]interface{}{}

				if metric.MetricConfig.Operator != nil {
					metricConfigMap["operator"] = metric.MetricConfig.Operator
				}

				if metric.MetricConfig.Period != nil {
					metricConfigMap["period"] = metric.MetricConfig.Period
				}

				if metric.MetricConfig.ContinuePeriod != nil {
					metricConfigMap["continue_period"] = metric.MetricConfig.ContinuePeriod
				}

				metricMap["metric_config"] = []interface{}{metricConfigMap}
			}

			if metric.IsAdvanced != nil {
				metricMap["is_advanced"] = metric.IsAdvanced
			}

			if metric.IsOpen != nil {
				metricMap["is_open"] = metric.IsOpen
			}

			if metric.ProductId != nil {
				metricMap["product_id"] = metric.ProductId
			}

			if metric.Operators != nil {
				operatorsList := []interface{}{}
				for _, operators := range metric.Operators {
					operatorsMap := map[string]interface{}{}

					if operators.Id != nil {
						operatorsMap["id"] = operators.Id
					}

					if operators.Name != nil {
						operatorsMap["name"] = operators.Name
					}

					operatorsList = append(operatorsList, operatorsMap)
				}

				metricMap["operators"] = operatorsList
			}

			if metric.Periods != nil {
				metricMap["periods"] = metric.Periods
			}

			ids = append(ids, *metric.MetricName)
			tmpList = append(tmpList, metricMap)
		}

		_ = d.Set("metrics", tmpList)
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
