package apm

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	apmv20210622 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/apm/v20210622"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudApmApplicationConfig() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudApmApplicationConfigCreate,
		Read:   resourceTencentCloudApmApplicationConfigRead,
		Update: resourceTencentCloudApmApplicationConfigUpdate,
		Delete: resourceTencentCloudApmApplicationConfigDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Business system ID.",
			},

			"service_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Application name.",
			},

			"url_convergence_switch": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "URL convergence switch. 0: Off; 1: On.",
			},

			"url_convergence_threshold": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "URL convergence threshold.",
			},

			"exception_filter": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Regex rules for exception filtering, separated by commas.",
			},

			"url_convergence": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Regex rules for URL convergence, separated by commas.",
			},

			"error_code_filter": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Error code filtering, separated by commas.",
			},

			"url_exclude": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Regex rules for URL exclusion, separated by commas.",
			},

			"is_related_log": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "Log switch. 0: Off; 1: On.",
			},

			"log_region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Log region.",
			},

			"log_topic_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Log topic ID.",
			},

			"log_set": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "CLS log set/ES cluster ID.",
			},

			"log_source": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Log source: CLS or ES.",
			},

			"ignore_operation_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "APIs to be filtered.",
			},

			"enable_snapshot": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: "Whether thread profiling is enabled.",
			},

			"snapshot_timeout": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "Timeout threshold for thread profiling.",
			},

			"agent_enable": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: "Whether agent is enabled.",
			},

			"trace_squash": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: "Whether link compression is enabled.",
			},

			"event_enable": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: "Switch for enabling application diagnosis.",
			},

			"instrument_list": {
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				Description: "Component List.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "Component name.",
						},
						"enable": {
							Type:        schema.TypeBool,
							Optional:    true,
							Computed:    true,
							Description: "Component switch.",
						},
					},
				},
			},

			"agent_operation_config_view": {
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				MaxItems:    1,
				Description: "Related configurations of the probe APIs.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"retention_valid": {
							Type:        schema.TypeBool,
							Optional:    true,
							Computed:    true,
							Description: "Whether allowlist configuration is enabled for the current API.\nNote: This field may return null, indicating that no valid values can be obtained.",
						},
						"ignore_operation": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "Effective when RetentionValid is false. It indicates blocklist configuration in API settings. The APIs specified in the configuration do not support collection.\nNote: This field may return null, indicating that no valid values can be obtained.",
						},
						"retention_operation": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "Effective when RetentionValid is true. It indicates allowlist configuration in API settings. Only the APIs specified in the configuration support collection.\nNote: This field may return null, indicating that no valid values can be obtained.",
						},
					},
				},
			},

			"enable_log_config": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: "Whether to enable application log configuration.",
			},

			"enable_dashboard_config": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: "Whether to enable the dashboard configuration for applications. false: disabled (consistent with the business system configuration); true: enabled (application-level configuration).",
			},

			"is_related_dashboard": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "Whether to associate with Dashboard. 0: disabled; 1: enabled.",
			},

			"dashboard_topic_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "dashboard ID.",
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
				Computed:    true,
				Description: "Index key of traceId. It is valid when the CLS index type is key-value index.",
			},

			"enable_security_config": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: "Whether to enable application security configuration.",
			},
			"is_sql_injection_analysis": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "Whether to enable SQL injection analysis.",
			},

			"is_instrumentation_vulnerability_scan": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "Whether to enable detection of component vulnerability.",
			},

			"is_remote_command_execution_analysis": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "Whether remote command detection is enabled.",
			},

			"is_memory_hijacking_analysis": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "Whether to enable detection of Java webshell.",
			},

			"is_delete_any_file_analysis": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "Whether to enable the detection of deleting arbitrary files. (0 - disabled; 1: enabled.).",
			},

			"is_read_any_file_analysis": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "Whether to enable the detection of reading arbitrary files. (0 - disabled; 1 - enabled.).",
			},

			"is_upload_any_file_analysis": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "Whether to enable the detection of uploading arbitrary files. (0 - disabled; 1 - enabled.).",
			},

			"is_include_any_file_analysis": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "Whether to enable the detection of the inclusion of arbitrary files. (0: disabled, 1: enabled.).",
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
				Description: "Whether to enable template engine injection detection. (0: disabled; 1: enabled.).",
			},

			"is_script_engine_injection_analysis": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "Whether to enable script engine injection detection. (0 - disabled; 1 - enabled.).",
			},

			"is_expression_injection_analysis": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "Whether to enable expression injection detection. (0 - disabled; 1 - enabled.).",
			},

			"is_jndi_injection_analysis": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "Whether to enable JNDI injection detection. (0 - disabled; 1 - enabled.).",
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

			"url_auto_convergence_enable": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: "Automatic convergence switch for APIs. 0: disabled | 1: enabled.",
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

			"disable_memory_used": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "Specifies the memory threshold for probe fusing.",
			},

			"disable_cpu_used": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "Specifies the CPU threshold for probe fusing.",
			},
		},
	}
}

func resourceTencentCloudApmApplicationConfigCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_apm_application_config.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		instanceId  string
		serviceName string
	)

	if v, ok := d.GetOk("instance_id"); ok {
		instanceId = v.(string)
	}

	if v, ok := d.GetOk("service_name"); ok {
		serviceName = v.(string)
	}

	d.SetId(strings.Join([]string{instanceId, serviceName}, tccommon.FILED_SP))

	return resourceTencentCloudApmApplicationConfigUpdate(d, meta)
}

func resourceTencentCloudApmApplicationConfigRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_apm_application_config.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		service = ApmService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}

	instanceId := idSplit[0]
	serviceName := idSplit[1]

	respData, err := service.DescribeApmApplicationConfigById(ctx, instanceId, serviceName)
	if err != nil {
		return err
	}

	if respData == nil {
		log.Printf("[WARN]%s resource `tencentcloud_apm_application_config` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		d.SetId("")
		return nil
	}

	if respData.InstanceKey != nil {
		_ = d.Set("instance_id", respData.InstanceKey)
	}

	if respData.ServiceName != nil {
		_ = d.Set("service_name", respData.ServiceName)
	}

	if respData.UrlConvergenceSwitch != nil {
		_ = d.Set("url_convergence_switch", respData.UrlConvergenceSwitch)
	}

	if respData.UrlConvergenceThreshold != nil {
		_ = d.Set("url_convergence_threshold", respData.UrlConvergenceThreshold)
	}

	if respData.ExceptionFilter != nil {
		_ = d.Set("exception_filter", respData.ExceptionFilter)
	}

	if respData.UrlConvergence != nil {
		_ = d.Set("url_convergence", respData.UrlConvergence)
	}

	if respData.ErrorCodeFilter != nil {
		_ = d.Set("error_code_filter", respData.ErrorCodeFilter)
	}

	if respData.UrlExclude != nil {
		_ = d.Set("url_exclude", respData.UrlExclude)
	}

	if respData.IsRelatedLog != nil {
		_ = d.Set("is_related_log", respData.IsRelatedLog)
	}

	if respData.LogRegion != nil {
		_ = d.Set("log_region", respData.LogRegion)
	}

	if respData.LogTopicID != nil {
		_ = d.Set("log_topic_id", respData.LogTopicID)
	}

	if respData.LogSet != nil {
		_ = d.Set("log_set", respData.LogSet)
	}

	if respData.LogSource != nil {
		_ = d.Set("log_source", respData.LogSource)
	}

	if respData.IgnoreOperationName != nil {
		_ = d.Set("ignore_operation_name", respData.IgnoreOperationName)
	}

	if respData.EnableSnapshot != nil {
		_ = d.Set("enable_snapshot", respData.EnableSnapshot)
	}

	if respData.SnapshotTimeout != nil {
		_ = d.Set("snapshot_timeout", respData.SnapshotTimeout)
	}

	if respData.AgentEnable != nil {
		_ = d.Set("agent_enable", respData.AgentEnable)
	}

	if respData.TraceSquash != nil {
		_ = d.Set("trace_squash", respData.TraceSquash)
	}

	if respData.EventEnable != nil {
		_ = d.Set("event_enable", respData.EventEnable)
	}

	if respData.InstrumentList != nil && len(respData.InstrumentList) > 0 {
		instrumentListList := make([]map[string]interface{}, 0, len(respData.InstrumentList))
		for _, instrumentList := range respData.InstrumentList {
			instrumentListMap := map[string]interface{}{}
			if instrumentList.Name != nil {
				instrumentListMap["name"] = instrumentList.Name
			}

			if instrumentList.Enable != nil {
				instrumentListMap["enable"] = instrumentList.Enable
			}

			instrumentListList = append(instrumentListList, instrumentListMap)
		}

		_ = d.Set("instrument_list", instrumentListList)
	}

	if respData.AgentOperationConfigView != nil {
		agentOperationConfigViewMap := map[string]interface{}{}
		if respData.AgentOperationConfigView.RetentionValid != nil {
			agentOperationConfigViewMap["retention_valid"] = respData.AgentOperationConfigView.RetentionValid
		}

		if respData.AgentOperationConfigView.IgnoreOperation != nil {
			agentOperationConfigViewMap["ignore_operation"] = respData.AgentOperationConfigView.IgnoreOperation
		}

		if respData.AgentOperationConfigView.RetentionOperation != nil {
			agentOperationConfigViewMap["retention_operation"] = respData.AgentOperationConfigView.RetentionOperation
		}

		_ = d.Set("agent_operation_config_view", []interface{}{agentOperationConfigViewMap})
	}

	if respData.EnableLogConfig != nil {
		_ = d.Set("enable_log_config", respData.EnableLogConfig)
	}

	if respData.EnableDashboardConfig != nil {
		_ = d.Set("enable_dashboard_config", respData.EnableDashboardConfig)
	}

	if respData.IsRelatedDashboard != nil {
		_ = d.Set("is_related_dashboard", respData.IsRelatedDashboard)
	}

	if respData.DashboardTopicID != nil {
		_ = d.Set("dashboard_topic_id", respData.DashboardTopicID)
	}

	if respData.LogIndexType != nil {
		_ = d.Set("log_index_type", respData.LogIndexType)
	}

	if respData.LogTraceIdKey != nil {
		_ = d.Set("log_trace_id_key", respData.LogTraceIdKey)
	}

	if respData.EnableSecurityConfig != nil {
		_ = d.Set("enable_security_config", respData.EnableSecurityConfig)
	}

	if respData.IsSqlInjectionAnalysis != nil {
		_ = d.Set("is_sql_injection_analysis", respData.IsSqlInjectionAnalysis)
	}

	if respData.IsInstrumentationVulnerabilityScan != nil {
		_ = d.Set("is_instrumentation_vulnerability_scan", respData.IsInstrumentationVulnerabilityScan)
	}

	if respData.IsRemoteCommandExecutionAnalysis != nil {
		_ = d.Set("is_remote_command_execution_analysis", respData.IsRemoteCommandExecutionAnalysis)
	}

	if respData.IsMemoryHijackingAnalysis != nil {
		_ = d.Set("is_memory_hijacking_analysis", respData.IsMemoryHijackingAnalysis)
	}

	if respData.IsDeleteAnyFileAnalysis != nil {
		_ = d.Set("is_delete_any_file_analysis", respData.IsDeleteAnyFileAnalysis)
	}

	if respData.IsReadAnyFileAnalysis != nil {
		_ = d.Set("is_read_any_file_analysis", respData.IsReadAnyFileAnalysis)
	}

	if respData.IsUploadAnyFileAnalysis != nil {
		_ = d.Set("is_upload_any_file_analysis", respData.IsUploadAnyFileAnalysis)
	}

	if respData.IsIncludeAnyFileAnalysis != nil {
		_ = d.Set("is_include_any_file_analysis", respData.IsIncludeAnyFileAnalysis)
	}

	if respData.IsDirectoryTraversalAnalysis != nil {
		_ = d.Set("is_directory_traversal_analysis", respData.IsDirectoryTraversalAnalysis)
	}

	if respData.IsTemplateEngineInjectionAnalysis != nil {
		_ = d.Set("is_template_engine_injection_analysis", respData.IsTemplateEngineInjectionAnalysis)
	}

	if respData.IsScriptEngineInjectionAnalysis != nil {
		_ = d.Set("is_script_engine_injection_analysis", respData.IsScriptEngineInjectionAnalysis)
	}

	if respData.IsExpressionInjectionAnalysis != nil {
		_ = d.Set("is_expression_injection_analysis", respData.IsExpressionInjectionAnalysis)
	}

	if respData.IsJNDIInjectionAnalysis != nil {
		_ = d.Set("is_jndi_injection_analysis", respData.IsJNDIInjectionAnalysis)
	}

	if respData.IsJNIInjectionAnalysis != nil {
		_ = d.Set("is_jni_injection_analysis", respData.IsJNIInjectionAnalysis)
	}

	if respData.IsWebshellBackdoorAnalysis != nil {
		_ = d.Set("is_webshell_backdoor_analysis", respData.IsWebshellBackdoorAnalysis)
	}

	if respData.IsDeserializationAnalysis != nil {
		_ = d.Set("is_deserialization_analysis", respData.IsDeserializationAnalysis)
	}

	if respData.UrlAutoConvergenceEnable != nil {
		_ = d.Set("url_auto_convergence_enable", respData.UrlAutoConvergenceEnable)
	}

	if respData.UrlLongSegmentThreshold != nil {
		_ = d.Set("url_long_segment_threshold", respData.UrlLongSegmentThreshold)
	}

	if respData.UrlNumberSegmentThreshold != nil {
		_ = d.Set("url_number_segment_threshold", respData.UrlNumberSegmentThreshold)
	}

	if respData.DisableMemoryUsed != nil {
		_ = d.Set("disable_memory_used", respData.DisableMemoryUsed)
	}

	if respData.DisableCpuUsed != nil {
		_ = d.Set("disable_cpu_used", respData.DisableCpuUsed)
	}

	return nil
}

func resourceTencentCloudApmApplicationConfigUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_apm_application_config.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId = tccommon.GetLogId(tccommon.ContextNil)
		ctx   = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
	)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}

	instanceId := idSplit[0]
	serviceName := idSplit[1]

	request := apmv20210622.NewModifyApmApplicationConfigRequest()
	if v, ok := d.GetOkExists("url_convergence_switch"); ok {
		request.UrlConvergenceSwitch = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOkExists("url_convergence_threshold"); ok {
		request.UrlConvergenceThreshold = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("exception_filter"); ok {
		request.ExceptionFilter = helper.String(v.(string))
	}

	if v, ok := d.GetOk("url_convergence"); ok {
		request.UrlConvergence = helper.String(v.(string))
	}

	if v, ok := d.GetOk("error_code_filter"); ok {
		request.ErrorCodeFilter = helper.String(v.(string))
	}

	if v, ok := d.GetOk("url_exclude"); ok {
		request.UrlExclude = helper.String(v.(string))
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

	if v, ok := d.GetOk("ignore_operation_name"); ok {
		request.IgnoreOperationName = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("enable_snapshot"); ok {
		request.EnableSnapshot = helper.Bool(v.(bool))
	}

	if v, ok := d.GetOkExists("snapshot_timeout"); ok {
		request.SnapshotTimeout = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOkExists("agent_enable"); ok {
		request.AgentEnable = helper.Bool(v.(bool))
	}

	if v, ok := d.GetOkExists("trace_squash"); ok {
		request.TraceSquash = helper.Bool(v.(bool))
	}

	if v, ok := d.GetOkExists("event_enable"); ok {
		request.EventEnable = helper.Bool(v.(bool))
	}

	if v, ok := d.GetOk("instrument_list"); ok {
		for _, item := range v.([]interface{}) {
			instrumentListMap := item.(map[string]interface{})
			instrument := apmv20210622.Instrument{}
			if v, ok := instrumentListMap["name"].(string); ok && v != "" {
				instrument.Name = helper.String(v)
			}

			if v, ok := instrumentListMap["enable"].(bool); ok {
				instrument.Enable = helper.Bool(v)
			}

			request.InstrumentList = append(request.InstrumentList, &instrument)
		}
	}

	if agentOperationConfigViewMap, ok := helper.InterfacesHeadMap(d, "agent_operation_config_view"); ok {
		agentOperationConfigView := apmv20210622.AgentOperationConfigView{}
		if v, ok := agentOperationConfigViewMap["retention_valid"].(bool); ok {
			agentOperationConfigView.RetentionValid = helper.Bool(v)
		}

		if v, ok := agentOperationConfigViewMap["ignore_operation"].(string); ok && v != "" {
			agentOperationConfigView.IgnoreOperation = helper.String(v)
		}

		if v, ok := agentOperationConfigViewMap["retention_operation"].(string); ok && v != "" {
			agentOperationConfigView.RetentionOperation = helper.String(v)
		}

		request.AgentOperationConfigView = &agentOperationConfigView
	}

	if v, ok := d.GetOkExists("enable_log_config"); ok {
		request.EnableLogConfig = helper.Bool(v.(bool))
	}

	if v, ok := d.GetOkExists("enable_dashboard_config"); ok {
		request.EnableDashboardConfig = helper.Bool(v.(bool))
	}

	if v, ok := d.GetOkExists("is_related_dashboard"); ok {
		request.IsRelatedDashboard = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("dashboard_topic_id"); ok {
		request.DashboardTopicID = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("log_index_type"); ok {
		request.LogIndexType = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("log_trace_id_key"); ok {
		request.LogTraceIdKey = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("enable_security_config"); ok {
		request.EnableSecurityConfig = helper.Bool(v.(bool))
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

	if v, ok := d.GetOkExists("url_auto_convergence_enable"); ok {
		request.UrlAutoConvergenceEnable = helper.Bool(v.(bool))
	}

	if v, ok := d.GetOkExists("url_long_segment_threshold"); ok {
		request.UrlLongSegmentThreshold = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOkExists("url_number_segment_threshold"); ok {
		request.UrlNumberSegmentThreshold = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOkExists("disable_memory_used"); ok {
		request.DisableMemoryUsed = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOkExists("disable_cpu_used"); ok {
		request.DisableCpuUsed = helper.IntInt64(v.(int))
	}

	request.InstanceId = &instanceId
	request.ServiceName = &serviceName
	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseApmClient().ModifyApmApplicationConfigWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s update apm application config failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	return resourceTencentCloudApmApplicationConfigRead(d, meta)
}

func resourceTencentCloudApmApplicationConfigDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_apm_application_config.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}
