---
subcategory: "Cloud Monitor(Monitor)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_monitor_alarm_metric"
sidebar_current: "docs-tencentcloud-datasource-monitor_alarm_metric"
description: |-
  Use this data source to query detailed information of monitor alarm_metric
---

# tencentcloud_monitor_alarm_metric

Use this data source to query detailed information of monitor alarm_metric

## Example Usage

```hcl
data "tencentcloud_monitor_alarm_metric" "alarm_metric" {
  module       = "monitor"
  monitor_type = "Monitoring"
  namespace    = "cvm_device"
}
```

## Argument Reference

The following arguments are supported:

* `module` - (Required, String) Fixed value, as `monitor`.
* `monitor_type` - (Required, String) Monitoring Type Filter MT_QCE=Cloud Product Monitoring.
* `namespace` - (Required, String) Alarm policy type, obtained from DescribeAllNamespaces, such as cvm_device.
* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `metrics` - Alarm indicator list.
  * `description` - Indicator display name.
  * `dimensions` - Dimension List.
  * `is_advanced` - Is it a high-level indicator. 1 Yes 0 No.
  * `is_open` - Is the advanced indicator activated. 1 Yes 0 No.
  * `max` - Maximum value.
  * `metric_config` - Indicator configuration.
    * `continue_period` - Number of allowed duration cycles for configuration.
    * `operator` - Allowed Operators.
    * `period` - The data period allowed for configuration, in seconds.
  * `metric_name` - Indicator Name.
  * `min` - Minimum value.
  * `namespace` - Alarm strategy type.
  * `operators` - Matching operator.
    * `id` - Operator identification.
    * `name` - Operator Display Name.
  * `periods` - Indicator trigger.
  * `product_id` - Integration Center Product ID.
  * `unit` - Unit.


