---
subcategory: "Intelligent Global Traffic Manager(IGTM)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_igtm_strategy_list"
sidebar_current: "docs-tencentcloud-datasource-igtm_strategy_list"
description: |-
  Use this data source to query detailed information of IGTM strategy list
---

# tencentcloud_igtm_strategy_list

Use this data source to query detailed information of IGTM strategy list

## Example Usage

```hcl
data "tencentcloud_igtm_strategy_list" "example" {
  instance_id = "gtm-uukztqtoaru"
  filters {
    name  = "StrategyName"
    value = ["tf-example"]
    fuzzy = true
  }
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, String) Instance ID.
* `filters` - (Optional, List) Strategy filter conditions: StrategyName: strategy name.
* `result_output_file` - (Optional, String) Used to save results.

The `filters` object supports the following:

* `name` - (Required, String) Filter field name, supported list as follows:
- type: main resource type, CDN.
- instanceId: IGTM instance ID. This is a required parameter, failure to pass will cause interface query failure.
* `value` - (Required, Set) Filter field values.
* `fuzzy` - (Optional, Bool) Whether to enable fuzzy query, only supports filter field name as domain.
When fuzzy query is enabled, Value maximum length is 1, otherwise Value maximum length is 5. (Reserved field, currently unused).

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `strategy_set` - Strategy list.
  * `activate_level` - Current activated address pool level, 0 means fallback activated, null means unknown.
  * `activate_main_pool_id` - Activated main pool ID, null means unknown.
  * `active_pool_type` - Current activated address pool set type: main main pool; fallback fallback pool.
  * `active_traffic_strategy` - Current activated address pool traffic strategy: all resolve all; weight load balancing.
  * `created_on` - Creation time.
  * `instance_id` - Instance ID.
  * `is_enabled` - Whether enabled: ENABLED enabled; DISABLED disabled.
  * `keep_domain_records` - Whether to retain lines: enabled retain, disabled not retain, only retain default lines.
  * `monitor_num` - Monitor count.
  * `name` - Strategy name.
  * `source` - Address source.
    * `dns_line_id` - Resolution request source line ID.
    * `name` - Resolution request source line name.
  * `status` - Health status: ok healthy, warn risk, down failure.
  * `strategy_id` - Strategy ID.
  * `switch_pool_type` - Scheduling mode: AUTO default; PAUSE only pause without switching.
  * `updated_on` - Update time.


