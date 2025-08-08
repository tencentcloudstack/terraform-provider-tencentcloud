---
subcategory: "Data Lake Compute(DLC)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_dlc_describe_data_engine_events"
sidebar_current: "docs-tencentcloud-datasource-dlc_describe_data_engine_events"
description: |-
  Use this data source to query detailed information of DLC describe data engine events
---

# tencentcloud_dlc_describe_data_engine_events

Use this data source to query detailed information of DLC describe data engine events

## Example Usage

```hcl
data "tencentcloud_dlc_describe_data_engine_events" "example" {
  data_engine_name = "tf-example"
}
```

## Argument Reference

The following arguments are supported:

* `data_engine_name` - (Required, String) Data engine name.
* `result_output_file` - (Optional, String) Used to save results.
* `session_id` - (Optional, String) Session ID, standard spark engine required.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `events` - Event details.
  * `cluster_info` - Cluster information.
  * `events_action` - Event action.
  * `time` - Event time.


