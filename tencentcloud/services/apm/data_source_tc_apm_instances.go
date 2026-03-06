package apm

import (
	"context"
	"log"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	apm "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/apm/v20210622"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudApmInstances() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudApmInstancesRead,

		Schema: map[string]*schema.Schema{
			"instance_ids": {
				Type:        schema.TypeList,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Filter by instance ID list (exact match).",
			},
			"instance_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Filter by instance ID (fuzzy match).",
			},
			"instance_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Filter by instance name (fuzzy match).",
			},
			"tags": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "Filter by tags.",
			},
			"demo_instance_flag": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Whether to query official demo instances. 0: non-demo, 1: demo. Default is 0.",
			},
			"all_regions_flag": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Whether to query instances in all regions. 0: no, 1: yes. Default is 0.",
			},
			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results.",
			},
			"instance_list": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "APM instance list.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"instance_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Instance ID.",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Instance name.",
						},
						"description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Instance description.",
						},
						"region": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Region.",
						},
						"app_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "App ID.",
						},
						"status": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Instance status.",
						},
						"create_uin": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Creator UIN.",
						},
						"trace_duration": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Trace data retention duration.",
						},
						"span_daily_counters": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Daily span count quota.",
						},
						"pay_mode": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Billing mode.",
						},
						"free": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Whether it is free edition.",
						},
						"tags": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Tag list.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"key": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Tag key.",
									},
									"value": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Tag value.",
									},
								},
							},
						},
						"err_rate_threshold": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Error rate threshold.",
						},
						"sample_rate": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Sampling rate.",
						},
						"error_sample": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Error sampling switch.",
						},
						"service_count": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Service count.",
						},
						"amount_of_used_storage": {
							Type:        schema.TypeFloat,
							Computed:    true,
							Description: "Storage usage in MB.",
						},
						"count_of_report_span_per_day": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Daily reported span count.",
						},
						"billing_instance": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Whether billing is enabled. 0: not enabled, 1: enabled.",
						},
						"slow_request_saved_threshold": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Slow request saved threshold in ms.",
						},
						"log_region": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "CLS log region.",
						},
						"log_source": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Log source.",
						},
						"is_related_log": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Log feature switch. 0: off, 1: on.",
						},
						"log_topic_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Log topic ID.",
						},
						"client_count": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Client application count.",
						},
						"total_count": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Active application count in recent 2 days.",
						},
						"log_set": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "CLS log set.",
						},
						"metric_duration": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Metric data retention duration in days.",
						},
						"custom_show_tags": {
							Type:        schema.TypeList,
							Computed:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: "Custom display tag list.",
						},
						"pay_mode_effective": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Whether pay mode is effective.",
						},
						"response_duration_warning_threshold": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Response duration warning threshold in ms.",
						},
						"default_tsf": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Whether it is the default TSF instance. 0: no, 1: yes.",
						},
						"is_related_dashboard": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Whether dashboard is associated. 0: off, 1: on.",
						},
						"dashboard_topic_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Associated dashboard ID.",
						},
						"is_instrumentation_vulnerability_scan": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Whether instrumentation vulnerability scan is enabled. 0: off, 1: on.",
						},
						"is_sql_injection_analysis": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Whether SQL injection analysis is enabled. 0: off, 1: on.",
						},
						"stop_reason": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Throttling reason. 1: official version quota, 2: trial version quota, 4: trial expired, 8: account overdue.",
						},
						"is_remote_command_execution_analysis": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Whether remote command execution detection is enabled. 0: off, 1: on.",
						},
						"is_memory_hijacking_analysis": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Whether memory hijacking detection is enabled. 0: off, 1: on.",
						},
						"log_index_type": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "CLS index type. 0: full-text index, 1: key-value index.",
						},
						"log_trace_id_key": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "TraceId index key, effective when CLS index type is key-value.",
						},
						"is_delete_any_file_analysis": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Whether delete any file detection is enabled. 0: off, 1: on.",
						},
						"is_read_any_file_analysis": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Whether read any file detection is enabled. 0: off, 1: on.",
						},
						"is_upload_any_file_analysis": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Whether upload any file detection is enabled. 0: off, 1: on.",
						},
						"is_include_any_file_analysis": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Whether include any file detection is enabled. 0: off, 1: on.",
						},
						"is_directory_traversal_analysis": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Whether directory traversal detection is enabled. 0: off, 1: on.",
						},
						"is_template_engine_injection_analysis": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Whether template engine injection detection is enabled. 0: off, 1: on.",
						},
						"is_script_engine_injection_analysis": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Whether script engine injection detection is enabled. 0: off, 1: on.",
						},
						"is_expression_injection_analysis": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Whether expression injection detection is enabled. 0: off, 1: on.",
						},
						"is_jndi_injection_analysis": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Whether JNDI injection detection is enabled. 0: off, 1: on.",
						},
						"is_jni_injection_analysis": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Whether JNI injection detection is enabled. 0: off, 1: on.",
						},
						"is_webshell_backdoor_analysis": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Whether webshell backdoor detection is enabled. 0: off, 1: on.",
						},
						"is_deserialization_analysis": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Whether deserialization detection is enabled. 0: off, 1: on.",
						},
						"token": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Instance authentication token.",
						},
						"url_long_segment_threshold": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "URL long segment convergence threshold.",
						},
						"url_number_segment_threshold": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "URL number segment convergence threshold.",
						},
					},
				},
			},
		},
	}
}

func dataSourceTencentCloudApmInstancesRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_apm_instances.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
	service := ApmService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	params := make(map[string]interface{})

	if v, ok := d.GetOk("instance_ids"); ok {
		params["instance_ids"] = v.([]interface{})
	}
	if v, ok := d.GetOk("instance_id"); ok {
		params["instance_id"] = v.(string)
	}
	if v, ok := d.GetOk("instance_name"); ok {
		params["instance_name"] = v.(string)
	}
	if v, ok := d.GetOk("tags"); ok {
		params["tags"] = v.(map[string]interface{})
	}
	if v, ok := d.GetOkExists("demo_instance_flag"); ok {
		params["demo_instance_flag"] = v.(int)
	}
	if v, ok := d.GetOkExists("all_regions_flag"); ok {
		params["all_regions_flag"] = v.(int)
	}

	var instances []*apm.ApmInstanceDetail
	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		results, e := service.DescribeApmInstances(ctx, params)
		if e != nil {
			return tccommon.RetryError(e)
		}
		instances = results
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s read APM instances failed, reason:%+v", logId, err)
		return err
	}

	instanceList := make([]map[string]interface{}, 0, len(instances))
	ids := make([]string, 0, len(instances))

	for _, instance := range instances {
		instanceMap := map[string]interface{}{}

		if instance.InstanceId != nil {
			instanceMap["instance_id"] = instance.InstanceId
			ids = append(ids, *instance.InstanceId)
		}
		if instance.Name != nil {
			instanceMap["name"] = instance.Name
		}
		if instance.Description != nil {
			instanceMap["description"] = instance.Description
		}
		if instance.Region != nil {
			instanceMap["region"] = instance.Region
		}
		if instance.AppId != nil {
			instanceMap["app_id"] = instance.AppId
		}
		if instance.Status != nil {
			instanceMap["status"] = instance.Status
		}
		if instance.CreateUin != nil {
			instanceMap["create_uin"] = instance.CreateUin
		}
		if instance.TraceDuration != nil {
			instanceMap["trace_duration"] = instance.TraceDuration
		}
		if instance.SpanDailyCounters != nil {
			instanceMap["span_daily_counters"] = instance.SpanDailyCounters
		}
		if instance.PayMode != nil {
			instanceMap["pay_mode"] = instance.PayMode
		}
		if instance.Free != nil {
			instanceMap["free"] = instance.Free
		}
		if instance.Tags != nil {
			tags := make([]map[string]interface{}, 0, len(instance.Tags))
			for _, tag := range instance.Tags {
				tagMap := map[string]interface{}{}
				if tag.Key != nil {
					tagMap["key"] = tag.Key
				}
				if tag.Value != nil {
					tagMap["value"] = tag.Value
				}
				tags = append(tags, tagMap)
			}
			instanceMap["tags"] = tags
		}
		if instance.ErrRateThreshold != nil {
			instanceMap["err_rate_threshold"] = instance.ErrRateThreshold
		}
		if instance.SampleRate != nil {
			instanceMap["sample_rate"] = instance.SampleRate
		}
		if instance.ErrorSample != nil {
			instanceMap["error_sample"] = instance.ErrorSample
		}
		if instance.ServiceCount != nil {
			instanceMap["service_count"] = instance.ServiceCount
		}
		if instance.AmountOfUsedStorage != nil {
			instanceMap["amount_of_used_storage"] = instance.AmountOfUsedStorage
		}
		if instance.CountOfReportSpanPerDay != nil {
			instanceMap["count_of_report_span_per_day"] = instance.CountOfReportSpanPerDay
		}
		if instance.BillingInstance != nil {
			instanceMap["billing_instance"] = instance.BillingInstance
		}
		if instance.SlowRequestSavedThreshold != nil {
			instanceMap["slow_request_saved_threshold"] = instance.SlowRequestSavedThreshold
		}
		if instance.LogRegion != nil {
			instanceMap["log_region"] = instance.LogRegion
		}
		if instance.LogSource != nil {
			instanceMap["log_source"] = instance.LogSource
		}
		if instance.IsRelatedLog != nil {
			instanceMap["is_related_log"] = instance.IsRelatedLog
		}
		if instance.LogTopicID != nil {
			instanceMap["log_topic_id"] = instance.LogTopicID
		}
		if instance.ClientCount != nil {
			instanceMap["client_count"] = instance.ClientCount
		}
		if instance.TotalCount != nil {
			instanceMap["total_count"] = instance.TotalCount
		}
		if instance.LogSet != nil {
			instanceMap["log_set"] = instance.LogSet
		}
		if instance.MetricDuration != nil {
			instanceMap["metric_duration"] = instance.MetricDuration
		}
		if instance.CustomShowTags != nil {
			instanceMap["custom_show_tags"] = helper.StringsInterfaces(instance.CustomShowTags)
		}
		if instance.PayModeEffective != nil {
			instanceMap["pay_mode_effective"] = instance.PayModeEffective
		}
		if instance.ResponseDurationWarningThreshold != nil {
			instanceMap["response_duration_warning_threshold"] = instance.ResponseDurationWarningThreshold
		}
		if instance.DefaultTSF != nil {
			instanceMap["default_tsf"] = instance.DefaultTSF
		}
		if instance.IsRelatedDashboard != nil {
			instanceMap["is_related_dashboard"] = instance.IsRelatedDashboard
		}
		if instance.DashboardTopicID != nil {
			instanceMap["dashboard_topic_id"] = instance.DashboardTopicID
		}
		if instance.IsInstrumentationVulnerabilityScan != nil {
			instanceMap["is_instrumentation_vulnerability_scan"] = instance.IsInstrumentationVulnerabilityScan
		}
		if instance.IsSqlInjectionAnalysis != nil {
			instanceMap["is_sql_injection_analysis"] = instance.IsSqlInjectionAnalysis
		}
		if instance.StopReason != nil {
			instanceMap["stop_reason"] = instance.StopReason
		}
		if instance.IsRemoteCommandExecutionAnalysis != nil {
			instanceMap["is_remote_command_execution_analysis"] = instance.IsRemoteCommandExecutionAnalysis
		}
		if instance.IsMemoryHijackingAnalysis != nil {
			instanceMap["is_memory_hijacking_analysis"] = instance.IsMemoryHijackingAnalysis
		}
		if instance.LogIndexType != nil {
			instanceMap["log_index_type"] = instance.LogIndexType
		}
		if instance.LogTraceIdKey != nil {
			instanceMap["log_trace_id_key"] = instance.LogTraceIdKey
		}
		if instance.IsDeleteAnyFileAnalysis != nil {
			instanceMap["is_delete_any_file_analysis"] = instance.IsDeleteAnyFileAnalysis
		}
		if instance.IsReadAnyFileAnalysis != nil {
			instanceMap["is_read_any_file_analysis"] = instance.IsReadAnyFileAnalysis
		}
		if instance.IsUploadAnyFileAnalysis != nil {
			instanceMap["is_upload_any_file_analysis"] = instance.IsUploadAnyFileAnalysis
		}
		if instance.IsIncludeAnyFileAnalysis != nil {
			instanceMap["is_include_any_file_analysis"] = instance.IsIncludeAnyFileAnalysis
		}
		if instance.IsDirectoryTraversalAnalysis != nil {
			instanceMap["is_directory_traversal_analysis"] = instance.IsDirectoryTraversalAnalysis
		}
		if instance.IsTemplateEngineInjectionAnalysis != nil {
			instanceMap["is_template_engine_injection_analysis"] = instance.IsTemplateEngineInjectionAnalysis
		}
		if instance.IsScriptEngineInjectionAnalysis != nil {
			instanceMap["is_script_engine_injection_analysis"] = instance.IsScriptEngineInjectionAnalysis
		}
		if instance.IsExpressionInjectionAnalysis != nil {
			instanceMap["is_expression_injection_analysis"] = instance.IsExpressionInjectionAnalysis
		}
		if instance.IsJNDIInjectionAnalysis != nil {
			instanceMap["is_jndi_injection_analysis"] = instance.IsJNDIInjectionAnalysis
		}
		if instance.IsJNIInjectionAnalysis != nil {
			instanceMap["is_jni_injection_analysis"] = instance.IsJNIInjectionAnalysis
		}
		if instance.IsWebshellBackdoorAnalysis != nil {
			instanceMap["is_webshell_backdoor_analysis"] = instance.IsWebshellBackdoorAnalysis
		}
		if instance.IsDeserializationAnalysis != nil {
			instanceMap["is_deserialization_analysis"] = instance.IsDeserializationAnalysis
		}
		if instance.Token != nil {
			instanceMap["token"] = instance.Token
		}
		if instance.UrlLongSegmentThreshold != nil {
			instanceMap["url_long_segment_threshold"] = instance.UrlLongSegmentThreshold
		}
		if instance.UrlNumberSegmentThreshold != nil {
			instanceMap["url_number_segment_threshold"] = instance.UrlNumberSegmentThreshold
		}

		instanceList = append(instanceList, instanceMap)
	}

	_ = d.Set("instance_list", instanceList)
	d.SetId(helper.DataResourceIdsHash(ids))

	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), instanceList); e != nil {
			return e
		}
	}

	return nil
}
