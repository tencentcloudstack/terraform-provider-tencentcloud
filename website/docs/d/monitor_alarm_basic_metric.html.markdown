---
subcategory: "Cloud Monitor(Monitor)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_monitor_alarm_basic_metric"
sidebar_current: "docs-tencentcloud-datasource-monitor_alarm_basic_metric"
description: |-
  Use this data source to query detailed information of monitor basic_metric
---

# tencentcloud_monitor_alarm_basic_metric

Use this data source to query detailed information of monitor basic_metric

## Example Usage

```hcl
data "tencentcloud_monitor_alarm_basic_metric" "alarm_metric" {
  namespace   = "qce/cvm"
  metric_name = "WanOuttraffic"
  dimensions  = ["uuid"]
}
```

## Argument Reference

The following arguments are supported:

* `namespace` - (Required, String) The business namespace is different for each cloud product. To obtain the business namespace, please go to the product monitoring indicator documents, such as the namespace of the cloud server, which can be found in [Cloud Server Monitoring Indicators](https://cloud.tencent.com/document/product/248/6843 ).
* `dimensions` - (Optional, Set: [`String`]) Optional parameters, filtered by dimension.
* `metric_name` - (Optional, String) Indicator names are different for each cloud product. To obtain indicator names, please go to the monitoring indicator documents of each product, such as the indicator names of cloud servers, which can be found in [Cloud Server Monitoring Indicators]( https://cloud.tencent.com/document/product/248/6843).
* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `metric_set` - List of indicator descriptions obtained from query.
  * `dimensions` - Dimension description information.
    * `dimensions` - Dimension name array.
  * `meaning` - Explanation of the meaning of statistical indicators.
    * `en` - Explanation of indicators in English.
    * `zh` - Chinese interpretation of indicators.
  * `metric_c_name` - Indicator Chinese Name.
  * `metric_e_name` - Indicator English name.
  * `metric_name` - Indicator Name.
  * `namespace` - Namespaces, each cloud product will have a namespace.
  * `period` - The statistical period supported by the indicator, in seconds, such as 60, 300.
  * `periods` - Indicator method within the statistical cycle.
    * `period` - Cycle.
    * `stat_type` - Statistical methods.
  * `unit_cname` - Units used for indicators.
  * `unit` - Units used for indicators.


