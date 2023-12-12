Use this data source to query detailed information of monitor alarm_conditions_template

Example Usage

```hcl
data "tencentcloud_monitor_alarm_conditions_template" "alarm_conditions_template" {
  module             = "monitor"
  view_name          = "cvm_device"
  group_name         = "keep-template"
  group_id           = "7803070"
  update_time_order  = "desc=descending"
  policy_count_order = "asc=ascending"
}
```