package teo

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	teo "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/ratelimit"
)

func DataSourceTencentCloudTeoInferenceServiceMonitorData() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudTeoInferenceServiceMonitorDataRead,
		Schema: map[string]*schema.Schema{
			"zone_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Site ID.",
			},

			"service_ids": {
				Type:        schema.TypeList,
				Required:    true,
				MaxItems:    10,
				Description: "Inference service ID list, up to 10 inference service IDs.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},

			"metric_names": {
				Type:        schema.TypeList,
				Required:    true,
				MaxItems:    10,
				Description: "Metric list, up to 10 metrics. Valid values: `cpu_usage_average`, `cpu_usage_max`, `gpu_usage_average`, `gpu_usage_max`, `instance_num_average`, `instance_num_max`, `gpu_memory_usage_max`, `memory_usage_average`, `memory_usage_max`.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},

			"start_time": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Start time.",
			},

			"end_time": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "End time. The query time range (`EndTime` - `StartTime`) must be less than or equal to 30 days.",
			},

			"interval": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Query time granularity. Valid values: `min`, `5min`, `hour`, `day`.",
			},

			"inference_service_monitor_records": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Inference service monitor record list.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"service_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Inference service ID.",
						},
						"metric_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Metric name.",
						},
						"inference_service_monitor_items": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Detailed inference service monitor data.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"timestamp": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Timestamp of the monitor data point.",
									},
									"value": {
										Type:        schema.TypeFloat,
										Computed:    true,
										Description: "Numeric value of the monitor data point.",
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

func dataSourceTencentCloudTeoInferenceServiceMonitorDataRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_teo_inference_service_monitor_data.read")()
	defer tccommon.InconsistentCheck(d, meta)

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
	client := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTeoClient()

	request := teo.NewDescribeInferenceServiceMonitorDataRequest()
	if v, ok := d.GetOk("zone_id"); ok {
		request.ZoneId = helper.String(v.(string))
	}
	if v, ok := d.GetOk("service_ids"); ok {
		serviceIdsSet := v.([]interface{})
		request.ServiceIds = helper.InterfacesStringsPoint(serviceIdsSet)
	}
	if v, ok := d.GetOk("metric_names"); ok {
		metricNamesSet := v.([]interface{})
		request.MetricNames = helper.InterfacesStringsPoint(metricNamesSet)
	}
	if v, ok := d.GetOk("start_time"); ok {
		request.StartTime = helper.String(v.(string))
	}
	if v, ok := d.GetOk("end_time"); ok {
		request.EndTime = helper.String(v.(string))
	}
	if v, ok := d.GetOk("interval"); ok {
		request.Interval = helper.String(v.(string))
	}

	ratelimit.Check(request.GetAction())

	var response *teo.DescribeInferenceServiceMonitorDataResponse
	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		resp, e := client.DescribeInferenceServiceMonitorDataWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		}
		response = resp
		if response == nil || response.Response == nil || len(response.Response.InferenceServiceMonitorRecords) == 0 {
			log.Printf("[DATASOURCE] read empty, skip SetId")
			return resource.NonRetryableError(e)
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s read tencentcloud_teo_inference_service_monitor_data failed, reason:%s\n", logId, err.Error())
		return err
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	recordsList := make([]map[string]interface{}, 0)
	if response.Response.InferenceServiceMonitorRecords != nil {
		for _, record := range response.Response.InferenceServiceMonitorRecords {
			recordMap := map[string]interface{}{}
			if record.ServiceId != nil {
				recordMap["service_id"] = record.ServiceId
			}
			if record.MetricName != nil {
				recordMap["metric_name"] = record.MetricName
			}
			itemsList := make([]map[string]interface{}, 0)
			if record.InferenceServiceMonitorItems != nil {
				for _, item := range record.InferenceServiceMonitorItems {
					itemMap := map[string]interface{}{}
					if item.Timestamp != nil {
						itemMap["timestamp"] = item.Timestamp
					}
					if item.Value != nil {
						itemMap["value"] = *item.Value
					}
					itemsList = append(itemsList, itemMap)
				}
			}
			recordMap["inference_service_monitor_items"] = itemsList
			recordsList = append(recordsList, recordMap)
		}
	}
	_ = d.Set("inference_service_monitor_records", recordsList)

	d.SetId(helper.BuildToken())

	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), recordsList); e != nil {
			return e
		}
	}

	return nil
}
