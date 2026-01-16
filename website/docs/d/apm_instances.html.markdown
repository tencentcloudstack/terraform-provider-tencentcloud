---
subcategory: "Application Performance Management(APM)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_apm_instances"
sidebar_current: "docs-tencentcloud-datasource-apm_instances"
description: |-
  Use this data source to query APM (Application Performance Management) instances.
---

# tencentcloud_apm_instances

Use this data source to query APM (Application Performance Management) instances.

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
  * `app_id` - App ID.
  * `create_uin` - Creator UIN.
  * `description` - Instance description.
  * `err_rate_threshold` - Error rate threshold.
  * `error_sample` - Error sampling switch.
  * `free` - Whether it is free edition.
  * `instance_id` - Instance ID.
  * `name` - Instance name.
  * `pay_mode` - Billing mode.
  * `region` - Region.
  * `sample_rate` - Sampling rate.
  * `service_count` - Service count.
  * `span_daily_counters` - Daily span count quota.
  * `status` - Instance status.
  * `tags` - Tag list.
    * `key` - Tag key.
    * `value` - Tag value.
  * `trace_duration` - Trace data retention duration.


