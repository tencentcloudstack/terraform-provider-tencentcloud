---
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_monitor_data"
sidebar_current: "docs-tencentcloud-datasource-monitor_data"
description: |-
  Use this data source to query monitor data. for complex queries, use (https://github.com/tencentyun/tencentcloud-exporter)
---

# tencentcloud_monitor_data

Use this data source to query monitor data. for complex queries, use (https://github.com/tencentyun/tencentcloud-exporter)

## Example Usage

```hcl
data "tencentcloud_instances" "instances" {
}

#cvm
data "tencentcloud_monitor_data" "cvm_monitor_data" {
  namespace   = "QCE/CVM"
  metric_name = "CPUUsage"
  dimensions {
    name  = "InstanceId"
    value = data.tencentcloud_instances.instances.instance_list[0].instance_id
  }
  period     = 300
  start_time = "2020-04-28T18:45:00+08:00"
  end_time   = "2020-04-28T19:00:00+08:00"
}

#cos
data "tencentcloud_monitor_data" "cos_monitor_data" {
  namespace   = "QCE/COS"
  metric_name = "InternetTraffic"
  dimensions {
    name  = "appid"
    value = "1258798060"
  }
  dimensions {
    name  = "bucket"
    value = "test-1258798060"
  }

  period     = 300
  start_time = "2020-04-28T18:30:00+08:00"
  end_time   = "2020-04-28T19:00:00+08:00"
}
```

## Argument Reference

The following arguments are supported:

* `dimensions` - (Required) Dimensional composition of instance objects.
* `end_time` - (Required) End time for this query, eg:`2018-09-22T20:00:00+08:00`.
* `metric_name` - (Required) Metric name, please refer to the documentation of monitor interface of each product.
* `namespace` - (Required) Namespace of each cloud product in monitor system, refer to `data.tencentcloud_monitor_product_namespace`.
* `start_time` - (Required) Start time for this query, eg:`2018-09-22T19:51:23+08:00`.
* `period` - (Optional) Statistical period.
* `result_output_file` - (Optional) Used to store results.

The `dimensions` object supports the following:

* `name` - (Required) Instance dimension name, eg: `InstanceId` for cvm.
* `value` - (Required) Instance dimension value, eg: `ins-j0hk02zo` for cvm.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `list` - A list data point. Each element contains the following attributes:
  * `timestamp` - Statistical timestamp.
  * `value` - Statistical value.


