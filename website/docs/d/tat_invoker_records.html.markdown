---
subcategory: "TencentCloud Automation Tools(TAT)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_tat_invoker_records"
sidebar_current: "docs-tencentcloud-datasource-tat_invoker_records"
description: |-
  Use this data source to query detailed information of tat invoker_records
---

# tencentcloud_tat_invoker_records

Use this data source to query detailed information of tat invoker_records

## Example Usage

```hcl
data "tencentcloud_tat_invoker_records" "invoker_records" {
  invoker_ids = ["ivk-cas4upyf"]
}
```

## Argument Reference

The following arguments are supported:

* `invoker_ids` - (Optional, Set: [`String`]) List of invoker IDs. Up to 100 IDs are allowed.
* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `invoker_record_set` - Execution history of an invoker.
  * `invocation_id` - Command execution ID.
  * `invoke_time` - Execution time.
  * `invoker_id` - Invoker ID.
  * `reason` - Execution reason.
  * `result` - Trigger result.


