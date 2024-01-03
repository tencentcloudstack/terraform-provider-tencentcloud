---
subcategory: "Oceanus"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_oceanus_job_events"
sidebar_current: "docs-tencentcloud-datasource-oceanus_job_events"
description: |-
  Use this data source to query detailed information of oceanus job_events
---

# tencentcloud_oceanus_job_events

Use this data source to query detailed information of oceanus job_events

## Example Usage

```hcl
data "tencentcloud_oceanus_job_events" "example" {
  job_id          = "cql-6w8eab6f"
  start_timestamp = 1630932161
  end_timestamp   = 1631232466
  types           = ["1", "2"]
  work_space_id   = "space-6w8eab6f"
}
```

## Argument Reference

The following arguments are supported:

* `end_timestamp` - (Required, Int) Filter condition:End Unix timestamp (seconds).
* `job_id` - (Required, String) Job ID.
* `start_timestamp` - (Required, Int) Filter condition:Start Unix timestamp (seconds).
* `work_space_id` - (Required, String) Workspace SerialId.
* `result_output_file` - (Optional, String) Used to save results.
* `types` - (Optional, Set: [`String`]) Event types. If not passed, data of all types will be returned.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `events` - List of events within the specified range for this jobNote: This field may return null, indicating that no valid values can be obtained.
  * `description` - Description text of the event type.
  * `message` - Some optional explanations of the eventNote: This field may return null, indicating that no valid values can be obtained.
  * `running_order_id` - Running ID when the event occurredNote: This field may return null, indicating that no valid values can be obtained.
  * `solution_link` - Troubleshooting manual link for the abnormal eventNote: This field may return null, indicating that no valid values can be obtained.
  * `timestamp` - Unix timestamp (seconds) when the event occurred.
  * `type` - Internally defined event type.
* `running_order_ids` - Array of running instance IDs.


