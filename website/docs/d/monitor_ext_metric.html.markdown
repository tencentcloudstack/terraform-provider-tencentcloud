---
subcategory: "Cloud Monitor(Monitor)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_monitor_ext_metric"
sidebar_current: "docs-tencentcloud-data_source-monitor_ext_metric"
description: |-
  Use this data source to query Cloud Monitor external metrics of an external namespace.
---

# tencentcloud_monitor_ext_metric

Use this data source to query Cloud Monitor external metrics of an external namespace.

## Example Usage

```hcl
data "tencentcloud_monitor_ext_metric" "example" {
  ext_namespace = "QCE/CVM"
}
```

## Argument Reference

The following arguments are supported:

* `ext_namespace` - (Required, String) External namespace to query.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

  * `cn_meaning` - Chinese meaning of the metric. This field may return null, which means no valid value is available.
  * `dimension_flag` - Whether external dimensions are configured.
  * `en_meaning` - English meaning of the metric. This field may return null, which means no valid value is available.
  * `metric_c_name` - Chinese metric name. This field may return null, which means no valid value is available.
  * `metric_name` - Metric name. This field may return null, which means no valid value is available.
  * `unit` - Metric unit. This field may return null, which means no valid value is available.
* `ext_metric_list` - List of external metrics under the namespace.


