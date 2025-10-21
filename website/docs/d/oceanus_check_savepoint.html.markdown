---
subcategory: "Oceanus"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_oceanus_check_savepoint"
sidebar_current: "docs-tencentcloud-datasource-oceanus_check_savepoint"
description: |-
  Use this data source to query detailed information of oceanus check_savepoint
---

# tencentcloud_oceanus_check_savepoint

Use this data source to query detailed information of oceanus check_savepoint

## Example Usage

```hcl
data "tencentcloud_oceanus_check_savepoint" "example" {
  job_id         = "cql-314rw6w0"
  serial_id      = "svp-52xkpymp"
  record_type    = 1
  savepoint_path = "cosn://52xkpymp-12345/12345/10000/cql-12345/2/flink-savepoints/savepoint-000000-12334"
  work_space_id  = "space-2idq8wbr"
}
```

## Argument Reference

The following arguments are supported:

* `job_id` - (Required, String) Job id.
* `record_type` - (Required, Int) Snapshot type. 1:savepoint; 2:checkpoint; 3:cancelWithSavepoint.
* `savepoint_path` - (Required, String) Snapshot path, currently only supports COS path.
* `serial_id` - (Required, String) Snapshot resource ID.
* `work_space_id` - (Required, String) Workspace ID.
* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `savepoint_status` - 1=available, 2=unavailable.


