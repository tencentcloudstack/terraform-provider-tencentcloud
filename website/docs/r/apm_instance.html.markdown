---
subcategory: "Application Performance Management(APM)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_apm_instance"
sidebar_current: "docs-tencentcloud-resource-apm_instance"
description: |-
  Provides a resource to create a APM instance
---

# tencentcloud_apm_instance

Provides a resource to create a APM instance

~> **NOTE:** To use the field `pay_mode`, you need to contact official customer service to join the whitelist.

## Example Usage

```hcl
resource "tencentcloud_apm_instance" "example" {
  name                = "tf-example"
  description         = "desc."
  trace_duration      = 7
  span_daily_counters = 0
  tags = {
    createdBy = "Terraform"
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required, String) Name Of Instance.
* `custom_show_tags` - (Optional, Set: [`String`]) List of custom display tags.
* `dashboard_topic_id` - (Optional, String) Associated dashboard id, which takes effect after the associated dashboard is enabled.
* `description` - (Optional, String) Description Of Instance.
* `err_rate_threshold` - (Optional, Int) Error rate warning line. when the average error rate of the application exceeds this threshold, the system will give an abnormal note.
* `error_sample` - (Optional, Int) Error sampling switch (0: off, 1: on).
* `free` - (Optional, Int) Whether it is free (0 = paid edition; 1 = tsf restricted free edition; 2 = free edition), default 0.
* `is_delete_any_file_analysis` - (Optional, Int) Whether to enable the detection of deleting arbitrary files. (0 - disabled; 1: enabled).
* `is_deserialization_analysis` - (Optional, Int) Whether to enable deserialization detection. (0 - disabled; 1 - enabled).
* `is_directory_traversal_analysis` - (Optional, Int) Whether to enable traversal detection of the directory. (0 - disabled; 1 - enabled).
* `is_expression_injection_analysis` - (Optional, Int) Whether to enable expression injection detection. (0 - disabled; 1 - enabled).
* `is_include_any_file_analysis` - (Optional, Int) Whether to enable the detection of the inclusion of arbitrary files. (0: disabled, 1: enabled).
* `is_instrumentation_vulnerability_scan` - (Optional, Int) Whether to enable component vulnerability detection (0 = no, 1 = yes).
* `is_jndi_injection_analysis` - (Optional, Int) Whether to enable JNDI injection detection. (0 - disabled; 1 - enabled).
* `is_jni_injection_analysis` - (Optional, Int) Whether to enable JNI injection detection. (0 - disabled, 1 - enabled).
* `is_memory_hijacking_analysis` - (Optional, Int) Whether to enable detection of Java webshell.
* `is_read_any_file_analysis` - (Optional, Int) Whether to enable the detection of reading arbitrary files. (0 - disabled; 1 - enabled).
* `is_related_dashboard` - (Optional, Int) Whether to associate the dashboard (0 = off, 1 = on).
* `is_related_log` - (Optional, Int) Log feature switch (0: off; 1: on).
* `is_remote_command_execution_analysis` - (Optional, Int) Whether to enable detection of the remote command attack.
* `is_script_engine_injection_analysis` - (Optional, Int) Whether to enable script engine injection detection. (0 - disabled; 1 - enabled).
* `is_sql_injection_analysis` - (Optional, Int) SQL injection detection switch (0: off, 1: on).
* `is_template_engine_injection_analysis` - (Optional, Int) Whether to enable template engine injection detection. (0: disabled; 1: enabled).
* `is_upload_any_file_analysis` - (Optional, Int) Whether to enable the detection of uploading arbitrary files. (0 - disabled; 1 - enabled).
* `is_webshell_backdoor_analysis` - (Optional, Int) Whether to enable Webshell backdoor detection. (0 - disabled; 1 - enabled).
* `log_index_type` - (Optional, Int) CLS index type. (0 = full-text index; 1 = key-value index).
* `log_region` - (Optional, String) Log region, which takes effect after the log feature is enabled.
* `log_set` - (Optional, String) Logset, which takes effect only after the log feature is enabled.
* `log_source` - (Optional, String) Log source, which takes effect only after the log feature is enabled.
* `log_topic_id` - (Optional, String) CLS log topic id, which takes effect after the log feature is enabled.
* `log_trace_id_key` - (Optional, String) Index key of traceId. It is valid when the CLS index type is key-value index.
* `open_billing` - (Optional, Bool) Billing switch.
* `pay_mode` - (Optional, Int) Modify the billing mode: `1` means prepaid, `0` means pay-as-you-go, the default value is `0`.
* `response_duration_warning_threshold` - (Optional, Int) Response time warning line.
* `sample_rate` - (Optional, Int) Sampling rate (unit: %).
* `slow_request_saved_threshold` - (Optional, Int) Sampling slow call saving threshold (unit: ms).
* `span_daily_counters` - (Optional, Int) Quota Of Instance Reporting.
* `tags` - (Optional, Map) Tag description list.
* `trace_duration` - (Optional, Int) Duration Of Trace Data.
* `url_long_segment_threshold` - (Optional, Int) Convergence threshold for URL long segments.
* `url_number_segment_threshold` - (Optional, Int) Convergence threshold for URL numerical segments.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `instance_id` - APM instance ID.
* `public_collector_url` - External Network Reporting Address.
* `token` - Business system authentication token.


## Import

APM instance can be imported using the id, e.g.

```
terraform import tencentcloud_apm_instance.example apm-IMVrxXl1K
```

