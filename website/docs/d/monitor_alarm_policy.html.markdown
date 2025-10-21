---
subcategory: "Cloud Monitor(Monitor)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_monitor_alarm_policy"
sidebar_current: "docs-tencentcloud-datasource-monitor_alarm_policy"
description: |-
  Use this data source to query detailed information of monitor alarm_policy
---

# tencentcloud_monitor_alarm_policy

Use this data source to query detailed information of monitor alarm_policy

## Example Usage

```hcl
data "tencentcloud_monitor_alarm_policy" "alarm_policy" {
  module        = "monitor"
  policy_name   = "terraform"
  monitor_types = ["MT_QCE"]
  namespaces    = ["cvm_device"]
  project_ids   = [0]
  notice_ids    = ["notice-f2svbu3w"]
  rule_types    = ["STATIC"]
  enable        = [1]
}
```

## Argument Reference

The following arguments are supported:

* `module` - (Required, String) Value fixed at monitor.
* `dimensions` - (Optional, String) The alarm object list, which is a JSON string. The outer array corresponds to multiple instances, and the inner array is the dimension of an object.For example, 'CVM - Basic Monitor' can be written as: [ {Dimensions: {unInstanceId: ins-qr8d555g}}, {Dimensions: {unInstanceId: ins-qr8d555h}} ]You can also refer to the 'Example 2' below.For more information on the parameter samples of different Tencent Cloud services, see [Product Policy Type and Dimension Information](https://www.tencentcloud.com/document/product/248/39565?has_map=1).Note: If 1 is passed in for NeedCorrespondence, the relationship between a policy and an instance needs to be returned. You can pass in up to 20 alarm object dimensions to avoid request timeout.
* `enable` - (Optional, Set: [`Int`]) Filter by alarm status. Valid values: [1]: enabled; [0]: disabled; [0, 1]: all.
* `field` - (Optional, String) Sort by field. For example, to sort by the last modification time, use Field: UpdateTime.
* `instance_group_id` - (Optional, Int) Instance group ID.
* `monitor_types` - (Optional, Set: [`String`]) Filter by monitor type. Valid values: MT_QCE (Tencent Cloud service monitoring). If this parameter is left empty, all will be queried by default.
* `namespaces` - (Optional, Set: [`String`]) Filter by namespace. For the values of different policy types, please see:[Poicy Type List](https://www.tencentcloud.com/document/product/248/39565?has_map=1).
* `need_correspondence` - (Optional, Int) Whether the relationship between a policy and the input parameter filter dimension is required. 1: Yes. 0: No. Default value: 0.
* `not_bind_all` - (Optional, Int) Whether the returned result needs to filter policies associated with all objects. Valid values: 1 (Yes), 0 (No).
* `not_binding_notice_rule` - (Optional, Int) If 1 is passed in, alarm policies with no notification rules configured are queried. If it is left empty or other values are passed in, all alarm policies are queried.
* `not_instance_group` - (Optional, Int) Whether the returned result needs to filter policies associated with instance groups. Valid values: 1 (Yes), 0 (No).
* `notice_ids` - (Optional, Set: [`String`]) List of the notification template IDs, which can be obtained by querying the notification template list.It can be queried with the API [DescribeAlarmNotices](https://www.tencentcloud.com/document/product/248/39300).
* `one_click_policy_type` - (Optional, Set: [`String`]) Filter by quick alarm policy. If this parameter is left empty, all policies are displayed. ONECLICK: Display quick alarm policies; NOT_ONECLICK: Display non-quick alarm policies.
* `order` - (Optional, String) Sort order. Valid values: ASC (ascending), DESC (descending).
* `policy_name` - (Optional, String) Fuzzy search by policy name.
* `policy_type` - (Optional, Set: [`String`]) Filter by default policy. Valid values: DEFAULT (display default policy), NOT_DEFAULT (display non-default policies). If this parameter is left empty, all policies will be displayed.
* `project_ids` - (Optional, Set: [`Int`]) ID array of the policy project, which can be viewed on the following page: [Project Management](https://console.tencentcloud.com/project).
* `prom_ins_id` - (Optional, String) ID of the TencentCloud Managed Service for Prometheus instance, which is used for customizing a metric policy.
* `receiver_groups` - (Optional, Set: [`Int`]) Search by recipient group. You can get the user group list with the API [ListGroups](https://www.tencentcloud.com/document/product/598/34589?from_cn_redirect=1) in 'Cloud Access Management' or query the user group list where a sub-user is in with the API [ListGroupsForUser](https://www.tencentcloud.com/document/product/598/34588?from_cn_redirect=1). The GroupId field in the returned result should be entered here.
* `receiver_on_call_form_ids` - (Optional, Set: [`String`]) Search by schedule.
* `receiver_uids` - (Optional, Set: [`Int`]) Search by recipient. You can get the user list with the API [ListUsers](https://www.tencentcloud.com/document/product/598/34587?from_cn_redirect=1) in 'Cloud Access Management' or query the sub-user information with the API [GetUser](https://www.tencentcloud.com/document/product/598/34590?from_cn_redirect=1). The Uid field in the returned result should be entered here.
* `result_output_file` - (Optional, String) Used to save results.
* `rule_types` - (Optional, Set: [`String`]) Filter by trigger condition. Valid values: STATIC (display policies with static threshold), DYNAMIC (display policies with dynamic threshold). If this parameter is left empty, all policies will be displayed.
* `trigger_tasks` - (Optional, List) Filter alarm policy by triggered task (such as auto scaling task). Up to 10 tasks can be specified.

The `trigger_tasks` object supports the following:

* `task_config` - (Required, String) Configuration information in JSON format, such as {Key1:Value1,Key2:Value2}Note: this field may return null, indicating that no valid values can be obtained.
* `type` - (Required, String) Triggered task type. Valid value: AS (auto scaling)Note: this field may return null, indicating that no valid values can be obtained.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `policies` - Policy array.
  * `advanced_metric_number` - The number of advanced metrics.Note: This field may return null, indicating that no valid values can be obtained.
  * `can_set_default` - Whether the default policy can be set. Valid values: 1 (yes), 0 (no)Note: this field may return null, indicating that no valid values can be obtained.
  * `condition_template_id` - Trigger condition template IDNote: this field may return null, indicating that no valid values can be obtained.
  * `condition` - Metric trigger conditionNote: this field may return null, indicating that no valid values can be obtained.
    * `complex_expression` - The judgment expression of composite alarm trigger conditions, which is valid when the value of IsUnionRule is 2. This parameter is used to determine that an alarm condition is met only when the expression values are True for multiple trigger conditionsNote: This field may return null, indicating that no valid values can be obtained.
    * `is_union_rule` - Judgment condition of an alarm trigger condition (0: Any; 1: All; 2: Composite). When the value is set to 2 (i.e., composite trigger conditions), this parameter should be used together with ComplexExpression.Note: This field may return null, indicating that no valid values can be obtained.
    * `rules` - Alarm trigger condition listNote: this field may return null, indicating that no valid values can be obtained.
      * `continue_period` - Number of periods. 1: continue for one period; 2: continue for two periods; and so on. The valid values can be queried via [DescribeAlarmMetrics](https://www.tencentcloud.com/document/product/248/39322)Note: this field may return null, indicating that no valid value is obtained.
      * `description` - Metric display name, which is used in the output parameteNote: this field may return null, indicating that no valid values can be obtained.
      * `filter` - Filter condition for one single trigger rulNote: this field may return null, indicating that no valid values can be obtained.
        * `dimensions` - JSON string generated by serializing the AlarmPolicyDimension two-dimensional array. The one-dimensional arrays are in OR relationship, and the elements in a one-dimensional array are in AND relationshiNote: this field may return null, indicating that no valid values can be obtained.
        * `type` - Filter condition type. Valid values: DIMENSION (uses dimensions for filteringNote: this field may return null, indicating that no valid values can be obtained.
      * `hierarchical_value` - The configuration of alarm level thresholNote: This field may return null, indicating that no valid values can be obtained.
        * `remind` - Threshold for the Remind leveNote: This field may return null, indicating that no valid values can be obtained.
        * `serious` - Threshold for the Serious leveNote: This field may return null, indicating that no valid values can be obtained.
        * `warn` - Threshold for the Warn leveNote: This field may return null, indicating that no valid values can be obtained.
      * `is_advanced` - Whether it is an advanced metric. 0: No; 1: YesNote: This field may return null, indicating that no valid values can be obtained.
      * `is_open` - Whether the advanced metric feature is enabled. 0: No; 1: YesNote: This field may return null, indicating that no valid values can be obtained.
      * `is_power_notice` - Whether the alarm frequency increases exponentially. Valid values: 0 (no), 1 (yesNote: this field may return null, indicating that no valid values can be obtained.
      * `metric_name` - Metric name or event name. The supported metrics can be queried via [DescribeAlarmMetrics](https://www.tencentcloud.com/document/product/248/39322) and the supported events via [DescribeAlarmEvents](https://www.tencentcloud.com/document/product/248/39324).Note: this field may return null, indicating that no valid value is obtained.
      * `notice_frequency` - Alarm interval in seconds. Valid values: 0 (do not repeat), 300 (alarm once every 5 minutes), 600 (alarm once every 10 minutes), 900 (alarm once every 15 minutes), 1800 (alarm once every 30 minutes), 3600 (alarm once every hour), 7200 (alarm once every 2 hours), 10800 (alarm once every 3 hours), 21600 (alarm once every 6 hours), 43200 (alarm once every 12 hours), 86400 (alarm once every dayNote: this field may return null, indicating that no valid values can be obtained.
      * `operator` - Operatorintelligent = intelligent detection without thresholdeq = equal toge = greater than or equal togt = greater thanle = less than or equal tolt = less thanne = not equal today_increase = day-on-day increaseday_decrease = day-on-day decreaseday_wave = day-on-day fluctuationweek_increase = week-on-week increaseweek_decrease = week-on-week decreaseweek_wave = week-on-week fluctuationcycle_increase = cyclical increasecycle_decrease = cyclical decreasecycle_wave = cyclical fluctuationre = regex matchThe valid values can be queried via [DescribeAlarmMetrics](https://www.tencentcloud.com/document/product/248/39322)Note: this field may return null, indicating that no valid value is obtained.
      * `period` - Statistical period in seconds. The valid values can be queried via DescribeAlarmMetrics.Note: this field may return null, indicating that no valid value is obtained.
      * `product_id` - Integration center product IDNote: This field may return null, indicating that no valid values can be obtained.
      * `rule_type` - Trigger condition type. STATIC: static threshold; dynamic: dynamic threshold. If you do not specify this parameter when creating or editing a policy, STATIC is used by defaultNote: this field may return null, indicating that no valid value is obtained.
      * `unit` - Unit, which is used in the output parameteNote: this field may return null, indicating that no valid values can be obtained.
      * `value_max` - Maximum valuNote: This field may return null, indicating that no valid values can be obtained.
      * `value_min` - Minimum valuNote: This field may return null, indicating that no valid values can be obtained.
      * `value` - Threshold. The valid value range can be queried via [DescribeAlarmMetrics](https://www.tencentcloud.com/document/product/248/39322)Note: this field may return null, indicating that no valid value is obtained.
  * `conditions_temp` - Template policy groupNote: this field may return null, indicating that no valid values can be obtained.
    * `condition` - Metric trigger conditionNote: this field may return null, indicating that no valid values can be obtained.
      * `complex_expression` - The judgment expression of composite alarm trigger conditions, which is valid when the value of IsUnionRule is 2. This parameter is used to determine that an alarm condition is met only when the expression values are True for multiple trigger conditions.Note: This field may return null, indicating that no valid values can be obtained.
      * `is_union_rule` - Judgment condition of an alarm trigger condition (0: Any; 1: All; 2: Composite). When the value is set to 2 (i.e., composite trigger conditions), this parameter should be used together with ComplexExpression.Note: This field may return null, indicating that no valid values can be obtained.
      * `rules` - Alarm trigger condition listNote: this field may return null, indicating that no valid values can be obtained.
        * `continue_period` - Number of periods. 1: continue for one period; 2: continue for two periods; and so on. The valid values can be queried via DescribeAlarmMetrics.Note: this field may return null, indicating that no valid value is obtained.
        * `description` - Metric display name, which is used in the output parameterNote: this field may return null, indicating that no valid values can be obtained.
        * `filter` - Filter condition for one single trigger ruleNote: this field may return null, indicating that no valid values can be obtained.
          * `dimensions` - JSON string generated by serializing the AlarmPolicyDimension two-dimensional array. The one-dimensional arrays are in OR relationship, and the elements in a one-dimensional array are in AND relationshipNote: this field may return null, indicating that no valid values can be obtained.
          * `type` - Filter condition type. Valid values: DIMENSION (uses dimensions for filtering)Note: this field may return null, indicating that no valid values can be obtained.
        * `hierarchical_value` - The configuration of alarm level thresholdNote: This field may return null, indicating that no valid values can be obtained.
          * `remind` - Threshold for the Remind levelNote: This field may return null, indicating that no valid values can be obtained.
          * `serious` - Threshold for the Serious levelNote: This field may return null, indicating that no valid values can be obtained.
          * `warn` - Threshold for the Warn levelNote: This field may return null, indicating that no valid values can be obtained.
        * `is_advanced` - Whether it is an advanced metric. 0: No; 1: Yes.Note: This field may return null, indicating that no valid values can be obtained.
        * `is_open` - Whether the advanced metric feature is enabled. 0: No; 1: Yes.Note: This field may return null, indicating that no valid values can be obtained.
        * `is_power_notice` - Whether the alarm frequency increases exponentially. Valid values: 0 (no), 1 (yes)Note: this field may return null, indicating that no valid values can be obtained.
        * `metric_name` - Metric name or event name. The supported metrics can be queried via DescribeAlarmMetrics and the supported events via DescribeAlarmEvents.Note: this field may return null, indicating that no valid value is obtained.
        * `notice_frequency` - Alarm interval in seconds. Valid values: 0 (do not repeat), 300 (alarm once every 5 minutes), 600 (alarm once every 10 minutes), 900 (alarm once every 15 minutes), 1800 (alarm once every 30 minutes), 3600 (alarm once every hour), 7200 (alarm once every 2 hours), 10800 (alarm once every 3 hours), 21600 (alarm once every 6 hours), 43200 (alarm once every 12 hours), 86400 (alarm once every day)Note: this field may return null, indicating that no valid values can be obtained.
        * `operator` - Operatorintelligent = intelligent detection without thresholdeq = equal toge = greater than or equal togt = greater thanle = less than or equal tolt = less thanne = not equal today_increase = day-on-day increaseday_decrease = day-on-day decreaseday_wave = day-on-day fluctuationweek_increase = week-on-week increaseweek_decrease = week-on-week decreaseweek_wave = week-on-week fluctuationcycle_increase = cyclical increasecycle_decrease = cyclical decreasecycle_wave = cyclical fluctuationre = regex matchThe valid values can be queried via DescribeAlarmMetrics.Note: this field may return null, indicating that no valid value is obtained.
        * `period` - Statistical period in seconds. The valid values can be queried via DescribeAlarmMetrics.Note: this field may return null, indicating that no valid value is obtained.
        * `product_id` - Integration center product ID.Note: This field may return null, indicating that no valid values can be obtained.
        * `rule_type` - Trigger condition type. STATIC: static threshold; dynamic: dynamic threshold. If you do not specify this parameter when creating or editing a policy, STATIC is used by default.Note: this field may return null, indicating that no valid value is obtained.
        * `unit` - Unit, which is used in the output parameterNote: this field may return null, indicating that no valid values can be obtained.
        * `value_max` - Maximum valueNote: This field may return null, indicating that no valid values can be obtained.
        * `value_min` - Minimum valueNote: This field may return null, indicating that no valid values can be obtained.
        * `value` - Threshold. The valid value range can be queried via DescribeAlarmMetrics.Note: this field may return null, indicating that no valid value is obtained.
    * `event_condition` - Event trigger conditionNote: this field may return null, indicating that no valid values can be obtained.
      * `rules` - Alarm trigger condition listNote: this field may return null, indicating that no valid values can be obtained.
        * `continue_period` - Number of periods. 1: continue for one period; 2: continue for two periods; and so on. The valid values can be queried via DescribeAlarmMetrics.Note: this field may return null, indicating that no valid value is obtained.
        * `description` - Metric display name, which is used in the output parameterNote: this field may return null, indicating that no valid values can be obtained.
        * `filter` - Filter condition for one single trigger ruleNote: this field may return null, indicating that no valid values can be obtained.
          * `dimensions` - JSON string generated by serializing the AlarmPolicyDimension two-dimensional array. The one-dimensional arrays are in OR relationship, and the elements in a one-dimensional array are in AND relationshipNote: this field may return null, indicating that no valid values can be obtained.
          * `type` - Filter condition type. Valid values: DIMENSION (uses dimensions for filtering)Note: this field may return null, indicating that no valid values can be obtained.
        * `hierarchical_value` - The configuration of alarm level thresholdNote: This field may return null, indicating that no valid values can be obtained.
          * `remind` - Threshold for the Remind levelNote: This field may return null, indicating that no valid values can be obtained.
          * `serious` - Threshold for the Serious levelNote: This field may return null, indicating that no valid values can be obtained.
          * `warn` - Threshold for the Warn levelNote: This field may return null, indicating that no valid values can be obtained.
        * `is_advanced` - Whether it is an advanced metric. 0: No; 1: Yes.Note: This field may return null, indicating that no valid values can be obtained.
        * `is_open` - Whether the advanced metric feature is enabled. 0: No; 1: Yes.Note: This field may return null, indicating that no valid values can be obtained.
        * `is_power_notice` - Whether the alarm frequency increases exponentially. Valid values: 0 (no), 1 (yes)Note: this field may return null, indicating that no valid values can be obtained.
        * `metric_name` - Metric name or event name. The supported metrics can be queried via DescribeAlarmMetrics and the supported events via DescribeAlarmEvents.Note: this field may return null, indicating that no valid value is obtained.
        * `notice_frequency` - Alarm interval in seconds. Valid values: 0 (do not repeat), 300 (alarm once every 5 minutes), 600 (alarm once every 10 minutes), 900 (alarm once every 15 minutes), 1800 (alarm once every 30 minutes), 3600 (alarm once every hour), 7200 (alarm once every 2 hours), 10800 (alarm once every 3 hours), 21600 (alarm once every 6 hours), 43200 (alarm once every 12 hours), 86400 (alarm once every day)Note: this field may return null, indicating that no valid values can be obtained.
        * `operator` - Operatorintelligent = intelligent detection without thresholdeq = equal toge = greater than or equal togt = greater thanle = less than or equal tolt = less thanne = not equal today_increase = day-on-day increaseday_decrease = day-on-day decreaseday_wave = day-on-day fluctuationweek_increase = week-on-week increaseweek_decrease = week-on-week decreaseweek_wave = week-on-week fluctuationcycle_increase = cyclical increasecycle_decrease = cyclical decreasecycle_wave = cyclical fluctuationre = regex matchThe valid values can be queried via DescribeAlarmMetrics.Note: this field may return null, indicating that no valid value is obtained.
        * `period` - Statistical period in seconds. The valid values can be queried via DescribeAlarmMetrics.Note: this field may return null, indicating that no valid value is obtained.
        * `product_id` - Integration center product ID.Note: This field may return null, indicating that no valid values can be obtained.
        * `rule_type` - Trigger condition type. STATIC: static threshold; dynamic: dynamic threshold. If you do not specify this parameter when creating or editing a policy, STATIC is used by default.Note: this field may return null, indicating that no valid value is obtained.
        * `unit` - Unit, which is used in the output parameterNote: this field may return null, indicating that no valid values can be obtained.
        * `value_max` - Maximum valueNote: This field may return null, indicating that no valid values can be obtained.
        * `value_min` - Minimum valueNote: This field may return null, indicating that no valid values can be obtained.
        * `value` - Threshold. The valid value range can be queried via DescribeAlarmMetrics.Note: this field may return null, indicating that no valid value is obtained.
    * `template_name` - Template nameNote: u200dThis field may return null, indicating that no valid values can be obtained.
  * `enable` - Status. Valid values: 0 (disabled), 1 (enabled)Note: this field may return null, indicating that no valid values can be obtained.
  * `event_condition` - Event trigger conditioNote: this field may return null, indicating that no valid values can be obtained.
    * `rules` - Alarm trigger condition lisNote: this field may return null, indicating that no valid values can be obtained.
      * `continue_period` - Number of periods. 1: continue for one period; 2: continue for two periods; and so on. The valid values can be queried via DescribeAlarmMetrics.Note: this field may return null, indicating that no valid value is obtained.
      * `description` - Metric display name, which is used in the output parameterNote: this field may return null, indicating that no valid values can be obtained.
      * `filter` - Filter condition for one single trigger ruleNote: this field may return null, indicating that no valid values can be obtained.
        * `dimensions` - JSON string generated by serializing the AlarmPolicyDimension two-dimensional array. The one-dimensional arrays are in OR relationship, and the elements in a one-dimensional array are in AND relationshipNote: this field may return null, indicating that no valid values can be obtained.
        * `type` - Filter condition type. Valid values: DIMENSION (uses dimensions for filtering)Note: this field may return null, indicating that no valid values can be obtained.
      * `hierarchical_value` - The configuration of alarm level thresholdNote: This field may return null, indicating that no valid values can be obtained.
        * `remind` - Threshold for the Remind levelNote: This field may return null, indicating that no valid values can be obtained.
        * `serious` - Threshold for the Serious levelNote: This field may return null, indicating that no valid values can be obtained.
        * `warn` - Threshold for the Warn levelNote: This field may return null, indicating that no valid values can be obtained.
      * `is_advanced` - Whether it is an advanced metric. 0: No; 1: Yes.Note: This field may return null, indicating that no valid values can be obtained.
      * `is_open` - Whether the advanced metric feature is enabled. 0: No; 1: Yes.Note: This field may return null, indicating that no valid values can be obtained.
      * `is_power_notice` - Whether the alarm frequency increases exponentially. Valid values: 0 (no), 1 (yes)Note: this field may return null, indicating that no valid values can be obtained.
      * `metric_name` - Metric name or event name. The supported metrics can be queried via DescribeAlarmMetrics and the supported events via DescribeAlarmEventsNote: this field may return null, indicating that no valid value is obtained.
      * `notice_frequency` - Alarm interval in seconds. Valid values: 0 (do not repeat), 300 (alarm once every 5 minutes), 600 (alarm once every 10 minutes), 900 (alarm once every 15 minutes), 1800 (alarm once every 30 minutes), 3600 (alarm once every hour), 7200 (alarm once every 2 hours), 10800 (alarm once every 3 hours), 21600 (alarm once every 6 hours), 43200 (alarm once every 12 hours), 86400 (alarm once every day)Note: this field may return null, indicating that no valid values can be obtained.
      * `operator` - Statistical period in seconds. The valid values can be queried via DescribeAlarmMetrics.Note: this field may return null, indicating that no valid value is obtained.Operator	String	No	Operatorintelligent = intelligent detection without thresholdeq = equal toge = greater than or equal togt = greater thanle = less than or equal tolt = less thanne = not equal today_increase = day-on-day increaseday_decrease = day-on-day decreaseday_wave = day-on-day fluctuationweek_increase = week-on-week increaseweek_decrease = week-on-week decreaseweek_wave = week-on-week fluctuationcycle_increase = cyclical increasecycle_decrease = cyclical decreasecycle_wave = cyclical fluctuationre = regex matchThe valid values can be queried via DescribeAlarmMetrics.Note: this field may return null, indicating that no valid value is obtained.
      * `period` - Statistical period in seconds. The valid values can be queried via DescribeAlarmMetricsNote: this field may return null, indicating that no valid value is obtained.
      * `product_id` - Integration center product ID.Note: This field may return null, indicating that no valid values can be obtained.
      * `rule_type` - Trigger condition type. STATIC: static threshold; dynamic: dynamic threshold. If you do not specify this parameter when creating or editing a policy, STATIC is used by default.Note: this field may return null, indicating that no valid value is obtained.
      * `unit` - Unit, which is used in the output parameterNote: this field may return null, indicating that no valid values can be obtained.
      * `value_max` - Maximum valueNote: This field may return null, indicating that no valid values can be obtained.
      * `value_min` - Minimum valueNote: This field may return null, indicating that no valid values can be obtained.
      * `value` - Threshold. The valid value range can be queried via DescribeAlarmMetrics.Note: this field may return null, indicating that no valid value is obtained.
  * `filter_dimensions_param` - Information on the filter dimension associated with a policy.Note: This field may return null, indicating that no valid values can be obtained.
  * `insert_time` - Creation timeNote: this field may return null, indicating that no valid values can be obtained.
  * `instance_group_id` - Instance group IDNote: this field may return null, indicating that no valid values can be obtained.
  * `instance_group_name` - Instance group nameNote: this field may return null, indicating that no valid values can be obtained.
  * `instance_sum` - Total number of instances in instance groupNote: this field may return null, indicating that no valid values can be obtained.
  * `is_bind_all` - Whether the policy is associated with all objectsNote: This field may return null, indicating that no valid values can be obtained.
  * `is_default` - Whether it is the default policy. Valid values: 1 (yes), 0 (no)Note: this field may return null, indicating that no valid values can be obtained.
  * `is_one_click` - Whether it is a quick alarm policy.Note: This field may return null, indicating that no valid values can be obtained.
  * `last_edit_uin` - Uin of the last modifying userNote: this field may return null, indicating that no valid values can be obtained.
  * `monitor_type` - Monitor type. Valid values: MT_QCE (Tencent Cloud service monitoring)Note: this field may return null, indicating that no valid values can be obtained.
  * `namespace_show_name` - Namespace display nameNote: this field may return null, indicating that no valid values can be obtained.
  * `namespace` - Alarm policy typeNote: this field may return null, indicating that no valid values can be obtained.
  * `notice_ids` - Notification rule ID listNote: this field may return null, indicating that no valid values can be obtained.
  * `notices` - Notification rule listNote: this field may return null, indicating that no valid values can be obtained.
    * `amp_consumer_id` - Backend AMP consumer ID.Note: This field may return null, indicating that no valid values can be obtained.
    * `cls_notices` - Channel to push alarm notifications to CLS.Note: This field may return null, indicating that no valid values can be obtained.
      * `enable` - Status. Valid values: 0 (disabled), 1 (enabled). Default value: 1 (enabled). This parameter can be left empty.
      * `log_set_id` - Logset ID.
      * `region` - Region.
      * `topic_id` - Topic ID.
    * `id` - Alarm notification template IDNote: this field may return null, indicating that no valid values can be obtained.
    * `is_preset` - Whether it is the system default notification template. Valid values: 0 (no), 1 (yes)Note: this field may return null, indicating that no valid values can be obtained.
    * `name` - Alarm notification template nameNote: this field may return null, indicating that no valid values can be obtained.
    * `notice_language` - Notification language. Valid values: zh-CN (Chinese), en-US (English)Note: this field may return null, indicating that no valid values can be obtained.
    * `notice_type` - Alarm notification type. Valid values: ALARM (for unresolved alarms), OK (for resolved alarms), ALL (for all alarms)Note: this field may return null, indicating that no valid values can be obtained.
    * `policy_ids` - List of IDs of the alarm policies bound to alarm notification templateNote: this field may return null, indicating that no valid values can be obtained.
    * `tags` - Tags bound to a notification templateNote: This field may return null, indicating that no valid values can be obtained.
      * `key` - Tag key.
      * `value` - Tag value.
    * `updated_at` - Last modified timeNote: this field may return null, indicating that no valid values can be obtained.
    * `updated_by` - Last modified byNote: this field may return null, indicating that no valid values can be obtained.
    * `url_notices` - Callback notification listNote: this field may return null, indicating that no valid values can be obtained.
      * `end_time` - End time of the notification in seconds, which is calculated from 00:00:00.Note: this field may return null, indicating that no valid values can be obtained.
      * `is_valid` - Whether verification is passed. Valid values: 0 (no), 1 (yes)Note: this field may return null, indicating that no valid values can be obtained.
      * `start_time` - Start time of the notification in seconds, which is calculated from 00:00:00.Note: this field may return null, indicating that no valid values can be obtained.
      * `url` - Callback URL, which can contain up to 256 charactersNote: this field may return null, indicating that no valid values can be obtained.
      * `validation_code` - Verification codeNote: this field may return null, indicating that no valid values can be obtained.
      * `weekday` - Notification cycle. The values 1-7 indicate Monday to Sunday.Note: This field may return null, indicating that no valid values can be obtained.
    * `user_notices` - User notification listNote: this field may return null, indicating that no valid values can be obtained.
      * `end_time` - Notification end time, which is expressed by the number of seconds since 00:00:00. Value range: 0-86399Note: this field may return null, indicating that no valid values can be obtained.
      * `group_ids` - User group ID listNote: this field may return null, indicating that no valid values can be obtained.
      * `need_phone_arrive_notice` - Whether receipt notification is required. Valid values: 0 (no), 1 (yes)Note: this field may return null, indicating that no valid values can be obtained.
      * `notice_way` - Notification channel list. Valid values: EMAIL (email), SMS (SMS), CALL (phone), WECHAT (WeChat), RTX (WeCom)Note: This field may return null, indicating that no valid values can be obtained.
      * `on_call_form_ids` - List of schedule IDsNote: u200dThis field may return null, indicating that no valid values can be obtained.
      * `phone_call_type` - Dial type. SYNC (simultaneous dial), CIRCLE (polled dial). Default value: CIRCLE.Note: This field may return null, indicating that no valid values can be obtained.
      * `phone_circle_interval` - Polling interval in seconds. Value range: 60-900Note: this field may return null, indicating that no valid values can be obtained.
      * `phone_circle_times` - Number of phone pollings. Value range: 1-5Note: this field may return null, indicating that no valid values can be obtained.
      * `phone_inner_interval` - Call interval in seconds within one polling. Value range: 60-900Note: this field may return null, indicating that no valid values can be obtained.
      * `phone_order` - Phone polling listNote: this field may return null, indicating that no valid values can be obtained.
      * `receiver_type` - Recipient type. Valid values: USER (user), GROUP (user group)Note: this field may return null, indicating that no valid values can be obtained.
      * `start_time` - Notification start time, which is expressed by the number of seconds since 00:00:00. Value range: 0-86399Note: this field may return null, indicating that no valid values can be obtained.
      * `user_ids` - User uid listNote: this field may return null, indicating that no valid values can be obtained.
      * `weekday` - Notification cycle. The values 1-7 indicate Monday to Sunday.Note: This field may return null, indicating that no valid values can be obtained.
  * `one_click_status` - Whether the quick alarm policy is enabled.Note: This field may return null, indicating that no valid values can be obtained.
  * `origin_id` - Policy ID for instance/instance group binding and unbinding APIs (BindingPolicyObject, UnBindingAllPolicyObject, UnBindingPolicyObject)Note: this field may return null, indicating that no valid values can be obtained.
  * `policy_id` - Alarm policy IDNote: this field may return null, indicating that no valid values can be obtained.
  * `policy_name` - Alarm policy nameNote: this field may return null, indicating that no valid values can be obtained.
  * `project_id` - Project ID. Valid values: -1 (no project), 0 (default project)Note: this field may return null, indicating that no valid values can be obtained.
  * `project_name` - Project nameNote: this field may return null, indicating that no valid values can be obtained.
  * `region` - RegionNote: this field may return null, indicating that no valid values can be obtained.
  * `remark` - RemarksNote: this field may return null, indicating that no valid values can be obtained.
  * `rule_type` - Trigger condition type. Valid values: STATIC (static threshold), DYNAMIC (dynamic)Note: this field may return null, indicating that no valid values can be obtained.
  * `tag_instances` - TagNote: This field may return null, indicating that no valid values can be obtained.
    * `binding_status` - Binding status. 2: bound; 1: bindingNote: This field may return null, indicating that no valid values can be obtained.
    * `instance_sum` - Number of instancesNote: This field may return null, indicating that no valid values can be obtained.
    * `key` - Tag keyNote: This field may return null, indicating that no valid values can be obtained.
    * `region_id` - Region IDNote: This field may return null, indicating that no valid values can be obtained.
    * `service_type` - Service type, for example, CVMNote: This field may return null, indicating that no valid values can be obtained.
    * `tag_status` - Tag status. 2: existent; 1: nonexistentNote: This field may return null, indicating that no valid values can be obtained.
    * `value` - Tag valueNote: This field may return null, indicating that no valid values can be obtained.
  * `tags` - Policy tagNote: This field may return null, indicating that no valid values can be obtained.
    * `key` - Tag key.
    * `value` - Tag value.
  * `trigger_tasks` - Triggered task listNote: this field may return null, indicating that no valid values can be obtained.
    * `task_config` - Configuration information in JSON format, such as {Key1:Value1,Key2:Value2}Note: this field may return null, indicating that no valid values can be obtained.
    * `type` - Triggered task type. Valid value: AS (auto scaling)Note: this field may return null, indicating that no valid values can be obtained.
  * `update_time` - Update timeNote: this field may return null, indicating that no valid values can be obtained.
  * `use_sum` - Number of instances bound to policy groupNote: this field may return null, indicating that no valid values can be obtained.


