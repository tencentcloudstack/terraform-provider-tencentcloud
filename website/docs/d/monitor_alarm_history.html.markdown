---
subcategory: "Cloud Monitor(Monitor)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_monitor_alarm_history"
sidebar_current: "docs-tencentcloud-datasource-monitor_alarm_history"
description: |-
  Use this data source to query detailed information of monitor alarm_history
---

# tencentcloud_monitor_alarm_history

Use this data source to query detailed information of monitor alarm_history

## Example Usage

```hcl
data "tencentcloud_monitor_alarm_history" "alarm_history" {
  module        = "monitor"
  order         = "DESC"
  start_time    = 1696608000
  end_time      = 1697212799
  monitor_types = ["MT_QCE"]
  project_ids   = [0]
  namespaces {
    monitor_type = "CpuUsage"
    namespace    = "cvm_device"
  }
  policy_name = "terraform_test"
  content     = "CPU利用率 > 3%"
  policy_ids  = ["policy-iejtp4ue"]
}
```

## Argument Reference

The following arguments are supported:

* `module` - (Required, String) Value fixed at monitor.
* `alarm_levels` - (Optional, Set: [`String`]) Alarm levels.
* `alarm_object` - (Optional, String) Filter by alarm object. Fuzzy search with string is supported.
* `alarm_status` - (Optional, Set: [`String`]) Filter by alarm status. Valid values: ALARM (not resolved), OK (resolved), NO_CONF (expired), NO_DATA (insufficient data). If this parameter is left empty, all will be queried by default.
* `content` - (Optional, String) Fuzzy search by alarm content.
* `end_time` - (Optional, Int) End time, which is the current timestamp and the time when the alarm FirstOccurTime first occurs. An alarm record can be searched only if its FirstOccurTime is earlier than the EndTime.
* `instance_group_ids` - (Optional, Set: [`Int`]) Filter by instance group ID.
* `metric_names` - (Optional, Set: [`String`]) Filter by metric name.
* `monitor_types` - (Optional, Set: [`String`]) Filter by monitor type. Valid values: MT_QCE (Tencent Cloud service monitoring), MT_TAW (application performance monitoring), MT_RUM (frontend performance monitoring), MT_PROBE (cloud automated testing). If this parameter is left empty, all types will be queried by default.
* `namespaces` - (Optional, List) Filter by policy type. Monitoring type and policy type are first-level and second-level filters respectively and both need to be passed in. For example, [{MonitorType: MT_QCE, Namespace: cvm_device}].
* `order` - (Optional, String) Sort by the first occurrence time in descending order by default. Valid values: ASC (ascending), DESC (descending).
* `policy_ids` - (Optional, Set: [`String`]) Search by alarm policy ID list.
* `policy_name` - (Optional, String) Fuzzy search by policy name.
* `project_ids` - (Optional, Set: [`Int`]) Filter by project ID. Valid values: -1 (no project), 0 (default project).
* `receiver_groups` - (Optional, Set: [`Int`]) Search by recipient group.
* `receiver_uids` - (Optional, Set: [`Int`]) Search by recipient.
* `result_output_file` - (Optional, String) Used to save results.
* `start_time` - (Optional, Int) Start time, which is the timestamp one day ago by default and the time when the alarm FirstOccurTime first occurs. An alarm record can be searched only if its FirstOccurTime is later than the StartTime.

The `namespaces` object supports the following:

* `monitor_type` - (Required, String) Monitor type.
* `namespace` - (Required, String) Policy type.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `histories` - Alarm record list.
  * `alarm_id` - Alarm record ID.
  * `alarm_level` - Alarm level.Note: this field may return null, indicating that no valid values can be obtained.
  * `alarm_object` - Alarm object.
  * `alarm_status` - Alarm status. Valid values: ALARM (not resolved), OK (resolved), NO_CONF (expired), NO_DATA (insufficient data).
  * `alarm_type` - Alarm type.
  * `content` - Alarm content.
  * `dimensions` - Dimension information of an instance that triggered alarms.Note: this field may return null, indicating that no valid values can be obtained.
  * `event_id` - Event ID.
  * `first_occur_time` - Timestamp of the first occurrence.
  * `instance_group` - Instance group of alarm object.
    * `id` - Instance group ID.
    * `name` - Instance group name.
  * `last_occur_time` - Timestamp of the last occurrence.
  * `metrics_info` - Metric informationNote: this field may return null, indicating that no valid values can be obtained.
    * `description` - Metric display name.
    * `metric_name` - Metric name.
    * `period` - Statistical period.
    * `qce_namespace` - Namespace used to query data by Tencent Cloud service monitoring type.
    * `value` - Value triggering alarm.
  * `monitor_type` - Monitor type.
  * `namespace` - Policy type.
  * `notice_ways` - Alarm channel list. Valid values: SMS (SMS), EMAIL (email), CALL (phone), WECHAT (WeChat).
  * `origin_id` - Alarm policy ID, which can be used when you call APIs (BindingPolicyObject, UnBindingAllPolicyObject, UnBindingPolicyObject) to bind/unbind instances or instance groups to/from an alarm policy.
  * `policy_exists` - Whether the policy exists. Valid values: 0 (no), 1 (yes).
  * `policy_id` - Alarm policy ID.
  * `policy_name` - Policy name.
  * `project_id` - Project ID.
  * `project_name` - Project name.
  * `receiver_groups` - Recipient group list.
  * `receiver_uids` - Recipient list.
  * `region` - Region.
  * `vpc` - VPC of alarm object for basic product alarm.


