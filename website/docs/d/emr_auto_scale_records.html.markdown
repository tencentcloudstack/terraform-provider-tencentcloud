---
subcategory: "MapReduce(EMR)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_emr_auto_scale_records"
sidebar_current: "docs-tencentcloud-datasource-emr_auto_scale_records"
description: |-
  Use this data source to query detailed information of emr auto_scale_records
---

# tencentcloud_emr_auto_scale_records

Use this data source to query detailed information of emr auto_scale_records

## Example Usage

```hcl
data "tencentcloud_emr_auto_scale_records" "auto_scale_records" {
  instance_id = "emr-bpum4pad"
  filters {
    key   = "StartTime"
    value = "2006-01-02 15:04:05"
  }
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, String) EMR cluster ID.
* `filters` - (Optional, List) Record filtering parameters, currently only `StartTime`, `EndTime` and `StrategyName` are supported. `StartTime` and `EndTime` support the time format of 2006-01-02 15:04:05 or 2006/01/02 15:04:05.
* `result_output_file` - (Optional, String) Used to save results.

The `filters` object supports the following:

* `key` - (Required, String) Key. Note: This field may return null, indicating that no valid value can be obtained.
* `value` - (Required, String) Value. Note: This field may return null, indicating that no valid value can be obtained.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `record_list` - Record list.
  * `action_status` - `SUCCESS`, `FAILED`, `PART_SUCCESS`, `IN_PROCESS`.
  * `action_time` - Process Trigger Time.
  * `compensate_count` - Compensation Times Note: This field may return null, indicating that no valid value can be obtained.
  * `compensate_flag` - Compensation and expansion, 0 represents no start, 1 represents start. Note: This field may return null, indicating that no valid value can be obtained.
  * `end_time` - Process End Time.
  * `expect_scale_num` - Effective only when ScaleAction is SCALE_OUT.
  * `scale_action` - `SCALE_OUT` and `SCALE_IN` respectively represent expanding and shrinking capacity.
  * `scale_info` - Scalability-related Description.
  * `spec_info` - Specification information used when expanding capacity.
  * `strategy_name` - Rule name of expanding and shrinking capacity.
  * `strategy_type` - Strategy Type, 1 for Load scaling, 2 for Time scaling.


