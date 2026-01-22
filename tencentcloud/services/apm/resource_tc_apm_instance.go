package apm

import (
	"context"
	"fmt"
	"log"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	svctag "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/tag"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	apm "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/apm/v20210622"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudApmInstance() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudApmInstanceCreate,
		Read:   resourceTencentCloudApmInstanceRead,
		Update: resourceTencentCloudApmInstanceUpdate,
		Delete: resourceTencentCloudApmInstanceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name Of Instance.",
			},

			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Description Of Instance.",
			},

			"trace_duration": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "Duration Of Trace Data.",
			},

			"span_daily_counters": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "Quota Of Instance Reporting.",
			},

			"pay_mode": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "Modify the billing mode: `1` means prepaid, `0` means pay-as-you-go, the default value is `0`.",
			},

			"free": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "Whether it is free (0 = paid edition; 1 = tsf restricted free edition; 2 = free edition), default 0.",
			},

			"tags": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "Tag description list.",
			},

			"open_billing": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: "Billing switch.",
			},

			"err_rate_threshold": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "Error rate warning line. when the average error rate of the application exceeds this threshold, the system will give an abnormal note.",
			},

			"sample_rate": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "Sampling rate (unit: %).",
			},

			"error_sample": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "Error sampling switch (0: off, 1: on).",
			},

			"slow_request_saved_threshold": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "Sampling slow call saving threshold (unit: ms).",
			},

			"is_related_log": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "Log feature switch (0: off; 1: on).",
			},

			"log_region": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Log region, which takes effect after the log feature is enabled.",
			},

			"log_topic_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "CLS log topic id, which takes effect after the log feature is enabled.",
			},

			"log_set": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Logset, which takes effect only after the log feature is enabled.",
			},

			"log_source": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Log source, which takes effect only after the log feature is enabled.",
			},

			"custom_show_tags": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "List of custom display tags.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},

			"response_duration_warning_threshold": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "Response time warning line.",
			},

			"is_related_dashboard": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "Whether to associate the dashboard (0 = off, 1 = on).",
			},

			"dashboard_topic_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Associated dashboard id, which takes effect after the associated dashboard is enabled.",
			},

			"is_sql_injection_analysis": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "SQL injection detection switch (0: off, 1: on).",
			},

			"is_instrumentation_vulnerability_scan": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "Whether to enable component vulnerability detection (0 = no, 1 = yes).",
			},

			"is_remote_command_execution_analysis": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "Whether to enable detection of the remote command attack.",
			},

			"is_memory_hijacking_analysis": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "Whether to enable detection of Java webshell.",
			},

			"log_index_type": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "CLS index type. (0 = full-text index; 1 = key-value index).",
			},

			"log_trace_id_key": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Index key of traceId. It is valid when the CLS index type is key-value index.",
			},

			"is_delete_any_file_analysis": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "Whether to enable the detection of deleting arbitrary files. (0 - disabled; 1: enabled).",
			},

			"is_read_any_file_analysis": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "Whether to enable the detection of reading arbitrary files. (0 - disabled; 1 - enabled).",
			},

			"is_upload_any_file_analysis": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "Whether to enable the detection of uploading arbitrary files. (0 - disabled; 1 - enabled).",
			},

			"is_include_any_file_analysis": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "Whether to enable the detection of the inclusion of arbitrary files. (0: disabled, 1: enabled).",
			},

			"is_directory_traversal_analysis": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "Whether to enable traversal detection of the directory. (0 - disabled; 1 - enabled).",
			},

			"is_template_engine_injection_analysis": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "Whether to enable template engine injection detection. (0: disabled; 1: enabled).",
			},

			"is_script_engine_injection_analysis": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "Whether to enable script engine injection detection. (0 - disabled; 1 - enabled).",
			},

			"is_expression_injection_analysis": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "Whether to enable expression injection detection. (0 - disabled; 1 - enabled).",
			},

			"is_jndi_injection_analysis": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "Whether to enable JNDI injection detection. (0 - disabled; 1 - enabled).",
			},

			"is_jni_injection_analysis": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "Whether to enable JNI injection detection. (0 - disabled, 1 - enabled).",
			},

			"is_webshell_backdoor_analysis": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "Whether to enable Webshell backdoor detection. (0 - disabled; 1 - enabled).",
			},

			"is_deserialization_analysis": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "Whether to enable deserialization detection. (0 - disabled; 1 - enabled).",
			},

			"url_long_segment_threshold": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "Convergence threshold for URL long segments.",
			},

			"url_number_segment_threshold": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "Convergence threshold for URL numerical segments.",
			},

			// computed
			"instance_id": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "APM instance ID.",
			},

			"token": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Business system authentication token.",
			},

			"public_collector_url": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "External Network Reporting Address.",
			},
		},
	}
}

func resourceTencentCloudApmInstanceCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_apm_instance.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId      = tccommon.GetLogId(tccommon.ContextNil)
		request    = apm.NewCreateApmInstanceRequest()
		response   = apm.NewCreateApmInstanceResponse()
		instanceId string
	)

	if v, ok := d.GetOk("name"); ok {
		request.Name = helper.String(v.(string))
	}

	if v, ok := d.GetOk("description"); ok {
		request.Description = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("trace_duration"); ok {
		request.TraceDuration = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOkExists("span_daily_counters"); ok {
		request.SpanDailyCounters = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOkExists("pay_mode"); ok {
		request.PayMode = helper.IntInt64(v.(int))
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseApmClient().CreateApmInstance(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("Create apm instance failed, Response is nil."))
		}

		response = result
		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s create apm instance failed, reason:%+v", logId, err)
		return err
	}

	if response.Response.InstanceId == nil {
		return fmt.Errorf("InstanceId is nil.")
	}

	instanceId = *response.Response.InstanceId
	d.SetId(instanceId)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
	if tags := helper.GetTags(d, "tags"); len(tags) > 0 {
		tagService := svctag.NewTagService(meta.(tccommon.ProviderMeta).GetAPIV3Conn())
		region := meta.(tccommon.ProviderMeta).GetAPIV3Conn().Region
		resourceName := fmt.Sprintf("qcs::apm:%s:uin/:apm-instance/%s", region, d.Id())
		if err := tagService.ModifyTags(ctx, resourceName, tags, nil); err != nil {
			return err
		}
	}

	// set config
	configRequest := apm.NewModifyApmInstanceRequest()
	if v, ok := d.GetOk("name"); ok {
		configRequest.Name = helper.String(v.(string))
	}

	if v, ok := d.GetOk("description"); ok {
		configRequest.Description = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("trace_duration"); ok {
		configRequest.TraceDuration = helper.IntInt64(v.(int))
	}

	if tags := helper.GetTags(d, "tags"); len(tags) > 0 {
		for k, v := range tags {
			key := k
			value := v
			tag := &apm.ApmTag{
				Key:   &key,
				Value: &value,
			}

			configRequest.Tags = append(configRequest.Tags, tag)
		}
	}

	if v, ok := d.GetOkExists("open_billing"); ok {
		configRequest.OpenBilling = helper.Bool(v.(bool))
	}

	if v, ok := d.GetOkExists("span_daily_counters"); ok {
		configRequest.SpanDailyCounters = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOkExists("err_rate_threshold"); ok {
		configRequest.ErrRateThreshold = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOkExists("sample_rate"); ok {
		configRequest.SampleRate = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOkExists("error_sample"); ok {
		configRequest.ErrorSample = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOkExists("slow_request_saved_threshold"); ok {
		configRequest.SlowRequestSavedThreshold = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOkExists("is_related_log"); ok {
		configRequest.IsRelatedLog = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("log_region"); ok {
		configRequest.LogRegion = helper.String(v.(string))
	}

	if v, ok := d.GetOk("log_topic_id"); ok {
		configRequest.LogTopicID = helper.String(v.(string))
	}

	if v, ok := d.GetOk("log_set"); ok {
		configRequest.LogSet = helper.String(v.(string))
	}

	if v, ok := d.GetOk("log_source"); ok {
		configRequest.LogSource = helper.String(v.(string))
	}

	if v, ok := d.GetOk("custom_show_tags"); ok {
		customShowTagsSet := v.(*schema.Set).List()
		for i := range customShowTagsSet {
			customShowTags := customShowTagsSet[i].(string)
			configRequest.CustomShowTags = append(configRequest.CustomShowTags, helper.String(customShowTags))
		}
	}

	if v, ok := d.GetOkExists("pay_mode"); ok {
		configRequest.PayMode = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOkExists("response_duration_warning_threshold"); ok {
		configRequest.ResponseDurationWarningThreshold = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOkExists("free"); ok {
		configRequest.Free = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOkExists("is_related_dashboard"); ok {
		configRequest.IsRelatedDashboard = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("dashboard_topic_id"); ok {
		configRequest.DashboardTopicID = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("is_sql_injection_analysis"); ok {
		configRequest.IsSqlInjectionAnalysis = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOkExists("is_instrumentation_vulnerability_scan"); ok {
		configRequest.IsInstrumentationVulnerabilityScan = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOkExists("is_remote_command_execution_analysis"); ok {
		configRequest.IsRemoteCommandExecutionAnalysis = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOkExists("is_memory_hijacking_analysis"); ok {
		configRequest.IsMemoryHijackingAnalysis = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOkExists("log_index_type"); ok {
		configRequest.LogIndexType = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("log_trace_id_key"); ok {
		configRequest.LogTraceIdKey = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("is_delete_any_file_analysis"); ok {
		configRequest.IsDeleteAnyFileAnalysis = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOkExists("is_read_any_file_analysis"); ok {
		configRequest.IsReadAnyFileAnalysis = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOkExists("is_upload_any_file_analysis"); ok {
		configRequest.IsUploadAnyFileAnalysis = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOkExists("is_include_any_file_analysis"); ok {
		configRequest.IsIncludeAnyFileAnalysis = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOkExists("is_directory_traversal_analysis"); ok {
		configRequest.IsDirectoryTraversalAnalysis = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOkExists("is_template_engine_injection_analysis"); ok {
		configRequest.IsTemplateEngineInjectionAnalysis = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOkExists("is_script_engine_injection_analysis"); ok {
		configRequest.IsScriptEngineInjectionAnalysis = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOkExists("is_expression_injection_analysis"); ok {
		configRequest.IsExpressionInjectionAnalysis = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOkExists("is_jndi_injection_analysis"); ok {
		configRequest.IsJNDIInjectionAnalysis = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOkExists("is_jni_injection_analysis"); ok {
		configRequest.IsJNIInjectionAnalysis = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOkExists("is_webshell_backdoor_analysis"); ok {
		configRequest.IsWebshellBackdoorAnalysis = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOkExists("is_deserialization_analysis"); ok {
		configRequest.IsDeserializationAnalysis = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOkExists("url_long_segment_threshold"); ok {
		configRequest.UrlLongSegmentThreshold = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOkExists("url_number_segment_threshold"); ok {
		configRequest.UrlNumberSegmentThreshold = helper.IntInt64(v.(int))
	}

	configRequest.InstanceId = &instanceId
	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseApmClient().ModifyApmInstanceWithContext(ctx, configRequest)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success,  configRequest body [%s], response body [%s]\n", logId, configRequest.GetAction(), configRequest.ToJsonString(), result.ToJsonString())
		}

		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s modify apm instance failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	return resourceTencentCloudApmInstanceRead(d, meta)
}

func resourceTencentCloudApmInstanceRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_apm_instance.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId      = tccommon.GetLogId(tccommon.ContextNil)
		ctx        = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service    = ApmService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		instance   *apm.ApmInstanceDetail
		instanceId = d.Id()
	)

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, err := service.DescribeApmInstanceById(ctx, instanceId)
		if err != nil {
			return tccommon.RetryError(err)
		}

		instance = result
		return nil
	})

	if err != nil {
		return err
	}

	if instance == nil {
		log.Printf("[WARN]%s resource `tencentcloud_apm_instance` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		d.SetId("")
		return nil
	}

	if instance.Name != nil {
		_ = d.Set("name", instance.Name)
	}

	if instance.Description != nil {
		_ = d.Set("description", instance.Description)
	}

	if instance.TraceDuration != nil {
		_ = d.Set("trace_duration", instance.TraceDuration)
	}

	if instance.SpanDailyCounters != nil {
		_ = d.Set("span_daily_counters", instance.SpanDailyCounters)
	}

	if instance.PayMode != nil {
		_ = d.Set("pay_mode", instance.PayMode)
	}

	if instance.Free != nil {
		_ = d.Set("free", instance.Free)
	}

	if instance.BillingInstance != nil {
		if *instance.BillingInstance == 0 {
			_ = d.Set("open_billing", false)
		} else {
			_ = d.Set("open_billing", true)
		}
	}

	if instance.ErrRateThreshold != nil {
		_ = d.Set("err_rate_threshold", instance.ErrRateThreshold)
	}

	if instance.SampleRate != nil {
		_ = d.Set("sample_rate", instance.SampleRate)
	}

	if instance.ErrorSample != nil {
		_ = d.Set("error_sample", instance.ErrorSample)
	}

	if instance.SlowRequestSavedThreshold != nil {
		_ = d.Set("slow_request_saved_threshold", instance.SlowRequestSavedThreshold)
	}

	if instance.IsRelatedLog != nil {
		_ = d.Set("is_related_log", instance.IsRelatedLog)
	}

	if instance.LogRegion != nil {
		_ = d.Set("log_region", instance.LogRegion)
	}

	if instance.LogTopicID != nil {
		_ = d.Set("log_topic_id", instance.LogTopicID)
	}

	if instance.LogSet != nil {
		_ = d.Set("log_set", instance.LogSet)
	}

	if instance.LogSource != nil {
		_ = d.Set("log_source", instance.LogSource)
	}

	if instance.CustomShowTags != nil {
		_ = d.Set("custom_show_tags", instance.CustomShowTags)
	}

	if instance.ResponseDurationWarningThreshold != nil {
		_ = d.Set("response_duration_warning_threshold", instance.ResponseDurationWarningThreshold)
	}

	if instance.IsRelatedDashboard != nil {
		_ = d.Set("is_related_dashboard", instance.IsRelatedDashboard)
	}

	if instance.DashboardTopicID != nil {
		_ = d.Set("dashboard_topic_id", instance.DashboardTopicID)
	}

	if instance.IsSqlInjectionAnalysis != nil {
		_ = d.Set("is_sql_injection_analysis", instance.IsSqlInjectionAnalysis)
	}

	if instance.IsInstrumentationVulnerabilityScan != nil {
		_ = d.Set("is_instrumentation_vulnerability_scan", instance.IsInstrumentationVulnerabilityScan)
	}

	if instance.IsRemoteCommandExecutionAnalysis != nil {
		_ = d.Set("is_remote_command_execution_analysis", instance.IsRemoteCommandExecutionAnalysis)
	}

	if instance.IsMemoryHijackingAnalysis != nil {
		_ = d.Set("is_memory_hijacking_analysis", instance.IsMemoryHijackingAnalysis)
	}

	if instance.LogIndexType != nil {
		_ = d.Set("log_index_type", instance.LogIndexType)
	}

	if instance.LogTraceIdKey != nil {
		_ = d.Set("log_trace_id_key", instance.LogTraceIdKey)
	}

	if instance.IsDeleteAnyFileAnalysis != nil {
		_ = d.Set("is_delete_any_file_analysis", instance.IsDeleteAnyFileAnalysis)
	}

	if instance.IsReadAnyFileAnalysis != nil {
		_ = d.Set("is_read_any_file_analysis", instance.IsReadAnyFileAnalysis)
	}

	if instance.IsUploadAnyFileAnalysis != nil {
		_ = d.Set("is_upload_any_file_analysis", instance.IsUploadAnyFileAnalysis)
	}

	if instance.IsIncludeAnyFileAnalysis != nil {
		_ = d.Set("is_include_any_file_analysis", instance.IsIncludeAnyFileAnalysis)
	}

	if instance.IsDirectoryTraversalAnalysis != nil {
		_ = d.Set("is_directory_traversal_analysis", instance.IsDirectoryTraversalAnalysis)
	}

	if instance.IsTemplateEngineInjectionAnalysis != nil {
		_ = d.Set("is_template_engine_injection_analysis", instance.IsTemplateEngineInjectionAnalysis)
	}

	if instance.IsScriptEngineInjectionAnalysis != nil {
		_ = d.Set("is_script_engine_injection_analysis", instance.IsScriptEngineInjectionAnalysis)
	}

	if instance.IsExpressionInjectionAnalysis != nil {
		_ = d.Set("is_expression_injection_analysis", instance.IsExpressionInjectionAnalysis)
	}

	if instance.IsJNDIInjectionAnalysis != nil {
		_ = d.Set("is_jndi_injection_analysis", instance.IsJNDIInjectionAnalysis)
	}

	if instance.IsJNIInjectionAnalysis != nil {
		_ = d.Set("is_jni_injection_analysis", instance.IsJNIInjectionAnalysis)
	}

	if instance.IsWebshellBackdoorAnalysis != nil {
		_ = d.Set("is_webshell_backdoor_analysis", instance.IsWebshellBackdoorAnalysis)
	}

	if instance.IsDeserializationAnalysis != nil {
		_ = d.Set("is_deserialization_analysis", instance.IsDeserializationAnalysis)
	}

	if instance.UrlLongSegmentThreshold != nil {
		_ = d.Set("url_long_segment_threshold", instance.UrlLongSegmentThreshold)
	}

	if instance.UrlNumberSegmentThreshold != nil {
		_ = d.Set("url_number_segment_threshold", instance.UrlNumberSegmentThreshold)
	}

	if instance.InstanceId != nil {
		_ = d.Set("instance_id", instance.InstanceId)
	}

	if instance.Token != nil {
		_ = d.Set("token", instance.Token)
	}

	apmAgent := &apm.ApmAgentInfo{}
	err = resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, err := service.DescribeApmAgentById(ctx, instanceId)
		if err != nil {
			return tccommon.RetryError(err)
		}

		apmAgent = result
		return nil
	})

	if err == nil {
		if apmAgent != nil && apmAgent.PublicCollectorURL != nil {
			_ = d.Set("public_collector_url", apmAgent.PublicCollectorURL)
		}
	}

	tcClient := meta.(tccommon.ProviderMeta).GetAPIV3Conn()
	tagService := svctag.NewTagService(tcClient)
	tags, err := tagService.DescribeResourceTags(ctx, "apm", "apm-instance", tcClient.Region, d.Id())
	if err != nil {
		return err
	}

	_ = d.Set("tags", tags)

	return nil
}

func resourceTencentCloudApmInstanceUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_apm_instance.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId      = tccommon.GetLogId(tccommon.ContextNil)
		request    = apm.NewModifyApmInstanceRequest()
		instanceId = d.Id()
	)

	if v, ok := d.GetOk("name"); ok {
		request.Name = helper.String(v.(string))
	}

	if tags := helper.GetTags(d, "tags"); len(tags) > 0 {
		for k, v := range tags {
			key := k
			value := v
			tag := &apm.ApmTag{
				Key:   &key,
				Value: &value,
			}

			request.Tags = append(request.Tags, tag)
		}
	}

	if v, ok := d.GetOk("description"); ok {
		request.Description = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("trace_duration"); ok {
		request.TraceDuration = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOkExists("open_billing"); ok {
		request.OpenBilling = helper.Bool(v.(bool))
	}

	if v, ok := d.GetOkExists("span_daily_counters"); ok {
		request.SpanDailyCounters = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOkExists("err_rate_threshold"); ok {
		request.ErrRateThreshold = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOkExists("sample_rate"); ok {
		request.SampleRate = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOkExists("error_sample"); ok {
		request.ErrorSample = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOkExists("slow_request_saved_threshold"); ok {
		request.SlowRequestSavedThreshold = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOkExists("is_related_log"); ok {
		request.IsRelatedLog = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("log_region"); ok {
		request.LogRegion = helper.String(v.(string))
	}

	if v, ok := d.GetOk("log_topic_id"); ok {
		request.LogTopicID = helper.String(v.(string))
	}

	if v, ok := d.GetOk("log_set"); ok {
		request.LogSet = helper.String(v.(string))
	}

	if v, ok := d.GetOk("log_source"); ok {
		request.LogSource = helper.String(v.(string))
	}

	if v, ok := d.GetOk("custom_show_tags"); ok {
		customShowTagsSet := v.(*schema.Set).List()
		for i := range customShowTagsSet {
			customShowTags := customShowTagsSet[i].(string)
			request.CustomShowTags = append(request.CustomShowTags, helper.String(customShowTags))
		}
	}

	if v, ok := d.GetOkExists("pay_mode"); ok {
		request.PayMode = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOkExists("response_duration_warning_threshold"); ok {
		request.ResponseDurationWarningThreshold = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOkExists("free"); ok {
		request.Free = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOkExists("is_related_dashboard"); ok {
		request.IsRelatedDashboard = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("dashboard_topic_id"); ok {
		request.DashboardTopicID = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("is_sql_injection_analysis"); ok {
		request.IsSqlInjectionAnalysis = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOkExists("is_instrumentation_vulnerability_scan"); ok {
		request.IsInstrumentationVulnerabilityScan = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOkExists("is_remote_command_execution_analysis"); ok {
		request.IsRemoteCommandExecutionAnalysis = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOkExists("is_memory_hijacking_analysis"); ok {
		request.IsMemoryHijackingAnalysis = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOkExists("log_index_type"); ok {
		request.LogIndexType = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("log_trace_id_key"); ok {
		request.LogTraceIdKey = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("is_delete_any_file_analysis"); ok {
		request.IsDeleteAnyFileAnalysis = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOkExists("is_read_any_file_analysis"); ok {
		request.IsReadAnyFileAnalysis = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOkExists("is_upload_any_file_analysis"); ok {
		request.IsUploadAnyFileAnalysis = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOkExists("is_include_any_file_analysis"); ok {
		request.IsIncludeAnyFileAnalysis = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOkExists("is_directory_traversal_analysis"); ok {
		request.IsDirectoryTraversalAnalysis = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOkExists("is_template_engine_injection_analysis"); ok {
		request.IsTemplateEngineInjectionAnalysis = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOkExists("is_script_engine_injection_analysis"); ok {
		request.IsScriptEngineInjectionAnalysis = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOkExists("is_expression_injection_analysis"); ok {
		request.IsExpressionInjectionAnalysis = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOkExists("is_jndi_injection_analysis"); ok {
		request.IsJNDIInjectionAnalysis = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOkExists("is_jni_injection_analysis"); ok {
		request.IsJNIInjectionAnalysis = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOkExists("is_webshell_backdoor_analysis"); ok {
		request.IsWebshellBackdoorAnalysis = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOkExists("is_deserialization_analysis"); ok {
		request.IsDeserializationAnalysis = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOkExists("url_long_segment_threshold"); ok {
		request.UrlLongSegmentThreshold = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOkExists("url_number_segment_threshold"); ok {
		request.UrlNumberSegmentThreshold = helper.IntInt64(v.(int))
	}

	request.InstanceId = &instanceId
	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseApmClient().ModifyApmInstance(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s update apm instance failed, reason:%+v", logId, err)
		return err
	}
	return resourceTencentCloudApmInstanceRead(d, meta)
}

func resourceTencentCloudApmInstanceDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_apm_instance.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId      = tccommon.GetLogId(tccommon.ContextNil)
		ctx        = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service    = ApmService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		instanceId = d.Id()
	)

	if err := service.DeleteApmInstanceById(ctx, instanceId); err != nil {
		return err
	}

	return nil
}
