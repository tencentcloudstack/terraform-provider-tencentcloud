---
subcategory: "Oceanus"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_oceanus_savepoint_list"
sidebar_current: "docs-tencentcloud-datasource-oceanus_savepoint_list"
description: |-
  Use this data source to query detailed information of oceanus savepoint_list
---

# tencentcloud_oceanus_savepoint_list

Use this data source to query detailed information of oceanus savepoint_list

## Example Usage

```hcl
data "tencentcloud_oceanus_savepoint_list" "example" {
  job_id        = "cql-314rw6w0"
  work_space_id = "space-2idq8wbr"
}
```

## Argument Reference

The following arguments are supported:

* `job_id` - (Required, String) Job SerialId.
* `result_output_file` - (Optional, String) Used to save results.
* `work_space_id` - (Optional, String) Workspace SerialId.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `savepoint` - Snapshot listNote: This field may return null, indicating that no valid value was found.
  * `create_time` - Creation timeNote: This field may return null, indicating that no valid value was found.
  * `description` - DescriptionNote: This field may return null, indicating that no valid value was found.
  * `id` - Primary keyNote: This field may return null, indicating that no valid value was found.
  * `job_runtime_id` - Sequential ID of the running job instanceNote: This field may return null, indicating that no valid value was found.
  * `path_status` - Snapshot path status: 1=available; 2=unavailable;Note: This field may return null, indicating that no valid value was found.
  * `path` - PathNote: This field may return null, indicating that no valid value was found.
  * `record_type` - Snapshot type: 1=savepoint; 2=checkpoint; 3=cancelWithSavepointNote: This field may return null, indicating that no valid value was found.
  * `serial_id` - Snapshot SerialIdNote: This field may return null, indicating that no valid value was found.
  * `size` - SizeNote: This field may return null, indicating that no valid value was found.
  * `status` - Status: 1=Active; 2=Expired; 3=InProgress; 4=Failed; 5=TimeoutNote: This field may return null, indicating that no valid value was found.
  * `time_consuming` - Time consumptionNote: This field may return null, indicating that no valid value was found.
  * `timeout` - Fixed timeoutNote: This field may return null, indicating that no valid value was found.
  * `update_time` - Update timeNote: This field may return null, indicating that no valid value was found.
  * `version_id` - Version numberNote: This field may return null, indicating that no valid value was found.


