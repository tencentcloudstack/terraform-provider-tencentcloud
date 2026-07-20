---
subcategory: "TencentCloud EdgeOne(TEO)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_teo_billing_data"
sidebar_current: "docs-tencentcloud-datasource-teo_billing_data"
description: |-
  Use this data source to query TEO billing data, such as traffic, bandwidth, and request metrics.
---

# tencentcloud_teo_billing_data

Use this data source to query TEO billing data, such as traffic, bandwidth, and request metrics.

## Example Usage

```hcl
data "tencentcloud_teo_billing_data" "example" {
  start_time  = "2025-01-01T00:00:00+08:00"
  end_time    = "2025-01-02T00:00:00+08:00"
  zone_ids    = ["zone-2qtuhspy7cr6"]
  metric_name = "acc_flux"
  interval    = "hour"
  filters {
    type  = "host"
    value = "test.example.com"
  }
  group_by = ["zone-id", "host"]
}
```

## Argument Reference

The following arguments are supported:

* `end_time` - (Required, String) End time. The query time range (`EndTime` - `StartTime`) must be less than or equal to 31 days.
* `metric_name` - (Required, String) Billing metric name, such as `acc_flux`, `acc_bandwidth`, etc.
* `start_time` - (Required, String) Start time.
* `zone_ids` - (Required, List: [`String`]) Site ID collection, up to 100 site IDs. Use `*` to query account-level data of all sites under the current TencentCloud main account.
* `filters` - (Optional, List) Filter conditions. Each item contains `type` and `value`. Valid `type` values: `host`, `proxy-id`, `region-id`.
* `group_by` - (Optional, List: [`String`]) Grouping and aggregation dimensions, up to two dimensions. Valid values: `zone-id`, `host`, `proxy-id`, `region-id`.
* `interval` - (Optional, String) Query time granularity. Valid values: `5min`, `hour`, `day`.
* `result_output_file` - (Optional, String) Used to save results.

The `filters` object supports the following:

* `type` - (Required, String) Parameter name. Valid values: `host`, `proxy-id`, `region-id`.
* `value` - (Required, String) Parameter value.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `data` - Billing data point list.
  * `host` - Domain name to which the data point belongs.
  * `proxy_id` - Layer 4 proxy instance ID to which the data point belongs.
  * `region_id` - Billing region ID to which the data point belongs.
  * `time` - Data timestamp.
  * `value` - Numeric value.
  * `zone_id` - Site ID to which the data point belongs.


