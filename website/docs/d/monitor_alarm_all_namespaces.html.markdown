---
subcategory: "Cloud Monitor(Monitor)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_monitor_alarm_all_namespaces"
sidebar_current: "docs-tencentcloud-datasource-monitor_alarm_all_namespaces"
description: |-
  Use this data source to query detailed information of monitor alarm_all_namespaces
---

# tencentcloud_monitor_alarm_all_namespaces

Use this data source to query detailed information of monitor alarm_all_namespaces

## Example Usage

```hcl
data "tencentcloud_monitor_alarm_all_namespaces" "alarm_all_namespaces" {
  scene_type    = "ST_ALARM"
  module        = "monitor"
  monitor_types = ["MT_QCE"]
  ids           = ["qaap_tunnel_l4_listeners"]
}
```

## Argument Reference

The following arguments are supported:

* `module` - (Required, String) Fixed value, as `monitor`.
* `scene_type` - (Required, String) Currently, only ST_ALARM=alarm type is filtered based on usage scenarios.
* `ids` - (Optional, Set: [`String`]) Filter based on the Id of the namespace without filling in the default query for all.
* `monitor_types` - (Optional, Set: [`String`]) Filter based on monitoring type, do not fill in default, check all types MT_QCE=cloud product monitoring.
* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `common_namespaces` - General alarm strategy types (including: application performance monitoring, front-end performance monitoring, cloud dial testing).
  * `dimensions` - Dimension Information.
    * `can_filter_history` - Can it be used to filter alarm history.
    * `can_filter_policy` - Can it be used to filter the policy list.
    * `can_group_by` - Can it be used as an aggregation dimension.
    * `is_multiple` - Do you support multiple selections.
    * `is_mutable` - Can I modify it after creation.
    * `is_required` - Required or not.
    * `is_visible` - Whether to display to users.
    * `key` - Dimension key identifier, backend English name.
    * `must_group_by` - Must it be used as an aggregation dimension.
    * `name` - Dimension key name, Chinese and English frontend display name.
    * `operators` - List of supported operators.
      * `id` - Operator identification.
      * `name` - Operator Display Name.
    * `show_value_replace` - Key to replace in front-end translation.
  * `id` - Namespace labeling.
  * `monitor_type` - Monitoring type.
  * `name` - Namespace name.
* `custom_namespaces_new` - Other alarm strategy types are currently not supported.
  * `available_regions` - List of supported regions.
  * `config` - Configuration information.
  * `dashboard_id` - Unique representation in dashboard.
  * `id` - Namespace labeling.
  * `name` - Namespace name.
  * `product_name` - Product Name.
  * `sort_id` - Sort Id.
  * `value` - Namespace value.
* `qce_namespaces_new` - Types of alarm strategies for cloud products.
  * `available_regions` - List of supported regions.
  * `config` - Configuration information.
  * `dashboard_id` - Unique representation in dashboard.
  * `id` - Namespace labeling.
  * `name` - Namespace name.
  * `product_name` - Product Name.
  * `sort_id` - Sort Id.
  * `value` - Namespace value.


