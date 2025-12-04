---
subcategory: "Application Performance Management(APM)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_apm_application_config"
sidebar_current: "docs-tencentcloud-resource-apm_application_config"
description: |-
  Provides a resource to create a APM application config
---

# tencentcloud_apm_application_config

Provides a resource to create a APM application config

## Example Usage

```hcl
resource "tencentcloud_apm_application_config" "example" {
  instance_id                           = tencentcloud_apm_instance.example.id
  service_name                          = "java-order-serive"
  url_convergence_switch                = 1
  agent_enable                          = true
  disable_cpu_used                      = 90
  disable_memory_used                   = 90
  enable_dashboard_config               = false
  enable_log_config                     = false
  enable_security_config                = false
  enable_snapshot                       = false
  event_enable                          = false
  is_delete_any_file_analysis           = 0
  is_deserialization_analysis           = 0
  is_directory_traversal_analysis       = 0
  is_expression_injection_analysis      = 0
  is_include_any_file_analysis          = 0
  is_instrumentation_vulnerability_scan = 1
  is_jndi_injection_analysis            = 0
  is_jni_injection_analysis             = 0
  is_memory_hijacking_analysis          = 0
  is_read_any_file_analysis             = 0
  is_related_dashboard                  = 0
  is_related_log                        = 0
  is_remote_command_execution_analysis  = 0
  is_script_engine_injection_analysis   = 0
  is_sql_injection_analysis             = 0
  is_template_engine_injection_analysis = 0
  is_upload_any_file_analysis           = 0
  is_webshell_backdoor_analysis         = 0
  log_index_type                        = 0
  log_source                            = "CLS"
  snapshot_timeout                      = 2000
  trace_squash                          = true
  url_auto_convergence_enable           = false
  url_convergence_threshold             = 1000
  url_long_segment_threshold            = 40
  url_number_segment_threshold          = 5

  agent_operation_config_view {
    retention_valid = false
  }

  instrument_list {
    enable = true
    name   = "apm-spring-annotations"
  }

  instrument_list {
    enable = true
    name   = "dubbo"
  }

  instrument_list {
    enable = true
    name   = "googlehttpclient"
  }

  instrument_list {
    enable = true
    name   = "grpc"
  }

  instrument_list {
    enable = true
    name   = "httpclient3"
  }

  instrument_list {
    enable = true
    name   = "httpclient4"
  }

  instrument_list {
    enable = true
    name   = "hystrix"
  }

  instrument_list {
    enable = true
    name   = "lettuce"
  }

  instrument_list {
    enable = true
    name   = "mongodb"
  }

  instrument_list {
    enable = true
    name   = "mybatis"
  }

  instrument_list {
    enable = true
    name   = "mysql"
  }

  instrument_list {
    enable = true
    name   = "okhttp"
  }

  instrument_list {
    enable = true
    name   = "redis"
  }

  instrument_list {
    enable = true
    name   = "rxjava"
  }

  instrument_list {
    enable = true
    name   = "spring-webmvc"
  }

  instrument_list {
    enable = true
    name   = "tomcat"
  }
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, String, ForceNew) Business system ID.
* `service_name` - (Required, String, ForceNew) Application name.
* `url_convergence_switch` - (Required, Int) URL convergence switch. 0: Off; 1: On.
* `agent_enable` - (Optional, Bool) Whether agent is enabled.
* `agent_operation_config_view` - (Optional, List) Related configurations of the probe APIs.
* `dashboard_topic_id` - (Optional, String) dashboard ID.
* `disable_cpu_used` - (Optional, Int) Specifies the CPU threshold for probe fusing.
* `disable_memory_used` - (Optional, Int) Specifies the memory threshold for probe fusing.
* `enable_dashboard_config` - (Optional, Bool) Whether to enable the dashboard configuration for applications. false: disabled (consistent with the business system configuration); true: enabled (application-level configuration).
* `enable_log_config` - (Optional, Bool) Whether to enable application log configuration.
* `enable_security_config` - (Optional, Bool) Whether to enable application security configuration.
* `enable_snapshot` - (Optional, Bool) Whether thread profiling is enabled.
* `error_code_filter` - (Optional, String) Error code filtering, separated by commas.
* `event_enable` - (Optional, Bool) Switch for enabling application diagnosis.
* `exception_filter` - (Optional, String) Regex rules for exception filtering, separated by commas.
* `ignore_operation_name` - (Optional, String) APIs to be filtered.
* `instrument_list` - (Optional, List) Component List.
* `is_delete_any_file_analysis` - (Optional, Int) Whether to enable the detection of deleting arbitrary files. (0 - disabled; 1: enabled.).
* `is_deserialization_analysis` - (Optional, Int) Whether to enable deserialization detection. (0 - disabled; 1 - enabled).
* `is_directory_traversal_analysis` - (Optional, Int) Whether to enable traversal detection of the directory. (0 - disabled; 1 - enabled).
* `is_expression_injection_analysis` - (Optional, Int) Whether to enable expression injection detection. (0 - disabled; 1 - enabled.).
* `is_include_any_file_analysis` - (Optional, Int) Whether to enable the detection of the inclusion of arbitrary files. (0: disabled, 1: enabled.).
* `is_instrumentation_vulnerability_scan` - (Optional, Int) Whether to enable detection of component vulnerability.
* `is_jndi_injection_analysis` - (Optional, Int) Whether to enable JNDI injection detection. (0 - disabled; 1 - enabled.).
* `is_jni_injection_analysis` - (Optional, Int) Whether to enable JNI injection detection. (0 - disabled, 1 - enabled).
* `is_memory_hijacking_analysis` - (Optional, Int) Whether to enable detection of Java webshell.
* `is_read_any_file_analysis` - (Optional, Int) Whether to enable the detection of reading arbitrary files. (0 - disabled; 1 - enabled.).
* `is_related_dashboard` - (Optional, Int) Whether to associate with Dashboard. 0: disabled; 1: enabled.
* `is_related_log` - (Optional, Int) Log switch. 0: Off; 1: On.
* `is_remote_command_execution_analysis` - (Optional, Int) Whether remote command detection is enabled.
* `is_script_engine_injection_analysis` - (Optional, Int) Whether to enable script engine injection detection. (0 - disabled; 1 - enabled.).
* `is_sql_injection_analysis` - (Optional, Int) Whether to enable SQL injection analysis.
* `is_template_engine_injection_analysis` - (Optional, Int) Whether to enable template engine injection detection. (0: disabled; 1: enabled.).
* `is_upload_any_file_analysis` - (Optional, Int) Whether to enable the detection of uploading arbitrary files. (0 - disabled; 1 - enabled.).
* `is_webshell_backdoor_analysis` - (Optional, Int) Whether to enable Webshell backdoor detection. (0 - disabled; 1 - enabled).
* `log_index_type` - (Optional, Int) CLS index type. (0 = full-text index; 1 = key-value index).
* `log_region` - (Optional, String) Log region.
* `log_set` - (Optional, String) CLS log set/ES cluster ID.
* `log_source` - (Optional, String) Log source: CLS or ES.
* `log_topic_id` - (Optional, String) Log topic ID.
* `log_trace_id_key` - (Optional, String) Index key of traceId. It is valid when the CLS index type is key-value index.
* `snapshot_timeout` - (Optional, Int) Timeout threshold for thread profiling.
* `trace_squash` - (Optional, Bool) Whether link compression is enabled.
* `url_auto_convergence_enable` - (Optional, Bool) Automatic convergence switch for APIs. 0: disabled | 1: enabled.
* `url_convergence_threshold` - (Optional, Int) URL convergence threshold.
* `url_convergence` - (Optional, String) Regex rules for URL convergence, separated by commas.
* `url_exclude` - (Optional, String) Regex rules for URL exclusion, separated by commas.
* `url_long_segment_threshold` - (Optional, Int) Convergence threshold for URL long segments.
* `url_number_segment_threshold` - (Optional, Int) Convergence threshold for URL numerical segments.

The `agent_operation_config_view` object supports the following:

* `ignore_operation` - (Optional, String) Effective when RetentionValid is false. It indicates blocklist configuration in API settings. The APIs specified in the configuration do not support collection.
Note: This field may return null, indicating that no valid values can be obtained.
* `retention_operation` - (Optional, String) Effective when RetentionValid is true. It indicates allowlist configuration in API settings. Only the APIs specified in the configuration support collection.
Note: This field may return null, indicating that no valid values can be obtained.
* `retention_valid` - (Optional, Bool) Whether allowlist configuration is enabled for the current API.
Note: This field may return null, indicating that no valid values can be obtained.

The `instrument_list` object supports the following:

* `enable` - (Optional, Bool) Component switch.
* `name` - (Optional, String) Component name.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

APM application config can be imported using the instanceId#serviceName, e.g.

```
terraform import tencentcloud_apm_application_config.example apm-jPr5iQL77#java-order-serive
```

