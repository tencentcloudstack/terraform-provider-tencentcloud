---
subcategory: "Application Performance Management(APM)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_apm_instances"
sidebar_current: "docs-tencentcloud-datasource-apm_instances"
description: |-
  Use this data source to query APM (Application Performance Management) instances. It returns all fields from the DescribeApmInstances API, including instance basic info, billing, log configuration, security detection settings, and more.
---

# tencentcloud_apm_instances

Use this data source to query APM (Application Performance Management) instances. It returns all fields from the DescribeApmInstances API, including instance basic info, billing, log configuration, security detection settings, and more.

## Example Usage

### Query all APM instances

```hcl
data "tencentcloud_apm_instances" "all" {
}

output "instances" {
  value = data.tencentcloud_apm_instances.all.instance_list
}
```

### Query APM instances by IDs

```hcl
data "tencentcloud_apm_instances" "by_ids" {
  instance_ids = ["apm-xxxxxxxx", "apm-yyyyyyyy"]
}

output "instances" {
  value = data.tencentcloud_apm_instances.by_ids.instance_list
}
```

### Query APM instances by name

```hcl
data "tencentcloud_apm_instances" "by_name" {
  instance_name = "test"
}

output "instances" {
  value = data.tencentcloud_apm_instances.by_name.instance_list
}
```

### Query APM instances by tags

```hcl
data "tencentcloud_apm_instances" "by_tags" {
  tags = {
    "Environment" = "Production"
    "Team"        = "DevOps"
  }
}

output "instances" {
  value = data.tencentcloud_apm_instances.by_tags.instance_list
}
```

## Argument Reference

The following arguments are supported:

* `all_regions_flag` - (Optional, Int) Whether to query instances in all regions. 0: no, 1: yes. Default is 0.
* `demo_instance_flag` - (Optional, Int) Whether to query official demo instances. 0: non-demo, 1: demo. Default is 0.
* `instance_id` - (Optional, String) Filter by instance ID (fuzzy match).
* `instance_ids` - (Optional, List: [`String`]) Filter by instance ID list (exact match).
* `instance_name` - (Optional, String) Filter by instance name (fuzzy match).
* `result_output_file` - (Optional, String) Used to save results.
* `tags` - (Optional, Map) Filter by tags.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `instance_list` - APM instance list.
  * `amount_of_used_storage` - Storage usage in MB.
  * `app_id` - App ID.
  * `billing_instance` - Whether billing is enabled. 0: not enabled, 1: enabled.
  * `client_count` - Client application count.
  * `count_of_report_span_per_day` - Daily reported span count.
  * `create_uin` - Creator UIN.
  * `custom_show_tags` - Custom display tag list.
  * `dashboard_topic_id` - Associated dashboard ID.
  * `default_tsf` - Whether it is the default TSF instance. 0: no, 1: yes.
  * `description` - Instance description.
  * `err_rate_threshold` - Error rate threshold.
  * `error_sample` - Error sampling switch.
  * `free` - Whether it is free edition.
  * `instance_id` - Instance ID.
  * `is_delete_any_file_analysis` - Whether delete any file detection is enabled. 0: off, 1: on.
  * `is_deserialization_analysis` - Whether deserialization detection is enabled. 0: off, 1: on.
  * `is_directory_traversal_analysis` - Whether directory traversal detection is enabled. 0: off, 1: on.
  * `is_expression_injection_analysis` - Whether expression injection detection is enabled. 0: off, 1: on.
  * `is_include_any_file_analysis` - Whether include any file detection is enabled. 0: off, 1: on.
  * `is_instrumentation_vulnerability_scan` - Whether instrumentation vulnerability scan is enabled. 0: off, 1: on.
  * `is_jndi_injection_analysis` - Whether JNDI injection detection is enabled. 0: off, 1: on.
  * `is_jni_injection_analysis` - Whether JNI injection detection is enabled. 0: off, 1: on.
  * `is_memory_hijacking_analysis` - Whether memory hijacking detection is enabled. 0: off, 1: on.
  * `is_read_any_file_analysis` - Whether read any file detection is enabled. 0: off, 1: on.
  * `is_related_dashboard` - Whether dashboard is associated. 0: off, 1: on.
  * `is_related_log` - Log feature switch. 0: off, 1: on.
  * `is_remote_command_execution_analysis` - Whether remote command execution detection is enabled. 0: off, 1: on.
  * `is_script_engine_injection_analysis` - Whether script engine injection detection is enabled. 0: off, 1: on.
  * `is_sql_injection_analysis` - Whether SQL injection analysis is enabled. 0: off, 1: on.
  * `is_template_engine_injection_analysis` - Whether template engine injection detection is enabled. 0: off, 1: on.
  * `is_upload_any_file_analysis` - Whether upload any file detection is enabled. 0: off, 1: on.
  * `is_webshell_backdoor_analysis` - Whether webshell backdoor detection is enabled. 0: off, 1: on.
  * `log_index_type` - CLS index type. 0: full-text index, 1: key-value index.
  * `log_region` - CLS log region.
  * `log_set` - CLS log set.
  * `log_source` - Log source.
  * `log_topic_id` - Log topic ID.
  * `log_trace_id_key` - TraceId index key, effective when CLS index type is key-value.
  * `metric_duration` - Metric data retention duration in days.
  * `name` - Instance name.
  * `pay_mode_effective` - Whether pay mode is effective.
  * `pay_mode` - Billing mode.
  * `region` - Region.
  * `response_duration_warning_threshold` - Response duration warning threshold in ms.
  * `sample_rate` - Sampling rate.
  * `service_count` - Service count.
  * `slow_request_saved_threshold` - Slow request saved threshold in ms.
  * `span_daily_counters` - Daily span count quota.
  * `status` - Instance status.
  * `stop_reason` - Throttling reason. 1: official version quota, 2: trial version quota, 4: trial expired, 8: account overdue.
  * `tags` - Tag list.
    * `key` - Tag key.
    * `value` - Tag value.
  * `token` - Instance authentication token.
  * `total_count` - Active application count in recent 2 days.
  * `trace_duration` - Trace data retention duration.
  * `url_long_segment_threshold` - URL long segment convergence threshold.
  * `url_number_segment_threshold` - URL number segment convergence threshold.


