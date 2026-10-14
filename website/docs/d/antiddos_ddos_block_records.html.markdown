---
subcategory: "Anti-DDoS(antiddos)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_antiddos_ddos_block_records"
sidebar_current: "docs-tencentcloud-datasource-antiddos_ddos_block_records"
description: |-
  Use this data source to query AntiDDoS DDoS block/unblock records and unblock quota info.
---

# tencentcloud_antiddos_ddos_block_records

Use this data source to query AntiDDoS DDoS block/unblock records and unblock quota info.

## Example Usage

```hcl
data "tencentcloud_antiddos_ddos_block_records" "example" {
  start_time = "2026-02-04T11:30:00+08:00"
  end_time   = "2026-03-04T11:30:00+08:00"

  filters {
    name   = "Status"
    values = ["Blocked"]
  }
}
```

## Argument Reference

The following arguments are supported:

* `end_time` - (Required, String) Query end time. (EndTime - StartTime) must be <= 31 days. Parameter format: 2026-03-04T11:30:00+08:00.
* `start_time` - (Required, String) Query start time. Parameter format: 2026-02-04T11:30:00+08:00.
* `filters` - (Optional, List) Filter conditions. The upper limit of Filters.Values is 20. If not filled in, it returns the list of all blocked resources under the current appid.
- Resource: filter by blocked IP or resource six-segment style.
- Status: filter by block status (Blocked/Unblocking/Unblocked).
* `result_output_file` - (Optional, String) Used to save results.

The `filters` object supports the following:

* `name` - (Required, String) Filter field name. Supported: Resource, Status.
* `values` - (Required, Set) Filter field value.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `block_records` - Block/unblock records.
  * `block_time` - Block time.
  * `resource` - Blocked resource, public network IP.
  * `status` - Block/unblock status. Enum: Blocked, Unblocking, Unblocked.
* `unblock_quota_info` - Unblock quota info.
  * `quota_end_time` - Quota effective end time.
  * `quota_start_time` - Quota effective start time.
  * `total_quota` - Total unblock quota.
  * `used_quota` - Used unblock quota.


