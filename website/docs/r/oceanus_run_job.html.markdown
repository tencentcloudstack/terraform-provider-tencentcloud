---
subcategory: "Oceanus"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_oceanus_run_job"
sidebar_current: "docs-tencentcloud-resource-oceanus_run_job"
description: |-
  Provides a resource to create a oceanus run_job
---

# tencentcloud_oceanus_run_job

Provides a resource to create a oceanus run_job

## Example Usage

```hcl
resource "tencentcloud_oceanus_run_job" "example" {
  run_job_descriptions {
    job_id                   = "cql-4xwincyn"
    run_type                 = 1
    start_mode               = "LATEST"
    job_config_version       = 10
    use_old_system_connector = false
  }
  work_space_id = "space-2idq8wbr"
}
```

## Argument Reference

The following arguments are supported:

* `run_job_descriptions` - (Required, List, ForceNew) The description information for batch job startup.
* `work_space_id` - (Optional, String, ForceNew) Workspace SerialId.

The `run_job_descriptions` object supports the following:

* `job_id` - (Required, String) Job ID.
* `run_type` - (Required, Int) The type of the run. 1 indicates start, and 2 indicates resume.
* `custom_timestamp` - (Optional, Int) Custom timestamp.
* `job_config_version` - (Optional, Int) A certain version of the current job(Not passed by default as a non-draft job version).
* `savepoint_id` - (Optional, String) Savepoint ID.
* `savepoint_path` - (Optional, String) Savepoint path.
* `start_mode` - (Optional, String) Compatible with the startup parameters of the old SQL type job: specify the start time point of data source consumption (recommended to pass the value)Ensure that the parameter is LATEST, EARLIEST, T+Timestamp (example: T1557394288000).
* `use_old_system_connector` - (Optional, Bool) Use the historical version of the system dependency.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



